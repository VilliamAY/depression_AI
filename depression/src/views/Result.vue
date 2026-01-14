<template>
  <div class="result">
    <el-card v-if="result" class="result-card">
      <template #header>
        <div class="card-header">
          <span>评估结果</span>
          <el-button type="primary" @click="handleRetake">重新测试</el-button>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="12">
          <div class="score-section">
            <div class="score-title">总体评分</div>
            <el-progress
              type="dashboard"
              :percentage="scorePercentage"
              :color="getScoreColor(scorePercentage)"
            >
              <template #default="{ percentage }">
                <span class="score-value">{{ result.score }}</span>
                <span class="score-label">分</span>
              </template>
            </el-progress>
            <div class="risk-level" :style="{ color: getRiskColor(result.level) }">
              风险等级：{{ riskLevelText }}
            </div>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="description-section">
            <h3>评估描述</h3>
            <div class="description-content">
              {{ result.description }}
            </div>
          </div>
        </el-col>
      </el-row>

      <div class="suggestion-section">
        <h3>专业建议</h3>
        <div class="suggestion-content">
          <div 
            v-for="(item, index) in formattedSuggestions" 
            :key="index" 
            class="suggestion-item"
          >
            <el-icon><el-icon-warning /></el-icon>
            <span>{{ item }}</span>
          </div>
        </div>
      </div>
    </el-card>

    <div v-else class="loading-container">
      <el-empty description="暂无评估结果" v-if="!loading">
        <el-button type="primary" @click="router.push('/home/questionnaire')">
          开始评估
        </el-button>
      </el-empty>
      <el-skeleton :rows="5" animated v-else />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElIcon } from 'element-plus'
import { Warning as ElIconWarning } from '@element-plus/icons-vue'
import { useAssessmentStore } from '@/store/assessment'

const assessmentStore = useAssessmentStore()
const router = useRouter()
const route = useRoute()
const result = ref(null)
const loading = ref(false)

// 计算属性
const scorePercentage = computed(() => {
  // 满分是100分
  return Math.min(Math.round((result.value?.score || 0) / 100 * 100), 100)
})

const riskLevelText = computed(() => {
  const levels = {
    severe: '高风险',
    moderate: '中等风险',
    mild: '轻度风险',
    normal: '正常'
  }
  return levels[result.value?.level] || result.value?.level
})

const formattedSuggestions = computed(() => {
  return result.value?.suggestions?.split('\n') || []
})

// 颜色计算
const getScoreColor = (percentage) => {
  if (percentage < 30) return '#67C23A'
  if (percentage < 60) return '#E6A23C'
  return '#F56C6C'
}

const getRiskColor = (level) => {
  const colors = {
    severe: '#F56C6C',
    moderate: '#E6A23C',
    mild: '#409EFF',
    normal: '#67C23A'
  }
  return colors[level] || '#909399'
}

// 重新测试
const handleRetake = () => {
  router.push('/home/questionnaire')
}

// 初始化加载结果
onMounted(() => {
  // 从Pinia加载结果
  result.value = assessmentStore.result
  
  // 如果直接访问页面，尝试从本地存储加载
  if (!result.value) {
    assessmentStore.loadFromStorage()
    result.value = assessmentStore.result
  }
})
</script>

<style scoped>
.result {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.score-section {
  text-align: center;
  padding: 20px;
}

.score-title {
  font-size: 18px;
  color: #303133;
  margin-bottom: 20px;
}

.score-value {
  font-size: 28px;
  font-weight: bold;
}

.score-label {
  font-size: 14px;
  margin-left: 5px;
}

.risk-level {
  margin-top: 20px;
  font-size: 16px;
  font-weight: 500;
}

.description-section {
  padding: 20px;
}

.description-content {
  margin-top: 15px;
  line-height: 1.6;
  color: #606266;
}

.suggestion-section {
  margin-top: 30px;
  padding: 0 20px;
}

.suggestion-content {
  margin-top: 15px;
}

.suggestion-item {
  display: flex;
  align-items: flex-start;
  margin-bottom: 15px;
  padding: 15px;
  background-color: #f8f8f8;
  border-radius: 4px;
}

.suggestion-item i {
  margin-right: 10px;
  color: #e6a23c;
}

.loading-container {
  padding: 40px;
  text-align: center;
}

h3 {
  margin: 0 0 20px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}
</style>