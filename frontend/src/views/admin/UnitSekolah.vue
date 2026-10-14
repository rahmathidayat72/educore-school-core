<template>
  <div class="container-fluid py-4">
    <!-- Header Page -->
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
      <div>
        <h4 class="fw-bold text-dark mb-1">
          <i class="fas fa-city text-primary me-2"></i>Unit Sekolah Yayasan
        </h4>
        <p class="text-secondary small mb-0">
          Kelola unit pendidikan di bawah naungan yayasan lintas jenjang (TK, SD, SMP, SMA, SMK).
        </p>
      </div>

      <div class="d-flex align-items-center gap-2">
        <button @click="fetchUnits" class="btn btn-outline-secondary btn-sm" :disabled="loading">
          <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
          <span class="d-none d-sm-inline ms-1">Segarkan</span>
        </button>
        <button @click="openCreateModal" class="btn btn-primary btn-sm d-flex align-items-center gap-1 shadow-sm">
          <i class="fas fa-plus"></i>
          <span>Tambah Unit Sekolah</span>
        </button>
      </div>
    </div>

    <!-- Alert Success / Error -->
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

    <!-- Filter & Search Bar -->
    <div class="card card-custom border-0 shadow-sm p-3 mb-4">
      <div class="row g-2 align-items-center">
        <div class="col-md-6 col-lg-4">
          <div class="input-group input-group-sm">
            <span class="input-group-text bg-light border-end-0">
              <i class="fas fa-search text-muted"></i>
            </span>
            <input
              v-model="searchQuery"
              @input="onSearchInput"
              type="text"
              class="form-control form-control-sm border-start-0"
              placeholder="Cari nama unit, kode, atau jenjang..."
            />
            <button v-if="searchQuery" @click="clearSearch" class="btn btn-outline-secondary" type="button">
              <i class="fas fa-times"></i>
            </button>
          </div>
        </div>
        <div class="col-md-6 col-lg-8 text-md-end text-muted small">
          Total Unit: <strong class="text-dark">{{ units.length }}</strong> sekolah terdaftar
        </div>
      </div>
    </div>

    <!-- Table List Unit Sekolah -->
    <div class="card card-custom border-0 shadow-sm overflow-hidden">
      <div class="table-responsive">
        <table class="table table-hover align-middle mb-0">
          <thead class="table-light">
            <tr>
              <th style="width: 50px;" class="text-center">No</th>
              <th>Nama Unit Sekolah</th>
              <th>Jenjang</th>
              <th>Kode Unit</th>
              <th>Kontak & Email</th>
              <th>Alamat</th>
              <th style="width: 120px;" class="text-center">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="7" class="text-center py-5">
                <i class="fas fa-spinner fa-spin fs-4 text-primary mb-2"></i>
                <div class="text-secondary small">Memuat data unit sekolah...</div>
              </td>
            </tr>

            <tr v-else-if="filteredUnits.length === 0">
              <td colspan="7" class="text-center py-5">
                <i class="fas fa-school text-muted fs-1 mb-2"></i>
                <div class="fw-bold text-dark">Tidak Ada Data Unit Sekolah</div>
                <div class="text-secondary small">Belum ada unit sekolah yang sesuai dengan pencarian.</div>
              </td>
            </tr>

            <tr v-else v-for="(unit, index) in filteredUnits" :key="unit.id">
              <td class="text-center text-muted fw-semibold">{{ index + 1 }}</td>
              <td>
                <div class="d-flex align-items-center gap-2">
                  <div class="rounded-circle bg-primary-subtle text-primary d-flex align-items-center justify-content-center fw-bold" style="width: 34px; height: 34px; font-size: 0.85rem;">
                    {{ unit.nama ? unit.nama.substring(0, 2).toUpperCase() : 'US' }}
                  </div>
                  <div>
                    <div class="fw-bold text-dark">{{ unit.nama }}</div>
                    <div class="text-muted" style="font-size: 0.75rem;">ID: {{ unit.id }}</div>
                  </div>
                </div>
              </td>
              <td>
                <span class="badge" :class="getJenjangBadge(unit.jenjang)">
                  {{ unit.jenjang }}
                </span>
              </td>
              <td>
                <span class="badge bg-light text-dark border font-monospace">
                  {{ unit.kode || '-' }}
                </span>
              </td>
              <td>
                <div class="small">
                  <div v-if="unit.telepon"><i class="fas fa-phone-alt me-1 text-muted"></i>{{ unit.telepon }}</div>
                  <div v-if="unit.email" class="text-muted"><i class="fas fa-envelope me-1 text-muted"></i>{{ unit.email }}</div>
                  <div v-if="!unit.telepon && !unit.email" class="text-muted fst-italic">-</div>
                </div>
              </td>
              <td class="small text-secondary" style="max-width: 250px;">
                {{ unit.alamat || '-' }}
              </td>
              <td class="text-center">
                <div class="btn-group btn-group-sm">
                  <button @click="openEditModal(unit)" class="btn btn-outline-primary" title="Ubah Profil Unit">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button @click="openDeleteModal(unit)" class="btn btn-outline-danger" title="Hapus Unit Sekolah">
                    <i class="fas fa-trash-alt"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form Tambah / Ubah Unit -->
    <div class="modal fade" id="unitModal" tabindex="-1" ref="unitModalRef" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content border-0 shadow">
          <div class="modal-header border-bottom pb-3">
            <h5 class="modal-title fw-bold text-dark">
              <i class="fas" :class="isEditMode ? 'fa-edit text-warning' : 'fa-plus-circle text-primary'"></i>
              <span class="ms-2">{{ isEditMode ? 'Ubah Profil Unit Sekolah' : 'Tambah Unit Sekolah Baru' }}</span>
            </h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>

          <form @submit.prevent="submitForm">
            <div class="modal-body py-3">
              <div v-if="modalError" class="alert alert-danger small py-2 d-flex align-items-center gap-2 mb-3">
                <i class="fas fa-exclamation-circle"></i>
                <div>{{ modalError }}</div>
              </div>

              <div class="mb-3">
                <label class="form-label small fw-semibold">Nama Unit Sekolah <span class="text-danger">*</span></label>
                <input
                  v-model="form.nama"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder="Contoh: SMA IT Al-Hikmah"
                  required
                />
              </div>

              <div class="row g-2 mb-3">
                <div class="col-md-6">
                  <label class="form-label small fw-semibold">Jenjang Pendidikan <span class="text-danger">*</span></label>
                  <select v-model="form.jenjang" class="form-select form-select-sm" required>
                    <option value="" disabled>Pilih Jenjang</option>
                    <option value="TK">TK (Pendidikan Anak Usia Dini)</option>
                    <option value="SD">SD (Sekolah Dasar)</option>
                    <option value="SMP">SMP (Sekolah Menengah Pertama)</option>
                    <option value="SMA">SMA (Sekolah Menengah Atas)</option>
                    <option value="SMK">SMK (Sekolah Menengah Kejuruan)</option>
                  </select>
                </div>

                <div class="col-md-6">
                  <label class="form-label small fw-semibold">Kode Singkatan Unit <span class="text-danger">*</span></label>
                  <input
                    v-model="form.kode"
                    type="text"
                    class="form-control form-control-sm text-uppercase font-monospace"
                    placeholder="Contoh: SMA-HIKMAH"
                    required
                  />
                </div>
              </div>

              <div class="row g-2 mb-3">
                <div class="col-md-6">
                  <label class="form-label small fw-semibold">Nomor Telepon Kantor</label>
                  <input
                    v-model="form.telepon"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder="021-xxxxxxx"
                  />
                </div>

                <div class="col-md-6">
                  <label class="form-label small fw-semibold">Email Resmi Sekolah</label>
                  <input
                    v-model="form.email"
                    type="email"
                    class="form-control form-control-sm"
                    placeholder="sma@alhikmah.sch.id"
                  />
                </div>
              </div>

              <div class="mb-2">
                <label class="form-label small fw-semibold">Alamat Lengkap</label>
                <textarea
                  v-model="form.alamat"
                  rows="3"
                  class="form-control form-control-sm"
                  placeholder="Jl. Raya Pendidikan No. 45, Jakarta Selatan"
                ></textarea>
              </div>
            </div>

            <div class="modal-footer border-top pt-2">
              <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal" :disabled="submitting">
                Batal
              </button>
              <button type="submit" class="btn btn-primary btn-sm d-flex align-items-center gap-1" :disabled="submitting">
                <i v-if="submitting" class="fas fa-spinner fa-spin"></i>
                <i v-else class="fas fa-save"></i>
                <span>{{ submitting ? 'Menyimpan...' : (isEditMode ? 'Simpan Perubahan' : 'Tambahkan Unit') }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Modal Konfirmasi Hapus Unit -->
    <div class="modal fade" id="deleteModal" tabindex="-1" ref="deleteModalRef" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered modal-sm">
        <div class="modal-content border-0 shadow text-center p-3">
          <div class="modal-body">
            <div class="text-danger mb-3">
              <i class="fas fa-exclamation-triangle fs-1"></i>
            </div>
            <h6 class="fw-bold text-dark mb-1">Hapus Unit Sekolah?</h6>
            <p class="text-secondary small mb-3">
              Apakah Anda yakin ingin menonaktifkan unit <strong>{{ selectedUnit?.nama }}</strong>?
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
import { ref, computed, onMounted } from 'vue'
import { Modal } from 'bootstrap'
import { schoolUnitsApi } from '@/api/schoolUnitsApi'

const units = ref([])
const loading = ref(false)
const submitting = ref(false)
const searchQuery = ref('')
const successMessage = ref('')
const errorMessage = ref('')
const modalError = ref('')

// Modal references
const unitModalRef = ref(null)
const deleteModalRef = ref(null)
let unitModalInstance = null
let deleteModalInstance = null

// Form State
const isEditMode = ref(false)
const selectedUnit = ref(null)
const form = ref({
  nama: '',
  jenjang: '',
  kode: '',
  telepon: '',
  email: '',
  alamat: ''
})

const filteredUnits = computed(() => {
  if (!searchQuery.value) return units.value
  const q = searchQuery.value.toLowerCase()
  return units.value.filter(u =>
    (u.nama && u.nama.toLowerCase().includes(q)) ||
    (u.kode && u.kode.toLowerCase().includes(q)) ||
    (u.jenjang && u.jenjang.toLowerCase().includes(q))
  )
})

const getJenjangBadge = (jenjang) => {
  switch (jenjang?.toUpperCase()) {
    case 'TK': return 'bg-warning-subtle text-warning-emphasis border border-warning-subtle'
    case 'SD': return 'bg-danger-subtle text-danger border border-danger-subtle'
    case 'SMP': return 'bg-primary-subtle text-primary border border-primary-subtle'
    case 'SMA': return 'bg-info-subtle text-info border border-info-subtle'
    case 'SMK': return 'bg-success-subtle text-success border border-success-subtle'
    default: return 'bg-secondary-subtle text-secondary border'
  }
}

const fetchUnits = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await schoolUnitsApi.getAll()
    units.value = res.data?.data || []
  } catch (err) {
    errorMessage.value = err.response?.data?.message || 'Gagal mengambil data unit sekolah yayasan'
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  isEditMode.value = false
  selectedUnit.value = null
  modalError.value = ''
  form.value = {
    nama: '',
    jenjang: '',
    kode: '',
    telepon: '',
    email: '',
    alamat: ''
  }
  unitModalInstance?.show()
}

const openEditModal = (unit) => {
  isEditMode.value = true
  selectedUnit.value = unit
  modalError.value = ''
  form.value = {
    nama: unit.nama || '',
    jenjang: unit.jenjang || '',
    kode: unit.kode || '',
    telepon: unit.telepon || '',
    email: unit.email || '',
    alamat: unit.alamat || ''
  }
  unitModalInstance?.show()
}

const openDeleteModal = (unit) => {
  selectedUnit.value = unit
  deleteModalInstance?.show()
}

const submitForm = async () => {
  modalError.value = ''
  submitting.value = true
  try {
    if (isEditMode.value) {
      await schoolUnitsApi.update(selectedUnit.value.id, form.value)
      successMessage.value = `Berhasil memperbarui data unit ${form.value.nama}`
    } else {
      await schoolUnitsApi.create(form.value)
      successMessage.value = `Berhasil mendaftarkan unit sekolah baru ${form.value.nama}`
    }
    unitModalInstance?.hide()
    await fetchUnits()
  } catch (err) {
    modalError.value = err.response?.data?.message || 'Gagal menyimpan data unit sekolah'
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async () => {
  if (!selectedUnit.value) return
  submitting.value = true
  try {
    await schoolUnitsApi.delete(selectedUnit.value.id)
    successMessage.value = `Unit sekolah ${selectedUnit.value.nama} berhasil dihapus`
    deleteModalInstance?.hide()
    await fetchUnits()
  } catch (err) {
    errorMessage.value = err.response?.data?.message || 'Gagal menghapus unit sekolah'
    deleteModalInstance?.hide()
  } finally {
    submitting.value = false
  }
}

const clearSearch = () => {
  searchQuery.value = ''
}

const onSearchInput = () => {
  // Reactive through computed
}

onMounted(() => {
  if (unitModalRef.value) {
    unitModalInstance = new Modal(unitModalRef.value)
  }
  if (deleteModalRef.value) {
    deleteModalInstance = new Modal(deleteModalRef.value)
  }
  fetchUnits()
})
</script>

<style scoped>
.card-custom {
  border-radius: 0.75rem;
}
.table > :not(caption) > * > * {
  padding: 0.85rem 0.75rem;
}
</style>
