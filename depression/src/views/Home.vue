<template>
  <div class="home-container">
    <el-container>
      <el-aside width="200px">
        <el-menu
          :router="true"
          :default-active="activeMenu"
          class="menu-vertical"
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409EFF"
        >
          <el-menu-item index="/home/face-detection">
            <el-icon><Camera /></el-icon>
            <span>人脸心情识别</span>
          </el-menu-item>
          <el-menu-item index="/home/questionnaire">
            <el-icon><Document /></el-icon>
            <span>抑郁症筛查</span>
          </el-menu-item>
          <el-menu-item index="/home/result">
            <el-icon><DataAnalysis /></el-icon>
            <span>筛查结果分析</span>
          </el-menu-item>
          <el-menu-item index="/home/combined-result">
            <el-icon><DataAnalysis /></el-icon>
            <span>综合评估分析</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-container>
        <el-header>
          <div class="header-content">
            <h2>抑郁症筛查系统</h2>
            <el-dropdown @command="handleCommand">
              <span class="user-dropdown">
                <el-icon><User /></el-icon>
                <span>{{ username }}</span>
                <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        <el-main>
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Camera,
  Document,
  DataAnalysis,
  User,
  ArrowDown
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const username = ref('')

onMounted(() => {
  username.value = localStorage.getItem('username') || '用户'
})

const activeMenu = computed(() => route.path)

const handleCommand = (command) => {
  if (command === 'logout') {
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    localStorage.removeItem('userInfo')
    router.push('/login')
    ElMessage.success('已退出登录')
  }
}
</script>

<style scoped>
.home-container {
  height: 100vh;
}

.el-container {
  height: 100%;
}

.el-aside {
  background-color: #304156;
}

.menu-vertical {
  border-right: none;
  height: 100%;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
}

.header-content {
  height: 60px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #606266;
}

.user-dropdown .el-icon {
  margin: 0 5px;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}

h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}
</style>
