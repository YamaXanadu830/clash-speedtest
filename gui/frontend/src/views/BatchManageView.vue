<template>
  <div class="flex flex-col h-full p-6">
    <div class="max-w-6xl mx-auto w-full">
      <!-- 页面标题 -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-white mb-2">📦 批量管理</h1>
        <p class="text-gray-400">批量导入配置文件，对比测试结果，筛选最优节点</p>
      </div>

      <!-- 功能卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- 批量导入 -->
        <div class="card hover:shadow-xl transition-shadow">
          <div class="card-header">
            <h3 class="text-lg font-semibold text-white flex items-center space-x-2">
              <span>📥</span>
              <span>批量导入配置</span>
            </h3>
          </div>
          <div class="card-body">
            <p class="text-gray-400 mb-4">同时导入多个Clash配置文件进行批量测试</p>
            <button class="btn btn-primary w-full">
              选择配置文件
            </button>
          </div>
        </div>

        <!-- 批量测试 -->
        <div class="card hover:shadow-xl transition-shadow">
          <div class="card-header">
            <h3 class="text-lg font-semibold text-white flex items-center space-x-2">
              <span>🚀</span>
              <span>批量测试任务</span>
            </h3>
          </div>
          <div class="card-body">
            <p class="text-gray-400 mb-4">创建测试任务队列，自动执行批量测试</p>
            <button class="btn btn-primary w-full" :disabled="!configInfo">
              创建测试任务
            </button>
          </div>
        </div>

        <!-- 结果对比 -->
        <div class="card hover:shadow-xl transition-shadow">
          <div class="card-header">
            <h3 class="text-lg font-semibold text-white flex items-center space-x-2">
              <span>📊</span>
              <span>结果对比分析</span>
            </h3>
          </div>
          <div class="card-body">
            <p class="text-gray-400 mb-4">对比不同配置的测试结果，找出最优节点</p>
            <button class="btn btn-primary w-full">
              查看对比结果
            </button>
          </div>
        </div>

        <!-- 节点筛选 -->
        <div class="card hover:shadow-xl transition-shadow">
          <div class="card-header">
            <h3 class="text-lg font-semibold text-white flex items-center space-x-2">
              <span>🎯</span>
              <span>最优节点筛选</span>
            </h3>
          </div>
          <div class="card-body">
            <p class="text-gray-400 mb-4">根据设定条件自动筛选最优质的节点</p>
            <button class="btn btn-primary w-full">
              开始筛选
            </button>
          </div>
        </div>

        <!-- 配置合并 -->
        <div class="card hover:shadow-xl transition-shadow">
          <div class="card-header">
            <h3 class="text-lg font-semibold text-white flex items-center space-x-2">
              <span>🔗</span>
              <span>配置文件合并</span>
            </h3>
          </div>
          <div class="card-body">
            <p class="text-gray-400 mb-4">将多个配置文件的优质节点合并为一个</p>
            <button class="btn btn-primary w-full">
              合并配置
            </button>
          </div>
        </div>

        <!-- 定时任务 -->
        <div class="card hover:shadow-xl transition-shadow">
          <div class="card-header">
            <h3 class="text-lg font-semibold text-white flex items-center space-x-2">
              <span>⏰</span>
              <span>定时测试任务</span>
            </h3>
          </div>
          <div class="card-body">
            <p class="text-gray-400 mb-4">设置定时任务，自动执行节点测试</p>
            <button class="btn btn-primary w-full">
              配置定时任务
            </button>
          </div>
        </div>
      </div>

      <!-- 任务队列 -->
      <div class="mt-8">
        <h2 class="text-xl font-semibold text-white mb-4">任务队列</h2>
        <div class="bg-gray-800 rounded-lg p-6">
          <div v-if="taskQueue.length === 0" class="text-center py-12">
            <div class="text-gray-400 text-lg mb-4">📋</div>
            <p class="text-gray-400">暂无任务</p>
          </div>
          <div v-else class="space-y-4">
            <!-- 任务列表 -->
            <div v-for="task in taskQueue" :key="task.id" class="bg-gray-700 rounded-lg p-4">
              <div class="flex items-center justify-between">
                <div>
                  <h4 class="text-white font-medium">{{ task.name }}</h4>
                  <p class="text-sm text-gray-400">{{ task.description }}</p>
                </div>
                <div class="flex items-center space-x-2">
                  <span :class="getTaskStatusClass(task.status)">
                    {{ task.status }}
                  </span>
                  <button class="btn btn-sm btn-secondary">
                    取消
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  name: 'BatchManageView',
  props: {
    configInfo: Object
  },
  setup(props) {
    const taskQueue = ref([])
    
    const getTaskStatusClass = (status) => {
      const classes = {
        '待执行': 'text-gray-400',
        '执行中': 'text-yellow-400',
        '已完成': 'text-green-400',
        '失败': 'text-red-400'
      }
      return classes[status] || 'text-gray-400'
    }
    
    return {
      taskQueue,
      getTaskStatusClass
    }
  }
}
</script>

<style scoped>
.card {
  @apply bg-gray-800 rounded-lg border border-gray-700;
}

.card:hover {
  @apply border-gray-600;
}

.card-header {
  @apply px-6 py-4 border-b border-gray-700;
}

.card-body {
  @apply px-6 py-4;
}

.btn-sm {
  @apply px-3 py-1 text-sm;
}
</style>