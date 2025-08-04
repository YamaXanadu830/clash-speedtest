package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/YamaXanadu830/clash-speedtest/speedtester"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 应用结构体
type App struct {
	ctx            context.Context
	speedTester    *speedtester.SpeedTester
	isRunning      bool
	cancelTest     context.CancelFunc
	historyManager *HistoryManager
}

// NewApp 创建新的应用实例
func NewApp() *App {
	return &App{}
}

// OnStartup 应用启动时调用
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
	
	// 初始化历史数据管理器
	homeDir, _ := os.UserHomeDir()
	dataFile := filepath.Join(homeDir, ".clash-speedtest", "history.json")
	a.historyManager = NewHistoryManager(dataFile)
}

// OnDomReady DOM准备好时调用
func (a *App) OnDomReady(ctx context.Context) {
	// 这里可以添加DOM准备好后的初始化逻辑
}

// OnBeforeClose 应用关闭前调用
func (a *App) OnBeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// OnShutdown 应用关闭时调用
func (a *App) OnShutdown(ctx context.Context) {
	// 清理资源
}

// LoadConfig 加载配置文件
func (a *App) LoadConfig(configPath string, filterRegex string, blockKeywords string) (*ConfigInfo, error) {
	config := &speedtester.Config{
		ConfigPaths:      configPath,
		FilterRegex:      filterRegex,
		BlockRegex:       blockKeywords,
		ServerURL:        "https://speed.cloudflare.com",
		DownloadSize:     50 * 1024 * 1024, // 50MB
		UploadSize:       20 * 1024 * 1024, // 20MB
		Timeout:          5 * time.Second,
		Concurrent:       4,
		MaxLatency:       800 * time.Millisecond,
		MinDownloadSpeed: 5 * 1024 * 1024,   // 5MB/s
		MinUploadSpeed:   2 * 1024 * 1024,   // 2MB/s
		FastMode:         false,
	}

	a.speedTester = speedtester.New(config)
	
	proxies, err := a.speedTester.LoadProxies(false)
	if err != nil {
		return nil, fmt.Errorf("加载代理配置失败: %v", err)
	}

	return &ConfigInfo{
		ProxyCount: len(proxies),
		ConfigPath: configPath,
		Filter:     filterRegex,
		Block:      blockKeywords,
	}, nil
}

// StartSpeedTest 开始速度测试
func (a *App) StartSpeedTest(config TestConfig) error {
	if a.speedTester == nil {
		return fmt.Errorf("请先加载配置文件")
	}

	if a.isRunning {
		return fmt.Errorf("测试已在进行中")
	}

	a.isRunning = true
	
	// 创建可取消的上下文
	testCtx, cancel := context.WithCancel(a.ctx)
	a.cancelTest = cancel
	
	go func() {
		defer func() {
			// 只有在正常完成时才重置状态（StopTest已经处理了取消的情况）
			if a.isRunning {
				a.isRunning = false
				a.cancelTest = nil
			}
		}()

		// 更新配置
		a.speedTester = speedtester.New(&speedtester.Config{
			ConfigPaths:      config.ConfigPath,
			FilterRegex:      config.FilterRegex,
			BlockRegex:       config.BlockRegex,
			ServerURL:        config.ServerURL,
			DownloadSize:     config.DownloadSize,
			UploadSize:       config.UploadSize,
			Timeout:          time.Duration(config.Timeout) * time.Second,
			Concurrent:       config.Concurrent,
			MaxLatency:       time.Duration(config.MaxLatency) * time.Millisecond,
			MinDownloadSpeed: config.MinDownloadSpeed * 1024 * 1024,
			MinUploadSpeed:   config.MinUploadSpeed * 1024 * 1024,
			FastMode:         config.FastMode,
		})

		// 检查是否已被取消
		select {
		case <-testCtx.Done():
			wailsruntime.EventsEmit(a.ctx, "test-stopped", nil)
			return
		default:
		}

		proxies, err := a.speedTester.LoadProxies(false)
		if err != nil {
			wailsruntime.EventsEmit(a.ctx, "test-error", fmt.Sprintf("加载代理失败: %v", err))
			return
		}

		// 再次检查是否已被取消
		select {
		case <-testCtx.Done():
			wailsruntime.EventsEmit(a.ctx, "test-stopped", nil)
			return
		default:
		}

		// 发送测试开始事件
		wailsruntime.EventsEmit(a.ctx, "test-start", map[string]interface{}{
			"total": len(proxies),
		})

		var results []*speedtester.Result
		
		a.speedTester.TestProxies(proxies, func(result *speedtester.Result) {
			// 在每个结果回调中检查取消状态
			select {
			case <-testCtx.Done():
				return
			default:
			}
			
			results = append(results, result)
			
			// 转换结果格式并发送进度事件
			guiResult := &TestResult{
				ProxyName:     result.ProxyName,
				ProxyType:     result.ProxyType,
				Latency:       int64(result.Latency / time.Millisecond),
				Jitter:        int64(result.Jitter / time.Millisecond),
				PacketLoss:    result.PacketLoss,
				DownloadSpeed: result.DownloadSpeed,
				UploadSpeed:   result.UploadSpeed,
				DownloadSize:  result.DownloadSize,
				UploadSize:    result.UploadSize,
				DownloadTime:  int64(result.DownloadTime / time.Millisecond),
				UploadTime:    int64(result.UploadTime / time.Millisecond),
				Status:        "完成",
			}
			
			// 异步保存到历史记录，避免阻塞测试进程
			if a.historyManager != nil {
				go func(result *TestResult) {
					if err := a.historyManager.AddResult(result); err != nil {
						fmt.Printf("保存测试结果到历史记录失败: %v\n", err)
					}
				}(guiResult)
			}
			
			wailsruntime.EventsEmit(a.ctx, "test-progress", map[string]interface{}{
				"current": len(results),
				"total":   len(proxies),
				"result":  guiResult,
			})
		})

		// 最后检查是否已被取消
		select {
		case <-testCtx.Done():
			wailsruntime.EventsEmit(a.ctx, "test-stopped", nil)
			return
		default:
		}

		// 发送测试完成事件
		wailsruntime.EventsEmit(a.ctx, "test-complete", map[string]interface{}{
			"results": results,
			"total":   len(results),
		})
	}()

	return nil
}

// StartMonitor 开始监控模式
func (a *App) StartMonitor(config MonitorConfig) error {
	if a.speedTester == nil {
		return fmt.Errorf("请先加载配置文件")
	}

	if a.isRunning {
		return fmt.Errorf("监控已在进行中")
	}

	a.isRunning = true

	// 创建可取消的上下文
	monitorCtx, cancel := context.WithCancel(a.ctx)
	a.cancelTest = cancel

	go func() {
		defer func() {
			// 只有在正常完成时才重置状态（StopTest已经处理了取消的情况）
			if a.isRunning {
				a.isRunning = false
				a.cancelTest = nil
			}
		}()

		// 检查是否已被取消
		select {
		case <-monitorCtx.Done():
			wailsruntime.EventsEmit(a.ctx, "test-stopped", nil)
			return
		default:
		}

		proxies, err := a.speedTester.LoadProxies(false)
		if err != nil {
			wailsruntime.EventsEmit(a.ctx, "monitor-error", fmt.Sprintf("加载代理失败: %v", err))
			return
		}

		// 根据监控类型设置目标URL
		targetURL := config.TargetURL
		if targetURL == "" {
			if config.Type == "websocket" {
				targetURL = "wss://stream.binance.com:9443/ws/btcusdt@ticker"
			} else {
				targetURL = "https://cn.tradingview.com/chart/"
			}
		}

		monitorConfig := &speedtester.MonitorConfig{
			Duration:  time.Duration(config.Duration) * time.Second,
			Interval:  time.Duration(config.Interval) * time.Second,
			TargetURL: targetURL,
			Type:      config.Type,
		}

		// 发送监控开始事件
		wailsruntime.EventsEmit(a.ctx, "monitor-start", map[string]interface{}{
			"total":    len(proxies),
			"duration": config.Duration,
			"type":     config.Type,
		})

		statusMap := make(map[string]*speedtester.MonitorStatus)
		
		results := a.speedTester.MonitorProxiesWithContext(monitorCtx, proxies, monitorConfig, func(status *speedtester.MonitorStatus) {
			statusMap[status.ProxyName] = status
			
			// 转换状态格式并发送更新事件
			guiStatus := &MonitorStatus{
				ProxyName:       status.ProxyName,
				IsAlive:         status.IsAlive,
				OnlineDuration:  int64(status.OnlineDuration / time.Second),
				TotalDuration:   int64(status.TotalDuration / time.Second),
				DisconnectCount: status.DisconnectCount,
				DataPacketCount: status.DataPacketCount,
				TotalDataBytes:  status.TotalDataBytes,
				LastPacketTime:  status.LastPacketTime.Format("15:04:05"),
				LastUpdate:      time.Now().Format("15:04:05"),
			}
			
			wailsruntime.EventsEmit(a.ctx, "monitor-update", guiStatus)
		})

		// 最后检查是否已被取消
		select {
		case <-monitorCtx.Done():
			wailsruntime.EventsEmit(a.ctx, "test-stopped", nil)
			return
		default:
		}

		// 发送监控完成事件
		wailsruntime.EventsEmit(a.ctx, "monitor-complete", map[string]interface{}{
			"results": results,
			"total":   len(results),
		})
	}()

	return nil
}

// StopTest 停止当前测试
func (a *App) StopTest() error {
	if !a.isRunning {
		return fmt.Errorf("没有正在运行的测试")
	}
	
	// 调用取消函数停止goroutine
	if a.cancelTest != nil {
		a.cancelTest()
		// 立即重置运行状态，确保可以立即开始新的测试
		a.isRunning = false
		a.cancelTest = nil
	}
	
	return nil
}

// GetSystemInfo 获取系统信息
func (a *App) GetSystemInfo() *SystemInfo {
	// 获取操作系统和架构信息
	platform := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	
	return &SystemInfo{
		Platform:  platform,
		Version:   "1.0.0",
		BuildTime: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// GetAnalysisStats 获取分析统计数据
func (a *App) GetAnalysisStats() (*AnalysisStats, error) {
	if a.historyManager == nil {
		return nil, fmt.Errorf("历史数据管理器未初始化")
	}
	
	return a.historyManager.GetStats(), nil
}

// ClearHistory 清空历史数据
func (a *App) ClearHistory() error {
	if a.historyManager == nil {
		return fmt.Errorf("历史数据管理器未初始化")
	}
	
	return a.historyManager.ClearHistory()
}

// SelectConfigFile 选择配置文件
func (a *App) SelectConfigFile() (string, error) {
	filters := []wailsruntime.FileFilter{
		{
			DisplayName: "YAML配置文件 (*.yaml, *.yml)",
			Pattern:     "*.yaml;*.yml",
		},
		{
			DisplayName: "所有文件 (*.*)",
			Pattern:     "*.*",
		},
	}
	
	options := wailsruntime.OpenDialogOptions{
		Title:   "选择Clash配置文件",
		Filters: filters,
	}
	
	selectedFile, err := wailsruntime.OpenFileDialog(a.ctx, options)
	return selectedFile, err
}

// SaveReport 保存测试报告
func (a *App) SaveReport(results interface{}, format string) error {
	filters := []wailsruntime.FileFilter{
		{
			DisplayName: "YAML文件 (*.yaml)",
			Pattern:     "*.yaml",
		},
		{
			DisplayName: "JSON文件 (*.json)",
			Pattern:     "*.json",
		},
	}
	
	options := wailsruntime.SaveDialogOptions{
		Title:          "保存测试报告",
		Filters:        filters,
		DefaultFilename: fmt.Sprintf("speedtest-report-%s", time.Now().Format("20060102-150405")),
	}
	
	selectedFile, err := wailsruntime.SaveFileDialog(a.ctx, options)
	if err != nil {
		return err
	}
	
	if selectedFile == "" {
		return nil // 用户取消了保存
	}
	
	if format == "json" {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return err
		}
		return os.WriteFile(selectedFile, data, 0644)
	} else {
		// 默认YAML格式，这里需要使用yaml包
		return fmt.Errorf("YAML格式保存功能待实现")
	}
}

// 数据结构定义

// HistoryManager 历史数据管理器
type HistoryManager struct {
	dataFile string
	history  []HistoricalTestResult
	mutex    sync.RWMutex
}

// HistoricalTestResult 历史测试结果（包含时间戳）
type HistoricalTestResult struct {
	TestResult
	Timestamp time.Time `json:"timestamp"`
}

// AnalysisStats 分析统计数据
type AnalysisStats struct {
	TotalTests  int     `json:"totalTests"`
	TotalNodes  int     `json:"totalNodes"`
	AvgLatency  float64 `json:"avgLatency"`
	AvgSpeed    float64 `json:"avgSpeed"`
	LastUpdated string  `json:"lastUpdated"`
}

type ConfigInfo struct {
	ProxyCount int    `json:"proxyCount"`
	ConfigPath string `json:"configPath"`
	Filter     string `json:"filter"`
	Block      string `json:"block"`
}

type TestConfig struct {
	ConfigPath       string  `json:"configPath"`
	FilterRegex      string  `json:"filterRegex"`
	BlockRegex       string  `json:"blockRegex"`
	ServerURL        string  `json:"serverURL"`
	DownloadSize     int     `json:"downloadSize"`
	UploadSize       int     `json:"uploadSize"`
	Timeout          int     `json:"timeout"`
	Concurrent       int     `json:"concurrent"`
	MaxLatency       int     `json:"maxLatency"`
	MinDownloadSpeed float64 `json:"minDownloadSpeed"`
	MinUploadSpeed   float64 `json:"minUploadSpeed"`
	FastMode         bool    `json:"fastMode"`
}

type TestResult struct {
	ProxyName     string  `json:"proxyName"`
	ProxyType     string  `json:"proxyType"`
	Latency       int64   `json:"latency"`       // 毫秒
	Jitter        int64   `json:"jitter"`        // 毫秒
	PacketLoss    float64 `json:"packetLoss"`    // 百分比
	DownloadSpeed float64 `json:"downloadSpeed"` // 字节/秒
	UploadSpeed   float64 `json:"uploadSpeed"`   // 字节/秒
	DownloadSize  float64 `json:"downloadSize"`  // 字节
	UploadSize    float64 `json:"uploadSize"`    // 字节
	DownloadTime  int64   `json:"downloadTime"`  // 毫秒
	UploadTime    int64   `json:"uploadTime"`    // 毫秒
	Status        string  `json:"status"`
}

type MonitorConfig struct {
	Duration  int    `json:"duration"`  // 秒
	Interval  int    `json:"interval"`  // 秒
	TargetURL string `json:"targetURL"`
	Type      string `json:"type"` // "http" 或 "websocket"
}

type MonitorStatus struct {
	ProxyName       string `json:"proxyName"`
	IsAlive         bool   `json:"isAlive"`
	OnlineDuration  int64  `json:"onlineDuration"`  // 秒
	TotalDuration   int64  `json:"totalDuration"`   // 秒
	DisconnectCount int    `json:"disconnectCount"`
	DataPacketCount int64  `json:"dataPacketCount"`
	TotalDataBytes  int64  `json:"totalDataBytes"`
	LastPacketTime  string `json:"lastPacketTime"`
	LastUpdate      string `json:"lastUpdate"`
}

type SystemInfo struct {
	Platform  string `json:"platform"`
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
}

// HistoryManager 方法定义

// NewHistoryManager 创建新的历史数据管理器
func NewHistoryManager(dataFile string) *HistoryManager {
	hm := &HistoryManager{
		dataFile: dataFile,
		history:  make([]HistoricalTestResult, 0),
	}
	
	// 确保目录存在
	dir := filepath.Dir(dataFile)
	os.MkdirAll(dir, 0755)
	
	// 加载历史数据
	hm.loadHistory()
	
	return hm
}

// loadHistory 从文件加载历史数据
func (hm *HistoryManager) loadHistory() {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	
	data, err := os.ReadFile(hm.dataFile)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("加载历史数据失败: %v\n", err)
		}
		return
	}
	
	if err := json.Unmarshal(data, &hm.history); err != nil {
		fmt.Printf("解析历史数据失败: %v\n", err)
		hm.history = make([]HistoricalTestResult, 0)
	}
}

// saveHistory 保存历史数据到文件
func (hm *HistoryManager) saveHistory() error {
	// 注意：这个方法是在AddResult中调用的，AddResult已经获取了写锁
	// 所以这里不需要再获取锁
	
	data, err := json.MarshalIndent(hm.history, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化历史数据失败: %v", err)
	}
	
	return os.WriteFile(hm.dataFile, data, 0644)
}

// AddResult 添加测试结果到历史记录
func (hm *HistoryManager) AddResult(result *TestResult) error {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	
	historicalResult := HistoricalTestResult{
		TestResult: *result,
		Timestamp:  time.Now(),
	}
	
	hm.history = append(hm.history, historicalResult)
	
	return hm.saveHistory()
}

// GetStats 计算并返回统计数据
func (hm *HistoryManager) GetStats() *AnalysisStats {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	
	if len(hm.history) == 0 {
		return &AnalysisStats{
			TotalTests:  0,
			TotalNodes:  0,
			AvgLatency:  0,
			AvgSpeed:    0,
			LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
		}
	}
	
	// 统计总测试次数
	totalTests := len(hm.history)
	
	// 统计不同节点数
	nodeSet := make(map[string]bool)
	var totalLatency int64
	var totalSpeed float64
	validLatencyCount := 0
	validSpeedCount := 0
	
	for _, result := range hm.history {
		nodeSet[result.ProxyName] = true
		
		if result.Latency > 0 {
			totalLatency += result.Latency
			validLatencyCount++
		}
		
		if result.DownloadSpeed > 0 {
			totalSpeed += result.DownloadSpeed
			validSpeedCount++
		}
	}
	
	var avgLatency float64
	var avgSpeed float64
	
	if validLatencyCount > 0 {
		avgLatency = float64(totalLatency) / float64(validLatencyCount)
	}
	
	if validSpeedCount > 0 {
		avgSpeed = totalSpeed / float64(validSpeedCount) / 1024 / 1024 // 转换为MB/s
	}
	
	return &AnalysisStats{
		TotalTests:  totalTests,
		TotalNodes:  len(nodeSet),
		AvgLatency:  avgLatency,
		AvgSpeed:    avgSpeed,
		LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// ClearHistory 清空历史数据
func (hm *HistoryManager) ClearHistory() error {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	
	hm.history = make([]HistoricalTestResult, 0)
	return hm.saveHistory()
}