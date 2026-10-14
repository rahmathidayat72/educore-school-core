import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request Interceptor: Menambahkan Authorization Bearer token otomatis
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response Interceptor: Tangani error 401 (token expired / unauthorized)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      // Token tidak valid atau habis masa berlaku
      if (window.location.pathname.startsWith('/admin')) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        window.location.href = '/login?expired=1'
      }
    }
    return Promise.reject(error)
  }
)

export default api
