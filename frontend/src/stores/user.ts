import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '../api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<any>(null)

  const login = async (username: string, password: string) => {
    const res = await authApi.login({ username, password })
    token.value = res.data.token
    localStorage.setItem('token', res.data.token)
    return res
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  const isLoggedIn = () => !!token.value

  return {
    token,
    userInfo,
    login,
    logout,
    isLoggedIn
  }
})
