<template>
  <div class="animate-fade-up">
    <!-- Header -->
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
      <div>
        <div class="d-flex align-items-center gap-2 mb-1">
          <h3 class="fw-bold text-dark mb-0">Data Rombongan Belajar (Kelas)</h3>
          <span class="badge bg-info-subtle text-info px-2 py-1 small">
            {{ filteredKelas.length }} Kelas
          </span>
        </div>
        <p class="text-secondary small mb-0">Kelola daftar rombongan belajar, wali kelas, ruang belajar, dan kapasitas penempatan siswa.</p>
      </div>
      <div class="d-flex align-items-center gap-2">
        <button @click="exportDataKelas" class="btn btn-outline-info d-flex align-items-center gap-2 shadow-sm" :disabled="loading || filteredKelas.length === 0" title="Unduh rekapitulasi kapasitas rombel ke Excel/CSV">
          <i class="fas fa-file-excel"></i>
          <span>Ekspor Rekap</span>
        </button>
        <button @click="openCreateModal" class="btn btn-info btn-shine text-white d-flex align-items-center gap-2 shadow-sm">
          <i class="fas fa-plus"></i>
          <span>Tambah Kelas Baru</span>
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
        <div class="col-md-5">
          <div class="input-group">
            <span class="input-group-text bg-white text-muted border-end-0"><i class="fas fa-search"></i></span>
            <input
              v-model="searchQuery"
              type="text"
              class="form-control border-start-0 ps-0"
              placeholder="Cari kelas, kode, wali kelas, jurusan, atau ruangan..."
            />
            <button v-if="searchQuery" class="btn btn-link text-secondary p-0 px-2" @click="searchQuery = ''" title="Bersihkan">
              <i class="fas fa-times-circle"></i>
            </button>
          </div>
        </div>
        <div class="col-md-3">
          <select v-model="filterStatus" class="form-select">
            <option value="">Semua Status Kelas</option>
            <option value="Aktif">Aktif</option>
            <option value="Nonaktif">Nonaktif</option>
          </select>
        </div>
        <div class="col-md-2">
          <select v-model="filterTingkat" class="form-select">
            <option value="">Semua Tingkat</option>
            <optgroup label="SD / MI">
              <option value="1">Kelas 1</option>
              <option value="2">Kelas 2</option>
              <option value="3">Kelas 3</option>
              <option value="4">Kelas 4</option>
              <option value="5">Kelas 5</option>
              <option value="6">Kelas 6</option>
            </optgroup>
            <optgroup label="SMP / MTs">
              <option value="7">Kelas 7</option>
              <option value="8">Kelas 8</option>
              <option value="9">Kelas 9</option>
            </optgroup>
            <optgroup label="SMA / SMK / MA">
              <option value="10">Kelas 10</option>
              <option value="11">Kelas 11</option>
              <option value="12">Kelas 12</option>
              <option value="13">Kelas 13 (SMK 4 Thn)</option>
            </optgroup>
          </select>
        </div>
        <div class="col-md-2 text-md-end">
          <button @click="fetchKelas" class="btn btn-outline-secondary w-100 d-flex align-items-center justify-content-center gap-2" :disabled="loading">
            <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
            <span>Segarkan</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Tabel Data Kelas -->
    <div class="card card-custom overflow-hidden border-0 shadow-sm">
      <div class="table-responsive">
        <table class="table table-custom table-hover mb-0">
          <thead>
            <tr>
              <th style="width: 50px;">No</th>
              <th>Nama Rombel / Kode</th>
              <th>Tingkat &amp; Jurusan</th>
              <th>Wali Kelas</th>
              <th>Ruang &amp; Thn Ajaran</th>
              <th>Kapasitas &amp; Siswa</th>
              <th>Status</th>
              <th class="text-end" style="width: 140px;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="8" class="text-center py-5 text-secondary">
                <div class="spinner-border spinner-border-sm text-info me-2" role="status"></div>
                Memuat data rombongan belajar...
              </td>
            </tr>
            <tr v-else-if="filteredKelas.length === 0">
              <td colspan="8" class="text-center py-5 text-muted">
                <i class="fas fa-school fs-2 mb-2 d-block text-secondary"></i>
                Belum ada data rombongan belajar yang sesuai.
              </td>
            </tr>
            <tr v-for="(k, index) in filteredKelas" :key="k.id">
              <td class="text-secondary fw-semibold">{{ index + 1 }}</td>
              <td>
                <div class="d-flex align-items-center gap-2">
                  <div class="avatar-class rounded-circle d-flex align-items-center justify-content-center flex-shrink-0">
                    <i class="fas fa-door-open small"></i>
                  </div>
                  <div>
                    <span class="fw-bold text-dark fs-6 d-block">{{ k.kelas }}</span>
                    <span v-if="k.kode_kelas" class="badge bg-light text-secondary border font-monospace small py-0">
                      {{ k.kode_kelas }}
                    </span>
                  </div>
                </div>
              </td>
              <td>
                <div class="small">
                  <span class="fw-semibold text-dark">{{ k.tingkat ? 'Kelas ' + k.tingkat : '-' }}</span>
                  <span v-if="k.fase" class="text-muted ms-1">({{ k.fase }})</span>
                  <div class="text-muted">{{ k.jurusan || 'Umum' }}</div>
                </div>
              </td>
              <td>
                <div v-if="k.nama_guru" class="d-flex align-items-center gap-2 text-dark small">
                  <span class="badge bg-light text-secondary border px-2 py-1">
                    <i class="fas fa-chalkboard-teacher me-1 text-primary"></i>
                    {{ k.nama_guru }}
                  </span>
                </div>
                <span v-else class="text-muted small fst-italic">Belum Ditentukan</span>
              </td>
              <td>
                <div class="small">
                  <span class="text-dark fw-semibold d-block">{{ k.ruangan || '-' }}</span>
                  <span class="text-muted">{{ k.tahun_ajaran || '2026/2027' }}</span>
                </div>
              </td>
              <td>
                <div class="small">
                  <div class="d-flex justify-content-between align-items-center mb-1">
                    <span class="fw-semibold text-dark">{{ k.total_siswa || 0 }} / {{ k.kapasitas || 36 }}</span>
                    <span class="text-muted small">{{ Math.round(((k.total_siswa || 0) / (k.kapasitas || 36)) * 100) }}%</span>
                  </div>
                  <div class="progress" style="height: 5px; width: 90px;">
                    <div
                      class="progress-bar"
                      :class="((k.total_siswa || 0) >= (k.kapasitas || 36)) ? 'bg-danger' : 'bg-info'"
                      role="progressbar"
                      :style="{ width: Math.min(100, Math.round(((k.total_siswa || 0) / (k.kapasitas || 36)) * 100)) + '%' }"
                    ></div>
                  </div>
                </div>
              </td>
              <td>
                <span class="badge" :class="k.status === 'Aktif' || !k.status ? 'bg-success-subtle text-success' : 'bg-secondary-subtle text-secondary'">
                  {{ k.status || 'Aktif' }}
                </span>
              </td>
              <td class="text-end">
                <div class="btn-group btn-group-sm">
                  <button @click="openDetailModal(k)" class="btn btn-outline-info" title="Lihat Profil Rombel Lengkap">
                    <i class="fas fa-eye"></i>
                  </button>
                  <button @click="openEditModal(k)" class="btn btn-outline-primary" title="Ubah Data Kelas">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button @click="openDeleteModal(k)" class="btn btn-outline-danger" title="Hapus Data Kelas">
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
              <div class="modal-header bg-info text-white">
                <h5 class="modal-title fw-bold d-flex align-items-center gap-2">
                  <i :class="isEditing ? 'fas fa-edit' : 'fas fa-plus'"></i>
                  <span>{{ isEditing ? 'Ubah Rombongan Belajar' : 'Tambah Rombongan Belajar Baru' }}</span>
                </h5>
                <button type="button" class="btn-close btn-close-white" @click="closeModal"></button>
              </div>
              <form @submit.prevent="handleSubmit">
                <div class="modal-body p-4">
                  <!-- Section 1: Identitas Rombel -->
                  <div class="mb-4">
                    <h6 class="fw-bold text-info border-bottom pb-2 mb-3">
                      <i class="fas fa-school me-2"></i>1. Identitas Rombongan Belajar
                    </h6>
                    <div class="row g-3">
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Nama Kelas / Rombel <span class="text-danger">*</span></label>
                        <input
                          v-model="form.kelas"
                          type="text"
                          class="form-control"
                          required
                          placeholder="Contoh: X MIPA 1, VII A, XII IPS 2"
                        />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Kode Rombel (Opsional)</label>
                        <input
                          v-model="form.kode_kelas"
                          type="text"
                          class="form-control font-monospace"
                          placeholder="Contoh: KLS-10-MIPA-1"
                        />
                      </div>
                      <div class="col-md-4">
                        <label class="form-label small fw-semibold text-secondary">Tingkat Pendidikan</label>
                        <select v-model="form.tingkat" class="form-select">
                          <option value="">-- Pilih Tingkat --</option>
                          <optgroup label="SD / MI">
                            <option value="1">Kelas 1 (SD/MI)</option>
                            <option value="2">Kelas 2 (SD/MI)</option>
                            <option value="3">Kelas 3 (SD/MI)</option>
                            <option value="4">Kelas 4 (SD/MI)</option>
                            <option value="5">Kelas 5 (SD/MI)</option>
                            <option value="6">Kelas 6 (SD/MI)</option>
                          </optgroup>
                          <optgroup label="SMP / MTs">
                            <option value="7">Kelas 7 (SMP/MTs)</option>
                            <option value="8">Kelas 8 (SMP/MTs)</option>
                            <option value="9">Kelas 9 (SMP/MTs)</option>
                          </optgroup>
                          <optgroup label="SMA / SMK / MA">
                            <option value="10">Kelas 10 (SMA/SMK)</option>
                            <option value="11">Kelas 11 (SMA/SMK)</option>
                            <option value="12">Kelas 12 (SMA/SMK)</option>
                            <option value="13">Kelas 13 (SMK 4 Tahun)</option>
                          </optgroup>
                        </select>
                      </div>
                      <div class="col-md-4">
                        <label class="form-label small fw-semibold text-secondary">Fase (Kurikulum Merdeka)</label>
                        <select v-model="form.fase" class="form-select">
                          <option value="">-- Pilih Fase --</option>
                          <option value="Fase A">Fase A (SD / Kelas 1-2)</option>
                          <option value="Fase B">Fase B (SD / Kelas 3-4)</option>
                          <option value="Fase C">Fase C (SD / Kelas 5-6)</option>
                          <option value="Fase D">Fase D (SMP / Kelas 7-9)</option>
                          <option value="Fase E">Fase E (SMA / Kelas 10)</option>
                          <option value="Fase F">Fase F (SMA / Kelas 11-12/13)</option>
                        </select>
                      </div>
                      <div class="col-md-4">
                        <label class="form-label small fw-semibold text-secondary">Jurusan / Peminatan</label>
                        <input
                          v-model="form.jurusan"
                          type="text"
                          class="form-control"
                          placeholder="Contoh: MIPA, IPS, TKJ, Umum"
                        />
                      </div>
                    </div>
                  </div>

                  <!-- Section 2: Penugasan & Sarana -->
                  <div class="mb-4">
                    <h6 class="fw-bold text-info border-bottom pb-2 mb-3">
                      <i class="fas fa-chalkboard-teacher me-2"></i>2. Pembimbing &amp; Sarana Ruang
                    </h6>
                    <div class="row g-3">
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Wali Kelas (Guru Pembimbing)</label>
                        <select v-model="form.id_guru" class="form-select">
                          <option value="">-- Belum ada / Kosongkan --</option>
                          <option v-for="guru in availableGurus" :key="guru.id" :value="guru.id">
                            {{ guru.nama }} {{ guru.nip ? '(' + guru.nip + ')' : '' }}
                          </option>
                        </select>
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Tahun Ajaran</label>
                        <input
                          v-model="form.tahun_ajaran"
                          type="text"
                          class="form-control"
                          placeholder="Contoh: 2026/2027"
                        />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Ruang Belajar</label>
                        <input
                          v-model="form.ruangan"
                          type="text"
                          class="form-control"
                          placeholder="Contoh: Ruang 101 Lantai 2, Lab 1"
                        />
                      </div>
                      <div class="col-md-3">
                        <label class="form-label small fw-semibold text-secondary">Kapasitas Maksimal Siswa</label>
                        <input
                          v-model.number="form.kapasitas"
                          type="number"
                          min="1"
                          max="60"
                          class="form-control"
                          placeholder="36"
                        />
                      </div>
                      <div class="col-md-3">
                        <label class="form-label small fw-semibold text-secondary">Status Rombel</label>
                        <select v-model="form.status" class="form-select">
                          <option value="Aktif">Aktif</option>
                          <option value="Nonaktif">Nonaktif</option>
                        </select>
                      </div>
                    </div>
                  </div>

                  <!-- Section 3: Catatan Rombel -->
                  <div>
                    <h6 class="fw-bold text-info border-bottom pb-2 mb-3">
                      <i class="fas fa-sticky-note me-2"></i>3. Keterangan Tambahan
                    </h6>
                    <div>
                      <label class="form-label small fw-semibold text-secondary">Catatan / Deskripsi Rombel</label>
                      <textarea
                        v-model="form.deskripsi"
                        rows="2"
                        class="form-control"
                        placeholder="Contoh: Rombel peminatan saintek reguler dengan kurikulum merdeka"
                      ></textarea>
                    </div>
                  </div>
                </div>

                <div class="modal-footer bg-light">
                  <button type="button" class="btn btn-light" @click="closeModal">Batal</button>
                  <button type="submit" class="btn btn-info btn-shine text-white d-flex align-items-center gap-2" :disabled="submitting">
                    <span v-if="submitting" class="spinner-border spinner-border-sm"></span>
                    <span>{{ isEditing ? 'Simpan Perubahan' : 'Tambah Kelas' }}</span>
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Modal Detail Lengkap Kelas -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showDetailModal" class="app-modal-overlay" tabindex="-1" @click.self="closeDetailModal">
          <div class="modal-dialog modal-dialog-centered modal-lg">
            <div class="modal-content border-0 shadow-lg bg-white">
              <div class="modal-header bg-info text-white">
                <h5 class="modal-title fw-bold d-flex align-items-center gap-2">
                  <i class="fas fa-school"></i>
                  <span>Profil Rombongan Belajar</span>
                </h5>
                <button type="button" class="btn-close btn-close-white" @click="closeDetailModal"></button>
              </div>
              <div v-if="selectedKelas" class="modal-body p-4">
                <!-- Header Profil -->
                <div class="d-flex flex-column flex-sm-row align-items-sm-center gap-3 p-3 bg-light rounded-3 mb-4">
                  <div class="avatar-detail rounded-circle d-flex align-items-center justify-content-center text-info fs-2 bg-white shadow-sm" style="width: 64px; height: 64px;">
                    <i class="fas fa-chalkboard"></i>
                  </div>
                  <div class="flex-grow-1">
                    <div class="d-flex flex-wrap align-items-center gap-2 mb-1">
                      <h4 class="fw-bold text-dark mb-0">{{ selectedKelas.kelas }}</h4>
                      <span v-if="selectedKelas.kode_kelas" class="badge bg-secondary font-monospace">
                        {{ selectedKelas.kode_kelas }}
                      </span>
                    </div>
                    <div class="d-flex flex-wrap gap-2 align-items-center">
                      <span class="badge" :class="selectedKelas.status === 'Aktif' || !selectedKelas.status ? 'bg-success' : 'bg-secondary'">
                        {{ selectedKelas.status || 'Aktif' }}
                      </span>
                      <span class="badge bg-info-subtle text-info">
                        Tahun Ajaran {{ selectedKelas.tahun_ajaran || '2026/2027' }}
                      </span>
                      <span class="text-muted small">
                        <i class="fas fa-map-marker-alt me-1"></i>{{ selectedKelas.ruangan || 'Ruangan belum diatur' }}
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Informasi Terperinci -->
                <div class="row g-3">
                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <h6 class="fw-bold text-info mb-3 small text-uppercase">
                        <i class="fas fa-info-circle me-2"></i>Struktur Rombel &amp; Kurikulum
                      </h6>
                      <dl class="row mb-0 small">
                        <dt class="col-sm-5 text-secondary">Tingkat</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedKelas.tingkat ? 'Kelas ' + selectedKelas.tingkat : '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Fase</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedKelas.fase || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Jurusan</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedKelas.jurusan || 'Umum' }}</dd>

                        <dt class="col-sm-5 text-secondary">Tahun Ajaran</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedKelas.tahun_ajaran || '2026/2027' }}</dd>
                      </dl>
                    </div>
                  </div>

                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <h6 class="fw-bold text-info mb-3 small text-uppercase">
                        <i class="fas fa-users me-2"></i>Kapasitas &amp; Wali Kelas
                      </h6>
                      <dl class="row mb-0 small">
                        <dt class="col-sm-5 text-secondary">Wali Kelas</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedKelas.nama_guru || 'Belum ditentukan' }}</dd>

                        <dt class="col-sm-5 text-secondary">Ruang Belajar</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedKelas.ruangan || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Kapasitas Maks</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedKelas.kapasitas || 36 }} Siswa</dd>

                        <dt class="col-sm-5 text-secondary">Siswa Terdaftar</dt>
                        <dd class="col-sm-7 fw-semibold text-info">{{ selectedKelas.total_siswa || 0 }} Siswa</dd>
                      </dl>
                    </div>
                  </div>

                  <!-- Catatan & Deskripsi -->
                  <div class="col-12" v-if="selectedKelas.deskripsi">
                    <div class="border rounded-3 p-3 bg-light">
                      <h6 class="fw-bold text-secondary mb-1 small text-uppercase">
                        <i class="fas fa-sticky-note me-2"></i>Catatan Rombel
                      </h6>
                      <p class="mb-0 text-dark small">{{ selectedKelas.deskripsi }}</p>
                    </div>
                  </div>
                </div>
              </div>
              <div class="modal-footer bg-light">
                <button type="button" class="btn btn-secondary" @click="closeDetailModal">Tutup</button>
                <button type="button" class="btn btn-info text-white" @click="openEditModal(selectedKelas); closeDetailModal()">
                  <i class="fas fa-edit me-1"></i> Ubah Data Ini
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Modal Konfirmasi Hapus Kelas -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showDeleteModal" class="app-modal-overlay" tabindex="-1" @click.self="closeDeleteModal">
          <div class="modal-dialog modal-dialog-centered modal-sm">
            <div class="modal-content border-0 shadow-lg text-center p-4 bg-white">
              <div class="bg-danger-subtle text-danger rounded-circle p-3 d-inline-flex align-items-center justify-content-center mx-auto mb-3" style="width: 58px; height: 58px;">
                <i class="fas fa-exclamation-triangle fs-4"></i>
              </div>
              <h5 class="fw-bold text-dark mb-1">Hapus Data Kelas?</h5>
              <p class="text-secondary small mb-3">
                Apakah Anda yakin ingin menghapus kelas <strong class="text-dark">{{ kelasToDelete?.kelas }}</strong>?
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
import { kelasApi } from '@/api/kelasApi'
import { guruApi } from '@/api/guruApi'
import { exportToCSV } from '@/utils/exportHelper'

const classes = ref([])
const availableGurus = ref([])
const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)

const searchQuery = ref('')
const filterStatus = ref('')
const filterTingkat = ref('')

const alertMessage = ref('')
const alertClass = ref('alert-success')
const alertIcon = ref('fas fa-check-circle')

const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref(null)

const showDetailModal = ref(false)
const selectedKelas = ref(null)

const showDeleteModal = ref(false)
const kelasToDelete = ref(null)

const form = ref({
  kelas: '',
  kode_kelas: '',
  tingkat: '',
  fase: '',
  jurusan: '',
  tahun_ajaran: '2026/2027',
  ruangan: '',
  kapasitas: 36,
  status: 'Aktif',
  deskripsi: '',
  id_guru: '',
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

const fetchKelas = async () => {
  loading.value = true
  try {
    const res = await kelasApi.getAll()
    if (res.data && res.data.data) {
      classes.value = res.data.data
    } else {
      classes.value = []
    }
  } catch (err) {
    showAlert('Gagal mengambil daftar rombongan belajar.', 'danger')
  } finally {
    loading.value = false
  }
}

const fetchGurus = async () => {
  try {
    const res = await guruApi.getAll()
    if (res.data && res.data.data) {
      availableGurus.value = res.data.data
    }
  } catch (err) {
    console.error('Error fetching gurus for kelas', err)
  }
}

const filteredKelas = computed(() => {
  return classes.value.filter((k) => {
    const query = searchQuery.value.toLowerCase().trim()
    const matchesSearch =
      !query ||
      k.kelas?.toLowerCase().includes(query) ||
      k.kode_kelas?.toLowerCase().includes(query) ||
      k.nama_guru?.toLowerCase().includes(query) ||
      k.jurusan?.toLowerCase().includes(query) ||
      k.ruangan?.toLowerCase().includes(query)

    const matchesStatus = !filterStatus.value || k.status === filterStatus.value
    const matchesTingkat = !filterTingkat.value || k.tingkat === filterTingkat.value

    return matchesSearch && matchesStatus && matchesTingkat
  })
})

const exportDataKelas = () => {
  const headers = [
    { label: 'No', key: '__index' },
    { label: 'Nama Rombongan Belajar', key: 'kelas' },
    { label: 'Kode Kelas', key: 'kode_kelas' },
    { label: 'Tingkat', key: 'tingkat', formatter: (k) => k.tingkat ? `Kelas ${k.tingkat}` : '-' },
    { label: 'Fase', key: 'fase' },
    { label: 'Jurusan / Peminatan', key: 'jurusan' },
    { label: 'Wali Kelas', key: 'nama_guru', formatter: (k) => k.nama_guru || '-' },
    { label: 'Ruang Belajar', key: 'ruangan' },
    { label: 'Tahun Ajaran', key: 'tahun_ajaran' },
    { label: 'Siswa Terdaftar', key: 'total_siswa', formatter: (k) => k.total_siswa || 0 },
    { label: 'Kapasitas Maksimal', key: 'kapasitas', formatter: (k) => k.kapasitas || 36 },
    { label: 'Persentase Keterisian (%)', key: 'persentase', formatter: (k) => `${Math.round(((k.total_siswa || 0) / (k.kapasitas || 36)) * 100)}%` },
    { label: 'Sisa Kuota Kursi', key: 'sisa', formatter: (k) => Math.max(0, (k.kapasitas || 36) - (k.total_siswa || 0)) },
    { label: 'Status Rombel', key: 'status', formatter: (k) => k.status || 'Aktif' }
  ]
  exportToCSV('Rekapitulasi_Kapasitas_Rombel', headers, filteredKelas.value)
}

const openCreateModal = () => {
  isEditing.value = false
  editingId.value = null
  form.value = {
    kelas: '',
    kode_kelas: '',
    tingkat: '',
    fase: '',
    jurusan: '',
    tahun_ajaran: '2026/2027',
    ruangan: '',
    kapasitas: 36,
    status: 'Aktif',
    deskripsi: '',
    id_guru: '',
  }
  showModal.value = true
}

const openEditModal = (k) => {
  isEditing.value = true
  editingId.value = k.id
  form.value = {
    kelas: k.kelas || '',
    kode_kelas: k.kode_kelas || '',
    tingkat: k.tingkat || '',
    fase: k.fase || '',
    jurusan: k.jurusan || '',
    tahun_ajaran: k.tahun_ajaran || '2026/2027',
    ruangan: k.ruangan || '',
    kapasitas: k.kapasitas || 36,
    status: k.status || 'Aktif',
    deskripsi: k.deskripsi || '',
    id_guru: k.id_guru || '',
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const openDetailModal = (k) => {
  selectedKelas.value = k
  showDetailModal.value = true
}

const closeDetailModal = () => {
  showDetailModal.value = false
  selectedKelas.value = null
}

const openDeleteModal = (k) => {
  kelasToDelete.value = k
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  kelasToDelete.value = null
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (isEditing.value) {
      await kelasApi.update(editingId.value, form.value)
      showAlert('Data kelas berhasil diperbarui!')
    } else {
      await kelasApi.create(form.value)
      showAlert('Kelas baru berhasil ditambahkan!')
    }
    closeModal()
    fetchKelas()
  } catch (err) {
    const errMsg = err.response?.data?.message || err.response?.data || err.message || 'Gagal memproses data kelas.'
    showAlert(errMsg, 'danger')
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async () => {
  if (!kelasToDelete.value) return
  deleting.value = true
  try {
    await kelasApi.delete(kelasToDelete.value.id)
    showAlert(`Data kelas "${kelasToDelete.value.kelas}" berhasil dihapus!`)
    closeDeleteModal()
    fetchKelas()
  } catch (err) {
    showAlert('Gagal menghapus data kelas.', 'danger')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchKelas()
  fetchGurus()
})
</script>

<style scoped>
.avatar-class {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, rgba(13, 202, 240, 0.15), rgba(2, 132, 199, 0.2));
  color: #0284c7;
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
tr:hover .avatar-class {
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
