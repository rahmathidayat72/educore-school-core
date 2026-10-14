import api from './axiosInstance'

export const mapelApi = {
  getAll: () => api.get('/mapel'),
  getById: (id) => api.get(`/mapel/mapelbyid?id=${id}`),
  create: (data) => api.post('/mapel/tambah', data),
  update: (id, data) => api.put(`/mapel/update?id=${id}`, data),
  delete: (id) => api.delete(`/mapel/deleted?id=${id}`),
}
