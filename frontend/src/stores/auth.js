import { defineStore } from 'pinia'
import { authApi } from '@/api/authApi'
import { userApi } from '@/api/userApi'

export const useAuthStore = defineStore('auth', {
  state: () => {
    let savedUser = null
    try {
      savedUser = JSON.parse(localStorage.getItem('user') || 'null')
      // Sinkronisasi pemulihan jika role masih kosong pada sesi lama
      if (savedUser && !savedUser.role) {
        if (savedUser.email === 'admin@sekolah.id' || savedUser.username === 'admin') {
          savedUser.role = 'admin'
        } else {
          savedUser.role = 'user'
        }
        localStorage.setItem('user', JSON.stringify(savedUser))
      }
    } catch {
      savedUser = null
    }

    return {
      token: localStorage.getItem('token') || '',
      user: savedUser,
      loading: false,
      error: null,
    }
  },

  getters: {
    isAuthenticated: (state) => !!state.token,
    currentUser: (state) => state.user,
    isSuperAdmin: (state) => {
      if (!state.user) return false
      const r = String(state.user.role || '').trim().toLowerCase()
      return r === 'super_admin'
    },
    isAdmin: (state) => {
      if (!state.user) return false
      // Jika field role ada, periksa nilai eksplisitnya
      if (state.user.role) {
        const r = String(state.user.role).trim().toLowerCase()
        return r === 'admin' || r === 'super_admin'
      }
      // Khusus penanganan akun admin default pada sesi lama
      return state.user.email === 'admin@sekolah.id' || state.user.username === 'admin'
    },
    isGuru: (state) => {
      if (!state.user || !state.user.role) return false
      return String(state.user.role).trim().toLowerCase() === 'guru'
    },
    canManageGuru: (state) => {
      if (!state.user) return false
      if (state.user.role) {
        const r = String(state.user.role).trim().toLowerCase()
        return r === 'admin' || r === 'super_admin'
      }
      return state.user.email === 'admin@sekolah.id' || state.user.username === 'admin'
    },
    canManageSiswa: (state) => {
      if (!state.user) return false
      if (state.user.role) {
        const r = String(state.user.role).trim().toLowerCase()
        return r === 'admin' || r === 'super_admin'
      }
      return state.user.email === 'admin@sekolah.id' || state.user.username === 'admin'
    },
    canManageUsers: (state) => {
      if (!state.user) return false
      if (state.user.role) {
        const r = String(state.user.role).trim().toLowerCase()
        return r === 'admin' || r === 'super_admin'
      }
      return state.user.email === 'admin@sekolah.id' || state.user.username === 'admin'
    },
    canManageKelas: (state) => {
      if (!state.user) return false
      const r = String(state.user.role || '').trim().toLowerCase()
      return r === 'admin' || r === 'super_admin' || r === 'guru' || state.user.email === 'admin@sekolah.id' || state.user.username === 'admin'
    },
    canManageMapel: (state) => {
      if (!state.user) return false
      const r = String(state.user.role || '').trim().toLowerCase()
      return r === 'admin' || r === 'super_admin' || r === 'guru' || state.user.email === 'admin@sekolah.id' || state.user.username === 'admin'
    },
    portalTitle: (state) => {
      if (!state.user) return 'Portal Sekolah'
      const r = String(state.user.role || '').trim().toLowerCase()
      if (r === 'admin') return 'Portal Administrator'
      if (r === 'guru') return 'Portal Guru'
      return 'Portal Akademik'
    },
  },

  actions: {
    async login(email, password) {
      this.loading = true
      this.error = null
      try {
        const response = await authApi.login({ email, password })
        const resData = response.data
        if (resData.success && resData.data) {
          const { token, id, username, email: userEmail, role } = resData.data
          let userRole = role
          if (!userRole) {
            if (userEmail === 'admin@sekolah.id' || username === 'admin') {
              userRole = 'admin'
            } else {
              userRole = 'user'
            }
          }
          this.token = token
          this.user = { id, username, email: userEmail, role: userRole }
          localStorage.setItem('token', token)
          localStorage.setItem('user', JSON.stringify(this.user))
          return { success: true }
        } else {
          throw new Error(resData.message || 'Login gagal')
        }
      } catch (err) {
        const message =
          err.response?.data?.message ||
          (typeof err.response?.data === 'string' ? err.response.data : null) ||
          err.message ||
          'Gagal melakukan login. Silakan cek email & password.'
        this.error = message
        return { success: false, message }
      } finally {
        this.loading = false
      }
    },

    async syncUserRole() {
      if (!this.user?.id || !this.token) return
      try {
        const response = await userApi.getById(this.user.id)
        const resData = response.data
        if (resData.success && resData.data) {
          const userData = Array.isArray(resData.data) ? resData.data[0] : resData.data
          if (userData && userData.role) {
            this.user.role = userData.role
            localStorage.setItem('user', JSON.stringify(this.user))
          }
        }
      } catch {
        // Tangani silent error jika gagal sync di background
      }
    },

    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },
  },
})
