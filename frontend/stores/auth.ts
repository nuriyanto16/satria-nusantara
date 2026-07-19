import { defineStore } from 'pinia'

interface User {
  id: string
  email: string
  nama_lengkap: string
  no_hp: string
  role_id: number
  role_name: string
  scope: string
  status: string
  cabang_id?: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: typeof window !== 'undefined' ? localStorage.getItem('token') : null,
    user: typeof window !== 'undefined' ? JSON.parse(localStorage.getItem('user') || 'null') : null as User | null,
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdminPusat: (state) => state.user?.scope === 'pusat',
    isPengurusCabang: (state) => state.user?.scope === 'cabang',
  },
  actions: {
    rehydrate() {
      if (typeof window !== 'undefined') {
        const storedToken = localStorage.getItem('token')
        const storedUser = localStorage.getItem('user')
        if (storedToken) {
          this.token = storedToken
        }
        if (storedUser) {
          try {
            this.user = JSON.parse(storedUser)
          } catch (e) {
            console.error(e)
          }
        }
      }
    },
    setAuth(token: string, user: User) {
      this.token = token
      this.user = user
      if (typeof window !== 'undefined') {
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(user))
      }
    },
    logout() {
      this.token = null
      this.user = null
      if (typeof window !== 'undefined') {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
      }
      navigateTo('/login')
    }
  }
})
