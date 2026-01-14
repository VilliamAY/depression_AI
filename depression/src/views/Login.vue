<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <h2>登录</h2>
      </template>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef">
        <el-form-item prop="username">
          <el-input v-model="loginForm.username" placeholder="用户名">
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="密码">
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" :loading="loading" class="login-button">登录</el-button>
        </el-form-item>
        <div class="register-link">
          还没有账号？<router-link to="/register">立即注册</router-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { login } from '../api'

const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
    loading.value = true

    const res = await login({ 
      username: loginForm.username, 
      password: loginForm.password 
    })

    if (res.code === 200) {
      // 保存token和用户信息
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('username', res.data.user.username)  // 从user对象获取username
      localStorage.setItem('userInfo', JSON.stringify(res.data.user)) // 存储完整用户信息

      ElMessage.success('登录成功')
      router.push('/home/face-detection')
    } else {
      ElMessage.error(res.message || '登录失败')
    }
  } catch (error) {
    ElMessage.error(error.message || '登录异常')
  } finally {
    loading.value = false
  }
}

</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
}

.login-card {
  width: 400px;
}

.login-button {
  width: 100%;
}

.register-link {
  text-align: center;
  margin-top: 16px;
  font-size: 14px;
}

.register-link a {
  color: #409eff;
  text-decoration: none;
}

:deep(.el-card__header) {
  text-align: center;
  padding: 15px;
}

h2 {
  margin: 0;
  font-weight: 500;
  color: #303133;
}
</style> 