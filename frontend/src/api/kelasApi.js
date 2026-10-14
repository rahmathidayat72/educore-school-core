import api from './axiosInstance'

export const kelasApi = {
  getAll: () => api.get('/kelas'),
  getById: (id) => api.get(`/kelas/kelasbyid?id=${id}`),
  create: (data) => api.post('/kelas/tambah', data),
  update: (id, data) => api.put(`/kelas/update?id=${id}`, data),
  delete: (id) => api.delete(`/kelas/deleted?id=${id}`),
}
