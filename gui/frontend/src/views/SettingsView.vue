<template>
  <div class="flex h-full">
    <!-- 设置内容 -->
    <div class="flex-1 p-6 overflow-y-auto">
      <div class="max-w-4xl mx-auto space-y-8">
        <!-- 页面标题 -->
        <div>
          <h1 class="text-2xl font-bold text-white mb-2">⚙️ 应用设置</h1>
          <p class="text-gray-400">配置应用的默认参数和行为</p>
        </div>

        <!-- 配置文件设置 -->
        <div class="card">
          <div class="card-header">
            <h2 class="text-xl font-semibold text-white">📁 配置文件</h2>
          </div>
          <div class="card-body space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">默认配置文件</label>
                <div class="flex space-x-2">
                  <input
                    v-model="settings.defaultConfigPath"
                    type="text"
                    class="input flex-1"
                    placeholder="选择默认配置文件..."
                    readonly
                  />
                  <button @click="selectDefaultConfig" class="btn btn-secondary whitespace-nowrap">
                    选择
                  </button>
                </div>
                <p class="text-xs text-gray-400 mt-1">启动时自动加载的配置文件</p>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">默认过滤规则</label>
                <input
                  v-model="settings.defaultFilter"
                  type="text"
                  class="input w-full"
                  placeholder="例如: HK|SG|JP"
                />
                <p class="text-xs text-gray-400 mt-1">默认的节点过滤正则表达式</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 测速默认设置 -->
        <div class="card">
          <div class="card-header">
            <h2 class="text-xl font-semibold text-white">🚀 测速设置</h2>
          </div>
          <div class="card-body space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">默认测速服务器</label>
                <select v-model="settings.defaultServerURL" class="input w-full">
                  <option value="https://speed.cloudflare.com">Cloudflare (推荐)</option>
                  <option value="https://fast.com">Fast.com</option>
                  <option value="https://speedtest.net">Speedtest.net</option>
                </select>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">默认并发数</label>
                <input
                  v-model.number="settings.defaultConcurrent"
                  type="number"
                  class="input w-full"
                  min="1"
                  max="16"
                />
              </div>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">下载测试大小 (MB)</label>
                <input
                  v-model.number="settings.defaultDownloadSize"
                  type="number"
                  class="input w-full"
                  min="1"
                  max="1000"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">上传测试大小 (MB)</label>
                <input
                  v-model.number="settings.defaultUploadSize"
                  type="number"
                  class="input w-full"
                  min="1"
                  max="100"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">超时时间 (秒)</label>
                <input
                  v-model.number="settings.defaultTimeout"
                  type="number"
                  class="input w-full"
                  min="1"
                  max="60"
                />
              </div>
            </div>
            
            <div class="flex items-center space-x-2">
              <input
                v-model="settings.defaultFastMode"
                type="checkbox"
                id="defaultFastMode"
                class="w-4 h-4 text-primary-600 rounded focus:ring-primary-500"
              />
              <label for="defaultFastMode" class="text-sm text-gray-300">默认启用快速模式</label>
            </div>
          </div>
        </div>

        <!-- 监控默认设置 -->
        <div class="card">
          <div class="card-header">
            <h2 class="text-xl font-semibold text-white">📡 监控设置</h2>
          </div>
          <div class="card-body space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">默认监控类型</label>
                <select v-model="settings.defaultMonitorType" class="input w-full">
                  <option value="http">HTTP 长连接</option>
                  <option value="websocket">WebSocket 数据流</option>
                </select>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">默认监控时长</label>
                <select v-model.number="settings.defaultMonitorDuration" class="input w-full">
                  <option :value="300">5 分钟</option>
                  <option :value="600">10 分钟</option>
                  <option :value="1800">30 分钟</option>
                  <option :value="3600">1 小时</option>
                  <option :value="14400">4 小时</option>
                  <option :value="86400">24 小时</option>
                </select>
              </div>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">心跳间隔 (秒)</label>
                <input
                  v-model.number="settings.defaultMonitorInterval"
                  type="number"
                  class="input w-full"
                  min="1"
                  max="60"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">HTTP目标地址</label>
                <input
                  v-model="settings.defaultHttpTarget"
                  type="text"
                  class="input w-full"
                  placeholder="https://cn.tradingview.com/chart/"
                />
              </div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">WebSocket目标地址</label>
              <input
                v-model="settings.defaultWebSocketTarget"
                type="text"
                class="input w-full"
                placeholder="wss://stream.binance.com:9443/ws/btcusdt@ticker"
              />
            </div>
          </div>
        </div>

        <!-- 界面设置 -->
        <div class="card">
          <div class="card-header">
            <h2 class="text-xl font-semibold text-white">🎨 界面设置</h2>
          </div>
          <div class="card-body space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">主题模式</label>
                <select v-model="settings.theme" class="input w-full">
                  <option value="dark">深色模式</option>
                  <option value="light">浅色模式</option>
                  <option value="auto">跟随系统</option>
                </select>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-300 mb-2">语言</label>
                <select v-model="settings.language" class="input w-full">
                  <option value="zh-CN">简体中文</option>
                  <option value="en-US">English</option>
                </select>
              </div>
            </div>
            
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <div>
                  <h3 class="text-sm font-medium text-gray-300">启动时最小化到托盘</h3>
                  <p class="text-xs text-gray-400">应用启动后自动最小化到系统托盘</p>
                </div>
                <input
                  v-model="settings.startMinimized"
                  type="checkbox"
                  class="w-4 h-4 text-primary-600 rounded focus:ring-primary-500"
                />
              </div>
              
              <div class="flex items-center justify-between">
                <div>
                  <h3 class="text-sm font-medium text-gray-300">自动保存测试结果</h3>
                  <p class="text-xs text-gray-400">测试完成后自动保存结果到文件</p>
                </div>
                <input
                  v-model="settings.autoSaveResults"
                  type="checkbox"
                  class="w-4 h-4 text-primary-600 rounded focus:ring-primary-500"
                />
              </div>
              
              <div class="flex items-center justify-between">
                <div>
                  <h3 class="text-sm font-medium text-gray-300">显示系统通知</h3>
                  <p class="text-xs text-gray-400">测试或监控完成时显示系统通知</p>
                </div>
                <input
                  v-model="settings.showNotifications"
                  type="checkbox"
                  class="w-4 h-4 text-primary-600 rounded focus:ring-primary-500"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 系统信息 -->
        <div class="card">
          <div class="card-header">
            <h2 class="text-xl font-semibold text-white">ℹ️ 系统信息</h2>
          </div>
          <div class="card-body">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-3">
                <div class="flex justify-between">
                  <span class="text-gray-300">应用版本:</span>
                  <span class="text-white font-medium">{{ systemInfo.version }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-300">构建时间:</span>
                  <span class="text-white font-medium">{{ systemInfo.buildTime }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-300">运行平台:</span>
                  <span class="text-white font-medium">{{ systemInfo.platform }}</span>
                </div>
              </div>
              
              <div class="space-y-3">
                <div class="flex justify-between">
                  <span class="text-gray-300">框架版本:</span>
                  <span class="text-white font-medium">Wails v2.8.0</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-300">Go版本:</span>
                  <span class="text-white font-medium">Go 1.24</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-300">Vue版本:</span>
                  <span class="text-white font-medium">Vue 3.4</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex justify-end space-x-4">
          <button @click="resetSettings" class="btn btn-secondary">
            🔄 重置默认
          </button>
          <button @click="saveSettings" class="btn btn-primary">
            💾 保存设置
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { SelectConfigFile } from '../../wailsjs/go/main/App'

export default {
  name: 'SettingsView',
  props: {
    systemInfo: Object
  },
  emits: ['config-changed'],
  setup(props, { emit }) {
    // 使用传入的系统信息，如果没有则使用默认值
    const systemInfo = computed(() => props.systemInfo || {
      version: '1.0.0',
      buildTime: '2025-01-20 12:00:00',
      platform: 'Unknown'
    })

    const settings = reactive({
      // 配置文件设置
      defaultConfigPath: '',
      defaultFilter: '.+',
      
      // 测速设置
      defaultServerURL: 'https://speed.cloudflare.com',
      defaultConcurrent: 4,
      defaultDownloadSize: 50,
      defaultUploadSize: 20,
      defaultTimeout: 5,
      defaultFastMode: false,
      
      // 监控设置
      defaultMonitorType: 'websocket',
      defaultMonitorDuration: 3600,
      defaultMonitorInterval: 1,
      defaultHttpTarget: 'https://cn.tradingview.com/chart/',
      defaultWebSocketTarget: 'wss://stream.binance.com:9443/ws/btcusdt@ticker',
      
      // 界面设置
      theme: 'dark',
      language: 'zh-CN',
      startMinimized: false,
      autoSaveResults: true,
      showNotifications: true
    })

    // 选择默认配置文件
    const selectDefaultConfig = async () => {
      try {
        const selectedFile = await SelectConfigFile()
        if (selectedFile) {
          settings.defaultConfigPath = selectedFile
        }
      } catch (error) {
        console.error('选择文件失败:', error)
      }
    }

    // 保存设置
    const saveSettings = () => {
      try {
        // 保存到localStorage
        localStorage.setItem('app-settings', JSON.stringify(settings))
        alert('设置已保存')
      } catch (error) {
        console.error('保存设置失败:', error)
        alert('保存设置失败: ' + error)
      }
    }

    // 重置设置
    const resetSettings = () => {
      if (confirm('确定要重置所有设置到默认值吗？')) {
        Object.assign(settings, {
          defaultConfigPath: '',
          defaultFilter: '.+',
          defaultServerURL: 'https://speed.cloudflare.com',
          defaultConcurrent: 4,
          defaultDownloadSize: 50,
          defaultUploadSize: 20,
          defaultTimeout: 5,
          defaultFastMode: false,
          defaultMonitorType: 'websocket',
          defaultMonitorDuration: 3600,
          defaultMonitorInterval: 1,
          defaultHttpTarget: 'https://cn.tradingview.com/chart/',
          defaultWebSocketTarget: 'wss://stream.binance.com:9443/ws/btcusdt@ticker',
          theme: 'dark',
          language: 'zh-CN',
          startMinimized: false,
          autoSaveResults: true,
          showNotifications: true
        })
      }
    }

    // 加载设置
    const loadSettings = () => {
      try {
        const savedSettings = localStorage.getItem('app-settings')
        if (savedSettings) {
          const parsed = JSON.parse(savedSettings)
          Object.assign(settings, parsed)
        }
      } catch (error) {
        console.error('加载设置失败:', error)
      }
    }

    onMounted(() => {
      loadSettings()
    })

    return {
      settings,
      systemInfo,
      selectDefaultConfig,
      saveSettings,
      resetSettings
    }
  }
}
</script>