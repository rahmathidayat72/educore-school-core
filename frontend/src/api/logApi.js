import api from './axiosInstance'

export const logApi = {
  getAll: (params = {}) => api.get('/logs', { params }),
  getById: (id) => api.get(`/logs/detail?id=${id}`),
}
