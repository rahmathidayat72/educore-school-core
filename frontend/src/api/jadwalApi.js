import api from './axiosInstance'

export const getJadwal = (params = {}) => {
  return api.get('/jadwal', { params })
}

export const getJadwalById = (id) => {
  return api.get('/jadwal/detail', { params: { id } })
}

export const createJadwal = (data) => {
  return api.post('/jadwal/tambah', data)
}

export const updateJadwal = (id, data) => {
  return api.put('/jadwal/update', data, { params: { id } })
}

export const deleteJadwal = (id) => {
  return api.delete('/jadwal/deleted', { params: { id } })
}

export default {
  getJadwal,
  getJadwalById,
  createJadwal,
  updateJadwal,
  deleteJadwal
}
