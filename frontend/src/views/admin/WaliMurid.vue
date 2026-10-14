<template>
  <div class="container-fluid py-4">
    <!-- Header Page -->
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
      <div>
        <h4 class="fw-bold text-dark mb-1">
          <i class="fas fa-user-friends text-primary me-2"></i>Direktori Wali Murid
        </h4>
        <p class="text-secondary small mb-0">
          Kelola profil orang tua / wali siswa, kontak darurat, serta tautan anak asuhan di sekolah.
        </p>
      </div>

      <div class="d-flex align-items-center gap-2">
        <button @click="fetchGuardians" class="btn btn-outline-secondary btn-sm" :disabled="loading">
          <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
          <span class="d-none d-sm-inline ms-1">Segarkan</span>
        </button>
        <button @click="openCreateModal" class="btn btn-primary btn-sm d-flex align-items-center gap-1 shadow-sm">
          <i class="fas fa-user-plus"></i>
          <span>Tambah Wali Murid</span>
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
              placeholder="Cari nama wali, NIK, telepon, atau nama anak..."
            />
            <button v-if="searchQuery" @click="clearSearch" class="btn btn-outline-secondary" type="button">
              <i class="fas fa-times"></i>
            </button>
          </div>
        </div>
        <div class="col-md-6 col-lg-8 text-md-end text-muted small">
          Total Wali Terdaftar: <strong class="text-dark">{{ guardians.length }}</strong> orang
        </div>
      </div>
    </div>

    <!-- Table List Wali Murid -->
    <div class="card card-custom border-0 shadow-sm overflow-hidden mb-4">
      <div class="table-responsive">
        <table class="table table-hover align-middle mb-0">
          <thead class="table-light">
            <tr>
              <th style="width: 50px;" class="text-center">No</th>
              <th>Nama Wali Murid</th>
              <th>Hubungan & NIK</th>
              <th>Kontak (Telepon / WA)</th>
              <th>Pekerjaan</th>
              <th>Anak Asuhan Terhubung</th>
              <th style="width: 140px;" class="text-center">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="7" class="text-center py-5">
                <i class="fas fa-spinner fa-spin fs-4 text-primary mb-2"></i>
                <div class="text-secondary small">Memuat data direktori wali murid...</div>
              </td>
            </tr>

            <tr v-else-if="filteredGuardians.length === 0">
              <td colspan="7" class="text-center py-5">
                <i class="fas fa-user-slash text-muted fs-1 mb-2"></i>
                <div class="fw-bold text-dark">Data Wali Murid Tidak Ditemukan</div>
                <div class="text-secondary small">Belum ada wali murid yang terdaftar sesuai pencarian.</div>
              </td>
            </tr>

            <tr v-else v-for="(g, index) in filteredGuardians" :key="g.id">
              <td class="text-center text-muted fw-semibold">{{ index + 1 }}</td>
              <td>
                <div class="d-flex align-items-center gap-2">
                  <div class="rounded-circle bg-success-subtle text-success d-flex align-items-center justify-content-center fw-bold" style="width: 34px; height: 34px; font-size: 0.85rem;">
                    {{ g.nama ? g.nama.substring(0, 1).toUpperCase() : 'W' }}
                  </div>
                  <div>
                    <div class="fw-bold text-dark">{{ g.nama }}</div>
                    <div class="text-muted" style="font-size: 0.75rem;">{{ g.email || '-' }}</div>
                  </div>
                </div>
              </td>
              <td>
                <div>
                  <span class="badge bg-primary-subtle text-primary border border-primary-subtle">
                    {{ g.hubungan || 'Wali' }}
                  </span>
                </div>
                <div class="text-muted font-monospace mt-1" style="font-size: 0.75rem;">
                  NIK: {{ g.nik || '-' }}
                </div>
              </td>
              <td>
                <div class="small">
                  <a v-if="g.telepon && g.telepon !== '-'" :href="'tel:' + g.telepon" class="text-decoration-none text-dark">
                    <i class="fas fa-phone-alt text-success me-1"></i>{{ g.telepon }}
                  </a>
                  <span v-else class="text-muted fst-italic">-</span>
                </div>
              </td>
              <td class="small text-secondary">
                {{ g.pekerjaan || '-' }}
              </td>
              <td>
                <div v-if="g.students && g.students.length > 0" class="d-flex flex-wrap gap-1">
                  <span
                    v-for="st in g.students"
                    :key="st.siswa_id || st.id"
                    class="badge bg-light text-dark border d-flex align-items-center gap-1"
                    style="font-size: 0.75rem;"
                  >
                    <i class="fas fa-user-graduate text-primary"></i>
                    <span>{{ st.nama_siswa || st.nama }}</span>
                    <span v-if="st.kelas" class="text-muted">({{ st.kelas }})</span>
                  </span>
                </div>
                <div v-else class="text-muted small fst-italic">
                  Belum ditautkan siswa
                </div>
              </td>
              <td class="text-center">
                <div class="btn-group btn-group-sm">
                  <button @click="openLinkModal(g)" class="btn btn-outline-success" title="Tautkan Anak / Siswa">
                    <i class="fas fa-link"></i>
                  </button>
                  <button @click="openEditModal(g)" class="btn btn-outline-primary" title="Ubah Data">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button @click="openDeleteModal(g)" class="btn btn-outline-danger" title="Hapus Wali">
                    <i class="fas fa-trash-alt"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form Tambah / Ubah Data Wali Murid -->
    <div class="modal fade" id="guardianModal" tabindex="-1" ref="guardianModalRef" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content border-0 shadow">
          <div class="modal-header border-bottom pb-3">
            <h5 class="modal-title fw-bold text-dark">
              <i class="fas" :class="isEditMode ? 'fa-edit text-warning' : 'fa-user-plus text-primary'"></i>
              <span class="ms-2">{{ isEditMode ? 'Ubah Data Wali Murid' : 'Tambah Wali Murid Baru' }}</span>
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
                <label class="form-label small fw-semibold">Nama Lengkap Wali Murid <span class="text-danger">*</span></label>
                <input
                  v-model="form.nama"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder="Nama sesuai KTP"
                  required
                />
              </div>

              <div class="row g-2 mb-3">
                <div class="col-md-6">
                  <label class="form-label small fw-semibold">Hubungan Keluarga <span class="text-danger">*</span></label>
                  <select v-model="form.hubungan" class="form-select form-select-sm" required>
                    <option value="Ayah Kandung">Ayah Kandung</option>
                    <option value="Ibu Kandung">Ibu Kandung</option>
                    <option value="Wali / Keluarga">Wali / Keluarga</option>
                    <option value="Kakek / Nenek">Kakek / Nenek</option>
                    <option value="Paman / Bibi">Paman / Bibi</option>
                  </select>
                </div>

                <div class="col-md-6">
                  <label class="form-label small fw-semibold">NIK (Nomor Induk Kependudukan)</label>
                  <input
                    v-model="form.nik"
                    type="text"
                    class="form-control form-control-sm font-monospace"
                    placeholder="16 Digit NIK"
                    maxlength="16"
                  />
                </div>
              </div>

              <div class="row g-2 mb-3">
                <div class="col-md-6">
                  <label class="form-label small fw-semibold">Nomor WhatsApp / HP <span class="text-danger">*</span></label>
                  <input
                    v-model="form.telepon"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder="0812xxxxxxxx"
                    required
                  />
                </div>

                <div class="col-md-6">
                  <label class="form-label small fw-semibold">Email Kontak</label>
                  <input
                    v-model="form.email"
                    type="email"
                    class="form-control form-control-sm"
                    placeholder="nama@email.com"
                  />
                </div>
              </div>

              <div class="row g-2 mb-3">
                <div class="col-md-6">
                  <label class="form-label small fw-semibold">Pekerjaan</label>
                  <input
                    v-model="form.pekerjaan"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder="Karyawan Swasta / Wiraswasta"
                  />
                </div>

                <div class="col-md-6">
                  <label class="form-label small fw-semibold">Pendidikan Terakhir</label>
                  <select v-model="form.pendidikan" class="form-select form-select-sm">
                    <option value="">Pilih Pendidikan</option>
                    <option value="SD">SD</option>
                    <option value="SMP">SMP</option>
                    <option value="SMA/SMK">SMA / SMK</option>
                    <option value="D3">Diploma 3 (D3)</option>
                    <option value="S1">Sarjana (S1)</option>
                    <option value="S2/S3">Pascasarjana (S2/S3)</option>
                  </select>
                </div>
              </div>

              <div class="mb-2">
                <label class="form-label small fw-semibold">Alamat Rumah Tinggal</label>
                <textarea
                  v-model="form.alamat"
                  rows="2"
                  class="form-control form-control-sm"
                  placeholder="Alamat domisili wali murid"
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
                <span>{{ submitting ? 'Menyimpan...' : (isEditMode ? 'Simpan Perubahan' : 'Tambahkan Wali') }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Modal Penautan Anak / Siswa (Link Student) -->
    <div class="modal fade" id="linkModal" tabindex="-1" ref="linkModalRef" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content border-0 shadow">
          <div class="modal-header border-bottom pb-3">
            <h5 class="modal-title fw-bold text-dark">
              <i class="fas fa-link text-success me-2"></i>Kelola Anak Asuhan
            </h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>

          <div class="modal-body py-3">
            <div class="p-3 bg-light rounded-3 mb-3 border">
              <div class="text-secondary small">Wali Murid Terpilih:</div>
              <div class="fw-bold text-dark fs-6">{{ selectedGuardian?.nama }}</div>
              <div class="text-muted small">Hubungan: {{ selectedGuardian?.hubungan || 'Wali' }} &bull; Kontak: {{ selectedGuardian?.telepon || '-' }}</div>
            </div>

            <div class="mb-3">
              <label class="form-label small fw-semibold">Pilih Siswa untuk Ditautkan</label>
              <div class="input-group input-group-sm">
                <select v-model="selectedStudentIdToLink" class="form-select form-select-sm" :disabled="loadingStudents">
                  <option value="" disabled>{{ loadingStudents ? 'Memuat daftar siswa...' : 'Pilih Siswa' }}</option>
                  <option v-for="s in availableStudents" :key="s.id" :value="s.id">
                    {{ s.nama }} ({{ s.nama_kelas || 'Tanpa Rombel' }}) - NIS: {{ s.nis || '-' }}
                  </option>
                </select>
                <button
                  @click="linkStudent"
                  class="btn btn-success"
                  type="button"
                  :disabled="!selectedStudentIdToLink || linkingStudent"
                >
                  <i v-if="linkingStudent" class="fas fa-spinner fa-spin"></i>
                  <i v-else class="fas fa-plus"></i>
                  <span class="ms-1">Tautkan</span>
                </button>
              </div>
            </div>

            <!-- Daftar Siswa yang Sudah Ditautkan -->
            <div>
              <div class="text-secondary small fw-semibold mb-2">Daftar Anak Asuhan Saat Ini:</div>
              <div v-if="linkedStudents.length === 0" class="text-center py-3 border rounded bg-white text-muted small fst-italic">
                Belum ada siswa yang terhubung dengan wali murid ini.
              </div>
              <div v-else class="d-flex flex-column gap-2">
                <div
                  v-for="st in linkedStudents"
                  :key="st.siswa_id || st.id"
                  class="p-2 border rounded bg-white d-flex justify-content-between align-items-center"
                >
                  <div class="d-flex align-items-center gap-2">
                    <i class="fas fa-user-graduate text-primary"></i>
                    <div>
                      <div class="fw-bold text-dark small">{{ st.nama_siswa || st.nama }}</div>
                      <div class="text-muted" style="font-size: 0.72rem;">Kelas: {{ st.kelas || '-' }}</div>
                    </div>
                  </div>
                  <button
                    @click="unlinkStudent(st.siswa_id || st.id)"
                    class="btn btn-outline-danger btn-sm py-0 px-2 small"
                    title="Lepaskan tautan siswa dari wali murid ini"
                    :disabled="unlinkingStudent"
                  >
                    <i class="fas fa-unlink me-1"></i>Lepas
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div class="modal-footer border-top pt-2">
            <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">
              Tutup
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Konfirmasi Hapus Wali -->
    <div class="modal fade" id="deleteModal" tabindex="-1" ref="deleteModalRef" aria-hidden="true">
      <div class="modal-dialog modal-dialog-centered modal-sm">
        <div class="modal-content border-0 shadow text-center p-3">
          <div class="modal-body">
            <div class="text-danger mb-3">
              <i class="fas fa-exclamation-triangle fs-1"></i>
            </div>
            <h6 class="fw-bold text-dark mb-1">Hapus Data Wali?</h6>
            <p class="text-secondary small mb-3">
              Apakah Anda yakin ingin menghapus data wali murid <strong>{{ selectedGuardian?.nama }}</strong>?
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
import { guardiansApi } from '@/api/guardiansApi'
import { siswaApi } from '@/api/siswaApi'

const guardians = ref([])
const availableStudents = ref([])
const linkedStudents = ref([])
const loading = ref(false)
const loadingStudents = ref(false)
const submitting = ref(false)
const linkingStudent = ref(false)
const unlinkingStudent = ref(false)

const searchQuery = ref('')
const successMessage = ref('')
const errorMessage = ref('')
const modalError = ref('')

const guardianModalRef = ref(null)
const linkModalRef = ref(null)
const deleteModalRef = ref(null)
let guardianModalInstance = null
let linkModalInstance = null
let deleteModalInstance = null

const isEditMode = ref(false)
const selectedGuardian = ref(null)
const selectedStudentIdToLink = ref('')

const form = ref({
  nama: '',
  hubungan: 'Ayah Kandung',
  nik: '',
  telepon: '',
  email: '',
  pekerjaan: '',
  pendidikan: '',
  alamat: ''
})

const filteredGuardians = computed(() => {
  if (!searchQuery.value) return guardians.value
  const q = searchQuery.value.toLowerCase()
  return guardians.value.filter(g =>
    (g.nama && g.nama.toLowerCase().includes(q)) ||
    (g.nik && g.nik.toLowerCase().includes(q)) ||
    (g.telepon && g.telepon.toLowerCase().includes(q)) ||
    (g.hubungan && g.hubungan.toLowerCase().includes(q))
  )
})

const fetchGuardians = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await guardiansApi.getAll()
    guardians.value = res.data?.data || []
  } catch (err) {
    errorMessage.value = err.response?.data?.message || 'Gagal mengambil data wali murid'
  } finally {
    loading.value = false
  }
}

const fetchAvailableStudents = async () => {
  loadingStudents.value = true
  try {
    const res = await siswaApi.getAll()
    availableStudents.value = res.data?.data || []
  } catch {
    availableStudents.value = []
  } finally {
    loadingStudents.value = false
  }
}

const openCreateModal = () => {
  isEditMode.value = false
  selectedGuardian.value = null
  modalError.value = ''
  form.value = {
    nama: '',
    hubungan: 'Ayah Kandung',
    nik: '',
    telepon: '',
    email: '',
    pekerjaan: '',
    pendidikan: '',
    alamat: ''
  }
  guardianModalInstance?.show()
}

const openEditModal = (g) => {
  isEditMode.value = true
  selectedGuardian.value = g
  modalError.value = ''
  form.value = {
    nama: g.nama || '',
    hubungan: g.hubungan || 'Wali',
    nik: g.nik || '',
    telepon: g.telepon || '',
    email: g.email || '',
    pekerjaan: g.pekerjaan || '',
    pendidikan: g.pendidikan || '',
    alamat: g.alamat || ''
  }
  guardianModalInstance?.show()
}

const openLinkModal = async (g) => {
  selectedGuardian.value = g
  selectedStudentIdToLink.value = ''
  linkedStudents.value = []
  linkModalInstance?.show()
  await refreshLinkedStudents(g.id)
  if (availableStudents.value.length === 0) {
    await fetchAvailableStudents()
  }
}

const refreshLinkedStudents = async (guardianId) => {
  try {
    const res = await guardiansApi.getStudentsByGuardianId(guardianId)
    linkedStudents.value = res.data?.data || []
  } catch {
    linkedStudents.value = []
  }
}

const linkStudent = async () => {
  if (!selectedGuardian.value || !selectedStudentIdToLink.value) return
  linkingStudent.value = true
  try {
    await guardiansApi.linkStudent(selectedGuardian.value.id, selectedStudentIdToLink.value, false)
    selectedStudentIdToLink.value = ''
    await refreshLinkedStudents(selectedGuardian.value.id)
    await fetchGuardians()
  } catch (err) {
    errorMessage.value = err.response?.data?.message || 'Gagal menautkan siswa'
  } finally {
    linkingStudent.value = false
  }
}

const unlinkStudent = async (siswaId) => {
  if (!selectedGuardian.value) return
  unlinkingStudent.value = true
  try {
    await guardiansApi.unlinkStudent(selectedGuardian.value.id, siswaId)
    await refreshLinkedStudents(selectedGuardian.value.id)
    await fetchGuardians()
  } catch (err) {
    errorMessage.value = err.response?.data?.message || 'Gagal melepas tautan siswa'
  } finally {
    unlinkingStudent.value = false
  }
}

const openDeleteModal = (g) => {
  selectedGuardian.value = g
  deleteModalInstance?.show()
}

const submitForm = async () => {
  modalError.value = ''
  submitting.value = true
  try {
    if (isEditMode.value) {
      await guardiansApi.update(selectedGuardian.value.id, form.value)
      successMessage.value = `Berhasil memperbarui data wali murid ${form.value.nama}`
    } else {
      await guardiansApi.create(form.value)
      successMessage.value = `Berhasil mendaftarkan wali murid ${form.value.nama}`
    }
    guardianModalInstance?.hide()
    await fetchGuardians()
  } catch (err) {
    modalError.value = err.response?.data?.message || 'Gagal menyimpan data wali murid'
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async () => {
  if (!selectedGuardian.value) return
  submitting.value = true
  try {
    await guardiansApi.delete(selectedGuardian.value.id)
    successMessage.value = `Data wali murid ${selectedGuardian.value.nama} berhasil dihapus`
    deleteModalInstance?.hide()
    await fetchGuardians()
  } catch (err) {
    errorMessage.value = err.response?.data?.message || 'Gagal menghapus data wali murid'
    deleteModalInstance?.hide()
  } finally {
    submitting.value = false
  }
}

const clearSearch = () => {
  searchQuery.value = ''
}

const onSearchInput = () => {
  // Computed filter
}

onMounted(() => {
  if (guardianModalRef.value) {
    guardianModalInstance = new Modal(guardianModalRef.value)
  }
  if (linkModalRef.value) {
    linkModalInstance = new Modal(linkModalRef.value)
  }
  if (deleteModalRef.value) {
    deleteModalInstance = new Modal(deleteModalRef.value)
  }
  fetchGuardians()
})
</script>

<style scoped>
.card-custom {
  border-radius: 0.75rem;
}
</style>
