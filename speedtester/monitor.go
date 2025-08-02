package speedtester

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/metacubex/mihomo/constant"
)

// MonitorConfig 监控配置
type MonitorConfig struct {
	Duration  time.Duration // 监控时长
	Interval  time.Duration // 心跳间隔
	TargetURL string        // 监控目标URL
	Type      string        // 监控类型: http 或 websocket
}

// MonitorResult 监控结果
type MonitorResult struct {
	ProxyName        string          // 代理名称
	ProxyType        string          // 代理类型
	StartTime        time.Time       // 开始时间
	EndTime          time.Time       // 结束时间
	TotalDuration    time.Duration   // 总监控时长
	OnlineDuration   time.Duration   // 在线时长
	DisconnectCount  int             // 断线次数
	DisconnectEvents []DisconnectEvent // 断线事件列表
	StabilityRate    float64         // 稳定率 (在线时长/总时长)
	MaxOnlineTime    time.Duration   // 最长连续在线时间
	IsAlive          bool            // 当前是否在线
}

// DisconnectEvent 断线事件
type DisconnectEvent struct {
	StartTime time.Time     // 断线开始时间
	EndTime   time.Time     // 断线结束时间（重连成功时间）
	Duration  time.Duration // 断线持续时间
	Error     string        // 断线原因
}

// MonitorSession 监控会话
type MonitorSession struct {
	proxy           constant.Proxy
	proxyName       string
	proxyType       string
	config          *MonitorConfig
	result          *MonitorResult
	client          *http.Client
	ctx             context.Context
	cancel          context.CancelFunc
	mu              sync.Mutex
	lastOnlineTime  time.Time
	currentOnlineStart time.Time
	wsConn          *websocket.Conn
	wsDialer        *websocket.Dialer
	// 数据包统计字段
	dataPacketCount    int64     // 接收到的数据包总数
	lastDataPacketTime time.Time // 最后一个数据包的时间
	totalDataBytes     int64     // 接收的总字节数
}

// MonitorProxies 监控代理节点稳定性
func (st *SpeedTester) MonitorProxies(proxies map[string]*CProxy, config *MonitorConfig, callback func(status *MonitorStatus)) []*MonitorResult {
	ctx, cancel := context.WithTimeout(context.Background(), config.Duration)
	defer cancel()

	var wg sync.WaitGroup
	results := make([]*MonitorResult, 0, len(proxies))
	resultChan := make(chan *MonitorResult, len(proxies))
	statusChan := make(chan *MonitorStatus, len(proxies)*10)

	// 启动状态收集器
	go func() {
		for status := range statusChan {
			callback(status)
		}
	}()

	// 为每个代理启动监控会话
	for name, proxy := range proxies {
		wg.Add(1)
		go func(proxyName string, p *CProxy) {
			defer wg.Done()
			session := &MonitorSession{
				proxy:     p.Proxy,
				proxyName: proxyName,
				proxyType: p.Type().String(),
				config:    config,
				client:    st.createClient(p.Proxy, 10*time.Second),
				ctx:       ctx,
			}
			result := session.Run(statusChan)
			resultChan <- result
		}(name, proxy)
	}

	// 等待所有监控完成
	wg.Wait()
	close(resultChan)
	close(statusChan)

	// 收集结果
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// Run 运行监控会话
func (s *MonitorSession) Run(statusChan chan<- *MonitorStatus) *MonitorResult {
	s.result = &MonitorResult{
		ProxyName:        s.proxyName,
		ProxyType:        s.proxyType,
		StartTime:        time.Now(),
		DisconnectEvents: make([]DisconnectEvent, 0),
		IsAlive:          false,
	}

	ticker := time.NewTicker(s.config.Interval)
	defer ticker.Stop()

	var currentDisconnect *DisconnectEvent
	var maxOnlineTime time.Duration

	for {
		select {
		case <-s.ctx.Done():
			// 监控时间到达
			s.result.EndTime = time.Now()
			s.result.TotalDuration = s.result.EndTime.Sub(s.result.StartTime)
			s.result.MaxOnlineTime = maxOnlineTime
			s.calculateStability()
			// 清理WebSocket连接
			if s.config.Type == "websocket" {
				s.closeWebSocket()
			}
			return s.result

		case <-ticker.C:
			// 执行心跳检测
			isAlive, err := s.heartbeat()
			
			s.mu.Lock()
			if isAlive {
				if !s.result.IsAlive {
					// 从断线恢复
					s.result.IsAlive = true
					s.currentOnlineStart = time.Now()
					if currentDisconnect != nil {
						currentDisconnect.EndTime = time.Now()
						currentDisconnect.Duration = currentDisconnect.EndTime.Sub(currentDisconnect.StartTime)
						s.result.DisconnectEvents = append(s.result.DisconnectEvents, *currentDisconnect)
						currentDisconnect = nil
					}
				}
				// 更新最长在线时间
				currentOnlineTime := time.Since(s.currentOnlineStart)
				if currentOnlineTime > maxOnlineTime {
					maxOnlineTime = currentOnlineTime
				}
			} else {
				if s.result.IsAlive {
					// 新的断线事件
					s.result.IsAlive = false
					s.result.DisconnectCount++
					currentDisconnect = &DisconnectEvent{
						StartTime: time.Now(),
						Error:     err.Error(),
					}
					// 累计之前的在线时间
					s.result.OnlineDuration += time.Since(s.currentOnlineStart)
				}
			}
			s.mu.Unlock()

			// 发送状态更新
			s.mu.Lock()
			status := &MonitorStatus{
				ProxyName:       s.proxyName,
				IsAlive:         isAlive,
				DisconnectCount: s.result.DisconnectCount,
				OnlineDuration:  s.getOnlineDuration(),
				TotalDuration:   time.Since(s.result.StartTime),
				// WebSocket流式数据统计
				DataPacketCount: s.dataPacketCount,
				TotalDataBytes:  s.totalDataBytes,
				LastPacketTime:  s.lastDataPacketTime,
			}
			s.mu.Unlock()
			
			select {
			case statusChan <- status:
			default:
				// 防止阻塞
			}
		}
	}
}

// heartbeat 执行心跳检测
func (s *MonitorSession) heartbeat() (bool, error) {
	if s.config.Type == "websocket" {
		return s.heartbeatWebSocket()
	}
	return s.heartbeatHTTP()
}

// heartbeatHTTP HTTP方式心跳检测
func (s *MonitorSession) heartbeatHTTP() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", s.config.TargetURL, nil)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 读取一部分响应确保连接正常
	_, err = io.CopyN(io.Discard, resp.Body, 1024)
	if err != nil && err != io.EOF {
		return false, err
	}

	return resp.StatusCode < 500, nil
}

// heartbeatWebSocket WebSocket流式数据监控
func (s *MonitorSession) heartbeatWebSocket() (bool, error) {
	// 如果连接不存在或已断开，尝试建立新连接
	if s.wsConn == nil {
		err := s.connectWebSocket()
		if err != nil {
			return false, err
		}
	}

	// 设置读取超时 - 10秒内必须收到数据包，否则视为断线
	s.wsConn.SetReadDeadline(time.Now().Add(10 * time.Second))
	
	// 读取实时数据包（Binance价格数据）
	_, message, err := s.wsConn.ReadMessage()
	if err != nil {
		// 读取错误或超时，连接已断开
		s.closeWebSocket()
		return false, err
	}

	// 记录数据包接收统计
	s.recordDataPacket(message)
	
	return true, nil
}

// connectWebSocket 建立WebSocket连接
func (s *MonitorSession) connectWebSocket() error {
	if s.wsDialer == nil {
		s.wsDialer = &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 10 * time.Second,
		}
	}

	// 使用代理的HTTP客户端传输
	if s.client.Transport != nil {
		s.wsDialer.NetDialContext = s.client.Transport.(*http.Transport).DialContext
	}

	// 使用Binance公开WebSocket API - 实时BTC/USDT价格数据流
	wsURL := "wss://stream.binance.com:9443/ws/btcusdt@ticker"

	// 添加适当的headers
	headers := http.Header{}
	headers.Set("User-Agent", "clash-speedtest/1.0")

	conn, _, err := s.wsDialer.Dial(wsURL, headers)
	if err != nil {
		return err
	}

	s.wsConn = conn
	return nil
}

// closeWebSocket 关闭WebSocket连接
func (s *MonitorSession) closeWebSocket() {
	if s.wsConn != nil {
		s.wsConn.Close()
		s.wsConn = nil
	}
}

// recordDataPacket 记录数据包接收统计
func (s *MonitorSession) recordDataPacket(message []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	now := time.Now()
	s.dataPacketCount++
	s.lastDataPacketTime = now
	s.totalDataBytes += int64(len(message))
	
	// 可选：记录详细的数据包信息用于调试
	// 注意：Binance返回的是JSON格式的价格数据，格式如：
	// {"e":"24hrTicker","E":1640995200000,"s":"BTCUSDT","p":"1000.00",...}
}

// getOnlineDuration 获取当前总在线时长
func (s *MonitorSession) getOnlineDuration() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	onlineDuration := s.result.OnlineDuration
	if s.result.IsAlive {
		onlineDuration += time.Since(s.currentOnlineStart)
	}
	return onlineDuration
}

// calculateStability 计算稳定率
func (s *MonitorSession) calculateStability() {
	// 如果当前还在线，需要加上最后一段在线时间
	if s.result.IsAlive {
		s.result.OnlineDuration += time.Since(s.currentOnlineStart)
	}
	
	if s.result.TotalDuration > 0 {
		s.result.StabilityRate = float64(s.result.OnlineDuration) / float64(s.result.TotalDuration) * 100
	}
}

// MonitorStatus 监控状态（用于实时更新）
type MonitorStatus struct {
	ProxyName       string
	IsAlive         bool
	DisconnectCount int
	OnlineDuration  time.Duration
	TotalDuration   time.Duration
	// WebSocket流式数据统计
	DataPacketCount int64     // 数据包接收总数
	TotalDataBytes  int64     // 接收总字节数
	LastPacketTime  time.Time // 最后数据包时间
}

// FormatStabilityRate 格式化稳定率
func (r *MonitorResult) FormatStabilityRate() string {
	return fmt.Sprintf("%.2f%%", r.StabilityRate)
}

// FormatOnlineDuration 格式化在线时长
func (r *MonitorResult) FormatOnlineDuration() string {
	hours := int(r.OnlineDuration.Hours())
	minutes := int(r.OnlineDuration.Minutes()) % 60
	seconds := int(r.OnlineDuration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}