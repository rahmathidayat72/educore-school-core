<template>
  <div class="animate-fade-up">
    <!-- Header -->
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
      <div>
        <div class="d-flex align-items-center gap-2 mb-1">
          <h3 class="fw-bold text-dark mb-0">Mata Pelajaran</h3>
          <span class="badge bg-warning-subtle text-warning px-2 py-1 small">
            {{ filteredMapel.length }} Mapel
          </span>
        </div>
        <p class="text-secondary small mb-0">Kelola kurikulum, beban jam pelajaran (JP), standar KKM/KKTP, guru pengampu, dan alokasi kelas.</p>
      </div>
      <button @click="openCreateModal" class="btn btn-warning btn-shine text-dark fw-semibold d-flex align-items-center gap-2 shadow-sm">
        <i class="fas fa-plus"></i>
        <span>Tambah Mapel Baru</span>
      </button>
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
        <div class="col-md-4">
          <div class="input-group">
            <span class="input-group-text bg-white text-muted border-end-0"><i class="fas fa-search"></i></span>
            <input
              v-model="searchQuery"
              type="text"
              class="form-control border-start-0 ps-0"
              placeholder="Cari mapel, kode, guru, kelas..."
            />
            <button v-if="searchQuery" class="btn btn-link text-secondary p-0 px-2" @click="searchQuery = ''" title="Bersihkan">
              <i class="fas fa-times-circle"></i>
            </button>
          </div>
        </div>
        <div class="col-md-3">
          <select v-model="filterKelompok" class="form-select">
            <option value="">Semua Kelompok</option>
            <option value="Umum">Umum (Wajib)</option>
            <option value="Peminatan">Peminatan / Kejuruan</option>
            <option value="Muatan Lokal">Muatan Lokal</option>
          </select>
        </div>
        <div class="col-md-3">
          <select v-model="filterStatus" class="form-select">
            <option value="">Semua Status Mapel</option>
            <option value="Aktif">Aktif</option>
            <option value="Nonaktif">Nonaktif</option>
          </select>
        </div>
        <div class="col-md-2 text-md-end">
          <button @click="fetchMapel" class="btn btn-outline-secondary w-100 d-flex align-items-center justify-content-center gap-2" :disabled="loading">
            <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
            <span>Segarkan</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Tabel Data Mapel -->
    <div class="card card-custom overflow-hidden border-0 shadow-sm">
      <div class="table-responsive">
        <table class="table table-custom table-hover mb-0">
          <thead>
            <tr>
              <th style="width: 50px;">No</th>
              <th>Mata Pelajaran / Kode</th>
              <th>Kelompok &amp; Kurikulum</th>
              <th>Beban Belajar &amp; KKM</th>
              <th>Guru Pengampu</th>
              <th>Alokasi Kelas</th>
              <th>Status</th>
              <th class="text-end" style="width: 140px;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="8" class="text-center py-5 text-secondary">
                <div class="spinner-border spinner-border-sm text-warning me-2" role="status"></div>
                Memuat mata pelajaran...
              </td>
            </tr>
            <tr v-else-if="filteredMapel.length === 0">
              <td colspan="8" class="text-center py-5 text-muted">
                <i class="fas fa-book-open fs-2 mb-2 d-block text-secondary"></i>
                Belum ada data mata pelajaran yang sesuai.
              </td>
            </tr>
            <tr v-for="(m, index) in filteredMapel" :key="m.id">
              <td class="text-secondary fw-semibold">{{ index + 1 }}</td>
              <td>
                <div class="d-flex align-items-center gap-2">
                  <div class="avatar-mapel rounded-circle d-flex align-items-center justify-content-center flex-shrink-0">
                    <i class="fas fa-book-open small"></i>
                  </div>
                  <div>
                    <span class="fw-bold text-dark fs-6 d-block">{{ m.mata_pelajaran }}</span>
                    <span v-if="m.kode_mapel" class="badge bg-light text-secondary border font-monospace small py-0">
                      {{ m.kode_mapel }}
                    </span>
                  </div>
                </div>
              </td>
              <td>
                <div class="small">
                  <span class="badge bg-warning-subtle text-dark border-0 px-2 py-0 mb-1 d-inline-block">
                    {{ m.kelompok || 'Umum' }}
                  </span>
                  <div class="text-muted">{{ m.kurikulum || 'Kurikulum Merdeka' }}</div>
                </div>
              </td>
              <td>
                <div class="small">
                  <div class="fw-semibold text-dark">
                    <i class="fas fa-clock text-warning me-1"></i>{{ m.jam_per_minggu || 2 }} JP / Minggu
                  </div>
                  <div class="text-muted">
                    KKM: <span class="fw-bold text-dark">{{ m.kkm ? Number(m.kkm).toFixed(0) : '75' }}</span>
                  </div>
                </div>
              </td>
              <td>
                <span v-if="m.guru" class="badge bg-light text-secondary border px-2 py-1">
                  <i class="fas fa-chalkboard-teacher me-1 text-primary"></i>
                  {{ m.guru }}
                </span>
                <span v-else class="text-muted small fst-italic">Belum ditentukan</span>
              </td>
              <td>
                <span v-if="m.nama_kelas" class="badge bg-secondary-subtle text-secondary px-2 py-1">
                  <i class="fas fa-school me-1"></i>{{ m.nama_kelas }}
                </span>
                <span v-else class="text-muted small fst-italic">Semua Kelas</span>
              </td>
              <td>
                <span class="badge" :class="m.status === 'Aktif' || !m.status ? 'bg-success-subtle text-success' : 'bg-secondary-subtle text-secondary'">
                  {{ m.status || 'Aktif' }}
                </span>
              </td>
              <td class="text-end">
                <div class="btn-group btn-group-sm">
                  <button @click="openDetailModal(m)" class="btn btn-outline-warning text-dark" title="Lihat Silabus / Detail Mapel">
                    <i class="fas fa-eye"></i>
                  </button>
                  <button @click="openEditModal(m)" class="btn btn-outline-primary" title="Ubah Data Mapel">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button @click="openDeleteModal(m)" class="btn btn-outline-danger" title="Hapus Data Mapel">
                    <i class="fas fa-trash-alt"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form (Create / Edit) -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showModal" class="app-modal-overlay" tabindex="-1" @click.self="closeModal">
          <div class="modal-dialog modal-dialog-centered modal-lg">
            <div class="modal-content border-0 shadow-lg bg-white">
              <div class="modal-header bg-warning text-dark">
                <h5 class="modal-title fw-bold d-flex align-items-center gap-2">
                  <i :class="isEditing ? 'fas fa-edit' : 'fas fa-plus'"></i>
                  <span>{{ isEditing ? 'Ubah Mata Pelajaran' : 'Tambah Mata Pelajaran Baru' }}</span>
                </h5>
                <button type="button" class="btn-close" @click="closeModal"></button>
              </div>
              <form @submit.prevent="handleSubmit">
                <div class="modal-body p-4">
                  <!-- Section 1: Identitas Mata Pelajaran -->
                  <div class="mb-4">
                    <h6 class="fw-bold text-warning border-bottom pb-2 mb-3 text-dark">
                      <i class="fas fa-book me-2"></i>1. Identitas Mata Pelajaran
                    </h6>
                    <div class="row g-3">
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Nama Mata Pelajaran <span class="text-danger">*</span></label>
                        <input
                          v-model="form.mata_pelajaran"
                          type="text"
                          class="form-control"
                          required
                          placeholder="Contoh: Matematika Wajib, Fisika, Biologi"
                        />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Kode Mata Pelajaran (Opsional)</label>
                        <input
                          v-model="form.kode_mapel"
                          type="text"
                          class="form-control font-monospace"
                          placeholder="Contoh: MP-MAT-01"
                        />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Kelompok Mapel</label>
                        <select v-model="form.kelompok" class="form-select">
                          <option value="Umum">Kelompok Umum (Wajib)</option>
                          <option value="Peminatan">Kelompok Peminatan / Kejuruan</option>
                          <option value="Muatan Lokal">Kelompok Muatan Lokal (Mulok)</option>
                        </select>
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Kurikulum Acuan</label>
                        <select v-model="form.kurikulum" class="form-select">
                          <option value="Kurikulum Merdeka">Kurikulum Merdeka</option>
                          <option value="Kurikulum 2013">Kurikulum 2013 (K13)</option>
                          <option value="Kurikulum Cambridge">Kurikulum Cambridge / Internasional</option>
                        </select>
                      </div>
                    </div>
                  </div>

                  <!-- Section 2: Beban Belajar & Penilaian -->
                  <div class="mb-4">
                    <h6 class="fw-bold text-dark border-bottom pb-2 mb-3">
                      <i class="fas fa-clock me-2 text-warning"></i>2. Beban Jam &amp; Standar Kelulusan
                    </h6>
                    <div class="row g-3">
                      <div class="col-md-3">
                        <label class="form-label small fw-semibold text-secondary">Tingkat</label>
                        <select v-model="form.tingkat" class="form-select">
                          <option value="">-- Semua Tingkat --</option>
                          <option value="7">Kelas 7</option>
                          <option value="8">Kelas 8</option>
                          <option value="9">Kelas 9</option>
                          <option value="10">Kelas 10</option>
                          <option value="11">Kelas 11</option>
                          <option value="12">Kelas 12</option>
                        </select>
                      </div>
                      <div class="col-md-3">
                        <label class="form-label small fw-semibold text-secondary">Fase</label>
                        <select v-model="form.fase" class="form-select">
                          <option value="">-- Semua Fase --</option>
                          <option value="Fase D">Fase D (SMP)</option>
                          <option value="Fase E">Fase E (Kelas 10)</option>
                          <option value="Fase F">Fase F (Kelas 11-12)</option>
                        </select>
                      </div>
                      <div class="col-md-3">
                        <label class="form-label small fw-semibold text-secondary">JP per Minggu</label>
                        <input
                          v-model.number="form.jam_per_minggu"
                          type="number"
                          min="1"
                          max="10"
                          class="form-control"
                          placeholder="2"
                        />
                      </div>
                      <div class="col-md-3">
                        <label class="form-label small fw-semibold text-secondary">KKM / KKTP</label>
                        <input
                          v-model.number="form.kkm"
                          type="number"
                          min="0"
                          max="100"
                          step="0.5"
                          class="form-control"
                          placeholder="75.0"
                        />
                      </div>
                    </div>
                  </div>

                  <!-- Section 3: Pengajar, Alokasi & Status -->
                  <div>
                    <h6 class="fw-bold text-dark border-bottom pb-2 mb-3">
                      <i class="fas fa-user-check me-2 text-warning"></i>3. Pengampu &amp; Alokasi Kelas
                    </h6>
                    <div class="row g-3">
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Guru Pengampu</label>
                        <select v-model="form.id_guru" class="form-select">
                          <option value="">-- Belum ada / Kosongkan --</option>
                          <option v-for="guru in availableGurus" :key="guru.id" :value="guru.id">
                            {{ guru.nama }} {{ guru.nip ? '(' + guru.nip + ')' : '' }}
                          </option>
                        </select>
                      </div>
                      <div class="col-md-4">
                        <label class="form-label small fw-semibold text-secondary">Alokasi Kelas Terpilih</label>
                        <select v-model="form.kelas_id" class="form-select">
                          <option value="">-- Umum / Semua Kelas --</option>
                          <option v-for="k in availableKelas" :key="k.id" :value="k.id">
                            {{ k.kelas }}
                          </option>
                        </select>
                      </div>
                      <div class="col-md-2">
                        <label class="form-label small fw-semibold text-secondary">Status</label>
                        <select v-model="form.status" class="form-select">
                          <option value="Aktif">Aktif</option>
                          <option value="Nonaktif">Nonaktif</option>
                        </select>
                      </div>
                      <div class="col-12">
                        <label class="form-label small fw-semibold text-secondary">Capaian Pembelajaran / Deskripsi</label>
                        <textarea
                          v-model="form.deskripsi"
                          rows="2"
                          class="form-control"
                          placeholder="Ringkasan kompetensi dasar dan materi pokok pembelajaran"
                        ></textarea>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="modal-footer bg-light">
                  <button type="button" class="btn btn-light" @click="closeModal">Batal</button>
                  <button type="submit" class="btn btn-warning btn-shine text-dark fw-semibold d-flex align-items-center gap-2" :disabled="submitting">
                    <span v-if="submitting" class="spinner-border spinner-border-sm"></span>
                    <span>{{ isEditing ? 'Simpan Perubahan' : 'Tambah Mapel' }}</span>
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Modal Detail Lengkap Mapel -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showDetailModal" class="app-modal-overlay" tabindex="-1" @click.self="closeDetailModal">
          <div class="modal-dialog modal-dialog-centered modal-lg">
            <div class="modal-content border-0 shadow-lg bg-white">
              <div class="modal-header bg-warning text-dark">
                <h5 class="modal-title fw-bold d-flex align-items-center gap-2">
                  <i class="fas fa-book-open"></i>
                  <span>Detail Silabus Mata Pelajaran</span>
                </h5>
                <button type="button" class="btn-close" @click="closeDetailModal"></button>
              </div>
              <div v-if="selectedMapel" class="modal-body p-4">
                <!-- Header Profil Mapel -->
                <div class="d-flex flex-column flex-sm-row align-items-sm-center gap-3 p-3 bg-light rounded-3 mb-4">
                  <div class="avatar-detail rounded-circle d-flex align-items-center justify-content-center text-dark fs-2 bg-white shadow-sm" style="width: 64px; height: 64px;">
                    <i class="fas fa-book-reader text-warning"></i>
                  </div>
                  <div class="flex-grow-1">
                    <div class="d-flex flex-wrap align-items-center gap-2 mb-1">
                      <h4 class="fw-bold text-dark mb-0">{{ selectedMapel.mata_pelajaran }}</h4>
                      <span v-if="selectedMapel.kode_mapel" class="badge bg-secondary font-monospace">
                        {{ selectedMapel.kode_mapel }}
                      </span>
                    </div>
                    <div class="d-flex flex-wrap gap-2 align-items-center">
                      <span class="badge" :class="selectedMapel.status === 'Aktif' || !selectedMapel.status ? 'bg-success' : 'bg-secondary'">
                        {{ selectedMapel.status || 'Aktif' }}
                      </span>
                      <span class="badge bg-warning-subtle text-dark border">
                        {{ selectedMapel.kelompok || 'Kelompok Umum' }}
                      </span>
                      <span class="text-muted small">
                        <i class="fas fa-layer-group me-1"></i>{{ selectedMapel.kurikulum || 'Kurikulum Merdeka' }}
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Informasi Terperinci -->
                <div class="row g-3">
                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <h6 class="fw-bold text-dark mb-3 small text-uppercase">
                        <i class="fas fa-award text-warning me-2"></i>Beban Jam &amp; Penilaian
                      </h6>
                      <dl class="row mb-0 small">
                        <dt class="col-sm-5 text-secondary">Alokasi Waktu</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedMapel.jam_per_minggu || 2 }} JP / Minggu</dd>

                        <dt class="col-sm-5 text-secondary">Standar KKM/KKTP</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedMapel.kkm ? Number(selectedMapel.kkm).toFixed(1) : '75.0' }}</dd>

                        <dt class="col-sm-5 text-secondary">Tingkat</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedMapel.tingkat ? 'Kelas ' + selectedMapel.tingkat : 'Semua Tingkat' }}</dd>

                        <dt class="col-sm-5 text-secondary">Fase</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedMapel.fase || '-' }}</dd>
                      </dl>
                    </div>
                  </div>

                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <h6 class="fw-bold text-dark mb-3 small text-uppercase">
                        <i class="fas fa-chalkboard-teacher text-warning me-2"></i>Pengajar &amp; Rombel
                      </h6>
                      <dl class="row mb-0 small">
                        <dt class="col-sm-5 text-secondary">Guru Pengampu</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedMapel.guru || 'Belum Ditentukan' }}</dd>

                        <dt class="col-sm-5 text-secondary">Alokasi Rombel</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedMapel.nama_kelas || 'Seluruh Kelas' }}</dd>

                        <dt class="col-sm-5 text-secondary">Kurikulum</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedMapel.kurikulum || 'Kurikulum Merdeka' }}</dd>

                        <dt class="col-sm-5 text-secondary">Status Mapel</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedMapel.status || 'Aktif' }}</dd>
                      </dl>
                    </div>
                  </div>

                  <!-- Silabus / Deskripsi Kompetensi -->
                  <div class="col-12" v-if="selectedMapel.deskripsi">
                    <div class="border rounded-3 p-3 bg-light">
                      <h6 class="fw-bold text-secondary mb-1 small text-uppercase">
                        <i class="fas fa-sticky-note me-2"></i>Capaian &amp; Silabus Pembelajaran
                      </h6>
                      <p class="mb-0 text-dark small">{{ selectedMapel.deskripsi }}</p>
                    </div>
                  </div>
                </div>
              </div>
              <div class="modal-footer bg-light">
                <button type="button" class="btn btn-secondary" @click="closeDetailModal">Tutup</button>
                <button type="button" class="btn btn-warning text-dark fw-semibold" @click="openEditModal(selectedMapel); closeDetailModal()">
                  <i class="fas fa-edit me-1"></i> Ubah Data Ini
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Modal Konfirmasi Hapus Mapel -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showDeleteModal" class="app-modal-overlay" tabindex="-1" @click.self="closeDeleteModal">
          <div class="modal-dialog modal-dialog-centered modal-sm">
            <div class="modal-content border-0 shadow-lg text-center p-4 bg-white">
              <div class="bg-danger-subtle text-danger rounded-circle p-3 d-inline-flex align-items-center justify-content-center mx-auto mb-3" style="width: 58px; height: 58px;">
                <i class="fas fa-exclamation-triangle fs-4"></i>
              </div>
              <h5 class="fw-bold text-dark mb-1">Hapus Mata Pelajaran?</h5>
              <p class="text-secondary small mb-3">
                Apakah Anda yakin ingin menghapus mapel <strong class="text-dark">{{ mapelToDelete?.mata_pelajaran }}</strong>?
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { mapelApi } from '@/api/mapelApi'
import { guruApi } from '@/api/guruApi'
import { kelasApi } from '@/api/kelasApi'

const mapels = ref([])
const availableGurus = ref([])
const availableKelas = ref([])
const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)

const searchQuery = ref('')
const filterKelompok = ref('')
const filterStatus = ref('')

const alertMessage = ref('')
const alertClass = ref('alert-success')
const alertIcon = ref('fas fa-check-circle')

const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref(null)

const showDetailModal = ref(false)
const selectedMapel = ref(null)

const showDeleteModal = ref(false)
const mapelToDelete = ref(null)

const form = ref({
  mata_pelajaran: '',
  kode_mapel: '',
  kelompok: 'Umum',
  kurikulum: 'Kurikulum Merdeka',
  tingkat: '',
  fase: '',
  jam_per_minggu: 2,
  kkm: 75.0,
  status: 'Aktif',
  id_guru: '',
  kelas_id: '',
  deskripsi: '',
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

const fetchMapel = async () => {
  loading.value = true
  try {
    const res = await mapelApi.getAll()
    if (res.data && res.data.data) {
      mapels.value = res.data.data
    } else {
      mapels.value = []
    }
  } catch (err) {
    showAlert('Gagal mengambil daftar mata pelajaran.', 'danger')
  } finally {
    loading.value = false
  }
}

const fetchGurusAndKelas = async () => {
  try {
    const [guruRes, kelasRes] = await Promise.allSettled([guruApi.getAll(), kelasApi.getAll()])
    if (guruRes.status === 'fulfilled' && guruRes.value.data?.data) {
      availableGurus.value = guruRes.value.data.data
    }
    if (kelasRes.status === 'fulfilled' && kelasRes.value.data?.data) {
      availableKelas.value = kelasRes.value.data.data
    }
  } catch (err) {
    console.error('Error fetching supporting data for mapel', err)
  }
}

const filteredMapel = computed(() => {
  return mapels.value.filter((m) => {
    const query = searchQuery.value.toLowerCase().trim()
    const matchesSearch =
      !query ||
      m.mata_pelajaran?.toLowerCase().includes(query) ||
      m.kode_mapel?.toLowerCase().includes(query) ||
      m.guru?.toLowerCase().includes(query) ||
      m.nama_kelas?.toLowerCase().includes(query) ||
      m.deskripsi?.toLowerCase().includes(query)

    const matchesKelompok = !filterKelompok.value || m.kelompok === filterKelompok.value
    const matchesStatus = !filterStatus.value || m.status === filterStatus.value

    return matchesSearch && matchesKelompok && matchesStatus
  })
})

const openCreateModal = () => {
  isEditing.value = false
  editingId.value = null
  form.value = {
    mata_pelajaran: '',
    kode_mapel: '',
    kelompok: 'Umum',
    kurikulum: 'Kurikulum Merdeka',
    tingkat: '',
    fase: '',
    jam_per_minggu: 2,
    kkm: 75.0,
    status: 'Aktif',
    id_guru: '',
    kelas_id: '',
    deskripsi: '',
  }
  showModal.value = true
}

const openEditModal = (m) => {
  isEditing.value = true
  editingId.value = m.id
  form.value = {
    mata_pelajaran: m.mata_pelajaran || '',
    kode_mapel: m.kode_mapel || '',
    kelompok: m.kelompok || 'Umum',
    kurikulum: m.kurikulum || 'Kurikulum Merdeka',
    tingkat: m.tingkat || '',
    fase: m.fase || '',
    jam_per_minggu: m.jam_per_minggu || 2,
    kkm: m.kkm || 75.0,
    status: m.status || 'Aktif',
    id_guru: m.id_guru || '',
    kelas_id: m.kelas_id || '',
    deskripsi: m.deskripsi || '',
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const openDetailModal = (m) => {
  selectedMapel.value = m
  showDetailModal.value = true
}

const closeDetailModal = () => {
  showDetailModal.value = false
  selectedMapel.value = null
}

const openDeleteModal = (m) => {
  mapelToDelete.value = m
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  mapelToDelete.value = null
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (isEditing.value) {
      await mapelApi.update(editingId.value, form.value)
      showAlert('Mata pelajaran berhasil diperbarui!')
    } else {
      await mapelApi.create(form.value)
      showAlert('Mata pelajaran baru berhasil ditambahkan!')
    }
    closeModal()
    fetchMapel()
  } catch (err) {
    const errMsg = err.response?.data?.message || err.response?.data || err.message || 'Gagal memproses mata pelajaran.'
    showAlert(errMsg, 'danger')
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async () => {
  if (!mapelToDelete.value) return
  deleting.value = true
  try {
    await mapelApi.delete(mapelToDelete.value.id)
    showAlert(`Mata pelajaran "${mapelToDelete.value.mata_pelajaran}" berhasil dihapus!`)
    closeDeleteModal()
    fetchMapel()
  } catch (err) {
    showAlert('Gagal menghapus data mata pelajaran.', 'danger')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchMapel()
  fetchGurusAndKelas()
})
</script>

<style scoped>
.avatar-mapel {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.15), rgba(217, 119, 6, 0.2));
  color: #b45309;
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
tr:hover .avatar-mapel {
  transform: scale(1.1) rotate(-6deg);
}
.app-modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1055;
  padding: 1rem;
}
.modal-content {
  background-color: #ffffff !important;
  color: #1e293b !important;
  border-radius: 12px;
  overflow: hidden;
}
</style>
