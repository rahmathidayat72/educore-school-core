<template>
  <div class="animate-fade-up">
    <!-- Header -->
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
      <div>
        <div class="d-flex align-items-center gap-2 mb-1">
          <h3 class="fw-bold text-dark mb-0">Data Tenaga Pendidik (Guru)</h3>
          <span class="badge bg-primary-subtle text-primary px-2 py-1 small">
            {{ filteredGuru.length }} Guru
          </span>
        </div>
        <p class="text-secondary small mb-0">Kelola informasi biodata lengkap guru, status kepegawaian, dan akun akses pendidik.</p>
      </div>
      <div class="d-flex align-items-center gap-2">
        <button @click="openImportModal" class="btn btn-outline-info d-flex align-items-center gap-2 shadow-sm" title="Import data guru secara masal dari Excel/CSV">
          <i class="fas fa-file-import"></i>
          <span>Import CSV</span>
        </button>
        <button @click="exportDataGuru" class="btn btn-outline-primary d-flex align-items-center gap-2 shadow-sm" :disabled="loading || filteredGuru.length === 0" title="Unduh data guru ke file Excel/CSV">
          <i class="fas fa-file-excel"></i>
          <span>Ekspor Excel</span>
        </button>
        <button @click="openCreateModal" class="btn btn-primary btn-shine d-flex align-items-center gap-2 shadow-sm">
          <i class="fas fa-plus"></i>
          <span>Tambah Guru</span>
        </button>
      </div>
    </div>

    <!-- Alert Notifikasi -->
    <Transition name="alert-pop">
      <div v-if="alertMessage" :class="['alert alert-dismissible fade show small py-2 d-flex align-items-center gap-2 mb-4 shadow-sm', alertClass]">
        <i :class="alertIcon"></i>
        <div>{{ alertMessage }}</div>
        <button type="button" class="btn-close" @click="alertMessage = ''"></button>
      </div>
    </Transition>

    <!-- Search & Filter Card -->
    <div class="card card-custom p-3 mb-4 border-0 shadow-sm">
      <div class="row g-2 align-items-center">
        <div class="col-md-6">
          <div class="input-group">
            <span class="input-group-text bg-white text-muted border-end-0"><i class="fas fa-search"></i></span>
            <input
              v-model="searchQuery"
              type="text"
              class="form-control border-start-0 ps-0"
              placeholder="Cari berdasarkan nama, NIP, NUPTK, email, atau telepon..."
            />
            <button v-if="searchQuery" class="btn btn-link text-secondary p-0 px-2" @click="searchQuery = ''" title="Bersihkan">
              <i class="fas fa-times-circle"></i>
            </button>
          </div>
        </div>
        <div class="col-md-4">
          <select v-model="filterStatus" class="form-select">
            <option value="">Semua Status Kepegawaian</option>
            <option value="PNS">PNS</option>
            <option value="PPPK">PPPK</option>
            <option value="Guru Tetap Yayasan">Guru Tetap Yayasan</option>
            <option value="Guru Tidak Tetap">Guru Tidak Tetap</option>
            <option value="Honorer">Honorer</option>
          </select>
        </div>
        <div class="col-md-2 text-md-end">
          <button @click="fetchGuru" class="btn btn-outline-secondary w-100 d-flex align-items-center justify-content-center gap-2" :disabled="loading">
            <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
            <span>Segarkan</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Tabel Data Guru -->
    <div class="card card-custom overflow-hidden border-0 shadow-sm">
      <div class="table-responsive">
        <table class="table table-custom table-hover mb-0">
          <thead>
            <tr>
              <th style="width: 50px;">No</th>
              <th>Nama &amp; Identitas</th>
              <th>JK / Agama</th>
              <th>Status Pegawai</th>
              <th>Kontak</th>
              <th class="text-end" style="width: 150px;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="6" class="text-center py-5 text-secondary">
                <div class="spinner-border spinner-border-sm text-primary me-2" role="status"></div>
                Memuat data guru...
              </td>
            </tr>
            <tr v-else-if="filteredGuru.length === 0">
              <td colspan="6" class="text-center py-5 text-muted">
                <i class="fas fa-chalkboard-teacher fs-2 mb-2 d-block text-secondary"></i>
                Belum ada data guru yang sesuai kriteria.
              </td>
            </tr>
            <tr v-for="(g, index) in filteredGuru" :key="g.id">
              <td class="text-secondary fw-semibold">{{ index + 1 }}</td>
              <td>
                <div class="d-flex align-items-center gap-2">
                  <div class="avatar-teacher rounded-circle d-flex align-items-center justify-content-center flex-shrink-0">
                    <i class="fas fa-user-tie small"></i>
                  </div>
                  <div>
                    <div class="fw-bold text-dark">{{ g.nama }}</div>
                    <div class="small text-muted d-flex align-items-center gap-2">
                      <span v-if="g.nip"><i class="fas fa-id-badge text-primary me-1"></i>NIP: {{ g.nip }}</span>
                      <span v-else-if="g.nuptk"><i class="fas fa-id-card text-info me-1"></i>NUPTK: {{ g.nuptk }}</span>
                      <span v-else class="text-secondary opacity-75">NIP/NUPTK: -</span>
                    </div>
                  </div>
                </div>
              </td>
              <td>
                <div class="small text-dark fw-semibold">{{ g.jenis_kelamin || '-' }}</div>
                <div class="small text-muted">{{ g.agama || '-' }}</div>
              </td>
              <td>
                <span class="badge" :class="getStatusBadgeClass(g.status_kepegawaian)">
                  {{ g.status_kepegawaian || 'Belum Diatur' }}
                </span>
                <div v-if="g.pendidikan_terakhir" class="small text-muted mt-1">
                  Pendidikan: {{ g.pendidikan_terakhir }}
                </div>
              </td>
              <td>
                <div class="small text-dark"><i class="fas fa-envelope text-secondary me-1"></i>{{ g.email }}</div>
                <div v-if="g.no_telepon" class="small text-muted"><i class="fas fa-phone text-success me-1"></i>{{ g.no_telepon }}</div>
              </td>
              <td class="text-end">
                <div class="btn-group btn-group-sm">
                  <button type="button" @click.stop="openDetailModal(g)" class="btn btn-outline-info" title="Lihat Profil Lengkap">
                    <i class="fas fa-eye"></i>
                  </button>
                  <button type="button" @click.stop="openEditModal(g)" class="btn btn-outline-primary" title="Ubah Biodata Guru">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button type="button" @click.stop="openDeleteModal(g)" class="btn btn-outline-danger" title="Hapus Data Guru">
                    <i class="fas fa-trash-alt"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form (Tambah / Ubah) -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showModal" class="app-modal-overlay" tabindex="-1" @click.self="closeModal">
          <div class="modal-dialog modal-dialog-centered modal-lg modal-dialog-scrollable">
            <div class="modal-content border-0 shadow-lg bg-white">
              <div class="modal-header bg-light">
                <h5 class="modal-title fw-bold text-dark d-flex align-items-center gap-2">
                  <i :class="isEditing ? 'fas fa-user-edit text-primary' : 'fas fa-user-plus text-primary'"></i>
                  <span>{{ isEditing ? 'Ubah Biodata Lengkap Guru' : 'Tambah Guru Baru' }}</span>
                </h5>
                <button type="button" class="btn-close" @click="closeModal"></button>
              </div>
              <form @submit.prevent="handleSubmit" class="d-flex flex-column overflow-hidden">
                <div class="modal-body p-4 overflow-y-auto" style="max-height: calc(85vh - 130px);">
                  <!-- 1. Akun Akses Pengguna -->
                  <div v-if="!isEditing" class="card p-3 mb-3 bg-light border-0">
                    <div class="fw-bold text-dark small mb-2 d-flex align-items-center gap-2">
                      <i class="fas fa-user-shield text-primary"></i>
                      <span>Akun Pengguna Aplikasi (Wajib)</span>
                    </div>
                    <select v-model="form.id_user" class="form-select form-select-sm" required>
                      <option value="" disabled>Pilih Akun Pengguna...</option>
                      <option v-for="user in availableUsers" :key="user.id" :value="user.id">
                        {{ user.username }} ({{ user.email }}) - {{ user.role }}
                      </option>
                    </select>
                    <div class="form-text small">Guru ini akan memiliki akses login sesuai akun pengguna yang dipilih.</div>
                  </div>

                  <!-- 2. Identitas Pokok -->
                  <div class="mb-3">
                    <div class="section-title text-primary small fw-bold text-uppercase mb-2">
                      <i class="fas fa-id-card me-1"></i>Identitas Pokok &amp; Kepegawaian
                    </div>
                    <div class="row g-3">
                      <div class="col-md-12">
                        <label class="form-label small fw-semibold text-secondary">Nama Lengkap &amp; Gelar <span class="text-danger">*</span></label>
                        <input v-model="form.nama" type="text" class="form-control" required placeholder="Contoh: Dra. Hj. Siti Aminah, M.Pd" />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">NIP (Nomor Induk Pegawai)</label>
                        <input v-model="form.nip" type="text" class="form-control" placeholder="19800101 200501 1 001" />
                        <div class="form-text small">Opsional jika belum memiliki NIP.</div>
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">NUPTK</label>
                        <input v-model="form.nuptk" type="text" class="form-control" placeholder="16 digit NUPTK..." />
                        <div class="form-text small">Nomor Unik Pendidik dan Tenaga Kependidikan.</div>
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Status Kepegawaian</label>
                        <select v-model="form.status_kepegawaian" class="form-select">
                          <option value="">Pilih Status Kepegawaian...</option>
                          <option value="PNS">PNS (Pegawai Negeri Sipil)</option>
                          <option value="PPPK">PPPK (Pegawai Pemerintah dg Perjanjian Kerja)</option>
                          <option value="Guru Tetap Yayasan">Guru Tetap Yayasan (GTY)</option>
                          <option value="Guru Tidak Tetap">Guru Tidak Tetap (GTT)</option>
                          <option value="Honorer">Guru Honorer Sekolah</option>
                        </select>
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Pendidikan Terakhir</label>
                        <select v-model="form.pendidikan_terakhir" class="form-select">
                          <option value="">Pilih Pendidikan...</option>
                          <option value="D3">Diploma 3 (D3)</option>
                          <option value="D4 / S1">Sarjana (S1 / D4)</option>
                          <option value="S2">Magister (S2)</option>
                          <option value="S3">Doktor (S3)</option>
                          <option value="Lainnya">Lainnya</option>
                        </select>
                      </div>
                    </div>
                  </div>

                  <hr class="my-3 opacity-25" />

                  <!-- 3. Biodata Diri -->
                  <div class="mb-3">
                    <div class="section-title text-primary small fw-bold text-uppercase mb-2">
                      <i class="fas fa-user-circle me-1"></i>Biodata Pribadi
                    </div>
                    <div class="row g-3">
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Jenis Kelamin</label>
                        <select v-model="form.jenis_kelamin" class="form-select">
                          <option value="">Pilih Jenis Kelamin...</option>
                          <option value="Laki-laki">Laki-laki</option>
                          <option value="Perempuan">Perempuan</option>
                        </select>
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Agama</label>
                        <select v-model="form.agama" class="form-select">
                          <option value="">Pilih Agama...</option>
                          <option value="Islam">Islam</option>
                          <option value="Kristen Protestan">Kristen Protestan</option>
                          <option value="Katolik">Katolik</option>
                          <option value="Hindu">Hindu</option>
                          <option value="Buddha">Buddha</option>
                          <option value="Khonghucu">Khonghucu</option>
                          <option value="Lainnya">Lainnya</option>
                        </select>
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Tempat Lahir</label>
                        <input v-model="form.tempat_lahir" type="text" class="form-control" placeholder="Contoh: Bandung" />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Tanggal Lahir</label>
                        <input v-model="form.tanggal_lahir" type="date" class="form-control" />
                      </div>
                    </div>
                  </div>

                  <hr class="my-3 opacity-25" />

                  <!-- 4. Kontak & Alamat -->
                  <div class="mb-2">
                    <div class="section-title text-primary small fw-bold text-uppercase mb-2">
                      <i class="fas fa-map-marker-alt me-1"></i>Kontak &amp; Tempat Tinggal
                    </div>
                    <div class="row g-3">
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Email Resmi / Utama <span class="text-danger">*</span></label>
                        <input v-model="form.email" type="email" class="form-control" required placeholder="siti.aminah@sekolah.id" />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">No. Telepon / WhatsApp</label>
                        <input v-model="form.no_telepon" type="tel" class="form-control" placeholder="Contoh: 08123456789" />
                      </div>
                      <div class="col-12">
                        <label class="form-label small fw-semibold text-secondary">Alamat Tempat Tinggal <span class="text-danger">*</span></label>
                        <textarea v-model="form.alamat" rows="2" class="form-control" required placeholder="Jl. Melati No. 12, Kelurahan X, Kota Y"></textarea>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="modal-footer bg-light">
                  <button type="button" class="btn btn-light" @click="closeModal">Batal</button>
                  <button type="submit" class="btn btn-primary btn-shine d-flex align-items-center gap-2" :disabled="submitting">
                    <span v-if="submitting" class="spinner-border spinner-border-sm"></span>
                    <span>{{ isEditing ? 'Simpan Perubahan' : 'Tambah Guru' }}</span>
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Modal Detail Lengkap Guru -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showDetailModal" class="app-modal-overlay" tabindex="-1" @click.self="closeDetailModal">
          <div class="modal-dialog modal-dialog-centered modal-lg">
            <div class="modal-content border-0 shadow-lg bg-white">
              <div class="modal-header bg-primary text-white">
                <h5 class="modal-title fw-bold d-flex align-items-center gap-2">
                  <i class="fas fa-id-badge"></i>
                  <span>Biodata Lengkap Tenaga Pendidik</span>
                </h5>
                <button type="button" class="btn-close btn-close-white" @click="closeDetailModal"></button>
              </div>
              <div v-if="selectedGuru" class="modal-body p-4">
                <!-- Profil Header Card -->
                <div class="d-flex flex-column flex-sm-row align-items-sm-center gap-3 p-3 bg-light rounded-3 mb-4">
                  <div class="avatar-detail rounded-circle d-flex align-items-center justify-content-center text-primary fs-2">
                    <i class="fas fa-chalkboard-teacher"></i>
                  </div>
                  <div class="flex-grow-1">
                    <h4 class="fw-bold text-dark mb-1">{{ selectedGuru.nama }}</h4>
                    <div class="d-flex flex-wrap gap-2 align-items-center">
                      <span class="badge" :class="getStatusBadgeClass(selectedGuru.status_kepegawaian)">
                        {{ selectedGuru.status_kepegawaian || 'Status Pegawai Belum Diatur' }}
                      </span>
                      <span v-if="selectedGuru.pendidikan_terakhir" class="badge bg-secondary-subtle text-secondary">
                        <i class="fas fa-graduation-cap me-1"></i>{{ selectedGuru.pendidikan_terakhir }}
                      </span>
                      <span class="text-muted small">
                        <i class="fas fa-envelope me-1"></i>{{ selectedGuru.email }}
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Informasi Terperinci -->
                <div class="row g-3">
                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <h6 class="fw-bold text-primary mb-3 small text-uppercase">
                        <i class="fas fa-id-card me-2"></i>Nomor &amp; Dokumen Identitas
                      </h6>
                      <dl class="row mb-0 small">
                        <dt class="col-sm-5 text-secondary">NIP</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedGuru.nip || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">NUPTK</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedGuru.nuptk || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Status Pegawai</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedGuru.status_kepegawaian || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Pendidikan Terakhir</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedGuru.pendidikan_terakhir || '-' }}</dd>
                      </dl>
                    </div>
                  </div>

                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <h6 class="fw-bold text-primary mb-3 small text-uppercase">
                        <i class="fas fa-user me-2"></i>Informasi Diri &amp; Kelahiran
                      </h6>
                      <dl class="row mb-0 small">
                        <dt class="col-sm-5 text-secondary">Jenis Kelamin</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedGuru.jenis_kelamin || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Agama</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedGuru.agama || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Tempat Lahir</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedGuru.tempat_lahir || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Tanggal Lahir</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ formatDateIndo(selectedGuru.tanggal_lahir) }}</dd>
                      </dl>
                    </div>
                  </div>

                  <div class="col-12">
                    <div class="border rounded-3 p-3 bg-white">
                      <h6 class="fw-bold text-primary mb-3 small text-uppercase">
                        <i class="fas fa-phone-alt me-2"></i>Kontak &amp; Domisili
                      </h6>
                      <div class="row g-2 small">
                        <div class="col-md-6">
                          <span class="text-secondary d-block">Nomor Telepon / WhatsApp:</span>
                          <span class="fw-semibold text-dark">{{ selectedGuru.no_telepon || '-' }}</span>
                        </div>
                        <div class="col-md-6">
                          <span class="text-secondary d-block">Email Resmi:</span>
                          <span class="fw-semibold text-dark">{{ selectedGuru.email }}</span>
                        </div>
                        <div class="col-12 mt-2">
                          <span class="text-secondary d-block">Alamat Lengkap:</span>
                          <span class="fw-semibold text-dark">{{ selectedGuru.alamat || '-' }}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="modal-footer bg-light">
                <button type="button" class="btn btn-outline-primary btn-sm d-flex align-items-center gap-1" @click="openEditFromDetail">
                  <i class="fas fa-edit"></i>
                  <span>Ubah Data</span>
                </button>
                <button type="button" class="btn btn-secondary btn-sm" @click="closeDetailModal">Tutup</button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Modal Konfirmasi Hapus Guru -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showDeleteModal" class="app-modal-overlay" tabindex="-1" @click.self="closeDeleteModal">
          <div class="modal-dialog modal-dialog-centered modal-sm">
            <div class="modal-content border-0 shadow-lg text-center p-4 bg-white">
              <div class="bg-danger-subtle text-danger rounded-circle p-3 d-inline-flex align-items-center justify-content-center mx-auto mb-3" style="width: 58px; height: 58px;">
                <i class="fas fa-exclamation-triangle fs-4"></i>
              </div>
              <h5 class="fw-bold text-dark mb-1">Hapus Data Guru?</h5>
              <p class="text-secondary small mb-3">
                Apakah Anda yakin ingin menghapus data guru <strong class="text-dark">{{ guruToDelete?.nama }}</strong>? Tindakan ini tidak dapat dibatalkan.
              </p>
              <div class="d-flex gap-2 justify-content-center">
                <button type="button" class="btn btn-light btn-sm px-3" @click="closeDeleteModal" :disabled="deleting">
                  Batal
                </button>
                <button type="button" class="btn btn-danger btn-sm px-3 d-flex align-items-center gap-1 shadow-sm" @click="confirmDelete" :disabled="deleting">
                  <span v-if="deleting" class="spinner-border spinner-border-sm"></span>
                  <i v-else class="fas fa-trash-alt"></i>
                  <span>{{ deleting ? 'Menghapus...' : 'Ya, Hapus' }}</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
    <!-- Modal Import Data Masal Guru -->
    <ImportModal
      :is-open="showImportModal"
      title="Tenaga Pendidik (Guru)"
      template-filename="Template_Import_Guru"
      :schema="guruImportSchema"
      :sample-rows="guruSampleRows"
      :on-import-row="handleImportGuruRow"
      @close="showImportModal = false"
      @success="onImportSuccess"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { guruApi } from '@/api/guruApi'
import { userApi } from '@/api/userApi'
import { exportToCSV } from '@/utils/exportHelper'
import ImportModal from '@/components/ImportModal.vue'

const router = useRouter()
const authStore = useAuthStore()

const gurus = ref([])
const availableUsers = ref([])
const loading = ref(false)
const submitting = ref(false)
const searchQuery = ref('')
const filterStatus = ref('')

const alertMessage = ref('')
const alertClass = ref('alert-success')
const alertIcon = ref('fas fa-check-circle')

const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref(null)

const showDetailModal = ref(false)
const selectedGuru = ref(null)

const showDeleteModal = ref(false)
const guruToDelete = ref(null)
const deleting = ref(false)

const form = ref({
  id_user: '',
  nama: '',
  nip: '',
  nuptk: '',
  jenis_kelamin: '',
  tempat_lahir: '',
  tanggal_lahir: '',
  no_telepon: '',
  agama: '',
  pendidikan_terakhir: '',
  status_kepegawaian: '',
  email: '',
  alamat: '',
})

const showAlert = (message, type = 'success') => {
  alertMessage.value = message
  if (type === 'success') {
    alertClass.value = 'alert-success'
    alertIcon.value = 'fas fa-check-circle'
  } else {
    alertClass.value = 'alert-danger'
    alertIcon.value = 'fas fa-exclamation-triangle'
  }
}

const getStatusBadgeClass = (status) => {
  switch (status) {
    case 'PNS':
      return 'bg-success-subtle text-success'
    case 'PPPK':
      return 'bg-primary-subtle text-primary'
    case 'Guru Tetap Yayasan':
      return 'bg-info-subtle text-info'
    case 'Honorer':
    case 'Guru Tidak Tetap':
      return 'bg-warning-subtle text-warning'
    default:
      return 'bg-secondary-subtle text-secondary'
  }
}

const formatDateIndo = (dateStr) => {
  if (!dateStr) return '-'
  try {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    return d.toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    })
  } catch {
    return dateStr
  }
}

const fetchGuru = async () => {
  loading.value = true
  try {
    const res = await guruApi.getAll()
    if (res.data && res.data.data) {
      gurus.value = res.data.data
    } else {
      gurus.value = []
    }
  } catch (err) {
    showAlert('Gagal mengambil daftar guru.', 'danger')
  } finally {
    loading.value = false
  }
}

const fetchUsers = async () => {
  try {
    const res = await userApi.getAll()
    if (res.data && res.data.data) {
      availableUsers.value = res.data.data
    }
  } catch (err) {
    console.error('Error fetching users for guru selection', err)
  }
}

const filteredGuru = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  return gurus.value.filter((g) => {
    const matchQuery =
      !query ||
      g.nama?.toLowerCase().includes(query) ||
      g.email?.toLowerCase().includes(query) ||
      g.nip?.toLowerCase().includes(query) ||
      g.nuptk?.toLowerCase().includes(query) ||
      g.no_telepon?.toLowerCase().includes(query) ||
      g.alamat?.toLowerCase().includes(query)

    const matchStatus = !filterStatus.value || g.status_kepegawaian === filterStatus.value
    return matchQuery && matchStatus
  })
})

const exportDataGuru = () => {
  const headers = [
    { label: 'No', key: '__index' },
    { label: 'Nama Lengkap Pendidik', key: 'nama' },
    { label: 'NIP', key: 'nip' },
    { label: 'NUPTK', key: 'nuptk' },
    { label: 'Status Kepegawaian', key: 'status_kepegawaian' },
    { label: 'Pendidikan Terakhir', key: 'pendidikan_terakhir' },
    { label: 'Jenis Kelamin', key: 'jenis_kelamin' },
    { label: 'Tempat Lahir', key: 'tempat_lahir' },
    { label: 'Tanggal Lahir', key: 'tanggal_lahir', formatter: (g) => g.tanggal_lahir ? g.tanggal_lahir.substring(0, 10) : '-' },
    { label: 'Agama', key: 'agama' },
    { label: 'No Telepon/WA', key: 'no_telepon' },
    { label: 'Email Resmi', key: 'email' },
    { label: 'Alamat', key: 'alamat' }
  ]
  exportToCSV('Data_Guru_Pendidik', headers, filteredGuru.value)
}

const showImportModal = ref(false)

const guruImportSchema = [
  { key: 'nama', label: 'Nama Lengkap', example: 'Dra. Hj. Siti Aminah, M.Pd' },
  { key: 'nip', label: 'NIP', example: '197508122000032001' },
  { key: 'nuptk', label: 'NUPTK', example: '8456753655300023' },
  { key: 'status_kepegawaian', label: 'Status Pegawai', example: 'PNS' },
  { key: 'pendidikan_terakhir', label: 'Pendidikan Terakhir', example: 'S2' },
  { key: 'jenis_kelamin', label: 'Jenis Kelamin', example: 'Perempuan' },
  { key: 'tempat_lahir', label: 'Tempat Lahir', example: 'Bandung' },
  { key: 'tanggal_lahir', label: 'Tanggal Lahir', example: '1975-08-12' },
  { key: 'agama', label: 'Agama', example: 'Islam' },
  { key: 'email', label: 'Email Resmi', example: 'siti.aminah@sekolah.id' },
  { key: 'no_telepon', label: 'No Telepon', example: '081298765432' },
  { key: 'alamat', label: 'Alamat', example: 'Jl. Diponegoro No. 88' }
]

const guruSampleRows = [
  {
    nama: 'Dra. Hj. Siti Aminah, M.Pd',
    nip: '197508122000032001',
    nuptk: '8456753655300023',
    status_kepegawaian: 'PNS',
    pendidikan_terakhir: 'S2',
    jenis_kelamin: 'Perempuan',
    tempat_lahir: 'Bandung',
    tanggal_lahir: '1975-08-12',
    agama: 'Islam',
    email: 'siti.aminah@sekolah.id',
    no_telepon: '081298765432',
    alamat: 'Jl. Diponegoro No. 88'
  },
  {
    nama: 'Budi Hartono, S.Kom',
    nip: '198803152014021003',
    nuptk: '3456762663200012',
    status_kepegawaian: 'PPPK',
    pendidikan_terakhir: 'S1',
    jenis_kelamin: 'Laki-laki',
    tempat_lahir: 'Surabaya',
    tanggal_lahir: '1988-03-15',
    agama: 'Islam',
    email: 'budi.hartono@sekolah.id',
    no_telepon: '081345678901',
    alamat: 'Jl. Pahlawan No. 17'
  }
]

const openImportModal = () => {
  showImportModal.value = true
}

const handleImportGuruRow = async (row) => {
  const payload = {
    id_user: availableUsers.value[0]?.id || '',
    nama: row.nama || row.nama_lengkap || '',
    nip: String(row.nip || ''),
    nuptk: String(row.nuptk || ''),
    status_kepegawaian: row.status_kepegawaian || 'Guru Tetap Yayasan',
    pendidikan_terakhir: row.pendidikan_terakhir || 'S1',
    jenis_kelamin: row.jenis_kelamin || 'Laki-laki',
    tempat_lahir: row.tempat_lahir || '',
    tanggal_lahir: row.tanggal_lahir ? `${row.tanggal_lahir}T00:00:00Z` : '',
    agama: row.agama || 'Islam',
    email: row.email || row.email_resmi || '',
    no_telepon: String(row.no_telepon || ''),
    alamat: row.alamat || ''
  }

  if (!payload.nama || !payload.email) {
    throw new Error('Kolom Nama Lengkap dan Email Resmi wajib diisi.')
  }

  return await guruApi.create(payload)
}

const onImportSuccess = ({ total, success, failed }) => {
  fetchGuru()
  showAlert(`Berhasil mengimpor ${success} dari total ${total} data guru. ${failed > 0 ? `(${failed} gagal)` : ''}`)
}

const openCreateModal = () => {
  isEditing.value = false
  editingId.value = null
  form.value = {
    id_user: availableUsers.value[0]?.id || '',
    nama: '',
    nip: '',
    nuptk: '',
    jenis_kelamin: '',
    tempat_lahir: '',
    tanggal_lahir: '',
    no_telepon: '',
    agama: '',
    pendidikan_terakhir: '',
    status_kepegawaian: '',
    email: '',
    alamat: '',
  }
  showModal.value = true
}

const openEditModal = (guru) => {
  isEditing.value = true
  editingId.value = guru.id
  form.value = {
    id_user: guru.id_user || '',
    nama: guru.nama || '',
    nip: guru.nip || '',
    nuptk: guru.nuptk || '',
    jenis_kelamin: guru.jenis_kelamin || '',
    tempat_lahir: guru.tempat_lahir || '',
    tanggal_lahir: guru.tanggal_lahir || '',
    no_telepon: guru.no_telepon || '',
    agama: guru.agama || '',
    pendidikan_terakhir: guru.pendidikan_terakhir || '',
    status_kepegawaian: guru.status_kepegawaian || '',
    email: guru.email || '',
    alamat: guru.alamat || '',
  }
  showModal.value = true
}

const openDetailModal = (guru) => {
  selectedGuru.value = guru
  showDetailModal.value = true
}

const closeDetailModal = () => {
  showDetailModal.value = false
  selectedGuru.value = null
}

const openEditFromDetail = () => {
  if (selectedGuru.value) {
    const target = selectedGuru.value
    closeDetailModal()
    openEditModal(target)
  }
}

const closeModal = () => {
  showModal.value = false
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (isEditing.value) {
      await guruApi.update(editingId.value, {
        nama: form.value.nama,
        nip: form.value.nip,
        nuptk: form.value.nuptk,
        jenis_kelamin: form.value.jenis_kelamin,
        tempat_lahir: form.value.tempat_lahir,
        tanggal_lahir: form.value.tanggal_lahir,
        no_telepon: form.value.no_telepon,
        agama: form.value.agama,
        pendidikan_terakhir: form.value.pendidikan_terakhir,
        status_kepegawaian: form.value.status_kepegawaian,
        email: form.value.email,
        alamat: form.value.alamat,
      })
      showAlert('Biodata guru berhasil diperbarui!')
    } else {
      await guruApi.create(form.value)
      showAlert('Guru baru berhasil ditambahkan!')
    }
    closeModal()
    fetchGuru()
  } catch (err) {
    const errMsg = err.response?.data?.message || err.response?.data || err.message || 'Gagal memproses data guru.'
    showAlert(typeof errMsg === 'string' ? errMsg : 'Gagal memproses data guru.', 'danger')
  } finally {
    submitting.value = false
  }
}

const openDeleteModal = (guru) => {
  guruToDelete.value = guru
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  guruToDelete.value = null
  showDeleteModal.value = false
  deleting.value = false
}

const confirmDelete = async () => {
  if (!guruToDelete.value) return
  deleting.value = true
  try {
    const res = await guruApi.delete(guruToDelete.value.id)
    const successMsg = res.data?.message || `Data guru "${guruToDelete.value.nama}" berhasil dihapus!`
    showAlert(successMsg, 'success')
    closeDeleteModal()
    fetchGuru()
  } catch (err) {
    const errMsg = err.response?.data?.message || err.response?.data || err.message || 'Gagal menghapus data guru.'
    showAlert(typeof errMsg === 'string' ? errMsg : 'Gagal menghapus data guru.', 'danger')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  if (!authStore.isAdmin) {
    router.replace('/admin')
    return
  }
  fetchGuru()
  fetchUsers()
})
</script>

<style scoped>
.avatar-teacher {
  width: 38px;
  height: 38px;
  background: linear-gradient(135deg, rgba(13, 110, 253, 0.15), rgba(2, 132, 199, 0.2));
  color: #0d6efd;
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
tr:hover .avatar-teacher {
  transform: scale(1.1) rotate(-6deg);
}
.avatar-detail {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, rgba(13, 110, 253, 0.15), rgba(2, 132, 199, 0.25));
}
.section-title {
  letter-spacing: 0.5px;
}
.app-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(4px);
  z-index: 1055;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}
.app-modal-overlay .modal-content {
  background-color: #ffffff !important;
}
</style>
