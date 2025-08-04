<template>
  <div class="flex h-full">
    <!-- 左侧配置面板 -->
    <div class="w-80 bg-gray-800 border-r border-gray-700 p-6 overflow-y-auto">
      <div class="space-y-6">
        <!-- 测速服务器 -->
        <div class="card">
          <div class="card-header">
            <h3 class="text-lg font-semibold text-white">📊 测速服务器</h3>
          </div>
          <div class="card-body space-y-4">
            <label class="flex items-center space-x-3">
              <input type="radio" v-model="testConfig.serverURL" value="https://speed.cloudflare.com" class="text-primary-500" />
              <span class="text-sm text-gray-300">Cloudflare (推荐)</span>
            </label>
            <label class="flex items-center space-x-3">
              <input type="radio" v-model="testConfig.serverURL" value="https://proof.ovh.net/files/100Mb.dat" class="text-primary-500" />
              <span class="text-sm text-gray-300">OVH (100MB)</span>
            </label>
            <label class="flex items-center space-x-3">
              <input type="radio" v-model="testConfig.serverURL" value="http://speedtest.tele2.net/100MB.zip" class="text-primary-500" />
              <span class="text-sm text-gray-300">Tele2 (100MB)</span>
            </label>
            <label class="flex items-center space-x-3">
              <input type="radio" v-model="testConfig.serverURL" value="https://ash-speed.hetzner.com/100MB.bin" class="text-primary-500" />
              <span class="text-sm text-gray-300">Hetzner (100MB)</span>
            </label>
            <label class="flex items-center space-x-3">
              <input type="radio" v-model="testConfig.serverURL" value="custom" class="text-primary-500" />
              <span class="text-sm text-gray-300">自定义</span>
            </label>
            <input
              v-if="testConfig.serverURL === 'custom'"
              v-model="customServerURL"
              type="text"
              class="input w-full mt-2"
              placeholder="输入自定义测速服务器URL"
            />
          </div>
        </div>

        <!-- 测试参数 -->
        <div class="card">
          <div class="card-header">
            <h3 class="text-lg font-semibold text-white">⚙️ 测试参数</h3>
          </div>
          <div class="card-body space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">下载大小 (MB)</label>
              <input
                v-model.number="downloadSizeMB"
                type="number"
                class="input w-full"
                min="1"
                max="1000"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">上传大小 (MB)</label>
              <input
                v-model.number="uploadSizeMB"
                type="number"
                class="input w-full"
                min="1"
                max="100"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">超时时间 (秒)</label>
              <input
                v-model.number="testConfig.timeout"
                type="number"
                class="input w-full"
                min="1"
                max="60"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">并发数</label>
              <div class="relative">
                <input
                  v-model.number="testConfig.concurrent"
                  type="range"
                  class="w-full"
                  min="1"
                  max="16"
                />
                <div class="flex justify-between text-xs text-gray-500 mt-1">
                  <span>1</span>
                  <span class="text-primary-400 font-medium">{{ testConfig.concurrent }}</span>
                  <span>16</span>
                </div>
              </div>
            </div>
            
            <div class="flex items-center space-x-3 pt-2">
              <input
                v-model="testConfig.fastMode"
                type="checkbox"
                id="fastMode"
                class="w-4 h-4 text-primary-600 rounded focus:ring-primary-500"
              />
              <label for="fastMode" class="text-sm text-gray-300">☑ 快速模式 (仅测延迟)</label>
            </div>
          </div>
        </div>

        <!-- 测试统计 -->
        <div v-if="testProgress.total > 0 || testResults.length > 0" class="card">
          <div class="card-header">
            <h3 class="text-lg font-semibold text-white">📈 测试统计</h3>
          </div>
          <div class="card-body space-y-3">
            <div class="flex justify-between">
              <span class="text-gray-300">总节点数:</span>
              <span class="text-white font-medium">{{ testProgress.total || testResults.length }}</span>
            </div>
            <div class="flex justify-between" v-if="testProgress.total > 0">
              <span class="text-gray-300">已测试:</span>
              <span class="text-blue-400 font-medium">{{ testProgress.current }}</span>
            </div>
            <div class="flex justify-between" v-if="testResults.length > 0">
              <span class="text-gray-300">合格节点:</span>
              <span class="text-green-400 font-medium">{{ passedCount }}</span>
            </div>
            <div class="flex justify-between" v-if="testResults.length > 0">
              <span class="text-gray-300">平均延迟:</span>
              <span class="text-white font-medium">{{ avgLatency }}ms</span>
            </div>
            <div class="flex justify-between" v-if="testProgress.total > 0">
              <span class="text-gray-300">进度:</span>
              <span class="text-white font-medium">{{ Math.round(testProgress.current / testProgress.total * 100) }}%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧测试面板 -->
    <div class="flex-1 flex flex-col">
      <!-- 操作栏 -->
      <div class="bg-gray-800 border-b border-gray-700 p-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <button
              @click="startTest"
              :disabled="!configInfo || isRunning"
              class="btn btn-success"
            >
              <span v-if="!isRunning" class="flex items-center space-x-2">
                <span>🚀</span>
                <span>开始测试</span>
              </span>
              <span v-else class="flex items-center space-x-2">
                <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>测试中...</span>
              </span>
            </button>
            
            <button
              @click="stopTest"
              :disabled="!isRunning"
              class="btn btn-danger"
            >
              ⏹ 停止测试
            </button>
            
            <div v-if="testProgress.total > 0" class="flex items-center space-x-3 ml-6">
              <div class="text-sm text-gray-400 whitespace-nowrap">
                进度: <span class="text-white font-medium">{{ testProgress.current }} / {{ testProgress.total }}</span>
              </div>
              <div class="w-48 progress-bar">
                <div 
                  class="progress-fill"
                  :style="{ width: (testProgress.current / testProgress.total * 100) + '%' }"
                ></div>
              </div>
              <span class="text-sm text-primary-400">{{ Math.round(testProgress.current / testProgress.total * 100) }}%</span>
            </div>
          </div>
          
          <div class="flex items-center space-x-3">
            <div v-if="testResults.length > 0" class="text-sm text-gray-400">
              合格: <span class="text-green-400 font-medium">{{ passedCount }}</span> / {{ testResults.length }}
            </div>
            <div class="relative">
              <button
                @click="showExportMenu = !showExportMenu"
                :disabled="!testResults.length"
                class="btn btn-secondary btn-sm"
              >
                导出 ▼
              </button>
              <div v-if="showExportMenu" class="absolute right-0 mt-2 w-48 bg-gray-800 rounded-lg shadow-lg z-10">
                <button @click="exportResults('csv')" class="block w-full text-left px-4 py-2 text-sm hover:bg-gray-700">
                  📄 导出为 CSV
                </button>
                <button @click="exportResults('json')" class="block w-full text-left px-4 py-2 text-sm hover:bg-gray-700">
                  📋 导出为 JSON
                </button>
                <button @click="exportResults('clash')" class="block w-full text-left px-4 py-2 text-sm hover:bg-gray-700">
                  ⚡ 导出为 Clash 配置
                </button>
              </div>
            </div>
            <button
              @click="clearResults"
              :disabled="!testResults.length"
              class="btn btn-secondary btn-sm"
            >
              清空
            </button>
          </div>
        </div>
      </div>

      <!-- 结果列表 -->
      <div class="flex-1 overflow-y-auto">
        <div v-if="!configInfo" class="text-center py-12">
          <div class="text-gray-400 text-lg mb-4">📁</div>
          <p class="text-gray-400">请先选择并加载配置文件</p>
        </div>
        
        <div v-else-if="!testResults.length && !isRunning" class="text-center py-12">
          <div class="text-gray-400 text-lg mb-4">🚀</div>
          <p class="text-gray-400">点击"开始测试"开始速度测试</p>
        </div>
        
        <div v-else>
          <!-- 结果表格 -->
          <div class="bg-gray-800 rounded-lg overflow-hidden">
            <table class="table">
              <thead>
                <tr>
                  <th class="w-8">#</th>
                  <th>节点名称</th>
                  <th class="w-20">类型</th>
                  <th class="w-24">延迟</th>
                  <th v-if="!testConfig.fastMode" class="w-24">抖动</th>
                  <th v-if="!testConfig.fastMode" class="w-24">丢包率</th>
                  <th v-if="!testConfig.fastMode" class="w-32">下载速度</th>
                  <th v-if="!testConfig.fastMode" class="w-32">上传速度</th>
                  <th class="w-20">状态</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(result, index) in sortedResults" :key="index">
                  <td class="text-gray-400">{{ index + 1 }}</td>
                  <td class="font-medium">{{ result.proxyName }}</td>
                  <td class="text-gray-400">{{ result.proxyType }}</td>
                  <td :class="getLatencyClass(result.latency)">
                    {{ formatLatency(result.latency) }}
                  </td>
                  <td v-if="!testConfig.fastMode" :class="getLatencyClass(result.jitter)">
                    {{ formatLatency(result.jitter) }}
                  </td>
                  <td v-if="!testConfig.fastMode" :class="getPacketLossClass(result.packetLoss)">
                    {{ formatPacketLoss(result.packetLoss) }}
                  </td>
                  <td v-if="!testConfig.fastMode" :class="getSpeedClass(result.downloadSpeed)">
                    {{ formatSpeed(result.downloadSpeed) }}
                  </td>
                  <td v-if="!testConfig.fastMode" :class="getSpeedClass(result.uploadSpeed)">
                    {{ formatSpeed(result.uploadSpeed) }}
                  </td>
                  <td>
                    <span v-if="result.status === '完成'" class="text-green-400">✅</span>
                    <span v-else-if="result.status === '测试中'" class="text-yellow-400">⏳</span>
                    <span v-else class="text-red-400">❌</span>
                  </td>
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
import { SaveReport, ClearHistory } from '../../wailsjs/go/main/App'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'

export default {
  name: 'SpeedTestView',
  props: {
    configInfo: Object,
    isRunning: Boolean
  },
  emits: ['start-test', 'stop-test'],
  setup(props, { emit }) {
    const testResults = ref([])
    const testProgress = ref({ current: 0, total: 0 })
    const showExportMenu = ref(false)
    
    const downloadSizeMB = ref(50)
    const uploadSizeMB = ref(20)
    const customServerURL = ref('')
    
    const testConfig = ref({
      serverURL: 'https://speed.cloudflare.com',
      timeout: 5,
      concurrent: 4,
      maxLatency: 800,
      minDownloadSpeed: 5,
      minUploadSpeed: 2,
      fastMode: false
    })

    // 计算属性
    const sortedResults = computed(() => {
      return [...testResults.value].sort((a, b) => {
        if (testConfig.value.fastMode) {
          return a.latency - b.latency
        }
        return b.downloadSpeed - a.downloadSpeed
      })
    })
    
    const passedCount = computed(() => {
      return testResults.value.filter(r => {
        if (r.status !== '完成') return false
        if (r.latency > testConfig.value.maxLatency) return false
        if (!testConfig.value.fastMode) {
          const downloadMBps = r.downloadSpeed / (1024 * 1024)
          const uploadMBps = r.uploadSpeed / (1024 * 1024)
          if (downloadMBps < testConfig.value.minDownloadSpeed) return false
          if (uploadMBps < testConfig.value.minUploadSpeed) return false
        }
        return true
      }).length
    })

    const avgLatency = computed(() => {
      const validResults = testResults.value.filter(r => r.status === '完成' && r.latency > 0)
      if (validResults.length === 0) return 0
      const totalLatency = validResults.reduce((sum, r) => sum + r.latency, 0)
      return Math.round(totalLatency / validResults.length)
    })

    // 开始测试
    const startTest = () => {
      if (!props.configInfo) return
      
      const serverURL = testConfig.value.serverURL === 'custom' 
        ? customServerURL.value 
        : testConfig.value.serverURL
      
      const config = {
        configPath: props.configInfo.configPath,
        filterRegex: props.configInfo.filter || '.+',
        blockRegex: props.configInfo.block || '',
        serverURL: serverURL,
        downloadSize: downloadSizeMB.value * 1024 * 1024,
        uploadSize: uploadSizeMB.value * 1024 * 1024,
        timeout: testConfig.value.timeout,
        concurrent: testConfig.value.concurrent,
        maxLatency: testConfig.value.maxLatency,
        minDownloadSpeed: testConfig.value.minDownloadSpeed,
        minUploadSpeed: testConfig.value.minUploadSpeed,
        fastMode: testConfig.value.fastMode
      }
      
      testResults.value = []
      testProgress.value = { current: 0, total: 0 }
      emit('start-test', config)
    }

    // 停止测试
    const stopTest = () => {
      emit('stop-test')
    }
    
    // 清空结果
    const clearResults = async () => {
      const message = props.isRunning 
        ? '确定要清空所有测试结果吗？这将停止当前测试并清空历史数据。'
        : '确定要清空所有测试结果和历史数据吗？'
        
      if (confirm(message)) {
        try {
          // 如果测试正在进行，先停止测试
          if (props.isRunning) {
            emit('stop-test')
            // 等待一小段时间确保后端处理停止请求
            await new Promise(resolve => setTimeout(resolve, 500))
          }
          
          // 清空前端数据
          testResults.value = []
          testProgress.value = { current: 0, total: 0 }
          showExportMenu.value = false
          
          // 清空后端历史数据
          await ClearHistory()
          console.log('历史数据已清空')
          
        } catch (error) {
          console.error('清空失败:', error)
          alert('清空失败，请重试: ' + error)
        }
      }
    }

    // 导出结果
    const exportResults = async (format) => {
      showExportMenu.value = false
      
      if (!testResults.value.length) return
      
      try {
        const timestamp = new Date().toISOString().replace(/[:.]/g, '-')
        const filename = `clash-speedtest-${timestamp}.${format}`
        
        await SaveReport(testResults.value, filename)
        console.log(`导出成功: ${filename}`)
      } catch (error) {
        console.error('导出失败:', error)
        alert('导出失败: ' + error)
      }
    }

    // 格式化函数
    const formatLatency = (latency) => {
      if (latency <= 0) return 'N/A'
      return `${latency}ms`
    }

    const formatSpeed = (speed) => {
      if (speed <= 0) return 'N/A'
      const mbps = speed / (1024 * 1024)
      return `${mbps.toFixed(2)} MB/s`
    }

    const formatPacketLoss = (loss) => {
      return `${loss.toFixed(1)}%`
    }

    // 样式类函数
    const getLatencyClass = (latency) => {
      if (latency <= 0) return 'text-gray-400'
      if (latency < 100) return 'text-latency-good'
      if (latency < 300) return 'text-latency-fair'
      return 'text-latency-poor'
    }

    const getSpeedClass = (speed) => {
      if (speed <= 0) return 'text-gray-400'
      const mbps = speed / (1024 * 1024)
      if (mbps >= 10) return 'text-speed-fast'
      if (mbps >= 5) return 'text-speed-medium'
      return 'text-speed-slow'
    }

    const getPacketLossClass = (loss) => {
      if (loss < 1) return 'text-speed-fast'
      if (loss < 5) return 'text-speed-medium'
      return 'text-speed-slow'
    }

    // 事件监听
    const setupEventListeners = () => {
      EventsOn('test-start', (data) => {
        testProgress.value = { current: 0, total: data.total }
        testResults.value = []
      })

      EventsOn('test-progress', (data) => {
        testProgress.value = { current: data.current, total: data.total }
        
        // 更新或添加结果
        const existingIndex = testResults.value.findIndex(r => r.proxyName === data.result.proxyName)
        if (existingIndex >= 0) {
          testResults.value[existingIndex] = data.result
        } else {
          testResults.value.push(data.result)
        }
      })

      EventsOn('test-complete', (data) => {
        testProgress.value.current = testProgress.value.total
      })

      EventsOn('test-error', (error) => {
        console.error('测试错误:', error)
        alert('测试错误: ' + error)
      })
    }

    onMounted(() => {
      setupEventListeners()
      
      // 点击外部关闭导出菜单
      const handleClickOutside = (event) => {
        if (!event.target.closest('.relative')) {
          showExportMenu.value = false
        }
      }
      document.addEventListener('click', handleClickOutside)
      
      onUnmounted(() => {
        document.removeEventListener('click', handleClickOutside)
      })
    })

    onUnmounted(() => {
      EventsOff('test-start')
      EventsOff('test-progress')
      EventsOff('test-complete')
      EventsOff('test-error')
    })

    return {
      testResults,
      testProgress,
      showExportMenu,
      downloadSizeMB,
      uploadSizeMB,
      customServerURL,
      testConfig,
      sortedResults,
      passedCount,
      avgLatency,
      startTest,
      stopTest,
      clearResults,
      exportResults,
      formatLatency,
      formatSpeed,
      formatPacketLoss,
      getLatencyClass,
      getSpeedClass,
      getPacketLossClass
    }
  }
}
</script>

<style scoped>
/* 小按钮样式 */
.btn-sm {
  @apply px-3 py-1.5 text-sm;
}

/* 卡片样式调整 */
.card {
  @apply bg-gray-700/50 rounded-lg;
}

.card-header {
  @apply px-4 py-2 border-b border-gray-600;
}

.card-body {
  @apply px-4 py-3;
}

/* 滑块样式 */
input[type="range"] {
  @apply h-2 bg-gray-600 rounded-lg appearance-none cursor-pointer;
}

input[type="range"]::-webkit-slider-thumb {
  @apply appearance-none w-4 h-4 bg-primary-500 rounded-full cursor-pointer;
}

input[type="range"]::-moz-range-thumb {
  @apply w-4 h-4 bg-primary-500 rounded-full cursor-pointer border-0;
}
</style>