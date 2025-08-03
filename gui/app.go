package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/YamaXanadu830/clash-speedtest/speedtester"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 应用结构体
type App struct {
	ctx         context.Context
	speedTester *speedtester.SpeedTester
	isRunning   bool
}

// NewApp 创建新的应用实例
func NewApp() *App {
	return &App{}
}

// OnStartup 应用启动时调用
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
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
	
	go func() {
		defer func() {
			a.isRunning = false
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

		proxies, err := a.speedTester.LoadProxies(false)
		if err != nil {
			runtime.EventsEmit(a.ctx, "test-error", fmt.Sprintf("加载代理失败: %v", err))
			return
		}

		// 发送测试开始事件
		runtime.EventsEmit(a.ctx, "test-start", map[string]interface{}{
			"total": len(proxies),
		})

		var results []*speedtester.Result
		
		a.speedTester.TestProxies(proxies, func(result *speedtester.Result) {
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
			
			runtime.EventsEmit(a.ctx, "test-progress", map[string]interface{}{
				"current": len(results),
				"total":   len(proxies),
				"result":  guiResult,
			})
		})

		// 发送测试完成事件
		runtime.EventsEmit(a.ctx, "test-complete", map[string]interface{}{
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

	go func() {
		defer func() {
			a.isRunning = false
		}()

		proxies, err := a.speedTester.LoadProxies(false)
		if err != nil {
			runtime.EventsEmit(a.ctx, "monitor-error", fmt.Sprintf("加载代理失败: %v", err))
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
		runtime.EventsEmit(a.ctx, "monitor-start", map[string]interface{}{
			"total":    len(proxies),
			"duration": config.Duration,
			"type":     config.Type,
		})

		statusMap := make(map[string]*speedtester.MonitorStatus)
		
		results := a.speedTester.MonitorProxies(proxies, monitorConfig, func(status *speedtester.MonitorStatus) {
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
			
			runtime.EventsEmit(a.ctx, "monitor-update", guiStatus)
		})

		// 发送监控完成事件
		runtime.EventsEmit(a.ctx, "monitor-complete", map[string]interface{}{
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
	
	// 这里可以添加停止逻辑
	a.isRunning = false
	runtime.EventsEmit(a.ctx, "test-stopped", nil)
	
	return nil
}

// GetSystemInfo 获取系统信息
func (a *App) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Platform: "unknown", // 可以使用runtime包获取实际平台信息
		Version:  "1.0.0",
		BuildTime: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// SelectConfigFile 选择配置文件
func (a *App) SelectConfigFile() (string, error) {
	filters := []runtime.FileFilter{
		{
			DisplayName: "YAML配置文件 (*.yaml, *.yml)",
			Pattern:     "*.yaml;*.yml",
		},
		{
			DisplayName: "所有文件 (*.*)",
			Pattern:     "*.*",
		},
	}
	
	options := runtime.OpenDialogOptions{
		Title:   "选择Clash配置文件",
		Filters: filters,
	}
	
	selectedFile, err := runtime.OpenFileDialog(a.ctx, options)
	return selectedFile, err
}

// SaveReport 保存测试报告
func (a *App) SaveReport(results interface{}, format string) error {
	filters := []runtime.FileFilter{
		{
			DisplayName: "YAML文件 (*.yaml)",
			Pattern:     "*.yaml",
		},
		{
			DisplayName: "JSON文件 (*.json)",
			Pattern:     "*.json",
		},
	}
	
	options := runtime.SaveDialogOptions{
		Title:          "保存测试报告",
		Filters:        filters,
		DefaultFilename: fmt.Sprintf("speedtest-report-%s", time.Now().Format("20060102-150405")),
	}
	
	selectedFile, err := runtime.SaveFileDialog(a.ctx, options)
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