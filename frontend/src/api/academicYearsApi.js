import api from './axiosInstance'

export const academicYearsApi = {
  // Master Academic Years
  getYears: (unitId = '') => {
    const q = unitId ? `?unit_id=${encodeURIComponent(unitId)}` : ''
    return api.get(`/academic-years${q}`)
  },
  getAll: (unitId = '') => {
    const q = unitId ? `?unit_id=${encodeURIComponent(unitId)}` : ''
    return api.get(`/academic-years${q}`)
  },
  getActiveYear: (unitId = '') => {
    const q = unitId ? `?unit_id=${encodeURIComponent(unitId)}` : ''
    return api.get(`/academic-years/active${q}`)
  },
  createYear: (data, unitId = '') => {
    const q = unitId ? `?unit_id=${encodeURIComponent(unitId)}` : ''
    return api.post(`/academic-years/tambah${q}`, data)
  },
  create: (data, unitId = '') => {
    const q = unitId ? `?unit_id=${encodeURIComponent(unitId)}` : ''
    return api.post(`/academic-years/tambah${q}`, data)
  },
  updateYear: (id, data) => api.put(`/academic-years/update?id=${id}`, data),
  update: (id, data) => api.put(`/academic-years/update?id=${id}`, data),
  activateYear: (id) => api.post(`/academic-years/set-active?id=${id}`),
  setActive: (id) => api.post(`/academic-years/set-active?id=${id}`),
  deleteYear: (id) => api.delete(`/academic-years/deleted?id=${id}`),
  delete: (id) => api.delete(`/academic-years/deleted?id=${id}`),

  // Class Enrollments
  getEnrollmentsByClass: (kelasId, academicYearId = '') => {
    let q = `?kelas_id=${encodeURIComponent(kelasId)}`
    if (academicYearId) q += `&academic_year_id=${encodeURIComponent(academicYearId)}`
    return api.get(`/class-enrollments/by-class${q}`)
  },
  enrollStudent: (data) => api.post('/class-enrollments/enroll', data),
  updateEnrollment: (id, data) => api.put(`/class-enrollments/update?id=${id}`, data),
  deleteEnrollment: (id) => api.delete(`/class-enrollments/deleted?id=${id}`),
  getStudentHistory: (siswaId) => api.get(`/class-enrollments/student-history?siswa_id=${siswaId}`)
}
