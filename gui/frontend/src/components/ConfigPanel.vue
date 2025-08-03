<template>
  <div class="config-panel">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- 配置文件部分 -->
      <div class="config-section">
        <h3 class="text-sm font-medium text-gray-400 mb-2">配置文件</h3>
        <div class="flex items-center space-x-2">
          <input
            v-model="configPath"
            type="text"
            class="input input-sm flex-1"
            placeholder="选择Clash配置文件..."
            readonly
          />
          <button @click="selectConfigFile" class="btn btn-sm btn-secondary" title="选择文件">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
            </svg>
          </button>
          <button 
            @click="loadConfig" 
            :disabled="!configPath || loading"
            class="btn btn-sm btn-primary"
          >
            {{ loading ? '加载中...' : '加载' }}
          </button>
        </div>
      </div>

      <!-- 节点过滤部分 -->
      <div class="config-section">
        <h3 class="text-sm font-medium text-gray-400 mb-2">节点过滤</h3>
        <div class="flex items-center space-x-2">
          <input
            v-model="filterRegex"
            type="text"
            class="input input-sm flex-1"
            placeholder="例如: HK|SG|JP"
            @keyup.enter="applyFilter"
          />
          <div class="tooltip">
            <button class="text-gray-400 hover:text-white">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </button>
            <div class="tooltip-text">
              使用正则表达式过滤节点<br>
              例如: HK|SG 匹配香港或新加坡<br>
              ^(?!.*x1) 排除包含x1的节点
            </div>
          </div>
        </div>
      </div>

      <!-- 排除关键词部分 -->
      <div class="config-section">
        <h3 class="text-sm font-medium text-gray-400 mb-2">排除关键词</h3>
        <div class="flex items-center space-x-2">
          <input
            v-model="blockKeywords"
            type="text"
            class="input input-sm flex-1"
            placeholder="例如: rate|x1|1x"
            @keyup.enter="applyFilter"
          />
          <button 
            @click="applyFilter"
            :disabled="!configInfo"
            class="btn btn-sm btn-secondary"
          >
            应用
          </button>
        </div>
      </div>
    </div>

    <!-- 配置信息显示 -->
    <Transition name="fade">
      <div v-if="configInfo" class="mt-4 p-3 bg-gray-700/50 rounded-lg">
        <div class="flex items-center justify-between text-sm">
          <div class="flex items-center space-x-4">
            <span class="text-gray-400">当前配置:</span>
            <span class="text-white font-medium">{{ configFileName }}</span>
            <span class="text-gray-500">|</span>
            <span class="text-gray-300">
              总节点: <strong class="text-white">{{ totalNodes }}</strong>
            </span>
            <span class="text-gray-500">|</span>
            <span class="text-gray-300">
              过滤后: <strong class="text-green-400">{{ filteredNodes }}</strong>
            </span>
          </div>
          <div class="flex items-center space-x-2">
            <button @click="refreshConfig" class="text-gray-400 hover:text-white" title="刷新配置">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </button>
            <button @click="showProxyList" class="text-gray-400 hover:text-white" title="查看节点列表">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script>
import { ref, computed, watch } from 'vue'
import { SelectConfigFile } from '../../wailsjs/go/main/App'

export default {
  name: 'ConfigPanel',
  props: {
    configInfo: Object
  },
  emits: ['load-config', 'config-changed'],
  setup(props, { emit }) {
    const configPath = ref('')
    const filterRegex = ref('.+')
    const blockKeywords = ref('')
    const loading = ref(false)
    
    // 计算属性
    const configFileName = computed(() => {
      if (!configPath.value) return ''
      return configPath.value.split('/').pop()
    })
    
    const totalNodes = computed(() => {
      return props.configInfo?.proxyCount || 0
    })
    
    const filteredNodes = computed(() => {
      // 这里应该根据实际过滤计算
      return props.configInfo?.proxyCount || 0
    })
    
    // 选择配置文件
    const selectConfigFile = async () => {
      try {
        const selectedFile = await SelectConfigFile()
        if (selectedFile) {
          configPath.value = selectedFile
        }
      } catch (error) {
        console.error('选择文件失败:', error)
      }
    }
    
    // 加载配置
    const loadConfig = async () => {
      if (!configPath.value) return
      
      loading.value = true
      try {
        emit('load-config', {
          configPath: configPath.value,
          filterRegex: filterRegex.value,
          blockKeywords: blockKeywords.value
        })
      } finally {
        loading.value = false
      }
    }
    
    // 应用过滤
    const applyFilter = () => {
      if (!props.configInfo) return
      
      emit('config-changed', {
        filterRegex: filterRegex.value,
        blockKeywords: blockKeywords.value
      })
    }
    
    // 刷新配置
    const refreshConfig = () => {
      if (configPath.value) {
        loadConfig()
      }
    }
    
    // 显示代理列表
    const showProxyList = () => {
      // TODO: 实现代理列表显示
      console.log('显示代理列表')
    }
    
    // 监听属性变化
    watch(() => props.configInfo, (newVal) => {
      if (newVal?.configPath) {
        configPath.value = newVal.configPath
        filterRegex.value = newVal.filter || '.+'
        blockKeywords.value = newVal.block || ''
      }
    }, { immediate: true })
    
    return {
      configPath,
      filterRegex,
      blockKeywords,
      loading,
      configFileName,
      totalNodes,
      filteredNodes,
      selectConfigFile,
      loadConfig,
      applyFilter,
      refreshConfig,
      showProxyList
    }
  }
}
</script>

<style scoped>
.config-panel {
  @apply w-full;
}

.config-section {
  @apply flex flex-col;
}

.input-sm {
  @apply px-2 py-1 text-sm;
}

.btn-sm {
  @apply px-3 py-1 text-sm;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 工具提示样式 */
.tooltip {
  @apply relative;
}

.tooltip:hover .tooltip-text {
  @apply opacity-100 visible;
}

.tooltip-text {
  @apply absolute top-full right-0 mt-2 w-64;
  @apply px-3 py-2 text-xs bg-gray-900 text-gray-300 rounded-lg shadow-lg;
  @apply opacity-0 invisible transition-all duration-200;
  @apply z-50;
  white-space: pre-line;
}
</style>