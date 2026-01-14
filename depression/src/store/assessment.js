import { defineStore } from 'pinia'

export const useAssessmentStore = defineStore('assessment', {
  state: () => ({
    result: null
  }),
  actions: {
    setResult(data) {
      this.result = data
      // 可选：同时存入本地存储防止页面刷新丢失
      localStorage.setItem('assessmentResult', JSON.stringify(data))
    },
    loadFromStorage() {
      const saved = localStorage.getItem('assessmentResult')
      if (saved) this.result = JSON.parse(saved)
    }
  }
})