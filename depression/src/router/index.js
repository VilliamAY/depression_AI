import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('../views/Home.vue'),
    redirect: '/home/face-detection',
    children: [
      {
        path: 'face-detection',
        name: 'FaceDetection',
        component: () => import('../views/FaceDetection.vue')
      },
      {
        path: 'questionnaire',
        name: 'Questionnaire',
        component: () => import('../views/Questionnaire.vue')
      },
      {
        path: 'result',
        name: 'Result',
        component: () => import('../views/Result.vue')
      },
      {
        path: 'combined-result',
        name: 'CombinedResult',
        component: () => import('../views/CombinedResult.vue') //
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && to.path !== '/register' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
