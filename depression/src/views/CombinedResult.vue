<template>
  <div class="result">
    <el-card v-if="result" class="result-card">
      <template #header>
        <div class="card-header">
          <span>综合评估结果</span>
          <el-button type="primary" @click="handleRetake">重新评估</el-button>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="8">
          <div class="score-section">
            <div class="score-title">综合得分</div>
            <el-progress
              type="dashboard"
              :percentage="scorePercentage"
              :color="getScoreColor(scorePercentage)"
            >
              <template #default="{ percentage }">
                <span class="score-value">{{ result.combined_score }}</span>
                <span class="score-label">分</span>
              </template>
            </el-progress>
            <div class="risk-level" :style="{ color: getRiskColor(result.combined_level) }">
              风险等级：{{ riskLevelText }}
            </div>
          </div>
        </el-col>

        <el-col :span="8">
          <div class="score-box">
            <div class="score-title">问卷得分</div>
            <el-tag type="info">分数：{{ result.questionnaire.score }}</el-tag>
            <div style="margin-top: 8px;">等级：{{ result.questionnaire.level }}</div>
            <div style="margin-top: 8px; color: #909399;">时间：{{ formatTime(result.assessment_date) }}</div>
          </div>
        </el-col>

        <el-col :span="8">
          <div class="score-box">
            <div class="score-title">人脸检测</div>
            <el-tag type="info">分数：{{ result.face_detection.score }}</el-tag>
            <div style="margin-top: 8px;">等级：{{ result.face_detection.level }}</div>
            <div style="margin-top: 8px;">情绪：{{ result.face_detection.emotion }}</div>
            <div style="margin-top: 8px; color: #909399;">时间：{{ formatTime(result.detection_date) }}</div>
          </div>
        </el-col>
      </el-row>

      <div class="description-section">
        <h3>评估描述</h3>
        <div class="description-content">{{ result.description }}</div>
      </div>

      <div class="suggestion-section">
        <h3>建议方案</h3>
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
      <el-empty description="暂无综合评估结果" v-if="!loading">
        <el-button type="primary" @click="router.push('/home/questionnaire')">
          去评估
        </el-button>
      </el-empty>
      <el-skeleton :rows="5" animated v-else />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Warning as ElIconWarning } from '@element-plus/icons-vue'
import { getAssessmentTotal } from '../api'

const router = useRouter()
const result = ref(null)
const loading = ref(true)

const scorePercentage = computed(() => {
  return Math.min(Math.round((result.value?.combined_score || 0)), 100)
})

const riskLevelText = computed(() => {
  const levels = {
    severe: '高风险',
    moderate: '中等风险',
    mild: '轻度风险',
    normal: '正常'
  }
  return levels[result.value?.combined_level] || '未知'
})

const formattedSuggestions = computed(() => {
  return result.value?.suggestions?.split('\n') || []
})

const getScoreColor = (percentage) => {
  if (percentage < 40) return '#67C23A'
  if (percentage < 60) return '#409EFF'
  if (percentage < 80) return '#E6A23C'
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

const formatTime = (time) => {
  return new Date(time).toLocaleString()
}

const handleRetake = () => {
  router.push('/home/questionnaire')
}
console.log('44444');

onMounted(async () => {
  try {
    const res = await getAssessmentTotal()
    console.log(res);
    
    if (res.code === 200) {
      result.value = res.data
    } else {
      ElMessage.warning(res.message || '加载结果失败')
    }
  } 
  catch (error) {
    ElMessage.error('网络错误，请稍后再试')
  } finally {
    loading.value = false
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
  font-weight: 600;
  color: #303133;
  margin-bottom: 10px;
}

.score-value {
  font-size: 28px;
  font-weight: bold;
}

.score-label {
  font-size: 14px;
  margin-left: 4px;
}

.risk-level {
  margin-top: 15px;
  font-size: 16px;
  font-weight: 500;
}

.score-box {
  padding: 20px;
  background: #f9f9f9;
  border-radius: 4px;
}

.description-section {
  margin-top: 30px;
  padding: 0 20px;
}

.description-content {
  margin-top: 10px;
  line-height: 1.6;
  color: #606266;
}

.suggestion-section {
  margin-top: 30px;
  padding: 0 20px;
}

.suggestion-item {
  display: flex;
  align-items: flex-start;
  margin-bottom: 15px;
  padding: 10px;
  background-color: #fef6e6;
  border-left: 4px solid #f56c6c;
  border-radius: 4px;
}

.suggestion-item i {
  margin-right: 10px;
  color: #f56c6c;
}

.loading-container {
  padding: 40px;
  text-align: center;
}
</style>
