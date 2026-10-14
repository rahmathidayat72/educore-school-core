import api from './axiosInstance'

export const schoolUnitsApi = {
  getAll: (search = '') => api.get(`/school-units?search=${encodeURIComponent(search)}`),
  getById: (id) => api.get(`/school-units/detail?id=${id}`),
  create: (data) => api.post('/school-units', data),
  update: (id, data) => api.put(`/school-units/update?id=${id}`, data),
  delete: (id) => api.delete(`/school-units/deleted?id=${id}`)
}
