package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/YamaXanadu830/clash-speedtest/speedtester"
	"github.com/metacubex/mihomo/log"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v3"
)

var (
	configPathsConfig = flag.String("c", "", "config file path, also support http(s) url")
	filterRegexConfig = flag.String("f", ".+", "filter proxies by name, use regexp")
	blockKeywords     = flag.String("b", "", "block proxies by keywords, use | to separate multiple keywords (example: -b 'rate|x1|1x')")
	serverURL         = flag.String("server-url", "https://speed.cloudflare.com", "server url")
	downloadSize      = flag.Int("download-size", 50*1024*1024, "download size for testing proxies")
	uploadSize        = flag.Int("upload-size", 20*1024*1024, "upload size for testing proxies")
	timeout           = flag.Duration("timeout", time.Second*5, "timeout for testing proxies")
	concurrent        = flag.Int("concurrent", 4, "download concurrent size")
	outputPath        = flag.String("output", "", "output config file path")
	stashCompatible   = flag.Bool("stash-compatible", false, "enable stash compatible mode")
	maxLatency        = flag.Duration("max-latency", 800*time.Millisecond, "filter latency greater than this value")
	minDownloadSpeed  = flag.Float64("min-download-speed", 5, "filter download speed less than this value(unit: MB/s)")
	minUploadSpeed    = flag.Float64("min-upload-speed", 2, "filter upload speed less than this value(unit: MB/s)")
	renameNodes       = flag.Bool("rename", false, "rename nodes with IP location and speed")
	fastMode          = flag.Bool("fast", false, "fast mode, only test latency")
	// 监控模式参数
	monitorMode     = flag.Bool("monitor", false, "enable stability monitoring mode")
	monitorDuration = flag.Duration("monitor-duration", 24*time.Hour, "monitoring duration (default: 24h)")
	monitorInterval = flag.Duration("monitor-interval", time.Second, "heartbeat interval for monitoring (default: 1s)")
	monitorType     = flag.String("monitor-type", "http", "monitoring type: http or websocket (default: http)")
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

func main() {
	flag.Parse()
	log.SetLevel(log.SILENT)

	if *configPathsConfig == "" {
		log.Fatalln("please specify the configuration file")
	}

	// 监控模式和速度测试模式互斥检查
	if *monitorMode && *fastMode {
		log.Fatalln("--monitor and --fast modes cannot be used together")
	}

	speedTester := speedtester.New(&speedtester.Config{
		ConfigPaths:      *configPathsConfig,
		FilterRegex:      *filterRegexConfig,
		BlockRegex:       *blockKeywords,
		ServerURL:        *serverURL,
		DownloadSize:     *downloadSize,
		UploadSize:       *uploadSize,
		Timeout:          *timeout,
		Concurrent:       *concurrent,
		MaxLatency:       *maxLatency,
		MinDownloadSpeed: *minDownloadSpeed * 1024 * 1024,
		MinUploadSpeed:   *minUploadSpeed * 1024 * 1024,
		FastMode:         *fastMode,
	})

	allProxies, err := speedTester.LoadProxies(*stashCompatible)
	if err != nil {
		log.Fatalln("load proxies failed: %v", err)
	}

	if *monitorMode {
		// 监控模式
		fmt.Println("=== 连接稳定性监控模式 ===")
		if *monitorType == "websocket" {
			fmt.Printf("目标: Binance WebSocket API (实时BTC/USDT数据流)\n")
			fmt.Printf("端点: wss://stream.binance.com:9443/ws/btcusdt@ticker\n")
		} else {
			fmt.Printf("目标: https://cn.tradingview.com/chart/\n")
		}
		fmt.Printf("监控类型: %s\n", *monitorType)
		fmt.Printf("监控时长: %v | 心跳间隔: %v\n", *monitorDuration, *monitorInterval)
		fmt.Printf("监控节点数: %d\n\n", len(allProxies))

		// 创建监控配置
		monitorConfig := &speedtester.MonitorConfig{
			Duration:  *monitorDuration,
			Interval:  *monitorInterval,
			TargetURL: "https://cn.tradingview.com/chart/",
			Type:      *monitorType,
		}

		// 执行监控
		startTime := time.Now()
		lastUpdate := time.Now()
		statusMap := make(map[string]*speedtester.MonitorStatus)

		results := speedTester.MonitorProxies(allProxies, monitorConfig, func(status *speedtester.MonitorStatus) {
			statusMap[status.ProxyName] = status

			// 每秒更新一次显示
			if time.Since(lastUpdate) >= time.Second {
				lastUpdate = time.Now()
				printMonitorStatus(statusMap, startTime)
			}
		})

		// 打印最终报告
		printMonitorReport(results)

		// 保存报告
		if *outputPath != "" {
			err = saveMonitorReport(results, *outputPath)
			if err != nil {
				log.Fatalln("save monitor report failed: %v", err)
			}
			fmt.Printf("\n监控报告已保存到: %s\n", *outputPath)
		}
		return
	} else {
		// 速度测试模式
		bar := progressbar.Default(int64(len(allProxies)), "测试中...")
		results := make([]*speedtester.Result, 0)
		speedTester.TestProxies(allProxies, func(result *speedtester.Result) {
			bar.Add(1)
			bar.Describe(result.ProxyName)
			results = append(results, result)
		})

		sort.Slice(results, func(i, j int) bool {
			return results[i].DownloadSpeed > results[j].DownloadSpeed
		})

		printResults(results)

		if *outputPath != "" {
			err = saveConfig(results)
			if err != nil {
				log.Fatalln("save config file failed: %v", err)
			}
			fmt.Printf("\nsave config file to: %s\n", *outputPath)
		}
	}
}

func printResults(results []*speedtester.Result) {
	table := tablewriter.NewWriter(os.Stdout)

	var headers []string
	if *fastMode {
		headers = []string{
			"序号",
			"节点名称",
			"类型",
			"延迟",
		}
	} else {
		headers = []string{
			"序号",
			"节点名称",
			"类型",
			"延迟",
			"抖动",
			"丢包率",
			"下载速度",
			"上传速度",
		}
	}
	table.SetHeader(headers)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.SetColMinWidth(0, 4)  // 序号
	table.SetColMinWidth(1, 20) // 节点名称
	table.SetColMinWidth(2, 8)  // 类型
	table.SetColMinWidth(3, 8)  // 延迟
	if !*fastMode {
		table.SetColMinWidth(4, 8)  // 抖动
		table.SetColMinWidth(5, 8)  // 丢包率
		table.SetColMinWidth(6, 12) // 下载速度
		table.SetColMinWidth(7, 12) // 上传速度
	}

	for i, result := range results {
		idStr := fmt.Sprintf("%d.", i+1)

		// 延迟颜色
		latencyStr := result.FormatLatency()
		if result.Latency > 0 {
			if result.Latency < 800*time.Millisecond {
				latencyStr = colorGreen + latencyStr + colorReset
			} else if result.Latency < 1500*time.Millisecond {
				latencyStr = colorYellow + latencyStr + colorReset
			} else {
				latencyStr = colorRed + latencyStr + colorReset
			}
		} else {
			latencyStr = colorRed + latencyStr + colorReset
		}

		jitterStr := result.FormatJitter()
		if result.Jitter > 0 {
			if result.Jitter < 800*time.Millisecond {
				jitterStr = colorGreen + jitterStr + colorReset
			} else if result.Jitter < 1500*time.Millisecond {
				jitterStr = colorYellow + jitterStr + colorReset
			} else {
				jitterStr = colorRed + jitterStr + colorReset
			}
		} else {
			jitterStr = colorRed + jitterStr + colorReset
		}

		// 丢包率颜色
		packetLossStr := result.FormatPacketLoss()
		if result.PacketLoss < 10 {
			packetLossStr = colorGreen + packetLossStr + colorReset
		} else if result.PacketLoss < 20 {
			packetLossStr = colorYellow + packetLossStr + colorReset
		} else {
			packetLossStr = colorRed + packetLossStr + colorReset
		}

		// 下载速度颜色 (以MB/s为单位判断)
		downloadSpeed := result.DownloadSpeed / (1024 * 1024)
		downloadSpeedStr := result.FormatDownloadSpeed()
		if downloadSpeed >= 10 {
			downloadSpeedStr = colorGreen + downloadSpeedStr + colorReset
		} else if downloadSpeed >= 5 {
			downloadSpeedStr = colorYellow + downloadSpeedStr + colorReset
		} else {
			downloadSpeedStr = colorRed + downloadSpeedStr + colorReset
		}

		// 上传速度颜色
		uploadSpeed := result.UploadSpeed / (1024 * 1024)
		uploadSpeedStr := result.FormatUploadSpeed()
		if uploadSpeed >= 5 {
			uploadSpeedStr = colorGreen + uploadSpeedStr + colorReset
		} else if uploadSpeed >= 2 {
			uploadSpeedStr = colorYellow + uploadSpeedStr + colorReset
		} else {
			uploadSpeedStr = colorRed + uploadSpeedStr + colorReset
		}

		var row []string
		if *fastMode {
			row = []string{
				idStr,
				result.ProxyName,
				result.ProxyType,
				latencyStr,
			}
		} else {
			row = []string{
				idStr,
				result.ProxyName,
				result.ProxyType,
				latencyStr,
				jitterStr,
				packetLossStr,
				downloadSpeedStr,
				uploadSpeedStr,
			}
		}

		table.Append(row)
	}

	fmt.Println()
	table.Render()
	fmt.Println()
}

func saveConfig(results []*speedtester.Result) error {
	proxies := make([]map[string]any, 0)
	for _, result := range results {
		if *maxLatency > 0 && result.Latency > *maxLatency {
			continue
		}
		if *downloadSize > 0 && *minDownloadSpeed > 0 && result.DownloadSpeed < *minDownloadSpeed*1024*1024 {
			continue
		}
		if *uploadSize > 0 && *minUploadSpeed > 0 && result.UploadSpeed < *minUploadSpeed*1024*1024 {
			continue
		}

		proxyConfig := result.ProxyConfig
		if *renameNodes {
			location, err := getIPLocation(proxyConfig["server"].(string))
			if err != nil || location.CountryCode == "" {
				proxies = append(proxies, proxyConfig)
				continue
			}
			proxyConfig["name"] = generateNodeName(location.CountryCode, result.DownloadSpeed)
		}
		proxies = append(proxies, proxyConfig)
	}

	config := &speedtester.RawConfig{
		Proxies: proxies,
	}
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(*outputPath, yamlData, 0o644)
}

type IPLocation struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
}

var countryFlags = map[string]string{
	"US": "🇺🇸", "CN": "🇨🇳", "GB": "🇬🇧", "UK": "🇬🇧", "JP": "🇯🇵", "DE": "🇩🇪", "FR": "🇫🇷", "RU": "🇷🇺",
	"SG": "🇸🇬", "HK": "🇭🇰", "TW": "🇹🇼", "KR": "🇰🇷", "CA": "🇨🇦", "AU": "🇦🇺", "NL": "🇳🇱", "IT": "🇮🇹",
	"ES": "🇪🇸", "SE": "🇸🇪", "NO": "🇳🇴", "DK": "🇩🇰", "FI": "🇫🇮", "CH": "🇨🇭", "AT": "🇦🇹", "BE": "🇧🇪",
	"BR": "🇧🇷", "IN": "🇮🇳", "TH": "🇹🇭", "MY": "🇲🇾", "VN": "🇻🇳", "PH": "🇵🇭", "ID": "🇮🇩", "UA": "🇺🇦",
	"TR": "🇹🇷", "IL": "🇮🇱", "AE": "🇦🇪", "SA": "🇸🇦", "EG": "🇪🇬", "ZA": "🇿🇦", "NG": "🇳🇬", "KE": "🇰🇪",
	"RO": "🇷🇴", "PL": "🇵🇱", "CZ": "🇨🇿", "HU": "🇭🇺", "BG": "🇧🇬", "HR": "🇭🇷", "SI": "🇸🇮", "SK": "🇸🇰",
	"LT": "🇱🇹", "LV": "🇱🇻", "EE": "🇪🇪", "PT": "🇵🇹", "GR": "🇬🇷", "IE": "🇮🇪", "LU": "🇱🇺", "MT": "🇲🇹",
	"CY": "🇨🇾", "IS": "🇮🇸", "MX": "🇲🇽", "AR": "🇦🇷", "CL": "🇨🇱", "CO": "🇨🇴", "PE": "🇵🇪", "VE": "🇻🇪",
	"EC": "🇪🇨", "UY": "🇺🇾", "PY": "🇵🇾", "BO": "🇧🇴", "CR": "🇨🇷", "PA": "🇵🇦", "GT": "🇬🇹", "HN": "🇭🇳",
	"SV": "🇸🇻", "NI": "🇳🇮", "BZ": "🇧🇿", "JM": "🇯🇲", "TT": "🇹🇹", "BB": "🇧🇧", "GD": "🇬🇩", "LC": "🇱🇨",
	"VC": "🇻🇨", "AG": "🇦🇬", "DM": "🇩🇲", "KN": "🇰🇳", "BS": "🇧🇸", "CU": "🇨🇺", "DO": "🇩🇴", "HT": "🇭🇹",
	"PR": "🇵🇷", "VI": "🇻🇮", "GU": "🇬🇺", "AS": "🇦🇸", "MP": "🇲🇵", "PW": "🇵🇼", "FM": "🇫🇲", "MH": "🇲🇭",
	"KI": "🇰🇮", "TV": "🇹🇻", "NR": "🇳🇷", "WS": "🇼🇸", "TO": "🇹🇴", "FJ": "🇫🇯", "VU": "🇻🇺", "SB": "🇸🇧",
	"PG": "🇵🇬", "NC": "🇳🇨", "PF": "🇵🇫", "WF": "🇼🇫", "CK": "🇨🇰", "NU": "🇳🇺", "TK": "🇹🇰", "SC": "🇸🇨",
}

func getIPLocation(ip string) (*IPLocation, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=country,countryCode", ip))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get location for IP %s", ip)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var location IPLocation
	if err := json.Unmarshal(body, &location); err != nil {
		return nil, err
	}
	return &location, nil
}

func generateNodeName(countryCode string, downloadSpeed float64) string {
	flag, exists := countryFlags[strings.ToUpper(countryCode)]
	if !exists {
		flag = "🏳️"
	}

	speedMBps := downloadSpeed / (1024 * 1024)
	return fmt.Sprintf("%s %s | ⬇️ %.2f MB/s", flag, strings.ToUpper(countryCode), speedMBps)
}

// printMonitorStatus 打印监控实时状态
func printMonitorStatus(statusMap map[string]*speedtester.MonitorStatus, startTime time.Time) {
	// 清屏
	fmt.Print("\033[H\033[2J")

	fmt.Println("=== 连接稳定性监控 (实时) ===")
	fmt.Printf("运行时间: %v\n\n", time.Since(startTime).Round(time.Second))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"节点名称", "状态", "在线时长", "断线次数", "稳定率", "数据包数", "数据量"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")

	// 按节点名称排序
	var names []string
	for name := range statusMap {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		status := statusMap[name]

		statusIcon := "✅"
		if !status.IsAlive {
			statusIcon = "❌"
		}

		stabilityRate := float64(0)
		if status.TotalDuration > 0 {
			stabilityRate = float64(status.OnlineDuration) / float64(status.TotalDuration) * 100
		}

		// 格式化数据量显示
		dataSize := formatDataSize(status.TotalDataBytes)

		row := []string{
			name,
			statusIcon,
			formatDuration(status.OnlineDuration),
			fmt.Sprintf("%d", status.DisconnectCount),
			fmt.Sprintf("%.2f%%", stabilityRate),
			fmt.Sprintf("%d", status.DataPacketCount),
			dataSize,
		}

		table.Append(row)
	}

	table.Render()
}

// printMonitorReport 打印监控最终报告
func printMonitorReport(results []*speedtester.MonitorResult) {
	fmt.Println("\n=== 监控报告 ===")

	if len(results) == 0 {
		fmt.Println("没有监控数据")
		return
	}

	fmt.Printf("监控时间: %s - %s\n\n",
		results[0].StartTime.Format("2006-01-02 15:04:05"),
		results[0].EndTime.Format("2006-01-02 15:04:05"))

	// 按稳定率排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].StabilityRate > results[j].StabilityRate
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"排名", "节点名称", "类型", "稳定率", "断线次数", "最长在线", "断线明细"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for i, result := range results {
		rank := fmt.Sprintf("%d", i+1)

		// 稳定率着色
		stabilityStr := result.FormatStabilityRate()
		if result.StabilityRate >= 99.9 {
			stabilityStr = colorGreen + stabilityStr + colorReset
		} else if result.StabilityRate >= 99 {
			stabilityStr = colorYellow + stabilityStr + colorReset
		} else {
			stabilityStr = colorRed + stabilityStr + colorReset
		}

		// 断线明细
		var disconnectDetails string
		if len(result.DisconnectEvents) > 0 {
			var details []string
			for _, event := range result.DisconnectEvents {
				details = append(details, fmt.Sprintf("%s(%v)",
					event.StartTime.Format("15:04:05"),
					event.Duration.Round(time.Second)))
			}
			disconnectDetails = strings.Join(details, ", ")
			if len(disconnectDetails) > 50 {
				disconnectDetails = disconnectDetails[:47] + "..."
			}
		} else {
			disconnectDetails = "-"
		}

		row := []string{
			rank,
			result.ProxyName,
			result.ProxyType,
			stabilityStr,
			fmt.Sprintf("%d", result.DisconnectCount),
			formatDuration(result.MaxOnlineTime),
			disconnectDetails,
		}

		table.Append(row)
	}

	table.Render()
}

// saveMonitorReport 保存监控报告
func saveMonitorReport(results []*speedtester.MonitorResult, outputPath string) error {
	// 转换为YAML格式
	type MonitorReportItem struct {
		Name            string  `yaml:"name"`
		Type            string  `yaml:"type"`
		StabilityRate   float64 `yaml:"stability_rate"`
		DisconnectCount int     `yaml:"disconnect_count"`
		MaxOnlineTime   string  `yaml:"max_online_time"`
		StartTime       string  `yaml:"start_time"`
		EndTime         string  `yaml:"end_time"`
	}

	var report []MonitorReportItem
	for _, result := range results {
		item := MonitorReportItem{
			Name:            result.ProxyName,
			Type:            result.ProxyType,
			StabilityRate:   result.StabilityRate,
			DisconnectCount: result.DisconnectCount,
			MaxOnlineTime:   formatDuration(result.MaxOnlineTime),
			StartTime:       result.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:         result.EndTime.Format("2006-01-02 15:04:05"),
		}
		report = append(report, item)
	}

	data, err := yaml.Marshal(report)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0644)
}

// formatDuration 格式化时长
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// formatDataSize 格式化数据量显示
func formatDataSize(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	sizes := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), sizes[exp])
}
