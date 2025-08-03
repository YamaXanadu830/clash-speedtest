<template>
  <div id="app" class="flex flex-col h-screen bg-gray-900">
    <!-- 顶部状态栏 -->
    <header class="bg-gray-800 border-b border-gray-700 px-4 py-3">
      <div class="flex items-center justify-between">
        <!-- 左侧：Logo和标题 -->
        <div class="flex items-center space-x-4">
          <div class="flex items-center space-x-2">
            <span class="text-2xl">⚡</span>
            <h1 class="text-xl font-bold text-white">Clash SpeedTest</h1>
          </div>
          <span class="text-sm text-gray-400">v{{ systemInfo.version || '1.0.0' }}</span>
        </div>
        
        <!-- 中间：全局状态 -->
        <div class="flex items-center space-x-6 text-sm">
          <div v-if="configInfo" class="flex items-center space-x-2">
            <span class="text-gray-400">配置:</span>
            <span class="text-green-400">{{ configFileName }}</span>
            <span class="text-gray-500">({{ configInfo.proxyCount }} 节点)</span>
          </div>
          <div v-if="isRunning" class="flex items-center space-x-2">
            <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <span class="text-gray-300">{{ runningTask }}</span>
          </div>
        </div>
        
        <!-- 右侧：快捷操作 -->
        <div class="flex items-center space-x-3">
          <button @click="toggleConfigPanel" class="text-gray-400 hover:text-white transition-colors" title="配置面板 (Ctrl+,)">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </button>
          <button @click="showHelp" class="text-gray-400 hover:text-white transition-colors" title="帮助">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </button>
        </div>
      </div>
    </header>

    <!-- 主工作区 -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <!-- 统一配置面板（可折叠） -->
      <Transition name="slide-down">
        <div v-if="showConfigPanel" class="bg-gray-800 border-b border-gray-700 p-4">
          <div class="max-w-7xl mx-auto">
            <ConfigPanel 
              :config-info="configInfo"
              @load-config="handleLoadConfig"
              @config-changed="handleConfigChanged"
            />
          </div>
        </div>
      </Transition>

      <!-- 功能标签页 -->
      <div class="bg-gray-800 border-b border-gray-700">
        <div class="flex items-center px-4">
          <nav class="flex space-x-1">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              @click="activeTab = tab.id"
              :class="[
                'px-4 py-3 text-sm font-medium transition-all duration-200',
                'border-b-2 hover:text-white',
                activeTab === tab.id 
                  ? 'text-white border-primary-500' 
                  : 'text-gray-400 border-transparent hover:border-gray-600'
              ]"
            >
              <span class="flex items-center space-x-2">
                <span>{{ tab.icon }}</span>
                <span>{{ tab.name }}</span>
                <span v-if="tab.badge" class="ml-2 px-2 py-0.5 text-xs bg-red-600 text-white rounded-full">
                  {{ tab.badge }}
                </span>
              </span>
            </button>
          </nav>
          
          <!-- 标签页右侧快捷操作 -->
          <div class="ml-auto flex items-center space-x-3">
            <button
              v-if="activeTab === 'speedtest' && configInfo"
              @click="quickTest"
              class="text-sm text-primary-400 hover:text-primary-300"
            >
              🚀 快速测试
            </button>
            <button
              v-if="activeTab === 'monitor' && configInfo"
              @click="quickMonitor"
              class="text-sm text-primary-400 hover:text-primary-300"
            >
              📡 快速监控
            </button>
          </div>
        </div>
      </div>

      <!-- 动态内容区 -->
      <div class="flex-1 overflow-hidden">
        <!-- 速度测试页面 -->
        <keep-alive>
          <SpeedTestView 
            v-if="activeTab === 'speedtest'" 
            :config-info="configInfo"
            :is-running="isRunning"
            @start-test="handleStartTest"
            @stop-test="handleStopTest"
          />
        </keep-alive>
        
        <!-- 监控页面 -->
        <keep-alive>
          <MonitorView 
            v-if="activeTab === 'monitor'"
            :config-info="configInfo"
            :is-running="isRunning"
            @start-monitor="handleStartMonitor"
            @stop-monitor="handleStopTest"
          />
        </keep-alive>
        
        <!-- 批量管理页面 -->
        <keep-alive>
          <BatchManageView
            v-if="activeTab === 'batch'"
            :config-info="configInfo"
          />
        </keep-alive>
        
        <!-- 数据分析页面 -->
        <keep-alive>
          <AnalysisView
            v-if="activeTab === 'analysis'"
            :config-info="configInfo"
          />
        </keep-alive>
        
        <!-- 设置页面 -->
        <keep-alive>
          <SettingsView 
            v-if="activeTab === 'settings'"
            :system-info="systemInfo"
            @config-changed="handleSettingsChanged"
          />
        </keep-alive>
      </div>
    </main>

    <!-- 底部状态栏 -->
    <footer class="bg-gray-800 border-t border-gray-700 px-4 py-2">
      <div class="flex items-center justify-between text-xs text-gray-400">
        <div class="flex items-center space-x-4">
          <span>{{ currentTime }}</span>
          <span v-if="lastTestTime">上次测试: {{ lastTestTime }}</span>
        </div>
        <div class="flex items-center space-x-4">
          <span v-if="cpuUsage">CPU: {{ cpuUsage }}%</span>
          <span v-if="memoryUsage">内存: {{ memoryUsage }}MB</span>
          <span>{{ systemInfo.platform || 'Wails' }}</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import ConfigPanel from './components/ConfigPanel.vue'
import SpeedTestView from './views/SpeedTestView.vue'
import MonitorView from './views/MonitorView.vue'
import BatchManageView from './views/BatchManageView.vue'
import AnalysisView from './views/AnalysisView.vue'
import SettingsView from './views/SettingsView.vue'
import { GetSystemInfo, LoadConfig, StartSpeedTest, StartMonitor, StopTest } from '../wailsjs/go/main/App'
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime'

export default {
  name: 'App',
  components: {
    ConfigPanel,
    SpeedTestView,
    MonitorView,
    BatchManageView,
    AnalysisView,
    SettingsView
  },
  setup() {
    // 系统信息
    const systemInfo = ref({
      platform: '',
      version: '',
      buildTime: ''
    })
    
    // 配置信息
    const configInfo = ref(null)
    const configFileName = computed(() => {
      if (!configInfo.value?.configPath) return ''
      return configInfo.value.configPath.split('/').pop()
    })
    
    // 界面状态
    const showConfigPanel = ref(true)
    const activeTab = ref('speedtest')
    const isRunning = ref(false)
    const runningTask = ref('')
    
    // 时间显示
    const currentTime = ref('')
    const lastTestTime = ref('')
    
    // 系统资源
    const cpuUsage = ref(0)
    const memoryUsage = ref(0)
    
    // 标签页配置
    const tabs = ref([
      { id: 'speedtest', name: '速度测试', icon: '🚀' },
      { id: 'monitor', name: '稳定监控', icon: '📡' },
      { id: 'batch', name: '批量管理', icon: '📦', badge: null },
      { id: 'analysis', name: '数据分析', icon: '📊' },
      { id: 'settings', name: '设置', icon: '⚙️' }
    ])
    
    // 切换配置面板
    const toggleConfigPanel = () => {
      showConfigPanel.value = !showConfigPanel.value
    }
    
    // 显示帮助
    const showHelp = () => {
      console.log('显示帮助')
    }
    
    // 处理配置加载
    const handleLoadConfig = async (configData) => {
      try {
        const config = await LoadConfig(
          configData.configPath, 
          configData.filterRegex || '.+', 
          configData.blockKeywords || ''
        )
        configInfo.value = config
        lastTestTime.value = new Date().toLocaleTimeString()
        
        // 保存到本地存储
        localStorage.setItem('lastConfigPath', configData.configPath)
        localStorage.setItem('lastFilterRegex', configData.filterRegex || '.+')
        localStorage.setItem('lastBlockKeywords', configData.blockKeywords || '')
      } catch (error) {
        console.error('加载配置失败:', error)
        alert('加载配置失败: ' + error)
      }
    }
    
    // 处理配置变更
    const handleConfigChanged = async (changes) => {
      if (!configInfo.value) return
      
      try {
        const config = await LoadConfig(
          configInfo.value.configPath,
          changes.filterRegex || '.+',
          changes.blockKeywords || ''
        )
        configInfo.value = config
      } catch (error) {
        console.error('更新配置失败:', error)
      }
    }
    
    // 处理开始测试
    const handleStartTest = async (testConfig) => {
      try {
        await StartSpeedTest(testConfig)
        isRunning.value = true
        runningTask.value = '速度测试中'
      } catch (error) {
        console.error('启动测试失败:', error)
        alert('启动测试失败: ' + error)
      }
    }
    
    // 处理停止测试
    const handleStopTest = async () => {
      try {
        await StopTest()
        isRunning.value = false
        runningTask.value = ''
        lastTestTime.value = new Date().toLocaleTimeString()
      } catch (error) {
        console.error('停止失败:', error)
        alert('停止失败: ' + error)
      }
    }
    
    // 处理开始监控
    const handleStartMonitor = async (monitorConfig) => {
      try {
        await StartMonitor(monitorConfig)
        isRunning.value = true
        runningTask.value = '稳定性监控中'
      } catch (error) {
        console.error('启动监控失败:', error)
        alert('启动监控失败: ' + error)
      }
    }
    
    // 处理设置变更
    const handleSettingsChanged = (settings) => {
      console.log('设置变更:', settings)
    }
    
    // 快速测试
    const quickTest = async () => {
      if (!configInfo.value) return
      
      const defaultConfig = {
        configPath: configInfo.value.configPath,
        filterRegex: localStorage.getItem('lastFilterRegex') || '.+',
        blockRegex: localStorage.getItem('lastBlockKeywords') || '',
        serverURL: 'https://speed.cloudflare.com',
        downloadSize: 50 * 1024 * 1024,
        uploadSize: 20 * 1024 * 1024,
        timeout: 5,
        concurrent: 4,
        maxLatency: 800,
        minDownloadSpeed: 5,
        minUploadSpeed: 2,
        fastMode: false
      }
      
      await handleStartTest(defaultConfig)
    }
    
    // 快速监控
    const quickMonitor = async () => {
      if (!configInfo.value) return
      
      const defaultConfig = {
        type: 'websocket',
        targetURL: 'wss://stream.binance.com:9443/ws/btcusdt@ticker',
        duration: 3600,
        interval: 1
      }
      
      await handleStartMonitor(defaultConfig)
    }
    
    // 更新时间
    const updateTime = () => {
      currentTime.value = new Date().toLocaleTimeString()
    }
    
    // 设置键盘快捷键
    const setupKeyboardShortcuts = () => {
      const handleKeyPress = (event) => {
        // Ctrl+O - 打开配置文件
        if (event.ctrlKey && event.key === 'o') {
          event.preventDefault()
          showConfigPanel.value = true
        }
        // Ctrl+T - 开始/停止测试
        else if (event.ctrlKey && event.key === 't') {
          event.preventDefault()
          if (isRunning.value) {
            handleStopTest()
          } else if (activeTab.value === 'speedtest' && configInfo.value) {
            quickTest()
          }
        }
        // Ctrl+M - 开始/停止监控
        else if (event.ctrlKey && event.key === 'm') {
          event.preventDefault()
          if (isRunning.value) {
            handleStopTest()
          } else if (activeTab.value === 'monitor' && configInfo.value) {
            quickMonitor()
          }
        }
        // Ctrl+, - 打开设置
        else if (event.ctrlKey && event.key === ',') {
          event.preventDefault()
          toggleConfigPanel()
        }
        // F5 - 刷新配置
        else if (event.key === 'F5') {
          event.preventDefault()
          if (configInfo.value) {
            handleLoadConfig({
              configPath: configInfo.value.configPath,
              filterRegex: configInfo.value.filter || '.+',
              blockKeywords: configInfo.value.block || ''
            })
          }
        }
        // ESC - 停止当前操作
        else if (event.key === 'Escape' && isRunning.value) {
          handleStopTest()
        }
      }
      
      document.addEventListener('keydown', handleKeyPress)
      return () => document.removeEventListener('keydown', handleKeyPress)
    }
    
    // 设置事件监听
    const setupEventListeners = () => {
      const eventHandlers = {
        'test-start': () => {
          isRunning.value = true
          runningTask.value = '速度测试中'
        },
        'test-complete': () => {
          isRunning.value = false
          runningTask.value = ''
          lastTestTime.value = new Date().toLocaleTimeString()
        },
        'test-error': () => {
          isRunning.value = false
          runningTask.value = ''
        },
        'test-stopped': () => {
          isRunning.value = false
          runningTask.value = ''
        },
        'monitor-start': () => {
          isRunning.value = true
          runningTask.value = '稳定性监控中'
        },
        'monitor-complete': () => {
          isRunning.value = false
          runningTask.value = ''
          lastTestTime.value = new Date().toLocaleTimeString()
        },
        'monitor-error': () => {
          isRunning.value = false
          runningTask.value = ''
        }
      }
      
      // 注册所有事件监听器
      Object.entries(eventHandlers).forEach(([eventName, handler]) => {
        EventsOn(eventName, handler)
      })
      
      return () => {
        Object.keys(eventHandlers).forEach(eventName => {
          EventsOff(eventName)
        })
      }
    }
    
    // 初始化
    onMounted(async () => {
      // 获取系统信息
      try {
        systemInfo.value = await GetSystemInfo()
      } catch (error) {
        console.error('获取系统信息失败:', error)
      }
      
      // 恢复上次配置
      const lastConfigPath = localStorage.getItem('lastConfigPath')
      if (lastConfigPath) {
        handleLoadConfig({
          configPath: lastConfigPath,
          filterRegex: localStorage.getItem('lastFilterRegex') || '.+',
          blockKeywords: localStorage.getItem('lastBlockKeywords') || ''
        })
      }
      
      // 启动时间更新
      updateTime()
      const timer = setInterval(updateTime, 1000)
      
      // 设置键盘快捷键
      const cleanupKeyboard = setupKeyboardShortcuts()
      
      // 设置事件监听
      const cleanupEvents = setupEventListeners()
      
      onUnmounted(() => {
        clearInterval(timer)
        cleanupKeyboard()
        cleanupEvents()
      })
    })
    
    return {
      systemInfo,
      configInfo,
      configFileName,
      showConfigPanel,
      activeTab,
      isRunning,
      runningTask,
      currentTime,
      lastTestTime,
      cpuUsage,
      memoryUsage,
      tabs,
      toggleConfigPanel,
      showHelp,
      handleLoadConfig,
      handleConfigChanged,
      handleStartTest,
      handleStopTest,
      handleStartMonitor,
      handleSettingsChanged,
      quickTest,
      quickMonitor
    }
  }
}
</script>

<style>
/* 过渡动画 */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from {
  transform: translateY(-100%);
  opacity: 0;
}

.slide-down-leave-to {
  transform: translateY(-100%);
  opacity: 0;
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #1f2937;
}

::-webkit-scrollbar-thumb {
  background: #4b5563;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #6b7280;
}
</style>