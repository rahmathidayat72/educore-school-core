import api from './axiosInstance'

export const siswaApi = {
  getAll: () => api.get('/siswa'),
  getById: (id) => api.get(`/siswa/siswabyid?id=${id}`),
  create: (data) => api.post('/siswa/tambah', data),
  update: (id, data) => api.put(`/siswa/update?id=${id}`, data),
  delete: (id) => api.delete(`/siswa/deleted?id=${id}`),
}
