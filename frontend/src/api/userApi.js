import api from './axiosInstance'

export const userApi = {
  getAll: () => api.get('/users'),
  getById: (id) => api.get(`/users/userbyid?id=${id}`),
  create: (data) => api.post('/users/tambah', data),
  update: (id, data) => api.put(`/users/update?id=${id}`, data),
  delete: (id) => api.delete(`/users/deleted?id=${id}`),
}
