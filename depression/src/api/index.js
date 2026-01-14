import axios from 'axios'

const instance = axios.create({
  baseURL: '/api', // 使用相对路径，将通过Vue代理转发到后端
  timeout: 5000
})

// 请求拦截器
instance.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // 确保不覆盖已设置的Content-Type，特别是对FormData请求
    // if (!config.headers['Content-Type']) {
    //   // 只有在未明确设置Content-Type时才设置默认值
    //   config.headers['Content-Type'] = 'application/json'
    // }

    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
instance.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    return Promise.reject(error)
  }
)

// API接口
export const login = (data) => {
  return instance.post('/auth/login', data)
}

export const register = (data) => {
  return instance.post('/auth/register', data)
}

export const uploadFaceImage = (data) => {
  return instance.post('/face/upload', data)
}


export const getQuestions = () => {
  return instance.get('/questions')
}

export const submitAnswers = (data) => {
  return instance.post('/questionnaire/submit', data)
}

export const getResult = () => {
  return instance.get('/assessment/combined')
}

export const getFaceHistory = () => {
  return instance.get('/face/history')
}

export const getAssessmentTotal = () => {
  return instance.get('/assessment/total')
}


export default instance