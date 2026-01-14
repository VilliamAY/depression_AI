<template>
  <div class="face-detection">
    <el-card class="method-card">
      <template #header>
        <div class="card-header">
          <span>选择识别方式</span>
          <el-button type="primary" @click="fetchHistory" :loading="loadingHistory">
            <el-icon><Clock /></el-icon> 查看历史记录
          </el-button>
        </div>
      </template>
      <div class="method-options">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="摄像头拍照" name="camera">
            <div class="camera-container" v-if="activeTab === 'camera'">
              <div class="camera-controls">
                <el-button type="primary" @click="startCamera" v-if="!cameraActive">
                  开启摄像头
                </el-button>
                <div v-else>
                  <el-button type="success" @click="captureAndUpload">
                    <el-icon><Camera /></el-icon> 拍照并分析
                  </el-button>
                  <el-button type="danger" @click="stopCamera">
                    关闭摄像头
                  </el-button>
                </div>
              </div>
              <div class="video-container">
                <video ref="videoRef" class="camera-view" autoplay playsinline style="transform: scaleX(-1);"></video>
                <canvas ref="canvasRef" class="canvas-overlay" style="transform: scaleX(-1);"></canvas>
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane label="本地图片上传" name="upload">
            <div class="upload-container" v-if="activeTab === 'upload'">
              <el-upload
                class="image-uploader"
                action="#"
                :auto-upload="false"
                :show-file-list="true"
                :on-change="handleFileChange"
                :limit="1"
                accept="image/jpeg,image/png,image/jpg"
              >
                <el-icon class="upload-icon"><Plus /></el-icon>
                <div class="upload-text">点击上传图片</div>
              </el-upload>
              <div class="upload-preview" v-if="selectedFile">
                <img :src="selectedFileUrl" class="preview-image" />
                <el-button type="primary" @click="uploadLocalImage" :loading="uploading">
                  上传并分析
                </el-button>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-card>

    <el-dialog v-model="historyDialogVisible" title="人脸检测历史记录" width="70%">
      <el-table :data="historyData" border style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="emotion" label="主要情绪" width="120">
          <template #default="{ row }">
            <el-tag :type="getEmotionTagType(row.emotion)">
              {{ row.emotion }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="情绪置信度" width="180">
          <template #default="{ row }">
            <el-progress 
              :percentage="Math.round(row.confidence * 100)" 
              :color="getEmotionColor(row.emotion)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="score" label="情绪分数" width="120" />
        <el-table-column prop="level" label="情绪等级" width="120" />
        <el-table-column prop="created_at" label="检测时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="分析结果" prop="result" />
      </el-table>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" :title="`检测详情 - ID: ${currentDetail?.id}`" width="50%">
      <div v-if="currentDetail" class="detail-container">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="检测图片">
            <el-image 
              style="width: 300px; height: 200px"
              :src="currentDetail.image_url" 
              fit="cover"
            />
          </el-descriptions-item>
          <el-descriptions-item label="主要情绪">
            <el-tag :type="getEmotionTagType(currentDetail.emotion)">
              {{ currentDetail.emotion }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="情绪置信度">
            <el-progress 
              :percentage="Math.round(currentDetail.confidence * 100)" 
              :color="getEmotionColor(currentDetail.emotion)"
            />
          </el-descriptions-item>
          <el-descriptions-item label="情绪分数">
            {{ currentDetail.score }}
          </el-descriptions-item>
          <el-descriptions-item label="情绪等级">
            {{ currentDetail.level }}
          </el-descriptions-item>
          <el-descriptions-item label="分析结果">
            {{ currentDetail.result }}
          </el-descriptions-item>
          <el-descriptions-item label="检测时间">
            {{ formatDateTime(currentDetail.created_at) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <el-card class="result-card" v-if="detectionResult">
      <template #header>
        <div class="card-header">
          <span>识别结果</span>
        </div>
      </template>
      <div class="result-container">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="主要情绪">
            {{ detectionResult.emotion }}
          </el-descriptions-item>
          <el-descriptions-item label="情绪置信度">
            <el-progress 
              :percentage="Math.round(detectionResult.confidence * 100)" 
              :color="getEmotionColor(detectionResult.emotion)"
            />
          </el-descriptions-item>
          <el-descriptions-item label="情绪分数">
            {{ detectionResult.score }}
          </el-descriptions-item>
          <el-descriptions-item label="情绪等级">
            {{ detectionResult.level }}
          </el-descriptions-item>
          <el-descriptions-item label="分析结果">
            {{ detectionResult.result }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { uploadFaceImage, getFaceHistory } from '../api'
import { Camera, Plus } from '@element-plus/icons-vue'

const historyDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const historyData = ref([])
const loadingHistory = ref(false)
const currentDetail = ref(null)
const videoRef = ref(null)
const canvasRef = ref(null)
const cameraActive = ref(false)
const stream = ref(null)
const detectionResult = ref(null)
const activeTab = ref('camera')
const selectedFile = ref(null)
const selectedFileUrl = ref('')
const uploading = ref(false)

// 情绪颜色映射
const emotionColors = {
  'happy': '#67C23A',
  'sad': '#409EFF',
  'angry': '#F56C6C',
  'fear': '#E6A23C',
  'surprise': '#909399',
  'neutral': '#67C23A',
  'disgust': '#E6A23C'
}

// 情绪标签类型映射
const emotionTagTypes = {
  'happy': 'success',
  'sad': 'info',
  'angry': 'danger',
  'fear': 'warning',
  'surprise': '',
  'neutral': 'success',
  'disgust': 'warning'
}

const getEmotionTagType = (emotion) => {
  return emotionTagTypes[emotion] || 'info'
}

const getEmotionColor = (emotion) => {
  return emotionColors[emotion] || '#409EFF'
}

// 格式化日期时间
const formatDateTime = (dateTimeStr) => {
  if (!dateTimeStr) return ''
  const date = new Date(dateTimeStr)
  return date.toLocaleString()
}

// 获取历史记录
const fetchHistory = async () => {
  loadingHistory.value = true
  try {
    const response = await getFaceHistory()
    console.log(response)
    historyData.value = response.data.list
    historyDialogVisible.value = true
  } catch (error) {
    console.error('获取历史记录失败:', error)
    ElMessage.error('获取历史记录失败: ' + (error.message || '未知错误'))
  } finally {
    loadingHistory.value = false
  }
}

// 显示详情
const showDetail = (row) => {
  currentDetail.value = row
  detailDialogVisible.value = true
}

// 开启摄像头
const startCamera = async () => {
  try {
    stream.value = await navigator.mediaDevices.getUserMedia({
      video: {
        width: 640,
        height: 480,
        facingMode: 'user'
      }
    })
    
    if (videoRef.value) {
      videoRef.value.srcObject = stream.value
      cameraActive.value = true
    }
  } catch (error) {
    ElMessage.error('无法访问摄像头')
    console.error('Error accessing camera:', error)
  }
}

// 关闭摄像头
const stopCamera = () => {
  if (stream.value) {
    stream.value.getTracks().forEach(track => track.stop())
    stream.value = null
  }
  cameraActive.value = false
}

// 捕获摄像头图像
const captureImage = () => {
  if (!canvasRef.value || !videoRef.value) return

  const context = canvasRef.value.getContext('2d')
  canvasRef.value.width = videoRef.value.videoWidth
  canvasRef.value.height = videoRef.value.videoHeight
  
  // 先水平翻转上下文
  context.translate(canvasRef.value.width, 0)
  context.scale(-1, 1)
  
  // 然后绘制图像
  context.drawImage(videoRef.value, 0, 0)
  
  // 重置变换
  context.setTransform(1, 0, 0, 1, 0, 0)
  
  return canvasRef.value.toDataURL('image/jpeg', 0.8)
}

// 手动拍照并上传分析
const captureAndUpload = async () => {
  if (!cameraActive.value) return
  
  const imageData = captureImage()
  if (!imageData) {
    ElMessage.warning('无法获取图像，请确保摄像头正常工作')
    return
  }
  
  uploading.value = true
  try {
    const base64Data = imageData.split(',')[1]
    const blob = base64ToBlob(base64Data, 'image/jpeg')
    const formData = new FormData()
    formData.append('image', blob, 'camera_capture.jpg')
    
    const response = await uploadFaceImage(formData)
    detectionResult.value = {
      emotion: response.data.emotion,
      confidence: response.data.confidence,
      score: response.data.score,
      level: response.data.level,
      result: response.data.result
    }
    ElMessage.success('人脸分析完成')
  } catch (error) {
    console.error('Error during face detection:', error)
    ElMessage.error('人脸分析失败: ' + (error.message || '未知错误'))
  } finally {
    uploading.value = false
  }
}

// 处理本地文件选择
const handleFileChange = (file) => {
  if (!file) return
  
  const isImage = file.raw?.type?.startsWith('image/') || file.type?.startsWith('image/')
  if (!isImage) {
    ElMessage.error('请上传图片文件')
    return
  }
  
  const isLt5M = (file.raw?.size || file.size) / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过5MB')
    return
  }
  
  selectedFile.value = file.raw || file
  selectedFileUrl.value = URL.createObjectURL(file.raw || file)
}

// 上传本地图片
const uploadLocalImage = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择图片')
    return
  }

  uploading.value = true
  try {
    const formData = new FormData()
    const file = selectedFile.value.raw || selectedFile.value
    let fileName = file.name || 'upload.jpg'
    if (!fileName.toLowerCase().match(/\.(jpg|jpeg|png)$/)) {
      fileName += '.jpg'
    }

    formData.append('image', file, fileName)

    const response = await uploadFaceImage(formData)
    detectionResult.value = {
      emotion: response.data.emotion,
      confidence: response.data.confidence,
      score: response.data.score,
      level: response.data.level,
      result: response.data.result
    }
    ElMessage.success('人脸分析完成')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败: ' + (error.message || '未知错误'))
  } finally {
    uploading.value = false
  }
}

// Base64转Blob工具函数
const base64ToBlob = (base64, mimeType) => {
  const byteCharacters = atob(base64)
  const byteArrays = []
  
  for (let offset = 0; offset < byteCharacters.length; offset += 512) {
    const slice = byteCharacters.slice(offset, offset + 512)
    
    const byteNumbers = new Array(slice.length)
    for (let i = 0; i < slice.length; i++) {
      byteNumbers[i] = slice.charCodeAt(i)
    }
    
    const byteArray = new Uint8Array(byteNumbers)
    byteArrays.push(byteArray)
  }
  
  return new Blob(byteArrays, { type: mimeType })
}

onUnmounted(() => {
  stopCamera()
  if (selectedFileUrl.value) {
    URL.revokeObjectURL(selectedFileUrl.value)
  }
})
</script>

<style scoped>
.face-detection {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.method-card,
.result-card {
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 18px;
  font-weight: 600;
}

.camera-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
}

.camera-controls {
  margin-bottom: 20px;
  display: flex;
  justify-content: center;
  gap: 10px;
}

.video-container {
  position: relative;
  width: 100%;
  max-width: 640px;
  margin: 0 auto;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.camera-view {
  width: 100%;
  height: auto;
  display: block;
  background-color: #f5f7fa;
}

.canvas-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.upload-container {
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detail-container {
  padding: 10px;
}

.image-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  width: 300px;
  height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  transition: border-color 0.3s;
}

.image-uploader:hover {
  border-color: #409EFF;
}

.upload-icon {
  font-size: 28px;
  color: #8c939d;
}

.upload-text {
  color: #8c939d;
  margin-top: 10px;
  font-size: 14px;
}

.upload-preview {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
}

.preview-image {
  max-width: 300px;
  max-height: 300px;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.result-container {
  padding: 10px;
}

:deep(.el-progress-bar__outer) {
  background-color: #ebeef5;
}

:deep(.el-descriptions) {
  margin-top: 10px;
}

:deep(.el-descriptions__body) {
  background-color: #fff;
}

:deep(.el-descriptions__label) {
  width: 120px;
  font-weight: 600;
}

:deep(.el-tabs__nav) {
  margin-bottom: 20px;
}
</style>