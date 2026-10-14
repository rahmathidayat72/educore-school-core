<template>
  <div class="container-fluid py-4">
    <!-- Header Page -->
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
      <div>
        <h4 class="fw-bold text-dark mb-1">
          <i class="fas fa-calendar-alt text-primary me-2"></i>Tahun Ajaran & Semester
        </h4>
        <p class="text-secondary small mb-0">
          Atur periode kalender akademik, semester berjalan, dan aktivasi tahun ajaran KBM.
        </p>
      </div>

      <div class="d-flex align-items-center gap-2">
        <button @click="fetchYears" class="btn btn-outline-secondary btn-sm" :disabled="loading">
          <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
          <span class="d-none d-sm-inline ms-1">Segarkan</span>
        </button>
        <button @click="openCreateModal" class="btn btn-primary btn-sm d-flex align-items-center gap-1 shadow-sm">
          <i class="fas fa-plus"></i>
          <span>Tambah Periode Baru</span>
        </button>
      </div>
    </div>

    <!-- Alert Messages -->
    <div v-if="successMessage" class="alert alert-success alert-dismissible fade show small py-2 d-flex align-items-center gap-2 mb-3">
      <i class="fas fa-check-circle fs-5"></i>
      <div>{{ successMessage }}</div>
      <button type="button" class="btn-close" @click="successMessage = ''"></button>
    </div>

    <div v-if="errorMessage" class="alert alert-danger alert-dismissible fade show small py-2 d-flex align-items-center gap-2 mb-3">
      <i class="fas fa-exclamation-triangle fs-5"></i>
      <div>{{ errorMessage }}</div>
      <button type="button" class="btn-close" @click="errorMessage = ''"></button>
    </div>

    <!-- Table List Tahun Ajaran -->
    <div class="card card-custom border-0 shadow-sm overflow-hidden mb-4">
      <div class="table-responsive">
        <table class="table table-hover align-middle mb-0">
          <thead class="table-light">
            <tr>
              <th style="width: 50px;" class="text-center">No</th>
              <th>Tahun Ajaran</th>
              <th>Semester</th>
              <th>Rentang Waktu KBM</th>
              <th>Status Operasional</th>
              <th style="width: 180px;" class="text-center">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="6" class="text-center py-5">
                <i class="fas fa-spinner fa-spin fs-4 text-primary mb-2"></i>
                <div class="text-secondary small">Memuat data periode tahun ajaran...</div>
              </td>
            </tr>

            <tr v-else-if="academicYears.length === 0">
              <td colspan="6" class="text-center py-5">
                <i class="fas fa-calendar-times text-muted fs-1 mb-2"></i>
                <div class="fw-bold text-dark">Belum Ada Tahun Ajaran</div>
                <div class="text-secondary small">Silakan tambahkan periode tahun ajaran baru.</div>
              </td>
            </tr>

            <tr v-else v-for="(item, index) in academicYears" :key="item.id">
              <td class="text-center text-muted fw-semibold">{{ index + 1 }}</td>
              <td>
                <div class="fw-bold text-dark font-monospace fs-6">{{ item.tahun_ajaran }}</div>
                <div class="text-muted small">ID: {{ item.id }}</div>
              </td>
              <td>
                <span class="badge" :class="item.semester === 'Ganjil' ? 'bg-primary-subtle text-primary border border-primary-subtle' : 'bg-info-subtle text-info border border-info-subtle'">
                  Semester {{ item.semester }}
                </span>
              </td>
              <td>
                <div class="small">
                  <i class="far fa-clock me-1 text-muted"></i>
                  <span>{{ formatDate(item.tanggal_mulai) }}</span>
                  <span class="mx-1 text-muted">&rarr;</span>
                  <span>{{ formatDate(item.tanggal_selesai) }}</span>
                </div>
              </td>
              <td>
                <span v-if="item.is_active" class="badge bg-success text-white px-2 py-1">
                  <i class="fas fa-check-circle me-1"></i>Aktif Berjalan
                </span>
                <span v-else class="badge bg-light text-muted border px-2 py-1">
                  Arsip / Non-Aktif
                </span>
              </td>
              <td class="text-center">
                <div class="btn-group btn-group-sm">
                  <button
                    v-if="!item.is_active"
                    @click="setActiveYear(item)"
                    class="btn btn-outline-success"
                    title="Jadikan Tahun Ajaran Aktif Berjalan"
                    :disabled="submitting"
                  >
                    <i class="fas fa-bolt"></i>
                    <span class="d-none d-lg-inline ms-1">Aktifkan</span>
                  </button>
                  <button @click="openEditModal(item)" class="btn btn-outline-primary" title="Ubah Periode">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button @click="openDeleteModal(item)" class="btn btn-outline-danger" title="Hapus Periode" :disabled="item.is_active">
                    <i class="fas fa-trash-alt"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form Tambah / Ubah Tahun Ajaran -->
    <div class="modal fade" id="yearModal" tabindex="-1" ref="yearModalRef" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content border-0 shadow">
          <div class="modal-header border-bottom pb-3">
            <h5 class="modal-title fw-bold text-dark">
              <i class="fas" :class="isEditMode ? 'fa-edit text-warning' : 'fa-plus-circle text-primary'"></i>
              <span class="ms-2">{{ isEditMode ? 'Ubah Periode Tahun Ajaran' : 'Tambah Tahun Ajaran Baru' }}</span>
            </h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>

          <form @submit.prevent="submitForm">
            <div class="modal-body py-3">
              <div v-if="modalError" class="alert alert-danger small py-2 d-flex align-items-center gap-2 mb-3">
                <i class="fas fa-exclamation-circle"></i>
                <div>{{ modalError }}</div>
              </div>

              <div class="row g-2 mb-3">
                <div class="col-md-7">
                  <label class="form-label small fw-semibold">Tahun Ajaran <span class="text-danger">*</span></label>
                  <input
                    v-model="form.tahun_ajaran"
                    type="text"
                    class="form-control form-control-sm font-monospace"
                    placeholder="2026/2027"
                    required
                  />
                  <div class="text-muted" style="font-size: 0.72rem;">Format standar: YYYY/YYYY (contoh: 2026/2027)</div>
                </div>

                <div class="col-md-5">
                  <label class="form-label small fw-semibold">Semester <span class="text-danger">*</span></label>
                  <select v-model="form.semester" class="form-select form-select-sm" required>
                    <option value="Ganjil">Ganjil</option>
                    <option value="Genap">Genap</option>
                  </select>
                </div>
              </div>

              <div class="row g-2 mb-3">
                <div class="col-md-6">
                  <label class="form-label small fw-semibold">Tanggal Mulai KBM <span class="text-danger">*</span></label>
                  <input
                    v-model="form.tanggal_mulai"
                    type="date"
                    class="form-control form-control-sm"
                    required
                  />
                </div>

                <div class="col-md-6">
                  <label class="form-label small fw-semibold">Tanggal Selesai KBM <span class="text-danger">*</span></label>
                  <input
                    v-model="form.tanggal_selesai"
                    type="date"
                    class="form-control form-control-sm"
                    required
                  />
                </div>
              </div>

              <div class="form-check form-switch mt-2">
                <input
                  v-model="form.is_active"
                  class="form-check-input"
                  type="checkbox"
                  role="switch"
                  id="isActiveSwitch"
                />
                <label class="form-check-label small fw-semibold" for="isActiveSwitch">
                  Jadikan Sebagai Tahun Ajaran Aktif Berjalan
                </label>
                <div class="text-muted" style="font-size: 0.72rem;">
                  Mengaktifkan periode ini akan otomatis mengarsipkan periode aktif sebelumnya.
                </div>
              </div>
            </div>

            <div class="modal-footer border-top pt-2">
              <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal" :disabled="submitting">
                Batal
              </button>
              <button type="submit" class="btn btn-primary btn-sm d-flex align-items-center gap-1" :disabled="submitting">
                <i v-if="submitting" class="fas fa-spinner fa-spin"></i>
                <i v-else class="fas fa-save"></i>
                <span>{{ submitting ? 'Menyimpan...' : (isEditMode ? 'Simpan Perubahan' : 'Tambahkan') }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Modal Konfirmasi Hapus -->
    <div class="modal fade" id="deleteModal" tabindex="-1" ref="deleteModalRef" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered modal-sm">
        <div class="modal-content border-0 shadow text-center p-3">
          <div class="modal-body">
            <div class="text-danger mb-3">
              <i class="fas fa-exclamation-triangle fs-1"></i>
            </div>
            <h6 class="fw-bold text-dark mb-1">Hapus Tahun Ajaran?</h6>
            <p class="text-secondary small mb-3">
              Apakah Anda yakin ingin menghapus periode <strong>{{ selectedItem?.tahun_ajaran }} ({{ selectedItem?.semester }})</strong>?
            </p>
            <div class="d-flex justify-content-center gap-2">
              <button type="button" class="btn btn-outline-secondary btn-sm px-3" data-bs-dismiss="modal" :disabled="submitting">
                Batal
              </button>
              <button @click="confirmDelete" type="button" class="btn btn-danger btn-sm px-3" :disabled="submitting">
                <i v-if="submitting" class="fas fa-spinner fa-spin me-1"></i>
                Hapus
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Modal } from 'bootstrap'
import { academicYearsApi } from '@/api/academicYearsApi'

const academicYears = ref([])
const loading = ref(false)
const submitting = ref(false)
const successMessage = ref('')
const errorMessage = ref('')
const modalError = ref('')

const yearModalRef = ref(null)
const deleteModalRef = ref(null)
let yearModalInstance = null
let deleteModalInstance = null

const isEditMode = ref(false)
const selectedItem = ref(null)
const form = ref({
  tahun_ajaran: '',
  semester: 'Ganjil',
  tanggal_mulai: '',
  tanggal_selesai: '',
  is_active: false
})

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
  } catch {
    return dateStr
  }
}

const fetchYears = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await academicYearsApi.getAll()
    academicYears.value = res.data?.data || []
  } catch (err) {
    errorMessage.value = err.response?.data?.message || 'Gagal mengambil data tahun ajaran'
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  isEditMode.value = false
  selectedItem.value = null
  modalError.value = ''
  form.value = {
    tahun_ajaran: '',
    semester: 'Ganjil',
    tanggal_mulai: '',
    tanggal_selesai: '',
    is_active: false
  }
  yearModalInstance?.show()
}

const openEditModal = (item) => {
  isEditMode.value = true
  selectedItem.value = item
  modalError.value = ''
  form.value = {
    tahun_ajaran: item.tahun_ajaran || '',
    semester: item.semester || 'Ganjil',
    tanggal_mulai: item.tanggal_mulai ? item.tanggal_mulai.substring(0, 10) : '',
    tanggal_selesai: item.tanggal_selesai ? item.tanggal_selesai.substring(0, 10) : '',
    is_active: item.is_active || false
  }
  yearModalInstance?.show()
}

const openDeleteModal = (item) => {
  selectedItem.value = item
  deleteModalInstance?.show()
}

const setActiveYear = async (item) => {
  submitting.value = true
  try {
    await academicYearsApi.setActive(item.id)
    successMessage.value = `Tahun ajaran ${item.tahun_ajaran} (${item.semester}) berhasil diaktifkan!`
    await fetchYears()
  } catch (err) {
    errorMessage.value = err.response?.data?.message || 'Gagal mengaktifkan tahun ajaran'
  } finally {
    submitting.value = false
  }
}

const submitForm = async () => {
  modalError.value = ''
  submitting.value = true
  try {
    if (isEditMode.value) {
      await academicYearsApi.update(selectedItem.value.id, form.value)
      successMessage.value = `Berhasil memperbarui periode ${form.value.tahun_ajaran}`
    } else {
      await academicYearsApi.create(form.value)
      successMessage.value = `Berhasil mendaftarkan tahun ajaran baru ${form.value.tahun_ajaran}`
    }
    yearModalInstance?.hide()
    await fetchYears()
  } catch (err) {
    modalError.value = err.response?.data?.message || 'Gagal menyimpan periode tahun ajaran'
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async () => {
  if (!selectedItem.value) return
  submitting.value = true
  try {
    await academicYearsApi.delete(selectedItem.value.id)
    successMessage.value = `Periode ${selectedItem.value.tahun_ajaran} berhasil dihapus`
    deleteModalInstance?.hide()
    await fetchYears()
  } catch (err) {
    errorMessage.value = err.response?.data?.message || 'Gagal menghapus periode tahun ajaran'
    deleteModalInstance?.hide()
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  if (yearModalRef.value) {
    yearModalInstance = new Modal(yearModalRef.value)
  }
  if (deleteModalRef.value) {
    deleteModalInstance = new Modal(deleteModalRef.value)
  }
  fetchYears()
})
</script>

<style scoped>
.card-custom {
  border-radius: 0.75rem;
}
</style>
