import api from './axiosInstance'

export const guardiansApi = {
  getAll: (search = '') => api.get(`/guardians?search=${encodeURIComponent(search)}`),
  getById: (id) => api.get(`/guardians/detail?id=${id}`),
  getMyProfile: () => api.get('/guardians/my-profile'),
  create: (data) => api.post('/guardians/tambah', data),
  update: (id, data) => api.put(`/guardians/update?id=${id}`, data),
  delete: (id) => api.delete(`/guardians/deleted?id=${id}`),

  // Student-Guardian Link API
  linkStudent: (guardianId, siswaId, isPrimary = false) =>
    api.post('/guardians/link-student', {
      guardian_id: guardianId,
      siswa_id: siswaId,
      is_primary: isPrimary
    }),
  unlinkStudent: (guardianId, siswaId) =>
    api.delete(`/guardians/unlink-student?guardian_id=${guardianId}&siswa_id=${siswaId}`),
  getStudentsByGuardianId: (guardianId) =>
    api.get(`/guardians/students?guardian_id=${guardianId}`),
  getGuardiansBySiswaId: (siswaId) =>
    api.get(`/guardians/by-student?siswa_id=${siswaId}`)
}
