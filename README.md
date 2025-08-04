# ⚡ Clash-SpeedTest

基于 Clash/Mihomo 的测速和监控工具

## 🚀 功能

- **📊 节点测速** - 测试延迟、下载/上传速度  
- **🔍 稳定性监控** - 24小时连接监控，支持HTTP和WebSocket模式
- **📈 流式数据监控** - 连接Binance实时数据流，适用于交易平台
- **🖥️ 图形界面** - 基于 Wails 的现代化 GUI 界面

## 📁 项目结构

```
clash-speedtest-fork/
├── main.go                   # 主程序入口和命令行参数处理
├── speedtester/              # 核心引擎
│   ├── speedtester.go        # 节点测速逻辑
│   ├── monitor.go            # 稳定性监控逻辑（NEW）
│   └── zeroreader.go         # 上传数据生成器
├── download-server/          # 自建测速服务器
│   └── download-server.go    # HTTP测速服务器
└── gui/                      # 图形界面版本
    ├── main.go              # GUI 主程序
    ├── app.go               # 应用逻辑
    └── frontend/            # 前端界面
```

## 📋 支持协议
Shadowsocks、VMess、Trojan、Hysteria、WireGuard、Tuic等

## 💻 安装使用

```bash
# 方式1：直接安装
go install github.com/YamaXanadu830/clash-speedtest@latest

# 方式2：本地编译
git clone https://github.com/YamaXanadu830/clash-speedtest.git
cd clash-speedtest
go build -o clash-speedtest

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
```

## 🖥️ GUI界面版本

基于 Wails v2 开发的现代化桌面应用程序。

### 主界面布局
```
┌─────────────────────────────────────────────────────────────┐
│  ⚡ Clash SpeedTest v1.0.0                    [⚙️] [❓]     │
├─────────────────────────────────────────────────────────────┤
│ [🚀速度测试] [📡稳定监控] [📊数据分析] [⚙️设置]          │
├─────────────────────────────────┬───────────────────────────┤
│  📊 测速服务器                  │  ⏳ 测试中... (45/100)   │
│  ● Cloudflare (推荐)           │  ████████░░░░ 45%         │
│  ○ OVH (100MB)                 │                           │
│  ○ Tele2 (100MB)               │  已完成节点:              │
│  ○ Hetzner (100MB)             │  ✅ HK-01  256MB/s        │
│                                 │  ✅ JP-02  198MB/s        │
│  ⚙️ 测试参数                   │  ❌ SG-03  超时           │
│  下载大小: [50] MB             │  ⏳ US-04  测试中...      │
│  上传大小: [20] MB             │                           │
│  超时时间: [5] 秒              │  [🚀开始] [⏹停止]        │
│  并发数:   ████░░░░ 4          │  [📤导出] [🗑清空]        │
│  ☑ 快速模式                   │                           │
└─────────────────────────────────┴───────────────────────────┘
```

### 安装运行
```bash
cd gui
wails build  # 编译
./build/bin/clash-speedtest-gui  # 运行
```

### 核心功能
- **🚀 速度测试** - 节点速度和延迟测试，支持多种测速服务器
- **📡 稳定监控** - 24小时连接稳定性监控，支持HTTP和WebSocket
- **📊 数据分析** - 历史数据统计分析，节点性能排名
- **⚙️ 设置** - 应用配置和系统信息管理

### 测速服务器选项
```
支持的测速服务器：
• Cloudflare     - https://speed.cloudflare.com (推荐)
• OVH           - https://proof.ovh.net/files/100Mb.dat
• Tele2         - http://speedtest.tele2.net/100MB.zip  
• Hetzner       - https://ash-speed.hetzner.com/100MB.bin
• 自定义服务器   - 支持输入自定义URL
```

## 🧪 测试示例

```bash
# 香港节点10分钟WebSocket监控测试
./clash-speedtest -c config.yaml --monitor --monitor-type websocket --monitor-duration 10m -f "香港"

# 带报告输出
./clash-speedtest -c config.yaml --monitor --monitor-type websocket --monitor-duration 10m -f "香港" --output hk_report.yaml

# 备选方案
go run main.go -c config.yaml --monitor --monitor-type websocket --monitor-duration 10m -f "🇭🇰|香港"

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

## 📝 最近更新 (v2.0)

### 🎯 核心改进
- **🆕 GUI界面全面重构** - 基于 Wails v2 的现代化桌面应用
- **✨ 多种测速服务器** - 新增 OVH、Tele2、Hetzner 测速选项
- **🔧 状态管理优化** - 修复停止测试后无法立即重启的问题
- **🗑️ 数据管理完善** - 清空功能支持前后端数据同步清理

### 🎨 界面优化
- **📊 数据分析重构** - 移除占位图表，优化统计数据展示
- **🗂️ 功能菜单精简** - 从5个菜单优化为4个核心功能
- **🖥️ 系统信息修复** - 支持实际平台检测（macOS/Windows/Linux）
- **🎯 用户体验提升** - 修复界面文字显示和布局问题

### 🛠️ 技术改进
- **⚡ 实时状态同步** - 优化前后端状态管理机制
- **🔄 并发处理优化** - 改进测试停止和重启的时序控制
- **📈 历史数据管理** - 完善数据持久化和统计分析功能
- **🧹 代码清理** - 移除无效功能，专注核心测速和监控能力

## License

[GPL-3.0](LICENSE)

