<template>
  <div class="flex h-full">
    <!-- 左侧配置面板 -->
    <div class="w-80 bg-gray-800 border-r border-gray-700 p-6 overflow-y-auto">
      <div class="space-y-6">
        <!-- 监控参数 -->
        <div class="card">
          <div class="card-header">
            <h3 class="text-lg font-semibold text-white">📡 监控参数</h3>
          </div>
          <div class="card-body space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">监控类型</label>
              <select v-model="monitorConfig.type" class="input w-full">
                <option value="http">HTTP 长连接</option>
                <option value="websocket">WebSocket 数据流</option>
              </select>
              <p class="text-xs text-gray-400 mt-1">
                WebSocket模式连接Binance实时数据流，更准确检测断线
              </p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">目标地址</label>
              <input
                v-model="monitorConfig.targetURL"
                type="text"
                class="input w-full"
                :placeholder="monitorConfig.type === 'websocket' ? 'wss://stream.binance.com:9443/ws/btcusdt@ticker' : 'https://cn.tradingview.com/chart/'"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">监控时长</label>
              <select v-model.number="monitorConfig.duration" class="input w-full">
                <option :value="60">1 分钟 (测试)</option>
                <option :value="300">5 分钟</option>
                <option :value="600">10 分钟</option>
                <option :value="1800">30 分钟</option>
                <option :value="3600">1 小时</option>
                <option :value="7200">2 小时</option>
                <option :value="14400">4 小时</option>
                <option :value="28800">8 小时</option>
                <option :value="86400">24 小时</option>
              </select>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">心跳间隔 (秒)</label>
              <input
                v-model.number="monitorConfig.interval"
                type="number"
                class="input w-full"
                min="1"
                max="60"
              />
            </div>
          </div>
        </div>

        <!-- 监控统计 -->
        <div v-if="monitorStats.totalNodes > 0" class="card">
          <div class="card-header">
            <h3 class="text-lg font-semibold text-white">📊 监控统计</h3>
          </div>
          <div class="card-body space-y-3">
            <div class="flex justify-between">
              <span class="text-gray-300">总节点数:</span>
              <span class="text-white font-medium">{{ monitorStats.totalNodes }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-300">在线节点:</span>
              <span class="text-green-400 font-medium">{{ monitorStats.onlineNodes }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-300">离线节点:</span>
              <span class="text-red-400 font-medium">{{ monitorStats.offlineNodes }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-300">平均稳定率:</span>
              <span class="text-white font-medium">{{ monitorStats.avgStability }}%</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-300">运行时长:</span>
              <span class="text-white font-medium">{{ formatDuration(monitorStats.runningTime) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧监控面板 -->
    <div class="flex-1 flex flex-col">
      <!-- 操作栏 -->
      <div class="bg-gray-800 border-b border-gray-700 p-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <button
              @click="startMonitor"
              :disabled="!configInfo || isRunning"
              class="btn btn-success"
            >
              <span v-if="!isRunning">📡 开始监控</span>
              <span v-else class="flex items-center space-x-2">
                <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>监控中...</span>
              </span>
            </button>
            
            <button
              @click="stopMonitor"
              :disabled="!isRunning"
              class="btn btn-danger"
            >
              ⏹ 停止监控
            </button>
            
            <button
              @click="exportMonitorReport"
              :disabled="!monitorResults.length"
              class="btn btn-secondary"
            >
              📋 导出报告
            </button>
          </div>
          
          <div v-if="isRunning" class="text-sm text-gray-300">
            监控类型: {{ monitorConfig.type === 'websocket' ? 'WebSocket数据流' : 'HTTP长连接' }}
          </div>
        </div>
      </div>

      <!-- 监控结果 -->
      <div class="flex-1 overflow-y-auto p-4">
        <div v-if="!configInfo" class="text-center py-12">
          <div class="text-gray-400 text-lg mb-4">📁</div>
          <p class="text-gray-400">请先选择并加载配置文件</p>
        </div>
        
        <div v-else-if="!monitorResults.length && !isRunning" class="text-center py-12">
          <div class="text-gray-400 text-lg mb-4">📡</div>
          <p class="text-gray-400">点击"开始监控"开始稳定性监控</p>
        </div>
        
        <div v-else>
          <!-- 实时状态表格 -->
          <div class="bg-gray-800 rounded-lg overflow-hidden">
            <table class="table">
              <thead>
                <tr>
                  <th class="w-8">#</th>
                  <th>节点名称</th>
                  <th class="w-20">状态</th>
                  <th class="w-32">在线时长</th>
                  <th class="w-24">断线次数</th>
                  <th class="w-24">稳定率</th>
                  <th v-if="monitorConfig.type === 'websocket'" class="w-24">数据包</th>
                  <th v-if="monitorConfig.type === 'websocket'" class="w-24">数据量</th>
                  <th class="w-24">最后更新</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(status, index) in sortedMonitorResults" :key="status.proxyName">
                  <td class="text-gray-400">{{ index + 1 }}</td>
                  <td class="font-medium">{{ status.proxyName }}</td>
                  <td>
                    <div class="flex items-center space-x-2">
                      <div :class="[
                        'status-dot',
                        status.isAlive ? 'status-online' : 'status-offline'
                      ]"></div>
                      <span :class="status.isAlive ? 'text-green-400' : 'text-red-400'">
                        {{ status.isAlive ? '在线' : '离线' }}
                      </span>
                    </div>
                  </td>
                  <td>{{ formatDuration(status.onlineDuration) }}</td>
                  <td :class="getDisconnectClass(status.disconnectCount)">
                    {{ status.disconnectCount }}
                  </td>
                  <td :class="getStabilityClass(getStabilityRate(status))">
                    {{ getStabilityRate(status) }}%
                  </td>
                  <td v-if="monitorConfig.type === 'websocket'" class="text-gray-300">
                    {{ formatNumber(status.dataPacketCount) }}
                  </td>
                  <td v-if="monitorConfig.type === 'websocket'" class="text-gray-300">
                    {{ formatDataSize(status.totalDataBytes) }}
                  </td>
                  <td class="text-gray-400 text-xs">{{ status.lastUpdate }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'

export default {
  name: 'MonitorView',
  props: {
    configInfo: Object,
    isRunning: Boolean
  },
  emits: ['start-monitor', 'stop-monitor'],
  setup(props, { emit }) {
    
    const monitorResults = ref([])
    const monitorStats = ref({
      totalNodes: 0,
      onlineNodes: 0,
      offlineNodes: 0,
      avgStability: 0,
      runningTime: 0
    })
    
    const monitorConfig = ref({
      type: 'websocket',
      targetURL: 'wss://stream.binance.com:9443/ws/btcusdt@ticker',
      duration: 3600, // 1小时
      interval: 1 // 1秒
    })

    const startTime = ref(null)

    // 计算属性
    const sortedMonitorResults = computed(() => {
      return [...monitorResults.value].sort((a, b) => {
        const stabilityA = getStabilityRate(a)
        const stabilityB = getStabilityRate(b)
        return stabilityB - stabilityA
      })
    })


    // 开始监控
    const startMonitor = () => {
      if (!props.configInfo) return
      
      // 根据监控类型设置默认目标地址
      if (monitorConfig.value.type === 'websocket' && !monitorConfig.value.targetURL) {
        monitorConfig.value.targetURL = 'wss://stream.binance.com:9443/ws/btcusdt@ticker'
      } else if (monitorConfig.value.type === 'http' && !monitorConfig.value.targetURL) {
        monitorConfig.value.targetURL = 'https://cn.tradingview.com/chart/'
      }
      
      monitorResults.value = []
      monitorStats.value = {
        totalNodes: 0,
        onlineNodes: 0,
        offlineNodes: 0,
        avgStability: 0,
        runningTime: 0
      }
      startTime.value = Date.now()
      
      emit('start-monitor', monitorConfig.value)
    }

    // 停止监控
    const stopMonitor = () => {
      emit('stop-monitor')
      startTime.value = null
    }

    // 导出监控报告
    const exportMonitorReport = () => {
      console.log('导出监控报告', monitorResults.value)
    }

    // 计算稳定率
    const getStabilityRate = (status) => {
      if (status.totalDuration <= 0) return 100
      return ((status.onlineDuration / status.totalDuration) * 100).toFixed(2)
    }

    // 格式化函数
    const formatDuration = (seconds) => {
      const hours = Math.floor(seconds / 3600)
      const minutes = Math.floor((seconds % 3600) / 60)
      const secs = seconds % 60
      return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
    }

    const formatNumber = (num) => {
      return num.toLocaleString()
    }

    const formatDataSize = (bytes) => {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return (bytes / Math.pow(k, i)).toFixed(1) + ' ' + sizes[i]
    }

    // 样式类函数
    const getStabilityClass = (stability) => {
      const rate = parseFloat(stability)
      if (rate >= 99.5) return 'text-green-400'
      if (rate >= 98) return 'text-yellow-400'
      return 'text-red-400'
    }

    const getDisconnectClass = (count) => {
      if (count === 0) return 'text-green-400'
      if (count <= 3) return 'text-yellow-400'
      return 'text-red-400'
    }

    // 更新统计信息
    const updateStats = () => {
      if (!monitorResults.value.length) return

      const totalNodes = monitorResults.value.length
      const onlineNodes = monitorResults.value.filter(r => r.isAlive).length
      const offlineNodes = totalNodes - onlineNodes
      
      const totalStability = monitorResults.value.reduce((sum, r) => {
        return sum + parseFloat(getStabilityRate(r))
      }, 0)
      const avgStability = (totalStability / totalNodes).toFixed(2)
      
      const runningTime = startTime.value ? Math.floor((Date.now() - startTime.value) / 1000) : 0

      monitorStats.value = {
        totalNodes,
        onlineNodes,
        offlineNodes,
        avgStability: parseFloat(avgStability),
        runningTime
      }
    }

    // 事件监听
    const setupEventListeners = () => {
      EventsOn('monitor-start', (data) => {
        monitorStats.value.totalNodes = data.total
        startTime.value = Date.now()
      })

      EventsOn('monitor-update', (status) => {
        // 更新或添加监控状态
        const existingIndex = monitorResults.value.findIndex(r => r.proxyName === status.proxyName)
        if (existingIndex >= 0) {
          monitorResults.value[existingIndex] = status
        } else {
          monitorResults.value.push(status)
        }
        
        updateStats()
      })

      EventsOn('monitor-complete', (data) => {
        console.log('监控完成:', data)
      })

      EventsOn('monitor-error', (error) => {
        console.error('监控错误:', error)
        alert('监控错误: ' + error)
      })
    }

    onMounted(() => {
      setupEventListeners()
      
      // 定时更新运行时间
      const timer = setInterval(() => {
        if (startTime.value) {
          updateStats()
        }
      }, 1000)
      
      onUnmounted(() => {
        clearInterval(timer)
      })
    })

    onUnmounted(() => {
      EventsOff('monitor-start')
      EventsOff('monitor-update')
      EventsOff('monitor-complete')
      EventsOff('monitor-error')
    })

    return {
      monitorResults,
      monitorStats,
      monitorConfig,
      sortedMonitorResults,
      startMonitor,
      stopMonitor,
      exportMonitorReport,
      getStabilityRate,
      formatDuration,
      formatNumber,
      formatDataSize,
      getStabilityClass,
      getDisconnectClass
    }
  }
}
</script>