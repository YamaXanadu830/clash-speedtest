<template>
  <div class="flex flex-col h-full p-6">
    <div class="max-w-7xl mx-auto w-full">
      <!-- 页面标题 -->
      <div class="mb-6 flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-white mb-2">📊 数据分析</h1>
          <p class="text-gray-400">查看历史测试数据，分析节点性能趋势</p>
        </div>
        <button
          @click="fetchStats"
          :disabled="isLoading"
          class="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <span v-if="isLoading">🔄 刷新中...</span>
          <span v-else>🔄 刷新数据</span>
        </button>
      </div>

      <!-- 错误提示 -->
      <div v-if="error" class="mb-6 p-4 bg-red-900/50 border border-red-500 rounded-lg">
        <div class="flex items-center">
          <span class="text-2xl mr-3">❌</span>
          <div>
            <p class="text-red-400 font-medium">数据加载失败</p>
            <p class="text-red-300 text-sm">{{ error }}</p>
          </div>
        </div>
      </div>

      <!-- 统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div class="bg-gray-800 rounded-lg p-4 border border-gray-700">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-400 text-sm">总测试次数</p>
              <p class="text-2xl font-bold text-white mt-1">
                <span v-if="isLoading">--</span>
                <span v-else>{{ stats.totalTests }}</span>
              </p>
            </div>
            <div class="text-3xl">📈</div>
          </div>
        </div>
        
        <div class="bg-gray-800 rounded-lg p-4 border border-gray-700">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-400 text-sm">测试节点数</p>
              <p class="text-2xl font-bold text-white mt-1">
                <span v-if="isLoading">--</span>
                <span v-else>{{ stats.totalNodes }}</span>
              </p>
            </div>
            <div class="text-3xl">🌐</div>
          </div>
        </div>
        
        <div class="bg-gray-800 rounded-lg p-4 border border-gray-700">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-400 text-sm">平均延迟</p>
              <p class="text-2xl font-bold text-green-400 mt-1">
                <span v-if="isLoading">--</span>
                <span v-else>{{ stats.avgLatency }}ms</span>
              </p>
            </div>
            <div class="text-3xl">⚡</div>
          </div>
        </div>
        
        <div class="bg-gray-800 rounded-lg p-4 border border-gray-700">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-gray-400 text-sm">平均速度</p>
              <p class="text-2xl font-bold text-blue-400 mt-1">
                <span v-if="isLoading">--</span>
                <span v-else>{{ stats.avgSpeed }}MB/s</span>
              </p>
            </div>
            <div class="text-3xl">🚀</div>
          </div>
        </div>
      </div>

      <!-- 节点排名 -->
      <div class="bg-gray-800 rounded-lg border border-gray-700">
        <div class="px-6 py-4 border-b border-gray-700">
          <h3 class="text-lg font-semibold text-white">节点性能排名</h3>
        </div>
        <div class="p-6">
          <div class="flex justify-between items-center mb-4">
            <div class="flex space-x-4">
              <button
                v-for="metric in metrics"
                :key="metric.id"
                @click="selectedMetric = metric.id"
                :class="[
                  'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
                  selectedMetric === metric.id 
                    ? 'bg-primary-600 text-white' 
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                ]"
              >
                {{ metric.name }}
              </button>
            </div>
            <select class="input py-2 px-4 text-sm">
              <option>最近7天</option>
              <option>最近30天</option>
              <option>全部时间</option>
            </select>
          </div>
          
          <!-- 排名列表 -->
          <div class="space-y-2">
            <div v-for="(node, index) in topNodes" :key="node.name" 
                 class="flex items-center justify-between p-3 bg-gray-700/50 rounded-lg">
              <div class="flex items-center space-x-4">
                <span class="text-2xl font-bold text-gray-500">#{{ index + 1 }}</span>
                <div>
                  <p class="text-white font-medium">{{ node.name }}</p>
                  <p class="text-sm text-gray-400">{{ node.type }} · {{ node.region }}</p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-lg font-semibold text-green-400">{{ node.value }}</p>
                <p class="text-sm text-gray-400">{{ node.tests }} 次测试</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { GetAnalysisStats } from '../../wailsjs/go/main/App'

export default {
  name: 'AnalysisView',
  props: {
    configInfo: Object
  },
  setup(props) {
    const selectedMetric = ref('latency')
    
    const stats = ref({
      totalTests: 0,
      totalNodes: 0,
      avgLatency: 0,
      avgSpeed: 0
    })
    
    const isLoading = ref(true)
    const error = ref(null)
    
    // 获取统计数据
    const fetchStats = async () => {
      try {
        isLoading.value = true
        error.value = null
        const result = await GetAnalysisStats()
        if (result) {
          stats.value = {
            totalTests: result.totalTests || 0,
            totalNodes: result.totalNodes || 0,
            avgLatency: Math.round(result.avgLatency) || 0,
            avgSpeed: (result.avgSpeed || 0).toFixed(1)
          }
        }
      } catch (err) {
        console.error('获取统计数据失败:', err)
        error.value = '获取统计数据失败'
      } finally {
        isLoading.value = false
      }
    }
    
    // 组件挂载时获取数据
    onMounted(() => {
      fetchStats()
    })
    
    const metrics = ref([
      { id: 'latency', name: '最低延迟' },
      { id: 'speed', name: '最高速度' },
      { id: 'stability', name: '稳定性' },
      { id: 'overall', name: '综合评分' }
    ])
    
    const topNodes = computed(() => {
      // 模拟数据
      return [
        { name: 'HK-Premium-01', type: 'Shadowsocks', region: '香港', value: '23ms', tests: 145 },
        { name: 'SG-Fast-02', type: 'VMess', region: '新加坡', value: '28ms', tests: 132 },
        { name: 'JP-Tokyo-03', type: 'Trojan', region: '日本', value: '35ms', tests: 128 },
        { name: 'TW-Taipei-01', type: 'Shadowsocks', region: '台湾', value: '38ms', tests: 115 },
        { name: 'US-LA-05', type: 'VMess', region: '美国', value: '125ms', tests: 98 }
      ]
    })
    
    return {
      selectedMetric,
      stats,
      isLoading,
      error,
      fetchStats,
      metrics,
      topNodes
    }
  }
}
</script>

<style scoped>
/* 保持一致的样式 */
</style>