# ⚡ Clash-SpeedTest

基于 Clash/Mihomo 的测速和监控工具

## 🚀 功能

- **📊 节点测速** - 测试延迟、下载/上传速度  
- **🔍 稳定性监控** - 24小时连接监控，支持HTTP和WebSocket模式
- **📈 流式数据监控** - 连接Binance实时数据流，适用于交易平台

<img width="1332" alt="image" src="https://github.com/user-attachments/assets/fdc47ec5-b626-45a3-a38a-6d88c326c588">

## 📁 项目结构

```
clash-speedtest-fork/
├── main.go                   # 主程序入口和命令行参数处理
├── speedtester/              # 核心引擎
│   ├── speedtester.go        # 节点测速逻辑
│   ├── monitor.go            # 稳定性监控逻辑（NEW）
│   └── zeroreader.go         # 上传数据生成器
└── download-server/          # 自建测速服务器
    └── download-server.go    # HTTP测速服务器
```

## 📋 支持协议
Shadowsocks、VMess、Trojan、Hysteria、WireGuard、Tuic等

## 💻 安装使用

```bash
# 安装
go install github.com/YamaXanadu830/clash-speedtest@latest

# 🚀 基础测速
clash-speedtest -c config.yaml                    # 测试所有节点
clash-speedtest -c config.yaml -f 'HK|SG'         # 测试指定地区
clash-speedtest -c config.yaml --fast             # 快速测试(仅延迟)

# 📊 稳定性监控
clash-speedtest -c config.yaml --monitor --monitor-duration 1h                    # HTTP监控
clash-speedtest -c config.yaml --monitor --monitor-type websocket --monitor-duration 24h  # WebSocket监控

# 🔥 新功能：WebSocket流式数据监控
> clash-speedtest -c config.yaml --monitor --monitor-type websocket --monitor-duration 1h
# 连接Binance实时BTC/USDT数据流，监控真正的流式数据连接稳定性
# 特别适用于交易平台、实时数据应用的24小时稳定性测试

## 📋 监控参数说明
```bash
--monitor              # 启用稳定性监控模式
--monitor-duration     # 监控时长（默认24h，支持：1m, 30m, 1h, 12h, 24h）
--monitor-interval     # 心跳检测间隔（默认1s）
--monitor-type         # 监控类型：http 或 websocket（默认http）
```

## 📊 实时监控输出示例
```
=== 连接稳定性监控 (实时) ===
运行时间: 01:23:45

节点名称        状态    在线时长    断线次数    稳定率    数据包数    数据量
🇭🇰 香港 01      ✅      01:23:45    0          100.00%   4985        2.1 MB
🇸🇬 新加坡 01    ✅      01:23:42    1          99.97%    4968        2.1 MB  
🇯🇵 日本 01      ❌      01:20:15    3          96.58%    4512        1.9 MB
```

## 🌐 WebSocket流式数据监控特性
- **🔴 真实数据流监控**：连接到Binance实时BTC/USDT价格数据流（wss://stream.binance.com:9443/ws/btcusdt@ticker）
- **⚡ 精确断线检测**：基于数据包接收间隔的智能检测（>10秒无数据=断线）
- **📈 详细统计指标**：数据包接收数、总字节数、连接稳定性、断线事件时间线
- **🕐 24小时监控**：专为交易平台、金融应用等关键服务的长期稳定性测试而设计
- **📝 完整报告导出**：支持YAML格式的详细监控报告，包含所有统计数据和事件记录

## ⚙️ 工作原理

### 📊 节点测速
默认使用Cloudflare测速服务器，测试节点的延迟和带宽性能。

### 📡 稳定性监控
- **HTTP模式**：使用Keep-Alive长连接持续监控节点连接状态
- **WebSocket模式**：连接Binance实时数据流 `wss://stream.binance.com:9443/ws/btcusdt@ticker`，通过数据包接收间隔检测断线（>10秒无数据=断线）

## License

[GPL-3.0](LICENSE)

