<template>
  <div class="animate-fade-up">
    <!-- Header Halaman -->
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
      <div>
        <div class="d-flex align-items-center gap-2 mb-1">
          <h3 class="fw-bold text-dark mb-0">Riwayat Aktivitas Pengguna</h3>
          <span class="badge bg-primary-subtle text-primary px-2 py-1 small">20 Data per Halaman</span>
        </div>
        <p class="text-secondary small mb-0">
          Catatan kegiatan seluruh pengguna di sistem sekolah, dilengkapi navigasi halaman (20 per halaman) dan auto-refresh tiap 30 detik.
        </p>
      </div>
      <div class="d-flex align-items-center gap-2">
        <button
          @click="toggleAutoRefresh"
          :class="autoRefresh ? 'btn btn-sm btn-success-subtle text-success border border-success pulse-beacon' : 'btn btn-sm btn-outline-secondary'"
          class="d-flex align-items-center gap-1 shadow-sm btn-living"
          title="Nyalakan/matikan pembaruan otomatis tiap 30 detik"
        >
          <span class="spinner-grow spinner-grow-sm text-success" v-if="autoRefresh" style="width: 8px; height: 8px;"></span>
          <i class="fas fa-satellite-dish" v-else></i>
          <span>{{ autoRefresh ? 'Pembaruan Otomatis (30s)' : 'Otomatis Nonaktif' }}</span>
        </button>
        <button @click="fetchLogs(false)" class="btn btn-outline-secondary btn-living d-flex align-items-center gap-2" :disabled="loading">
          <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
          <span>Segarkan Data</span>
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

    <!-- Filter & Pencarian Ramah Pengguna -->
    <div class="card card-custom p-3 mb-4">
      <div class="row g-2 align-items-center">
        <!-- Pencarian -->
        <div class="col-lg-4 col-md-12">
          <div class="input-group">
            <span class="input-group-text bg-light text-muted border-end-0">
              <i class="fas fa-search"></i>
            </span>
            <input
              v-model="searchQuery"
              @keyup.enter="applyFilter"
              type="text"
              class="form-control border-start-0 ps-0"
              placeholder="Cari nama pengguna, keterangan kegiatan..."
            />
            <button v-if="searchQuery" class="btn btn-link text-secondary p-0 px-2" @click="resetFilter" title="Bersihkan">
              <i class="fas fa-times-circle"></i>
            </button>
          </div>
        </div>

        <!-- Filter Bagian / Modul -->
        <div class="col-lg-2 col-md-4 col-sm-6">
          <select v-model="moduleFilter" @change="applyFilter" class="form-select">
            <option value="">Semua Bagian</option>
            <option value="Siswa">Siswa</option>
            <option value="Guru">Guru</option>
            <option value="Kelas">Kelas</option>
            <option value="Mata Pelajaran">Mata Pelajaran</option>
            <option value="Pengguna">Akun Pengguna</option>
            <option value="Autentikasi">Masuk Sistem</option>
            <option value="Log Sistem">Riwayat Log</option>
          </select>
        </div>

        <!-- Filter Jenis Aktivitas -->
        <div class="col-lg-2 col-md-4 col-sm-6">
          <select v-model="activityTypeFilter" @change="applyFilter" class="form-select">
            <option value="">Semua Aktivitas</option>
            <option value="LOGIN">Masuk Sistem</option>
            <option value="CREATE">Tambah Data</option>
            <option value="UPDATE">Ubah Data</option>
            <option value="DELETE">Hapus Data</option>
            <option value="READ">Melihat Data</option>
          </select>
        </div>

        <!-- Filter Status Hasil -->
        <div class="col-lg-2 col-md-4 col-sm-6">
          <select v-model="resultFilter" @change="applyFilter" class="form-select">
            <option value="">Semua Status</option>
            <option value="Success">Berhasil</option>
            <option value="Failed">Gagal</option>
          </select>
        </div>

        <!-- Tombol Filter & Reset -->
        <div class="col-lg-2 col-md-12 col-sm-6 d-flex gap-2">
          <button @click="applyFilter" class="btn btn-primary flex-grow-1 d-flex align-items-center justify-content-center gap-1">
            <i class="fas fa-filter"></i>
            <span>Saring</span>
          </button>
          <button @click="resetFilter" class="btn btn-outline-secondary" title="Reset Filter">
            <i class="fas fa-undo"></i>
          </button>
        </div>
      </div>
    </div>

    <!-- Tabel Riwayat Aktivitas yang Mudah Dibaca -->
    <div class="card card-custom overflow-hidden">
      <div class="table-responsive">
        <table class="table table-custom table-hover mb-0">
          <thead>
            <tr>
              <th style="width: 150px;">Waktu</th>
              <th style="width: 200px;">Pengguna</th>
              <th>Aktivitas yang Dilakukan</th>
              <th style="width: 130px;">Bagian</th>
              <th style="width: 120px;" class="text-center">Status</th>
              <th style="width: 100px;" class="text-end">Rincian</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="6" class="text-center py-5 text-secondary">
                <div class="spinner-border spinner-border-sm text-primary me-2" role="status"></div>
                Memuat riwayat aktivitas...
              </td>
            </tr>
            <tr v-else-if="logsList.length === 0">
              <td colspan="6" class="text-center py-5 text-muted">
                <i class="fas fa-clipboard-list fs-2 mb-2 d-block text-secondary"></i>
                Tidak ada data aktivitas yang sesuai dengan kriteria pencarian.
              </td>
            </tr>
            <tr v-for="log in logsList" :key="log.id">
              <!-- Waktu -->
              <td class="text-nowrap small text-secondary">
                <div class="fw-bold text-dark">{{ formatDateFriendly(log.timestamp) }}</div>
                <div class="text-muted" style="font-size: 0.78rem;">
                  <i class="far fa-clock me-1"></i>{{ formatTime(log.timestamp) }} WIB
                </div>
              </td>

              <!-- Pengguna / Pelaku -->
              <td>
                <div class="d-flex align-items-center gap-2">
                  <div class="avatar-circle" :class="getAvatarClass(log.role)">
                    <i :class="getAvatarIcon(log.role)"></i>
                  </div>
                  <div class="overflow-hidden">
                    <div class="fw-bold text-dark text-truncate" style="max-width: 150px;">
                      {{ log.username || 'Pengguna' }}
                    </div>
                    <div class="d-flex align-items-center gap-1 mt-1">
                      <span class="badge" :class="getRoleBadgeClass(log.role)">
                        {{ formatRole(log.role) }}
                      </span>
                    </div>
                  </div>
                </div>
              </td>

              <!-- Aktivitas yang Dilakukan -->
              <td>
                <div class="d-flex align-items-center gap-2 mb-1 flex-wrap">
                  <span class="badge" :class="getActivityBadgeClass(log.activity_type)">
                    <i :class="getActivityIcon(log.activity_type)" class="me-1"></i>
                    {{ getActivityLabel(log.activity_type) }}
                  </span>
                </div>
                <div class="fw-semibold text-dark small">
                  {{ cleanActivityDescription(log.activity_description) }}
                </div>
              </td>

              <!-- Bagian / Modul -->
              <td>
                <span class="badge bg-light text-dark border px-2 py-1 small">
                  {{ log.module || 'Sistem' }}
                </span>
              </td>

              <!-- Status -->
              <td class="text-center">
                <span
                  class="badge px-2 py-1"
                  :class="{
                    'bg-success-subtle text-success': log.result === 'Success',
                    'bg-danger-subtle text-danger': log.result === 'Failed',
                  }"
                >
                  <i :class="log.result === 'Success' ? 'fas fa-check-circle me-1' : 'fas fa-times-circle me-1'"></i>
                  {{ log.result === 'Success' ? 'Berhasil' : 'Gagal' }}
                </span>
              </td>

              <!-- Tombol Rincian -->
              <td class="text-end">
                <button @click="openDetailModal(log)" class="btn btn-sm btn-outline-primary" title="Lihat Rincian Lengkap">
                  <i class="fas fa-eye me-1"></i> Rincian
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Footer Pagination: 20 Data per Halaman -->
      <div class="card-footer bg-white d-flex flex-column flex-md-row justify-content-between align-items-center gap-3 py-3 border-top">
        <div class="small text-muted d-flex align-items-center gap-2">
          <i class="fas fa-list-ol text-primary"></i>
          <span>
            Menampilkan <strong class="text-dark">{{ fromNumber }} - {{ toNumber }}</strong> dari <strong class="text-dark">{{ totalItems }}</strong> aktivitas
          </span>
        </div>

        <!-- Tombol Pagination -->
        <nav v-if="totalPages > 1" aria-label="Navigasi Halaman Log">
          <ul class="pagination pagination-sm mb-0 flex-wrap justify-content-center">
            <li class="page-item" :class="{ disabled: currentPage <= 1 || loading }">
              <button class="page-link" @click="changePage(currentPage - 1)" :disabled="currentPage <= 1 || loading">
                <i class="fas fa-chevron-left me-1"></i> Sebelumnya
              </button>
            </li>
            <li
              v-for="p in visiblePages"
              :key="p"
              class="page-item"
              :class="{ active: p === currentPage, disabled: p === '...' || loading }"
            >
              <button v-if="p !== '...'" class="page-link" @click="changePage(p)" :disabled="loading">{{ p }}</button>
              <span v-else class="page-link bg-light text-muted">...</span>
            </li>
            <li class="page-item" :class="{ disabled: currentPage >= totalPages || loading }">
              <button class="page-link" @click="changePage(currentPage + 1)" :disabled="currentPage >= totalPages || loading">
                Berikutnya <i class="fas fa-chevron-right ms-1"></i>
              </button>
            </li>
          </ul>
        </nav>

        <div class="small text-secondary d-none d-md-flex align-items-center gap-2">
          <span v-if="autoRefresh" class="badge bg-success-subtle text-success border border-success-subtle d-inline-flex align-items-center gap-1 py-1 px-2">
            <i class="fas fa-sync-alt fa-spin" v-if="loading"></i>
            <i class="far fa-clock" v-else></i>
            Pembaruan tiap 30s
          </span>
          <span v-else class="badge bg-secondary-subtle text-secondary border py-1 px-2">
            Pembaruan nonaktif
          </span>
        </div>
      </div>
    </div>

    <!-- Modal Rincian Aktivitas (Ramah Admin Sekolah) -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div
          class="app-modal-overlay"
          v-if="showDetailModal"
          tabindex="-1"
          @click.self="showDetailModal = false"
        >
          <div class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable">
            <div class="modal-content border-0 shadow">
              <div class="modal-header bg-light">
                <h5 class="modal-title fw-bold d-flex align-items-center gap-2">
                  <i class="fas fa-clipboard-check text-primary"></i>
                  Rincian Aktivitas Pengguna
                </h5>
                <button type="button" class="btn-close" @click="showDetailModal = false"></button>
              </div>
              
              <div class="modal-body p-4" v-if="selectedLog">
                <!-- 1. Ringkasan Aktivitas Utama -->
                <div class="p-3 mb-4 rounded-3 border" :class="selectedLog.result === 'Success' ? 'bg-success-subtle border-success-subtle' : 'bg-danger-subtle border-danger-subtle'">
                  <div class="d-flex align-items-center gap-2 mb-2">
                    <span class="badge px-2 py-1" :class="getActivityBadgeClass(selectedLog.activity_type)">
                      <i :class="getActivityIcon(selectedLog.activity_type)" class="me-1"></i>
                      {{ getActivityLabel(selectedLog.activity_type) }}
                    </span>
                    <span class="badge bg-white text-dark border">
                      Bagian: {{ selectedLog.module || 'Sistem' }}
                    </span>
                    <span
                      class="badge ms-auto px-2 py-1"
                      :class="selectedLog.result === 'Success' ? 'bg-success text-white' : 'bg-danger text-white'"
                    >
                      <i :class="selectedLog.result === 'Success' ? 'fas fa-check-circle me-1' : 'fas fa-times-circle me-1'"></i>
                      {{ selectedLog.result === 'Success' ? 'Berhasil Dilakukan' : 'Gagal Dilakukan' }}
                    </span>
                  </div>
                  <h5 class="fw-bold text-dark mb-1">
                    {{ cleanActivityDescription(selectedLog.activity_description) }}
                  </h5>
                  <div class="text-secondary small">
                    Waktu Pelaksanaan: <span class="fw-semibold text-dark">{{ fullFormatDate(selectedLog.timestamp) }} WIB</span>
                  </div>
                </div>

                <!-- 2. Informasi Pelaku (Siapa) -->
                <div class="card border mb-3">
                  <div class="card-header bg-white py-2">
                    <div class="fw-bold small text-uppercase text-secondary">
                      <i class="fas fa-user-circle me-1"></i> Informasi Pelaku
                    </div>
                  </div>
                  <div class="card-body p-3">
                    <div class="row g-3">
                      <div class="col-sm-4">
                        <div class="text-secondary small">Nama Pengguna</div>
                        <div class="fw-bold text-dark">{{ selectedLog.username || 'Pengguna / Anonim' }}</div>
                      </div>
                      <div class="col-sm-4">
                        <div class="text-secondary small">Alamat Email</div>
                        <div class="fw-semibold text-dark">{{ selectedLog.email || '-' }}</div>
                      </div>
                      <div class="col-sm-4">
                        <div class="text-secondary small">Peran di Sekolah</div>
                        <div>
                          <span class="badge" :class="getRoleBadgeClass(selectedLog.role)">
                            {{ formatRole(selectedLog.role) }}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- 3. Rincian Data yang Diproses (Human-friendly) -->
                <div class="card border mb-4">
                  <div class="card-header bg-white py-2 d-flex justify-content-between align-items-center">
                    <div class="fw-bold small text-uppercase text-secondary">
                      <i class="fas fa-info-circle me-1"></i> Rincian Data yang Diproses
                    </div>
                    <span class="badge bg-light text-muted border">Bagian {{ selectedLog.module || 'Sistem' }}</span>
                  </div>
                  <div class="card-body p-3">
                    <!-- Tampilkan data dalam tabel atribut jika ada payload -->
                    <div v-if="parsedDataFields(selectedLog.request_body).length > 0" class="table-responsive">
                      <table class="table table-sm table-bordered mb-0">
                        <thead class="table-light">
                          <tr>
                            <th style="width: 35%;">Informasi / Kolom</th>
                            <th>Nilai Data</th>
                          </tr>
                        </thead>
                        <tbody>
                          <tr v-for="field in parsedDataFields(selectedLog.request_body)" :key="field.label">
                            <td class="text-secondary fw-semibold">{{ field.label }}</td>
                            <td class="fw-bold text-dark">{{ field.value }}</td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                    <div v-else class="text-muted small py-2">
                      <i class="fas fa-info-circle me-1 text-primary"></i>
                      Aktivitas ini berupa pemeriksaan atau pembacaan data, tidak ada formulir data baru yang dikirimkan.
                    </div>
                  </div>
                </div>

                <!-- 4. Bagian Lipat: Informasi Teknis untuk IT (Default Tertutup) -->
                <div class="card border">
                  <div class="card-header bg-light py-2 cursor-pointer d-flex justify-content-between align-items-center" @click="showTechnicalDetails = !showTechnicalDetails" style="cursor: pointer;">
                    <div class="fw-bold small text-secondary d-flex align-items-center gap-2">
                      <i class="fas fa-terminal"></i>
                      <span>Informasi Teknis & Diagnostik Sistem (Khusus Tim IT)</span>
                    </div>
                    <i class="fas" :class="showTechnicalDetails ? 'fa-chevron-up' : 'fa-chevron-down'"></i>
                  </div>
                  
                  <div class="card-body p-3" v-if="showTechnicalDetails">
                    <div class="row g-2 mb-3 small">
                      <div class="col-sm-6">
                        <span class="text-muted">Endpoint:</span>
                        <span class="font-monospace fw-bold ms-1">{{ selectedLog.http_method }} {{ selectedLog.endpoint || selectedLog.service_name }}</span>
                      </div>
                      <div class="col-sm-6">
                        <span class="text-muted">Status Code:</span>
                        <span class="font-monospace fw-bold ms-1">HTTP {{ selectedLog.status_code || 200 }}</span>
                      </div>
                      <div class="col-sm-6">
                        <span class="text-muted">Alamat IP:</span>
                        <span class="font-monospace ms-1">{{ selectedLog.ip_address || '127.0.0.1' }}</span>
                      </div>
                      <div class="col-sm-6">
                        <span class="text-muted">Waktu Proses:</span>
                        <span class="font-monospace ms-1">{{ selectedLog.duration_ms || 0 }} ms</span>
                      </div>
                      <div class="col-12 mt-2">
                        <span class="text-muted">Perangkat / Browser:</span>
                        <div class="font-monospace small text-break text-muted mt-1 p-2 bg-light rounded">
                          {{ selectedLog.perangkat || '-' }}
                        </div>
                      </div>
                    </div>

                    <!-- Raw JSON payloads -->
                    <div class="mb-2">
                      <div class="d-flex justify-content-between align-items-center mb-1">
                        <span class="fw-semibold small text-secondary">Request Payload (JSON)</span>
                        <button class="btn btn-sm btn-link text-decoration-none p-0 d-flex align-items-center gap-1" @click="copyToClipboard(formatJson(selectedLog.request_body), 'req')">
                          <i :class="copiedKey === 'req' ? 'fas fa-check text-success' : 'fas fa-copy'"></i>
                          <span :class="{ 'text-success fw-bold': copiedKey === 'req' }">{{ copiedKey === 'req' ? 'Tersalin!' : 'Salin JSON' }}</span>
                        </button>
                      </div>
                      <pre class="code-box">{{ formatJson(selectedLog.request_body) }}</pre>
                    </div>

                    <div class="mb-2">
                      <div class="d-flex justify-content-between align-items-center mb-1">
                        <span class="fw-semibold small text-secondary">Response Payload (JSON)</span>
                        <button class="btn btn-sm btn-link text-decoration-none p-0 d-flex align-items-center gap-1" @click="copyToClipboard(formatJson(selectedLog.response_body), 'res')">
                          <i :class="copiedKey === 'res' ? 'fas fa-check text-success' : 'fas fa-copy'"></i>
                          <span :class="{ 'text-success fw-bold': copiedKey === 'res' }">{{ copiedKey === 'res' ? 'Tersalin!' : 'Salin JSON' }}</span>
                        </button>
                      </div>
                      <pre class="code-box">{{ formatJson(selectedLog.response_body) }}</pre>
                    </div>
                  </div>
                </div>

              </div>
              <div class="modal-footer bg-light">
                <button type="button" class="btn btn-primary px-4" @click="showDetailModal = false">Tutup</button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { logApi } from '@/api/logApi'

const authStore = useAuthStore()
const router = useRouter()

const logsList = ref([])
const totalItems = ref(0)
const loading = ref(false)
const currentPage = ref(1)
const limit = ref(20)

// Realtime auto-refresh interval (30 detik)
const autoRefresh = ref(true)
let refreshTimer = null

// Filters
const searchQuery = ref('')
const moduleFilter = ref('')
const activityTypeFilter = ref('')
const resultFilter = ref('')

const selectedLog = ref(null)
const showDetailModal = ref(false)
const showTechnicalDetails = ref(false)
const copiedKey = ref('')

const alertMessage = ref('')
const alertClass = ref('alert-success')
const alertIcon = ref('fas fa-check-circle')

// Kalkulasi Pagination
const totalPages = computed(() => Math.ceil(totalItems.value / limit.value) || 1)
const fromNumber = computed(() => (totalItems.value === 0 ? 0 : (currentPage.value - 1) * limit.value + 1))
const toNumber = computed(() => Math.min(currentPage.value * limit.value, totalItems.value))

const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i)
  } else {
    pages.push(1)
    if (current > 3) pages.push('...')
    const start = Math.max(2, current - 1)
    const end = Math.min(total - 1, current + 1)
    for (let i = start; i <= end; i++) pages.push(i)
    if (current < total - 2) pages.push('...')
    pages.push(total)
  }
  return pages
})

const fetchLogs = async (isBackground = false) => {
  // Cegah request ganda atau tumpang tindih
  if (loading.value) return

  // Jika background auto-refresh dan admin sedang membaca halaman > 1, lewati refresh agar tidak memotong bacaan
  if (isBackground && currentPage.value > 1) return

  loading.value = true
  if (!isBackground) {
    alertMessage.value = ''
  }

  try {
    const params = {
      page: currentPage.value,
      limit: limit.value,
    }

    // Jika background refresh, gunakan no_count agar sangat cepat
    if (isBackground) {
      params.no_count = 'true'
    }

    if (searchQuery.value.trim()) {
      params.search = searchQuery.value.trim()
    }
    if (moduleFilter.value) {
      params.module = moduleFilter.value
    }
    if (activityTypeFilter.value) {
      params.activity_type = activityTypeFilter.value
    }
    if (resultFilter.value) {
      params.result = resultFilter.value
    }

    const response = await logApi.getAll(params)
    const resData = response.data
    if (resData.success && resData.data) {
      logsList.value = resData.data.data || []
      // Perbarui total data jika request eksplisit (bukan background), atau jika belum pernah terisi
      if (!isBackground || totalItems.value === 0) {
        totalItems.value = resData.data.total || 0
      }
    }
  } catch (err) {
    if (!isBackground) {
      showAlert(err.response?.data?.message || err.message || 'Gagal mengambil data aktivitas sistem', 'danger')
    }
  } finally {
    loading.value = false
  }
}

const changePage = (page) => {
  if (page < 1 || page > totalPages.value || page === currentPage.value || loading.value) return
  currentPage.value = page
  fetchLogs(false)
}

const applyFilter = () => {
  currentPage.value = 1
  fetchLogs(false)
}

const resetFilter = () => {
  searchQuery.value = ''
  moduleFilter.value = ''
  activityTypeFilter.value = ''
  resultFilter.value = ''
  currentPage.value = 1
  fetchLogs(false)
}

const openDetailModal = (log) => {
  selectedLog.value = log
  showTechnicalDetails.value = false
  showDetailModal.value = true
}

const formatDateFriendly = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'

  // Konversi perbandingan tanggal ke zona waktu Asia/Jakarta (WIB)
  const dWIB = new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Jakarta' }).format(d)
  const todayWIB = new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Jakarta' }).format(new Date())

  const yesterdayDate = new Date()
  yesterdayDate.setDate(yesterdayDate.getDate() - 1)
  const yesterdayWIB = new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Jakarta' }).format(yesterdayDate)

  if (dWIB === todayWIB) {
    return 'Hari ini'
  }
  if (dWIB === yesterdayWIB) {
    return 'Kemarin'
  }
  return d.toLocaleDateString('id-ID', {
    timeZone: 'Asia/Jakarta',
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

const formatTime = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleTimeString('id-ID', {
    timeZone: 'Asia/Jakarta',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  })
}

const fullFormatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleString('id-ID', {
    timeZone: 'Asia/Jakarta',
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  })
}

const formatRole = (role) => {
  switch (role?.toLowerCase()) {
    case 'admin':
      return 'Admin Sekolah'
    case 'guru':
      return 'Guru'
    default:
      return 'Pengguna'
  }
}

const getActivityLabel = (type) => {
  switch (type) {
    case 'CREATE':
      return 'Tambah Data'
    case 'UPDATE':
      return 'Ubah Data'
    case 'DELETE':
      return 'Hapus Data'
    case 'READ':
      return 'Melihat Data'
    case 'LOGIN':
      return 'Masuk Sistem'
    default:
      return 'Aktivitas'
  }
}

const getActivityIcon = (type) => {
  switch (type) {
    case 'CREATE':
      return 'fas fa-plus-circle'
    case 'UPDATE':
      return 'fas fa-edit'
    case 'DELETE':
      return 'fas fa-trash-alt'
    case 'READ':
      return 'fas fa-eye'
    case 'LOGIN':
      return 'fas fa-sign-in-alt'
    default:
      return 'fas fa-info-circle'
  }
}

const getActivityBadgeClass = (type) => {
  switch (type) {
    case 'CREATE':
      return 'bg-success-subtle text-success'
    case 'UPDATE':
      return 'bg-warning-subtle text-warning-emphasis'
    case 'DELETE':
      return 'bg-danger-subtle text-danger'
    case 'READ':
      return 'bg-info-subtle text-info-emphasis'
    case 'LOGIN':
      return 'bg-primary-subtle text-primary'
    default:
      return 'bg-secondary-subtle text-secondary'
  }
}

const getRoleBadgeClass = (role) => {
  switch (role?.toLowerCase()) {
    case 'admin':
      return 'bg-danger-subtle text-danger'
    case 'guru':
      return 'bg-primary-subtle text-primary'
    default:
      return 'bg-secondary-subtle text-secondary'
  }
}

const getAvatarClass = (role) => {
  switch (role?.toLowerCase()) {
    case 'admin':
      return 'bg-danger text-white'
    case 'guru':
      return 'bg-primary text-white'
    default:
      return 'bg-secondary text-white'
  }
}

const getAvatarIcon = (role) => {
  switch (role?.toLowerCase()) {
    case 'admin':
      return 'fas fa-user-shield'
    case 'guru':
      return 'fas fa-chalkboard-teacher'
    default:
      return 'fas fa-user'
  }
}

// Membersihkan deskripsi agar tidak ada kurung siku teknis seperti [admin/admin]
const cleanActivityDescription = (desc) => {
  if (!desc) return 'Aktivitas sistem tercatat'
  return desc.replace(/^\[[^\]]+\]\s*/, '')
}

// Menampilkan data request body dalam tabel bahasa Indonesia ramah pengguna
const parsedDataFields = (raw) => {
  if (!raw) return []
  let obj = raw
  if (typeof raw === 'string') {
    try {
      obj = JSON.parse(raw)
    } catch {
      return []
    }
  }
  if (!obj || typeof obj !== 'object' || Array.isArray(obj)) return []

  const labelMapping = {
    nama: 'Nama Lengkap',
    username: 'Username Akun',
    email: 'Alamat Email',
    role: 'Peran Pengguna',
    kelas: 'Nama Kelas',
    nama_pelajaran: 'Mata Pelajaran',
    id_guru: 'ID Guru Pengajar',
    guru_id: 'ID Guru',
    kelas_id: 'ID Kelas',
    nip: 'Nomor Induk Pegawai (NIP)',
    nis: 'Nomor Induk Siswa (NIS)',
    alamat: 'Alamat Tempat Tinggal',
    telepon: 'No. Telepon / WhatsApp',
    password: 'Kata Sandi',
  }

  const fields = []
  for (const [key, val] of Object.entries(obj)) {
    if (val !== undefined && val !== null && val !== '') {
      const label = labelMapping[key] || key.charAt(0).toUpperCase() + key.slice(1).replace(/_/g, ' ')
      let displayValue = String(val)
      if (key === 'password') {
        displayValue = '•••••••• (Dirahasiakan)'
      }
      fields.push({ label, value: displayValue })
    }
  }
  return fields
}

const formatJson = (data) => {
  if (!data) return '{}'
  if (typeof data === 'object') {
    return JSON.stringify(data, null, 2)
  }
  try {
    const parsed = JSON.parse(data)
    return JSON.stringify(parsed, null, 2)
  } catch {
    return data
  }
}

const copyToClipboard = (text, key = '') => {
  navigator.clipboard.writeText(text)
  if (key) {
    copiedKey.value = key
    setTimeout(() => {
      if (copiedKey.value === key) copiedKey.value = ''
    }, 2000)
  }
  showAlert('Informasi teknis berhasil disalin ke clipboard!', 'success')
}

const showAlert = (msg, type = 'success') => {
  alertMessage.value = msg
  alertClass.value = type === 'success' ? 'alert-success text-success' : 'alert-danger text-danger'
  alertIcon.value = type === 'success' ? 'fas fa-check-circle' : 'fas fa-exclamation-triangle'
}

const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value
  if (autoRefresh.value) {
    startAutoRefresh()
    fetchLogs(false)
  } else {
    stopAutoRefresh()
  }
}

const startAutoRefresh = () => {
  stopAutoRefresh()
  refreshTimer = setInterval(() => {
    // Hanya polling jika modal detail tidak sedang dibuka dan tidak sedang loading
    if (!showDetailModal.value && !loading.value) {
      fetchLogs(true)
    }
  }, 30000) // Diperbarui per 30 detik agar beban database sangat ringan
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(() => {
  if (!authStore.isAdmin) {
    router.replace('/admin')
    return
  }
  fetchLogs()
  if (autoRefresh.value) {
    startAutoRefresh()
  }
})

onBeforeUnmount(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.card-custom {
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.table-custom th {
  background-color: #f8fafc;
  color: #475569;
  font-weight: 600;
  font-size: 0.85rem;
  border-bottom: 2px solid #e2e8f0;
  padding: 12px 14px;
}

.table-custom td {
  padding: 12px 14px;
  vertical-align: middle;
  border-bottom: 1px solid #f1f5f9;
}

.avatar-circle {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  flex-shrink: 0;
}

.code-box {
  background-color: #1e293b;
  color: #f8fafc;
  padding: 12px;
  border-radius: 8px;
  max-height: 200px;
  overflow: auto;
  font-size: 0.78rem;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}
</style>
