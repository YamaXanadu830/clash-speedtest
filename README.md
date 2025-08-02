# Clash-SpeedTest

基于 Clash/Mihomo 核心的测速工具，快速测试你的节点速度。

## 特性

1. 无需额外的配置，直接将 Clash/Mihomo 配置本地文件路径或者订阅地址作为参数传入即可
2. 支持 Proxies 和 Proxy Provider 中定义的全部类型代理节点，兼容性跟 Mihomo 一致
3. 不依赖额外的 Clash/Mihomo 进程实例，单一工具即可完成测试
4. 代码简单而且开源，不发布构建好的二进制文件，保证你的节点安全

<img width="1332" alt="image" src="https://github.com/user-attachments/assets/fdc47ec5-b626-45a3-a38a-6d88c326c588">

## 项目结构

```
clash-speedtest-fork/
├── main.go                   # 主程序入口
├── speedtester/              # 核心测速引擎
│   ├── speedtester.go        # 测速逻辑
│   └── zeroreader.go         # 上传数据生成器
└── download-server/          # 自建测速服务器
    └── download-server.go    # 服务器代码
```

## 功能实现

### 测试流程
1. **延迟测试** - 6次ping，统计平均延迟、抖动、丢包率
2. **下载测试** - 多线程并发下载，计算下载速度
3. **上传测试** - 并发上传，计算上传速度

### 支持协议
Shadowsocks/ShadowsocksR、VMess/VLess、Trojan、Hysteria/Hysteria2、WireGuard、Tuic、SSH等

## 使用方法

```bash
# 支持从源码安装，或从 Release 里下载由 Github Action 自动构建的二进制文件
> go install github.com/YamaXanadu830/clash-speedtest@latest

# 查看帮助
> clash-speedtest -h
Usage of clash-speedtest:
  -c string
        configuration file path, also support http(s) url
  -f string
        filter proxies by name, use regexp (default ".*")
  -b string
        block proxies by keywords, use | to separate multiple keywords (example: -b 'rate|x1|1x')
  -server-url string
        server url for testing proxies (default "https://speed.cloudflare.com")
  -download-size int
        download size for testing proxies (default 50MB)
  -upload-size int
        upload size for testing proxies (default 20MB)
  -timeout duration
        timeout for testing proxies (default 5s)
  -concurrent int
        download concurrent size (default 4)
  -output string
        output config file path (default "")
  -stash-compatible
        enable stash compatible mode
  -max-latency duration
        filter latency greater than this value (default 800ms)
  -min-download-speed float
        filter speed less than this value(unit: MB/s) (default 5)
  -min-upload-speed float
        filter upload speed less than this value(unit: MB/s) (default 2)
  -rename
        rename nodes with IP location and speed
  -fast
        enable fast mode, only test latency

# 演示：

# 1. 测试全部节点，使用 HTTP 订阅地址
# 请在订阅地址后面带上 flag=meta 参数，否则无法识别出节点类型
> clash-speedtest -c 'https://domain.com/api/v1/client/subscribe?token=secret&flag=meta'

# 2. 测试香港节点，使用正则表达式过滤，使用本地文件
> clash-speedtest -c ~/.config/clash/config.yaml -f 'HK|港'
节点                                        	带宽          	延迟
Premium|广港|IEPL|01                        	484.80KB/s  	815.00ms
Premium|广港|IEPL|02                        	N/A         	N/A
Premium|广港|IEPL|03                        	2.62MB/s    	333.00ms
Premium|广港|IEPL|04                        	1.46MB/s    	272.00ms
Premium|广港|IEPL|05                        	3.87MB/s    	249.00ms

# 3. 当然你也可以混合使用
> clash-speedtest -c "https://domain.com/api/v1/client/subscribe?token=secret&flag=meta,/home/.config/clash/config.yaml"

# 4. 筛选出延迟低于 800ms 且下载速度大于 5MB/s 的节点，并输出到 filtered.yaml
> clash-speedtest -c "https://domain.com/api/v1/client/subscribe?token=secret&flag=meta" -output filtered.yaml -max-latency 800ms -min-speed 5
# 筛选后的配置文件可以直接粘贴到 Clash/Mihomo 中使用，或是贴到 Github\Gist 上通过 Proxy Provider 引用。

# 5. 使用 -rename 选项按照 IP 地区和下载速度重命名节点
> clash-speedtest -c config.yaml -output result.yaml -rename
# 重命名后的节点名称格式：🇺🇸 US | ⬇️ 15.67 MB/s
# 包含国旗 emoji、国家代码和下载速度

# 6. 快速测试模式
> clash-speedtest -f 'HK' -fast -c ~/.config/clash/config.yaml
# 此命令将只测试节点延迟，跳过其他测试项目，适用于：
# - 快速检查节点是否可用
# - 只需要检查延迟的场景
# - 需要快速得到测试结果的场景
🇭🇰 香港 HK-10 100% |██████████████████| (20/20, 13 it/min)
序号    节点名称                类型            延迟
1.      🇭🇰 香港 HK-01           Trojan          657ms
2.      🇭🇰 香港 HK-20           Trojan          649ms
3.      🇭🇰 香港 HK-15           Trojan          674ms
4.      🇭🇰 香港 HK-19           Trojan          649ms
5.      🇭🇰 香港 HK-12           Trojan          667ms

# 7. 连接稳定性监控模式
> clash-speedtest -c config.yaml --monitor --monitor-duration 24h
# 监控节点的连接稳定性，特别适用于需要24小时不间断连接的场景（如交易平台）
# 功能特点：
# - 持续监控节点连接状态
# - 秒级精度的断线检测
# - 实时显示各节点状态
# - 生成详细的稳定性报告

# 监控1小时并导出报告
> clash-speedtest -c config.yaml -f 'HK|SG' --monitor --monitor-duration 1h --output monitor-report.yaml

# 使用WebSocket流式数据监控模式
> clash-speedtest -c config.yaml --monitor --monitor-type websocket --monitor-duration 1h
# WebSocket模式连接到Binance实时数据流，监控真正的流式数据连接稳定性
# 适用于需要24小时不间断数据流的交易应用

监控参数说明：
  --monitor              启用稳定性监控模式
  --monitor-duration     监控时长（默认24h）
  --monitor-interval     心跳检测间隔（默认1s）
  --monitor-type         监控类型：http或websocket（默认http）

监控输出示例：
=== 连接稳定性监控 (实时) ===
运行时间: 01:23:45

节点名称        状态    在线时长    断线次数    稳定率    数据包数    数据量
🇭🇰 香港 01      ✅      01:23:45    0          100.00%   4985        2.1 MB
🇸🇬 新加坡 01    ✅      01:23:42    1          99.97%    4968        2.1 MB  
🇯🇵 日本 01      ❌      01:20:15    3          96.58%    4512        1.9 MB

## WebSocket流式数据监控特性
- **真实数据流监控**：连接到Binance实时BTC/USDT价格数据流
- **精确断线检测**：基于数据包接收间隔（>10秒=断线）
- **详细统计指标**：数据包接收数、总字节数、连接稳定性
- **24小时监控**：适用于交易平台的长期稳定性测试

## 测速原理

通过 HTTP GET 请求下载指定大小的文件，默认使用 https://speed.cloudflare.com (50MB) 进行测试，计算下载时间得到下载速度。

测试结果：
1. 带宽 是指下载指定大小文件的速度，即一般理解中的下载速度。当这个数值越高时表明节点的出口带宽越大。
2. 延迟 是指 HTTP GET 请求拿到第一个字节的的响应时间，即一般理解中的 TTFB。当这个数值越低时表明你本地到达节点的延迟越低，可能意味着中转节点有 BGP 部署、出海线路是 IEPL、IPLC 等。

请注意带宽跟延迟是两个独立的指标，两者并不关联：
1. 可能带宽很高但是延迟也很高，这种情况下你下载速度很快但是打开网页的时候却很慢，可能是是中转节点没有 BGP 加速，但出海线路带宽很充足。
2. 可能带宽很低但是延迟也很低，这种情况下你打开网页的时候很快但是下载速度很慢，可能是中转节点有 BGP 加速，但出海线路的 IEPL、IPLC 带宽很小。

Cloudflare 是全球知名的 CDN 服务商，其提供的测速服务器到海外绝大部分的节点速度都很快，一般情况下都没有必要自建测速服务器。

如果你不想使用 Cloudflare 的测速服务器，可以自己搭建一个测速服务器。

```shell
# 在您需要进行测速的服务器上安装和启动测速服务器
> go install github.com/YamaXanadu830/clash-speedtest/download-server@latest
> download-server

# 此时在本地使用 http://your-server-ip:8080 作为 server-url 即可
> clash-speedtest --server-url "http://your-server-ip:8080"
```

## License

[GPL-3.0](LICENSE)

## 📝 功能开发任务 (TODO List)

### 连接稳定性监控功能

#### 1. 添加监控模式命令行参数
- [x] 在 main.go 添加 --monitor 参数（布尔值）
- [x] 添加 --monitor-duration 参数（默认24h）
- [x] 添加 --monitor-interval 参数（默认1s）
- [x] 添加监控模式与速度测试模式的互斥判断

#### 2. 实现监控核心功能
- [x] 创建 speedtester/monitor.go 文件
- [x] 实现 MonitorConfig 结构体（包含duration、interval、target）
- [x] 实现 MonitorResult 结构体（包含稳定率、断线次数、断线明细）
- [x] 实现 testStability() 方法，使用HTTP Keep-Alive长连接
- [x] 实现断线检测和自动重连逻辑

#### 3. 添加实时状态显示
- [x] 实现实时状态输出（类似进度条）
- [x] 显示：节点名称、连接状态、运行时长、断线次数
- [x] 每秒刷新显示，使用终端控制符清屏更新

#### 4. 生成监控报告
- [x] 实现监控数据统计（稳定率计算、最长在线时间等）
- [x] 添加监控报告输出格式（表格形式）
- [x] 支持导出监控报告到文件（--output参数复用）

#### 5. 更新文档
- [x] 在 README.md 添加监控模式使用说明
- [x] 添加监控模式命令示例
- [x] 说明监控报告的指标含义

### 使用示例
```bash
# 监控所有香港节点24小时稳定性
clash-speedtest -c config.yaml -f 'HK|香港' --monitor --monitor-duration 24h

# 监控并输出报告
clash-speedtest -c config.yaml --monitor --monitor-duration 1h --output stability-report.yaml
```
