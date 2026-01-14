<template>
  <div class="questionnaire">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>抑郁症筛查问卷</span>
          <div class="progress-info">
            已完成 {{ currentQuestionIndex + 1 }}/{{ questions.length }}
          </div>
        </div>
      </template>
      
      <el-progress
        :percentage="progressPercentage"
        :format="(percentage) => `${percentage}%`"
        class="progress-bar"
      />

      <div v-if="currentQuestion" class="question-container">
        <h3 class="question-title">{{ currentQuestion.title }}</h3>
        <p class="question-description" v-if="currentQuestion.description">{{ currentQuestion.description }}</p>
        <el-radio-group v-model="currentAnswer" class="answer-group">
          <el-radio
            v-for="option in currentQuestion.options"
            :key="option.value"
            :label="option.value"
            class="answer-option"
          >
            {{ option.label }}
          </el-radio>
        </el-radio-group>
      </div>

      <div class="button-group">
        <el-button
          type="primary"
          @click="handlePrevious"
          :disabled="currentQuestionIndex === 0"
        >
          上一题
        </el-button>
        <el-button
          type="primary"
          @click="handleNext"
          :disabled="!currentAnswer"
        >
          {{ isLastQuestion ? '提交' : '下一题' }}
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAssessmentStore } from '@/store/assessment'
import { ElMessage } from 'element-plus'
import { getQuestions, submitAnswers } from '../api'

const router = useRouter()
const questions = ref([])
const answers = ref({}) // 使用对象存储答案，键是问题索引
const currentQuestionIndex = ref(0)
const loading = ref(false)
const questionnaireId = ref(null)
const assessmentStore = useAssessmentStore()

// 当前问题
const currentQuestion = computed(() => questions.value[currentQuestionIndex.value])

// 当前答案的双向绑定
const currentAnswer = computed({
  get: () => answers.value[currentQuestionIndex.value],
  set: (value) => {
    answers.value[currentQuestionIndex.value] = value
  }
})

// 计算进度百分比
const progressPercentage = computed(() => {
  const total = questions.value.length
  const answered = Object.keys(answers.value).length
  return Math.round((answered / total) * 100)
})

// 是否是最后一题
const isLastQuestion = computed(() => {
  return currentQuestionIndex.value === questions.value.length - 1
})

// 上一题
const handlePrevious = () => {
  if (currentQuestionIndex.value > 0) {
    currentQuestionIndex.value--
  }
}

// 下一题或提交
const handleNext = async () => {
  if (!currentAnswer.value) {
    ElMessage.warning('请选择一个选项')
    return
  }

  if (isLastQuestion.value) {
    await handleSubmit()
  } else {
    currentQuestionIndex.value++
  }
}

// 提交问卷
const handleSubmit = async () => {
  try {
    loading.value = true
    
    const formattedAnswers = questions.value.map((q, index) => ({
      question_id: q.id,
      answer_value: answers.value[index]
    }))

    const response = await submitAnswers({
      assessment_id: questionnaireId.value || 1,
      answers: formattedAnswers
    })
    
    // 保存结果到Pinia
    assessmentStore.setResult(response.data)
    
    // 跳转到结果页（不再需要传参）
    router.push('/home/result')
    
  } catch (error) {
    ElMessage.error(`提交失败: ${error.message}`)
  } finally {
    loading.value = false
  }
}

// 初始化加载问题
onMounted(async () => {
  try {
    loading.value = true
    const response = await getQuestions()

    // 如果有问卷ID，保存起来
    if (response.questionnaire_id) {
      questionnaireId.value = response.questionnaire_id
    }

    // 预处理问题选项
    questions.value = response.data.map((q) => {
      const parsedOptions = JSON.parse(q.options)
      return {
        ...q,
        options: parsedOptions.map((label, index) => ({
          label,
          value: index + 1 // 选项值从1开始
        }))
      }
    })
  } catch (error) {
    console.error('获取问题失败:', error)
    ElMessage.error('获取问题失败，请刷新页面重试')
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.questionnaire {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.progress-info {
  font-size: 14px;
  color: #909399;
}

.progress-bar {
  margin: 20px 0;
}

.question-container {
  margin: 30px 0;
}

.question-title {
  font-size: 18px;
  color: #303133;
  margin-bottom: 10px;
}

.question-description {
  font-size: 14px;
  color: #606266;
  margin-bottom: 20px;
}

.answer-group {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.answer-option {
  padding: 12px 15px;
  border-radius: 4px;
  transition: all 0.3s;
  text-align: center !important;
}
.answer-option:nth-child(4) {
  padding: 12px 15px;
  border-radius: 4px;
  transition: all 0.3s;
  text-align: center !important;
  padding-right: 45px;
}

.answer-option.is-checked {
  background-color: #ecf5ff;
  border-color: #409eff;
}

.button-group {
  display: flex;
  justify-content: space-between;
  margin-top: 30px;
}

:deep(.el-radio__label) {
  font-size: 16px;
}

:deep(.el-progress-bar__outer) {
  background-color: #ebeef5;
}
</style>