<template>
  <div class="animate-fade-up">
    <!-- Header -->
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
      <div>
        <div class="d-flex align-items-center gap-2 mb-1">
          <h3 class="fw-bold text-dark mb-0">Data Siswa / Peserta Didik</h3>
          <span class="badge bg-success-subtle text-success px-2 py-1 small">
            {{ filteredSiswa.length }} Siswa
          </span>
        </div>
        <p class="text-secondary small mb-0">Kelola biodata komprehensif siswa, penempatan rombongan belajar (kelas), serta data wali murid.</p>
      </div>
      <div class="d-flex align-items-center gap-2">
        <button @click="openImportModal" class="btn btn-outline-primary d-flex align-items-center gap-2 shadow-sm" title="Import data siswa secara masal dari Excel/CSV">
          <i class="fas fa-file-import"></i>
          <span>Import CSV</span>
        </button>
        <button @click="exportDataSiswa" class="btn btn-outline-success d-flex align-items-center gap-2 shadow-sm" :disabled="loading || filteredSiswa.length === 0" title="Unduh data siswa ke file Excel/CSV">
          <i class="fas fa-file-excel"></i>
          <span>Ekspor Excel</span>
        </button>
        <button @click="openCreateModal" class="btn btn-success btn-shine d-flex align-items-center gap-2 shadow-sm">
          <i class="fas fa-user-plus"></i>
          <span>Tambah Siswa</span>
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
              placeholder="Cari nama, NIS, NISN, wali, email, atau alamat..."
            />
            <button v-if="searchQuery" class="btn btn-link text-secondary p-0 px-2" @click="searchQuery = ''" title="Bersihkan">
              <i class="fas fa-times-circle"></i>
            </button>
          </div>
        </div>
        <div class="col-md-3">
          <select v-model="filterKelas" class="form-select">
            <option value="">Semua Kelas ({{ availableKelas.length }} Kelas)</option>
            <option v-for="k in availableKelas" :key="k.id" :value="k.id">
              {{ k.kelas }}
            </option>
          </select>
        </div>
        <div class="col-md-2">
          <select v-model="filterStatus" class="form-select">
            <option value="">Semua Status</option>
            <option value="Aktif">Aktif</option>
            <option value="Lulus">Lulus</option>
            <option value="Pindah">Pindah</option>
            <option value="Nonaktif">Nonaktif</option>
          </select>
        </div>
        <div class="col-md-2 text-md-end">
          <button @click="fetchSiswa" class="btn btn-outline-secondary w-100 d-flex align-items-center justify-content-center gap-2" :disabled="loading">
            <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
            <span>Segarkan</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Tabel Data Siswa -->
    <div class="card card-custom overflow-hidden border-0 shadow-sm">
      <div class="table-responsive">
        <table class="table table-custom table-hover mb-0">
          <thead>
            <tr>
              <th style="width: 50px;">No</th>
              <th>Nama Siswa &amp; Nomor Induk</th>
              <th>Kelas</th>
              <th>JK / Agama</th>
              <th>Wali Murid</th>
              <th>Status</th>
              <th class="text-end" style="width: 150px;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="7" class="text-center py-5 text-secondary">
                <div class="spinner-border spinner-border-sm text-success me-2" role="status"></div>
                Memuat data siswa...
              </td>
            </tr>
            <tr v-else-if="filteredSiswa.length === 0">
              <td colspan="7" class="text-center py-5 text-muted">
                <i class="fas fa-user-graduate fs-2 mb-2 d-block text-secondary"></i>
                Belum ada data siswa yang sesuai filter.
              </td>
            </tr>
            <tr v-for="(s, index) in filteredSiswa" :key="s.id">
              <td class="text-secondary fw-semibold">{{ index + 1 }}</td>
              <td>
                <div class="d-flex align-items-center gap-2">
                  <div class="avatar-student rounded-circle d-flex align-items-center justify-content-center flex-shrink-0">
                    <i class="fas fa-user-graduate small"></i>
                  </div>
                  <div>
                    <div class="fw-bold text-dark">{{ s.nama }}</div>
                    <div class="small text-muted d-flex align-items-center gap-2">
                      <span v-if="s.nis"><i class="fas fa-id-card text-success me-1"></i>NIS: {{ s.nis }}</span>
                      <span v-if="s.nisn"><i class="fas fa-fingerprint text-info me-1"></i>NISN: {{ s.nisn }}</span>
                      <span v-if="!s.nis && !s.nisn" class="text-secondary opacity-75">NIS/NISN: -</span>
                    </div>
                  </div>
                </div>
              </td>
              <td>
                <span class="badge bg-info-subtle text-info px-2 py-1">
                  {{ s.nama_kelas || s.kelas_id || '-' }}
                </span>
                <div v-if="s.tahun_masuk" class="small text-muted mt-1">
                  Masuk: {{ s.tahun_masuk }}
                </div>
              </td>
              <td>
                <div class="small text-dark fw-semibold">{{ s.jenis_kelamin || '-' }}</div>
                <div class="small text-muted">{{ s.agama || '-' }}</div>
              </td>
              <td>
                <div class="small text-dark fw-semibold">{{ s.nama_wali || '-' }}</div>
                <div v-if="s.no_telepon_wali" class="small text-muted">
                  <i class="fas fa-phone text-success me-1"></i>{{ s.no_telepon_wali }}
                </div>
              </td>
              <td>
                <span class="badge" :class="getStatusBadgeClass(s.status)">
                  {{ s.status || 'Aktif' }}
                </span>
              </td>
              <td class="text-end">
                <div class="btn-group btn-group-sm">
                  <button type="button" @click="openDetailModal(s)" class="btn btn-outline-info" title="Lihat Biodata Lengkap">
                    <i class="fas fa-eye"></i>
                  </button>
                  <button type="button" @click="openEditModal(s)" class="btn btn-outline-primary" title="Ubah Data Siswa">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button type="button" @click="openDeleteModal(s)" class="btn btn-outline-danger" title="Hapus Data Siswa">
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
          <div class="modal-dialog modal-dialog-centered modal-lg modal-dialog-scrollable">
            <div class="modal-content border-0 shadow-lg bg-white">
              <div class="modal-header bg-light">
                <h5 class="modal-title fw-bold text-dark d-flex align-items-center gap-2">
                  <i :class="isEditing ? 'fas fa-user-edit text-primary' : 'fas fa-user-plus text-success'"></i>
                  <span>{{ isEditing ? 'Ubah Biodata Lengkap Siswa' : 'Tambah Siswa Baru' }}</span>
                </h5>
                <button type="button" class="btn-close" @click="closeModal"></button>
              </div>
              <form @submit.prevent="handleSubmit" class="d-flex flex-column overflow-hidden">
                <div class="modal-body p-4 overflow-y-auto" style="max-height: calc(85vh - 130px);">
                  <!-- 1. Identitas & Rombel -->
                  <div class="mb-3">
                    <div class="section-title text-success small fw-bold text-uppercase mb-2">
                      <i class="fas fa-school me-1"></i>Identitas Siswa &amp; Kelas
                    </div>
                    <div class="row g-3">
                      <div class="col-md-8">
                        <label class="form-label small fw-semibold text-secondary">Nama Lengkap Siswa <span class="text-danger">*</span></label>
                        <input v-model="form.nama" type="text" class="form-control" required placeholder="Contoh: Muhammad Farhan Pratama" />
                      </div>
                      <div class="col-md-4">
                        <label class="form-label small fw-semibold text-secondary">Rombongan Belajar (Kelas) <span class="text-danger">*</span></label>
                        <select v-model="form.kelas_id" class="form-select" required>
                          <option value="" disabled>Pilih Kelas...</option>
                          <option v-for="k in availableKelas" :key="k.id" :value="k.id">
                            {{ k.kelas }}
                          </option>
                        </select>
                      </div>
                      <div class="col-md-4">
                        <label class="form-label small fw-semibold text-secondary">NIS (Nomor Induk Siswa)</label>
                        <input v-model="form.nis" type="text" class="form-control" placeholder="Contoh: 2024101" />
                      </div>
                      <div class="col-md-4">
                        <label class="form-label small fw-semibold text-secondary">NISN</label>
                        <input v-model="form.nisn" type="text" class="form-control" placeholder="10 digit NISN..." />
                      </div>
                      <div class="col-md-2">
                        <label class="form-label small fw-semibold text-secondary">Tahun Masuk</label>
                        <input v-model="form.tahun_masuk" type="text" class="form-control" placeholder="2024" />
                      </div>
                      <div class="col-md-2">
                        <label class="form-label small fw-semibold text-secondary">Status Siswa</label>
                        <select v-model="form.status" class="form-select">
                          <option value="Aktif">Aktif</option>
                          <option value="Lulus">Lulus</option>
                          <option value="Pindah">Pindah</option>
                          <option value="Nonaktif">Nonaktif</option>
                        </select>
                      </div>
                    </div>
                  </div>

                  <hr class="my-3 opacity-25" />

                  <!-- 2. Biodata Pribadi -->
                  <div class="mb-3">
                    <div class="section-title text-success small fw-bold text-uppercase mb-2">
                      <i class="fas fa-user-circle me-1"></i>Biodata Pribadi &amp; Kelahiran
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
                        <input v-model="form.tempat_lahir" type="text" class="form-control" placeholder="Contoh: Sleman" />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Tanggal Lahir</label>
                        <input v-model="form.tanggal_lahir" type="date" class="form-control" />
                      </div>
                    </div>
                  </div>

                  <hr class="my-3 opacity-25" />

                  <!-- 3. Orang Tua / Wali -->
                  <div class="mb-3">
                    <div class="section-title text-success small fw-bold text-uppercase mb-2">
                      <i class="fas fa-users me-1"></i>Orang Tua / Wali Murid
                    </div>
                    <div class="row g-3">
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Nama Lengkap Orang Tua / Wali</label>
                        <input v-model="form.nama_wali" type="text" class="form-control" placeholder="Contoh: Bambang Sutedjo" />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">No. Telepon / WhatsApp Wali</label>
                        <input v-model="form.no_telepon_wali" type="tel" class="form-control" placeholder="Contoh: 081298765432" />
                      </div>
                    </div>
                  </div>

                  <hr class="my-3 opacity-25" />

                  <!-- 4. Kontak Siswa & Alamat -->
                  <div class="mb-2">
                    <div class="section-title text-success small fw-bold text-uppercase mb-2">
                      <i class="fas fa-map-marker-alt me-1"></i>Kontak Siswa &amp; Alamat
                    </div>
                    <div class="row g-3">
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Email Siswa / Akun Utama <span class="text-danger">*</span></label>
                        <input v-model="form.email" type="email" class="form-control" required placeholder="farhan@siswa.sekolah.id" />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">No. Telepon Siswa (Opsional)</label>
                        <input v-model="form.no_telepon" type="tel" class="form-control" placeholder="Contoh: 085712345678" />
                      </div>
                      <div class="col-12">
                        <label class="form-label small fw-semibold text-secondary">Alamat Lengkap Tempat Tinggal <span class="text-danger">*</span></label>
                        <textarea v-model="form.alamat" rows="2" class="form-control" required placeholder="Jl. Cempaka No. 8, Sleman, DI Yogyakarta"></textarea>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="modal-footer bg-light">
                  <button type="button" class="btn btn-light" @click="closeModal">Batal</button>
                  <button type="submit" class="btn btn-success btn-shine d-flex align-items-center gap-2" :disabled="submitting">
                    <span v-if="submitting" class="spinner-border spinner-border-sm"></span>
                    <span>{{ isEditing ? 'Simpan Perubahan' : 'Tambah Siswa' }}</span>
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Modal Detail Lengkap Siswa -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showDetailModal" class="app-modal-overlay" tabindex="-1" @click.self="closeDetailModal">
          <div class="modal-dialog modal-dialog-centered modal-lg">
            <div class="modal-content border-0 shadow-lg bg-white">
              <div class="modal-header bg-success text-white">
                <h5 class="modal-title fw-bold d-flex align-items-center gap-2">
                  <i class="fas fa-id-card-alt"></i>
                  <span>Biodata Lengkap Peserta Didik</span>
                </h5>
                <button type="button" class="btn-close btn-close-white" @click="closeDetailModal"></button>
              </div>
              <div v-if="selectedSiswa" class="modal-body p-4">
                <!-- Profil Header Card -->
                <div class="d-flex flex-column flex-sm-row align-items-sm-center gap-3 p-3 bg-light rounded-3 mb-4">
                  <div class="avatar-detail rounded-circle d-flex align-items-center justify-content-center text-success fs-2">
                    <i class="fas fa-user-graduate"></i>
                  </div>
                  <div class="flex-grow-1">
                    <h4 class="fw-bold text-dark mb-1">{{ selectedSiswa.nama }}</h4>
                    <div class="d-flex flex-wrap gap-2 align-items-center">
                      <span class="badge bg-info-subtle text-info">
                        <i class="fas fa-school me-1"></i>Kelas {{ selectedSiswa.nama_kelas || selectedSiswa.kelas_id }}
                      </span>
                      <span class="badge" :class="getStatusBadgeClass(selectedSiswa.status)">
                        {{ selectedSiswa.status || 'Aktif' }}
                      </span>
                      <span v-if="selectedSiswa.tahun_masuk" class="badge bg-secondary-subtle text-secondary">
                        Tahun Masuk: {{ selectedSiswa.tahun_masuk }}
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Informasi Terperinci -->
                <div class="row g-3">
                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <h6 class="fw-bold text-success mb-3 small text-uppercase">
                        <i class="fas fa-fingerprint me-2"></i>Nomor Induk &amp; Status Akademik
                      </h6>
                      <dl class="row mb-0 small">
                        <dt class="col-sm-5 text-secondary">NIS</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.nis || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">NISN</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.nisn || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Rombel / Kelas</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.nama_kelas || selectedSiswa.kelas_id || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Tahun Masuk</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.tahun_masuk || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Status Siswa</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.status || 'Aktif' }}</dd>
                      </dl>
                    </div>
                  </div>

                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <h6 class="fw-bold text-success mb-3 small text-uppercase">
                        <i class="fas fa-user me-2"></i>Informasi Diri &amp; Kelahiran
                      </h6>
                      <dl class="row mb-0 small">
                        <dt class="col-sm-5 text-secondary">Jenis Kelamin</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.jenis_kelamin || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Agama</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.agama || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Tempat Lahir</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.tempat_lahir || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Tanggal Lahir</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ formatDateIndo(selectedSiswa.tanggal_lahir) }}</dd>
                      </dl>
                    </div>
                  </div>

                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <div class="d-flex justify-content-between align-items-center mb-2">
                        <h6 class="fw-bold text-success mb-0 small text-uppercase">
                          <i class="fas fa-users me-2"></i>Orang Tua / Wali Murid
                        </h6>
                        <button 
                          type="button" 
                          class="btn btn-xs btn-outline-success py-0 px-2 small" 
                          style="font-size: 0.75rem;"
                          @click="openLinkGuardianModal(selectedSiswa)"
                        >
                          <i class="fas fa-link me-1"></i>Kelola Wali
                        </button>
                      </div>

                      <div v-if="loadingGuardians" class="text-center py-2 text-muted small">
                        <div class="spinner-border spinner-border-sm text-success me-1"></div>
                        Memuat data relasi wali...
                      </div>
                      <div v-else-if="studentGuardians.length > 0" class="d-flex flex-column gap-2 mt-2">
                        <div 
                          v-for="g in studentGuardians" 
                          :key="g.id" 
                          class="p-2 border rounded bg-light small d-flex justify-content-between align-items-center"
                        >
                          <div>
                            <span class="fw-bold text-dark d-block">
                              {{ g.nama_wali }}
                              <span v-if="g.is_primary" class="badge bg-success text-white ms-1" style="font-size: 0.65rem;">Utama</span>
                            </span>
                            <span class="text-muted" style="font-size: 0.75rem;">
                              <i class="fas fa-user-tag me-1"></i>{{ g.hubungan || 'Wali' }} &bull; 
                              <i class="fas fa-phone me-1"></i>{{ g.telepon_wali || '-' }}
                            </span>
                          </div>
                          <button 
                            type="button" 
                            class="btn btn-link text-danger p-0 ms-2" 
                            title="Putuskan relasi"
                            @click="handleUnlinkGuardian(g.guardian_id, selectedSiswa.id)"
                          >
                            <i class="fas fa-times"></i>
                          </button>
                        </div>
                      </div>
                      <dl v-else class="row mb-0 small mt-1">
                        <dt class="col-sm-5 text-secondary">Nama Wali</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.nama_wali || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Telepon Wali</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.no_telepon_wali || '-' }}</dd>
                      </dl>
                    </div>
                  </div>

                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <h6 class="fw-bold text-success mb-3 small text-uppercase">
                        <i class="fas fa-address-book me-2"></i>Kontak Siswa &amp; Alamat
                      </h6>
                      <dl class="row mb-0 small">
                        <dt class="col-sm-5 text-secondary">Email</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.email }}</dd>

                        <dt class="col-sm-5 text-secondary">Telepon Siswa</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.no_telepon || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Alamat</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedSiswa.alamat || '-' }}</dd>
                      </dl>
                    </div>
                  </div>
                </div>
              </div>
              <div class="modal-footer bg-light">
                <button type="button" class="btn btn-outline-success btn-sm d-flex align-items-center gap-1" @click="openEditFromDetail">
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

    <!-- Modal Konfirmasi Hapus Siswa -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showDeleteModal" class="app-modal-overlay" tabindex="-1" @click.self="closeDeleteModal">
          <div class="modal-dialog modal-dialog-centered modal-sm">
            <div class="modal-content border-0 shadow-lg text-center p-4 bg-white">
              <div class="bg-danger-subtle text-danger rounded-circle p-3 d-inline-flex align-items-center justify-content-center mx-auto mb-3" style="width: 58px; height: 58px;">
                <i class="fas fa-exclamation-triangle fs-4"></i>
              </div>
              <h5 class="fw-bold text-dark mb-1">Hapus Data Siswa?</h5>
              <p class="text-secondary small mb-3">
                Apakah Anda yakin ingin menghapus data <strong class="text-dark">{{ siswaToDelete?.nama }}</strong>?
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

    <!-- Modal Penautan Wali Murid (Task-04.4) -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showLinkModal" class="app-modal-overlay" tabindex="-1" @click.self="closeLinkGuardianModal">
          <div class="modal-dialog modal-dialog-centered">
            <div class="modal-content border-0 shadow-lg bg-white">
              <div class="modal-header bg-success text-white">
                <h5 class="modal-title fw-bold d-flex align-items-center gap-2">
                  <i class="fas fa-link"></i>
                  <span>Tautkan Wali Murid ke {{ targetSiswaForLink?.nama }}</span>
                </h5>
                <button type="button" class="btn-close btn-close-white" @click="closeLinkGuardianModal"></button>
              </div>
              <div class="modal-body p-4">
                <ul class="nav nav-pills nav-fill mb-3" role="tablist">
                  <li class="nav-item">
                    <button 
                      class="nav-link py-1 small" 
                      :class="{ active: linkTab === 'existing' }" 
                      @click="linkTab = 'existing'"
                    >
                      Pilih Wali Terdaftar
                    </button>
                  </li>
                  <li class="nav-item">
                    <button 
                      class="nav-link py-1 small" 
                      :class="{ active: linkTab === 'new' }" 
                      @click="linkTab = 'new'"
                    >
                      Daftarkan Wali Baru
                    </button>
                  </li>
                </ul>

                <!-- Tab 1: Pilih Existing -->
                <div v-if="linkTab === 'existing'">
                  <div class="mb-3">
                    <label class="form-label small fw-semibold text-secondary">Cari &amp; Pilih Wali Murid</label>
                    <select v-model="linkForm.guardian_id" class="form-select" required>
                      <option value="">-- Pilih Wali Murid --</option>
                      <option v-for="g in allGuardians" :key="g.id" :value="g.id">
                        {{ g.nama }} ({{ g.hubungan || 'Wali' }}) - Telp: {{ g.telepon }}
                      </option>
                    </select>
                  </div>
                  <div class="form-check mb-3">
                    <input v-model="linkForm.is_primary" class="form-check-input" type="checkbox" id="checkPrimary" />
                    <label class="form-check-label small" for="checkPrimary">
                      Tetapkan sebagai Wali Utama (Kontak Darurat Prioritas)
                    </label>
                  </div>
                </div>

                <!-- Tab 2: Daftarkan Baru -->
                <div v-else>
                  <div class="row g-2 mb-2">
                    <div class="col-md-7">
                      <label class="form-label small fw-semibold text-secondary">Nama Lengkap Wali <span class="text-danger">*</span></label>
                      <input v-model="newGuardianForm.nama" type="text" class="form-control" placeholder="Nama ayah/ibu/wali" required />
                    </div>
                    <div class="col-md-5">
                      <label class="form-label small fw-semibold text-secondary">Hubungan</label>
                      <select v-model="newGuardianForm.hubungan" class="form-select">
                        <option value="Ayah">Ayah</option>
                        <option value="Ibu">Ibu</option>
                        <option value="Wali">Wali</option>
                        <option value="Kakak">Kakak</option>
                      </select>
                    </div>
                  </div>
                  <div class="row g-2 mb-2">
                    <div class="col-md-6">
                      <label class="form-label small fw-semibold text-secondary">No. Telepon / WA <span class="text-danger">*</span></label>
                      <input v-model="newGuardianForm.telepon" type="tel" class="form-control" placeholder="08xxxxxxxxxx" required />
                    </div>
                    <div class="col-md-6">
                      <label class="form-label small fw-semibold text-secondary">NIK (Opsional)</label>
                      <input v-model="newGuardianForm.nik" type="text" class="form-control" placeholder="16 digit NIK" />
                    </div>
                  </div>
                  <div class="mb-2">
                    <label class="form-label small fw-semibold text-secondary">Alamat Rumah (Opsional)</label>
                    <textarea v-model="newGuardianForm.alamat" rows="2" class="form-control" placeholder="Alamat wali murid"></textarea>
                  </div>
                  <div class="form-check mb-2">
                    <input v-model="newGuardianForm.is_primary" class="form-check-input" type="checkbox" id="checkNewPrimary" />
                    <label class="form-check-label small" for="checkNewPrimary">
                      Tetapkan sebagai Wali Utama
                    </label>
                  </div>
                </div>
              </div>
              <div class="modal-footer bg-light">
                <button type="button" class="btn btn-light" @click="closeLinkGuardianModal">Batal</button>
                <button 
                  type="button" 
                  class="btn btn-success d-flex align-items-center gap-1"
                  :disabled="submittingLink"
                  @click="submitLinkGuardian"
                >
                  <span v-if="submittingLink" class="spinner-border spinner-border-sm"></span>
                  <i v-else class="fas fa-check"></i>
                  <span>{{ linkTab === 'existing' ? 'Tautkan Wali' : 'Simpan &amp; Tautkan' }}</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
    <!-- Modal Import Data Masal Siswa -->
    <ImportModal
      :is-open="showImportModal"
      title="Siswa / Peserta Didik"
      template-filename="Template_Import_Siswa"
      :schema="siswaImportSchema"
      :sample-rows="siswaSampleRows"
      :on-import-row="handleImportSiswaRow"
      @close="showImportModal = false"
      @success="onImportSuccess"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { siswaApi } from '@/api/siswaApi'
import { kelasApi } from '@/api/kelasApi'
import { guardiansApi } from '@/api/guardiansApi'
import { exportToCSV } from '@/utils/exportHelper'
import ImportModal from '@/components/ImportModal.vue'

const router = useRouter()
const authStore = useAuthStore()

const siswas = ref([])
const availableKelas = ref([])
const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const searchQuery = ref('')
const filterKelas = ref('')
const filterStatus = ref('')

// Guardians State (Task-04.4)
const studentGuardians = ref([])
const loadingGuardians = ref(false)
const showLinkModal = ref(false)
const linkTab = ref('existing')
const allGuardians = ref([])
const targetSiswaForLink = ref(null)
const submittingLink = ref(false)

const linkForm = ref({
  guardian_id: '',
  is_primary: false,
})

const newGuardianForm = ref({
  nama: '',
  hubungan: 'Ayah',
  telepon: '',
  nik: '',
  alamat: '',
  is_primary: false,
})

const alertMessage = ref('')
const alertClass = ref('alert-success')
const alertIcon = ref('fas fa-check-circle')

const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref(null)

const showDetailModal = ref(false)
const selectedSiswa = ref(null)

const showDeleteModal = ref(false)
const siswaToDelete = ref(null)

const form = ref({
  nama: '',
  kelas_id: '',
  nis: '',
  nisn: '',
  jenis_kelamin: '',
  tempat_lahir: '',
  tanggal_lahir: '',
  agama: '',
  no_telepon: '',
  nama_wali: '',
  no_telepon_wali: '',
  tahun_masuk: '',
  status: 'Aktif',
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
    case 'Aktif':
      return 'bg-success-subtle text-success'
    case 'Lulus':
      return 'bg-primary-subtle text-primary'
    case 'Pindah':
      return 'bg-warning-subtle text-warning'
    case 'Nonaktif':
      return 'bg-danger-subtle text-danger'
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

const fetchSiswa = async () => {
  loading.value = true
  try {
    const res = await siswaApi.getAll()
    if (res.data && res.data.data) {
      siswas.value = res.data.data
    } else {
      siswas.value = []
    }
  } catch (err) {
    showAlert('Gagal mengambil daftar siswa.', 'danger')
  } finally {
    loading.value = false
  }
}

const fetchKelas = async () => {
  try {
    const res = await kelasApi.getAll()
    if (res.data && res.data.data) {
      availableKelas.value = res.data.data
    }
  } catch (err) {
    console.error('Error fetching kelas for siswa', err)
  }
}

const filteredSiswa = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  return siswas.value.filter((s) => {
    const matchesSearch =
      !query ||
      s.nama?.toLowerCase().includes(query) ||
      s.nis?.toLowerCase().includes(query) ||
      s.nisn?.toLowerCase().includes(query) ||
      s.nama_wali?.toLowerCase().includes(query) ||
      s.email?.toLowerCase().includes(query) ||
      s.no_telepon?.toLowerCase().includes(query) ||
      s.alamat?.toLowerCase().includes(query)

    const matchesKelas = filterKelas.value ? s.kelas_id === filterKelas.value : true
    const matchesStatus = filterStatus.value ? (s.status || 'Aktif') === filterStatus.value : true
    return matchesSearch && matchesKelas && matchesStatus
  })
})

const exportDataSiswa = () => {
  const headers = [
    { label: 'No', key: '__index' },
    { label: 'Nama Lengkap Siswa', key: 'nama' },
    { label: 'NIS', key: 'nis' },
    { label: 'NISN', key: 'nisn' },
    { label: 'Rombel Kelas', key: 'nama_kelas', formatter: (s) => s.nama_kelas || s.kelas_id || '-' },
    { label: 'Jenis Kelamin', key: 'jenis_kelamin' },
    { label: 'Tempat Lahir', key: 'tempat_lahir' },
    { label: 'Tanggal Lahir', key: 'tanggal_lahir', formatter: (s) => s.tanggal_lahir ? s.tanggal_lahir.substring(0, 10) : '-' },
    { label: 'Agama', key: 'agama' },
    { label: 'Status', key: 'status', formatter: (s) => s.status || 'Aktif' },
    { label: 'Nama Wali Murid', key: 'nama_wali' },
    { label: 'Telepon Wali', key: 'no_telepon_wali' },
    { label: 'Email Siswa', key: 'email' },
    { label: 'Alamat', key: 'alamat' }
  ]
  exportToCSV('Data_Siswa_Sekolah', headers, filteredSiswa.value)
}

const showImportModal = ref(false)

const siswaImportSchema = [
  { key: 'nama', label: 'Nama Lengkap', example: 'Ahmad Fauzi' },
  { key: 'nis', label: 'NIS', example: '2026001' },
  { key: 'nisn', label: 'NISN', example: '0081234567' },
  { key: 'kelas', label: 'Nama Kelas', example: 'X RPL 1' },
  { key: 'jenis_kelamin', label: 'Jenis Kelamin', example: 'Laki-laki' },
  { key: 'tempat_lahir', label: 'Tempat Lahir', example: 'Jakarta' },
  { key: 'tanggal_lahir', label: 'Tanggal Lahir', example: '2010-05-14' },
  { key: 'agama', label: 'Agama', example: 'Islam' },
  { key: 'status', label: 'Status', example: 'Aktif' },
  { key: 'nama_wali', label: 'Nama Wali', example: 'Bambang Sudiro' },
  { key: 'no_telepon_wali', label: 'No Telepon Wali', example: '081234567890' },
  { key: 'email', label: 'Email', example: 'ahmad.fauzi@siswa.id' },
  { key: 'alamat', label: 'Alamat', example: 'Jl. Merdeka No. 45' }
]

const siswaSampleRows = [
  {
    nama: 'Ahmad Fauzi',
    nis: '2026001',
    nisn: '0081234567',
    kelas: 'X RPL 1',
    jenis_kelamin: 'Laki-laki',
    tempat_lahir: 'Jakarta',
    tanggal_lahir: '2010-05-14',
    agama: 'Islam',
    status: 'Aktif',
    nama_wali: 'Bambang Sudiro',
    no_telepon_wali: '081234567890',
    email: 'ahmad.fauzi@siswa.id',
    alamat: 'Jl. Merdeka No. 45'
  },
  {
    nama: 'Siti Rahmawati',
    nis: '2026002',
    nisn: '0087654321',
    kelas: 'X RPL 1',
    jenis_kelamin: 'Perempuan',
    tempat_lahir: 'Bandung',
    tanggal_lahir: '2010-09-21',
    agama: 'Islam',
    status: 'Aktif',
    nama_wali: 'Hendra Gunawan',
    no_telepon_wali: '081398765432',
    email: 'siti.rahma@siswa.id',
    alamat: 'Jl. Kenanga No. 12'
  }
]

const openImportModal = () => {
  showImportModal.value = true
}

const handleImportSiswaRow = async (row) => {
  // Resolusi ID Kelas berdasarkan nama/kode kelas yang diinput di CSV
  let matchedKelasId = ''
  if (row.kelas || row.nama_kelas) {
    const targetName = (row.kelas || row.nama_kelas).toLowerCase().trim()
    const found = availableKelas.value.find(k => 
      k.kelas?.toLowerCase().trim() === targetName ||
      k.kode_kelas?.toLowerCase().trim() === targetName
    )
    if (found) matchedKelasId = found.id
  }

  // Jika tidak ditemukan, default ke kelas pertama yang tersedia
  if (!matchedKelasId && availableKelas.value.length > 0) {
    matchedKelasId = availableKelas.value[0].id
  }

  const payload = {
    nama: row.nama || row.nama_lengkap || '',
    nis: String(row.nis || ''),
    nisn: String(row.nisn || ''),
    kelas_id: matchedKelasId,
    jenis_kelamin: row.jenis_kelamin || 'Laki-laki',
    tempat_lahir: row.tempat_lahir || '',
    tanggal_lahir: row.tanggal_lahir ? `${row.tanggal_lahir}T00:00:00Z` : '',
    agama: row.agama || 'Islam',
    status: row.status || 'Aktif',
    nama_wali: row.nama_wali || '',
    no_telepon_wali: String(row.no_telepon_wali || ''),
    email: row.email || '',
    alamat: row.alamat || '',
    tahun_masuk: new Date().getFullYear().toString()
  }

  if (!payload.nama || !payload.nis) {
    throw new Error('Kolom Nama dan NIS wajib diisi.')
  }

  return await siswaApi.create(payload)
}

const onImportSuccess = ({ total, success, failed }) => {
  fetchSiswa()
  showAlert(`Berhasil mengimpor ${success} dari total ${total} data siswa. ${failed > 0 ? `(${failed} gagal)` : ''}`)
}

const openCreateModal = () => {
  isEditing.value = false
  editingId.value = null
  form.value = {
    nama: '',
    kelas_id: availableKelas.value[0]?.id || '',
    nis: '',
    nisn: '',
    jenis_kelamin: '',
    tempat_lahir: '',
    tanggal_lahir: '',
    agama: '',
    no_telepon: '',
    nama_wali: '',
    no_telepon_wali: '',
    tahun_masuk: new Date().getFullYear().toString(),
    status: 'Aktif',
    email: '',
    alamat: '',
  }
  showModal.value = true
}

const openEditModal = (s) => {
  isEditing.value = true
  editingId.value = s.id
  form.value = {
    nama: s.nama || '',
    kelas_id: s.kelas_id || '',
    nis: s.nis || '',
    nisn: s.nisn || '',
    jenis_kelamin: s.jenis_kelamin || '',
    tempat_lahir: s.tempat_lahir || '',
    tanggal_lahir: s.tanggal_lahir || '',
    agama: s.agama || '',
    no_telepon: s.no_telepon || '',
    nama_wali: s.nama_wali || '',
    no_telepon_wali: s.no_telepon_wali || '',
    tahun_masuk: s.tahun_masuk || '',
    status: s.status || 'Aktif',
    email: s.email || '',
    alamat: s.alamat || '',
  }
  showModal.value = true
}

const fetchGuardiansBySiswa = async (siswaId) => {
  if (!siswaId) return
  loadingGuardians.value = true
  try {
    const res = await guardiansApi.getGuardiansBySiswaId(siswaId)
    if (res.data && res.data.data) {
      studentGuardians.value = res.data.data
    } else {
      studentGuardians.value = []
    }
  } catch (err) {
    console.error('Error fetching student guardians', err)
    studentGuardians.value = []
  } finally {
    loadingGuardians.value = false
  }
}

const fetchAllGuardians = async () => {
  try {
    const res = await guardiansApi.getAll()
    if (res.data && res.data.data) {
      allGuardians.value = res.data.data
    }
  } catch (err) {
    console.error('Error fetching all guardians', err)
  }
}

const openDetailModal = (s) => {
  selectedSiswa.value = s
  showDetailModal.value = true
  studentGuardians.value = []
  fetchGuardiansBySiswa(s.id)
}

const closeDetailModal = () => {
  showDetailModal.value = false
  selectedSiswa.value = null
  studentGuardians.value = []
}

const openLinkGuardianModal = (s) => {
  targetSiswaForLink.value = s
  linkTab.value = 'existing'
  linkForm.value = {
    guardian_id: '',
    is_primary: false,
  }
  newGuardianForm.value = {
    nama: '',
    hubungan: 'Ayah',
    telepon: '',
    nik: '',
    alamat: '',
    is_primary: false,
  }
  fetchAllGuardians()
  showLinkModal.value = true
}

const closeLinkGuardianModal = () => {
  showLinkModal.value = false
  targetSiswaForLink.value = null
}

const submitLinkGuardian = async () => {
  if (!targetSiswaForLink.value) return
  submittingLink.value = true
  try {
    let guardianId = linkForm.value.guardian_id
    let isPrimary = linkForm.value.is_primary

    if (linkTab.value === 'new') {
      if (!newGuardianForm.value.nama || !newGuardianForm.value.telepon) {
        showAlert('Nama dan nomor telepon wali wajib diisi!', 'danger')
        submittingLink.value = false
        return
      }
      const created = await guardiansApi.create({
        nama: newGuardianForm.value.nama,
        hubungan: newGuardianForm.value.hubungan,
        telepon: newGuardianForm.value.telepon,
        nik: newGuardianForm.value.nik,
        alamat: newGuardianForm.value.alamat,
      })
      guardianId = created.data?.data?.id
      isPrimary = newGuardianForm.value.is_primary
    }

    if (!guardianId) {
      showAlert('Silakan pilih atau daftarkan wali murid terlebih dahulu.', 'danger')
      submittingLink.value = false
      return
    }

    await guardiansApi.linkStudent(guardianId, targetSiswaForLink.value.id, isPrimary)
    showAlert('Wali murid berhasil ditautkan ke siswa!')
    closeLinkGuardianModal()
    if (selectedSiswa.value?.id === targetSiswaForLink.value.id) {
      fetchGuardiansBySiswa(selectedSiswa.value.id)
    }
  } catch (err) {
    const errMsg = err.response?.data?.message || err.response?.data || err.message || 'Gagal menautkan wali murid.'
    showAlert(typeof errMsg === 'string' ? errMsg : 'Gagal menautkan wali murid.', 'danger')
  } finally {
    submittingLink.value = false
  }
}

const handleUnlinkGuardian = async (guardianId, siswaId) => {
  if (!confirm('Apakah Anda yakin ingin memutuskan tautan wali murid ini dari siswa?')) return
  try {
    await guardiansApi.unlinkStudent(guardianId, siswaId)
    showAlert('Tautan wali murid berhasil diputuskan!')
    fetchGuardiansBySiswa(siswaId)
  } catch (err) {
    const errMsg = err.response?.data?.message || err.message || 'Gagal memutuskan relasi wali.'
    showAlert(typeof errMsg === 'string' ? errMsg : 'Gagal memutuskan relasi wali.', 'danger')
  }
}

const openEditFromDetail = () => {
  if (selectedSiswa.value) {
    const target = selectedSiswa.value
    closeDetailModal()
    openEditModal(target)
  }
}

const closeModal = () => {
  showModal.value = false
}

const openDeleteModal = (s) => {
  siswaToDelete.value = s
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  siswaToDelete.value = null
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    if (isEditing.value) {
      await siswaApi.update(editingId.value, form.value)
      showAlert('Biodata siswa berhasil diperbarui!')
    } else {
      await siswaApi.create(form.value)
      showAlert('Siswa baru berhasil ditambahkan!')
    }
    closeModal()
    fetchSiswa()
  } catch (err) {
    const errMsg = err.response?.data?.message || err.response?.data || err.message || 'Gagal memproses data siswa.'
    showAlert(typeof errMsg === 'string' ? errMsg : 'Gagal memproses data siswa.', 'danger')
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async () => {
  if (!siswaToDelete.value) return
  deleting.value = true
  try {
    await siswaApi.delete(siswaToDelete.value.id)
    showAlert(`Data siswa "${siswaToDelete.value.nama}" berhasil dihapus!`)
    closeDeleteModal()
    fetchSiswa()
  } catch (err) {
    const errMsg = err.response?.data?.message || err.response?.data || err.message || 'Gagal menghapus data siswa.'
    showAlert(typeof errMsg === 'string' ? errMsg : 'Gagal menghapus data siswa.', 'danger')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  if (!authStore.isAdmin) {
    router.replace('/admin')
    return
  }
  fetchSiswa()
  fetchKelas()
})
</script>

<style scoped>
.avatar-student {
  width: 38px;
  height: 38px;
  background: linear-gradient(135deg, rgba(25, 135, 84, 0.15), rgba(20, 184, 166, 0.2));
  color: #198754;
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
tr:hover .avatar-student {
  transform: scale(1.1) rotate(-6deg);
}
.avatar-detail {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, rgba(25, 135, 84, 0.15), rgba(20, 184, 166, 0.25));
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
