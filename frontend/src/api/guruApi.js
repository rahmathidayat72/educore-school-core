import api from './axiosInstance'

export const guruApi = {
  getAll: () => api.get('/guru'),
  getById: (id) => api.get(`/guru/gurubyid?id=${id}`),
  create: (data) => api.post('/guru/tambah', data),
  update: (id, data) => api.put(`/guru/update?id=${id}`, data),
  delete: (id) => api.delete(`/guru/deleted?id=${id}`),
}
