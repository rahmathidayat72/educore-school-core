import api from './axiosInstance'

export const authApi = {
  login: (data) => api.post('/login', data),
}
