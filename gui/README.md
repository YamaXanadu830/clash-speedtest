# 🖥️ Clash SpeedTest GUI

基于 Wails v2 的现代化 Clash 节点测速和监控工具图形界面。

## ✨ 功能特性

### 🚀 速度测试
- **节点测速** - 批量测试代理节点的延迟、下载/上传速度
- **实时进度** - 实时显示测试进度和结果
- **多种模式** - 支持完整测试和快速模式（仅测延迟）
- **灵活过滤** - 支持正则表达式过滤和关键词屏蔽
- **结果导出** - 支持导出测试结果到文件

### 📡 稳定性监控
- **长期监控** - 支持 5分钟 到 24小时 的连续监控
- **双重模式** - HTTP长连接 和 WebSocket数据流监控
- **实时统计** - 显示在线率、断线次数、数据传输量
- **Binance集成** - WebSocket模式连接Binance实时数据流
- **监控报告** - 生成详细的稳定性监控报告

### ⚙️ 灵活配置
- **参数调节** - 自定义测速服务器、并发数、超时时间
- **过滤条件** - 设置延迟、速度阈值过滤节点
- **界面设置** - 深色模式、通知提醒、自动保存
- **默认配置** - 保存常用设置为默认值

## 🏗️ 技术架构

### 后端 (Go)
- **Wails v2** - 现代化Go GUI框架
- **核心引擎** - 复用原项目speedtester包
- **事件驱动** - 实时数据推送到前端
- **跨平台** - 支持Windows/macOS/Linux

### 前端 (Vue 3)
- **Vue 3** - 现代化响应式框架
- **Tailwind CSS** - 实用优先的CSS框架
- **组件化** - 清晰的界面组件结构
- **实时更新** - WebSocket事件实时更新界面

## 📁 项目结构

```
gui/
├── main.go                 # 应用入口
├── app.go                  # Wails应用绑定层
├── wails.json             # Wails配置文件
├── go.mod                 # Go模块依赖
├── frontend/              # Vue前端代码
│   ├── src/
│   │   ├── App.vue        # 主应用组件
│   │   ├── views/         # 页面组件
│   │   │   ├── SpeedTestView.vue    # 速度测试页面
│   │   │   ├── MonitorView.vue      # 监控页面
│   │   │   └── SettingsView.vue     # 设置页面
│   │   ├── components/    # 公共组件
│   │   └── utils/         # 工具函数
│   ├── package.json       # Node.js依赖
│   └── vite.config.js     # Vite构建配置
└── build/                 # 构建输出目录
    ├── bin/               # 可执行文件
    └── pkg/               # 安装包
```

## 🛠️ 开发环境

### 依赖要求
- **Go 1.24+** - 后端开发语言
- **Node.js 18+** - 前端开发环境
- **Wails CLI** - GUI框架工具链

### 安装 Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 初始化项目
```bash
cd gui
go mod tidy
npm install
```

## 🚀 开发模式

### 启动开发服务器
```bash
# 在gui目录下
wails dev
```

这将启动：
- Go后端热重载
- Vue前端开发服务器 (localhost:5173)
- 自动打开应用窗口

### 前端独立开发
```bash
cd frontend
npm run dev
```

## 📦 构建部署

### 构建应用
```bash
# 开发构建
wails build

# 生产构建（优化）
wails build -clean -upx -s

# 跨平台构建
wails build -platform windows/amd64
wails build -platform darwin/amd64
wails build -platform linux/amd64
```

### 构建输出
- **可执行文件**: `build/bin/clash-speedtest-gui`
- **安装包**: `build/pkg/` (根据平台生成)

## 🎨 界面预览

### 速度测试界面
- 左侧配置面板：配置文件加载、测试参数、过滤条件
- 右侧结果面板：实时进度、测试结果表格、操作按钮

### 监控界面
- 监控参数配置：类型选择、时长设置、目标地址
- 实时状态显示：在线状态、稳定率、数据统计

### 设置界面
- 默认参数配置：测速、监控默认值
- 界面设置：主题、语言、通知
- 系统信息：版本、平台信息

## 🔧 配置说明

### 支持的配置文件
- **Clash配置** - 标准Clash YAML格式
- **节点订阅** - HTTP(S) URL订阅地址
- **本地文件** - 本地YAML配置文件

### 测速服务器
- **Cloudflare** - 推荐，全球CDN网络
- **Fast.com** - Netflix CDN网络
- **Speedtest.net** - 传统测速服务

### 监控模式
- **HTTP模式** - Keep-Alive长连接监控
- **WebSocket模式** - 实时数据流监控（推荐）

## 📊 数据格式

### 测速结果
```json
{
  "proxyName": "🇭🇰 香港 01",
  "proxyType": "shadowsocks",
  "latency": 45,
  "jitter": 12,
  "packetLoss": 0.5,
  "downloadSpeed": 52428800,
  "uploadSpeed": 10485760,
  "status": "完成"
}
```

### 监控状态
```json
{
  "proxyName": "🇭🇰 香港 01",
  "isAlive": true,
  "onlineDuration": 3600,
  "totalDuration": 3660,
  "disconnectCount": 2,
  "dataPacketCount": 3600,
  "totalDataBytes": 1048576,
  "lastUpdate": "14:30:25"
}
```

## 🐛 问题排查

### 常见问题
1. **构建失败** - 检查Go和Node.js版本
2. **前端无法连接** - 确认Wails绑定正确
3. **配置加载失败** - 检查文件路径和格式
4. **测速无结果** - 确认网络连接和代理配置

### 调试模式
```bash
# 启用调试日志
wails dev -debug

# 查看前端控制台
# 在应用中按F12打开开发者工具
```

## 📝 更新日志

### v1.0.0 (2025-01-20)
- 🎉 初始版本发布
- ✅ 完整的速度测试功能
- ✅ 稳定性监控功能
- ✅ 现代化界面设计
- ✅ 跨平台支持

## 📄 许可证

[GPL-3.0](../LICENSE)