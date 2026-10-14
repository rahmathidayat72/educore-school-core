<template>
  <div class="animate-fade-up">
    <!-- Header -->
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
      <div>
        <div class="d-flex align-items-center gap-2 mb-1">
          <h3 class="fw-bold text-dark mb-0">Jadwal Pelajaran</h3>
          <span class="badge bg-primary-subtle text-primary px-2 py-1 small">
            {{ filteredJadwal.length }} Sesi Pembelajaran
          </span>
        </div>
        <p class="text-secondary small mb-0">
          {{ authStore.isAdmin ? 'Kelola jadwal KBM seluruh rombel, alokasi jam tatap muka/praktikum, dan pencegahan bentrok ruang & guru.' : 'Kelola dan pantau jadwal kegiatan mengajar Anda di setiap kelas.' }}
        </p>
      </div>
      <div class="d-flex align-items-center gap-2">
        <!-- Switcher Mode Tampilan: Tabel vs Kalender Mingguan -->
        <div class="btn-group shadow-sm me-2">
          <button 
            type="button" 
            class="btn btn-sm" 
            :class="viewMode === 'table' ? 'btn-primary text-white fw-bold' : 'btn-outline-primary'"
            @click="viewMode = 'table'"
          >
            <i class="fas fa-list me-1"></i>Tabel Sesi
          </button>
          <button 
            type="button" 
            class="btn btn-sm" 
            :class="viewMode === 'calendar' ? 'btn-primary text-white fw-bold' : 'btn-outline-primary'"
            @click="viewMode = 'calendar'"
          >
            <i class="fas fa-calendar-alt me-1"></i>Matriks Mingguan
          </button>
        </div>

        <button @click="openCreateModal" class="btn btn-primary text-white fw-semibold d-flex align-items-center gap-2 shadow-sm">
          <i class="fas fa-plus"></i>
          <span>Tambah Jadwal Baru</span>
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
        <!-- Pencarian Multi-Kriteria -->
        <div class="col-lg-3 col-md-6">
          <div class="input-group">
            <span class="input-group-text bg-white text-muted border-end-0"><i class="fas fa-search"></i></span>
            <input
              v-model="searchQuery"
              type="text"
              class="form-control border-start-0 ps-0"
              placeholder="Cari mapel, kode, guru, kelas, ruang..."
            />
            <button v-if="searchQuery" class="btn btn-link text-secondary p-0 px-2" @click="searchQuery = ''" title="Bersihkan">
              <i class="fas fa-times-circle"></i>
            </button>
          </div>
        </div>

        <!-- Filter Hari -->
        <div class="col-lg-2 col-md-3">
          <select v-model="filterHari" class="form-select">
            <option value="">Semua Hari</option>
            <option v-for="(day, idx) in listHari" :key="idx" :value="idx + 1">
              {{ day }}
            </option>
          </select>
        </div>

        <!-- Filter Kelas -->
        <div class="col-lg-2 col-md-3">
          <select v-model="filterKelas" class="form-select">
            <option value="">Semua Rombel</option>
            <option v-for="k in listKelas" :key="k.id" :value="k.id">
              {{ k.kelas }}
            </option>
          </select>
        </div>

        <!-- Filter Jenis Kegiatan -->
        <div class="col-lg-2 col-md-4">
          <select v-model="filterKegiatan" class="form-select">
            <option value="">Semua Kegiatan</option>
            <option value="Tatap Muka">Tatap Muka</option>
            <option value="Praktikum">Praktikum / Lab</option>
            <option value="Projek (P5)">Projek (P5)</option>
            <option value="Remedial / Pengayaan">Remedial / Pengayaan</option>
          </select>
        </div>

        <!-- Filter Status -->
        <div class="col-lg-2 col-md-4">
          <select v-model="filterStatus" class="form-select">
            <option value="">Semua Status</option>
            <option value="Aktif">Aktif</option>
            <option value="Nonaktif">Nonaktif</option>
          </select>
        </div>

        <!-- Tombol Segarkan -->
        <div class="col-lg-1 col-md-4 text-md-end">
          <button @click="fetchData" class="btn btn-outline-secondary w-100 d-flex align-items-center justify-content-center p-2" :disabled="loading" title="Segarkan Data">
            <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
          </button>
        </div>
      </div>
    </div>

    <!-- Tabel Data Jadwal (View Mode: Table) -->
    <div v-if="viewMode === 'table'" class="card card-custom overflow-hidden border-0 shadow-sm">
      <div class="table-responsive">
        <table class="table table-custom table-hover mb-0">
          <thead>
            <tr>
              <th style="width: 50px;">No</th>
              <th>Hari &amp; Jam KBM</th>
              <th>Mata Pelajaran / Kode</th>
              <th>Rombongan Belajar</th>
              <th>Guru Pengampu</th>
              <th>Ruang &amp; Kegiatan</th>
              <th>Tahun / Semester</th>
              <th>Status</th>
              <th class="text-end" style="width: 130px;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="9" class="text-center py-5 text-secondary">
                <div class="spinner-border spinner-border-sm text-primary me-2" role="status"></div>
                Memuat data jadwal pelajaran...
              </td>
            </tr>
            <tr v-else-if="filteredJadwal.length === 0">
              <td colspan="9" class="text-center py-5 text-muted">
                <i class="fas fa-calendar-times fs-2 mb-2 d-block text-secondary"></i>
                Belum ada data sesi jadwal pelajaran yang sesuai.
              </td>
            </tr>
            <tr v-for="(j, index) in filteredJadwal" :key="j.id">
              <td class="text-secondary fw-semibold">{{ index + 1 }}</td>
              <td>
                <div class="d-flex flex-column">
                  <div class="d-flex align-items-center gap-1 mb-1">
                    <span class="badge bg-primary-subtle text-primary fw-bold">
                      {{ j.nama_hari }}
                    </span>
                    <span v-if="j.jam_ke_mulai" class="badge bg-light text-secondary border font-monospace small py-0">
                      Jam ke-{{ j.jam_ke_mulai }}{{ j.jam_ke_selesai && j.jam_ke_selesai !== j.jam_ke_mulai ? '-' + j.jam_ke_selesai : '' }}
                    </span>
                  </div>
                  <span class="fw-semibold text-dark small">
                    <i class="far fa-clock me-1 text-muted"></i>
                    {{ j.jam_mulai }} - {{ j.jam_selesai }}
                  </span>
                </div>
              </td>
              <td>
                <div class="d-flex align-items-center gap-2">
                  <div class="avatar-jadwal rounded-circle d-flex align-items-center justify-content-center flex-shrink-0">
                    <i class="fas fa-book-reader small text-primary"></i>
                  </div>
                  <div>
                    <span class="fw-bold text-dark d-block">{{ j.nama_pelajaran }}</span>
                    <span v-if="j.kode_jadwal" class="badge bg-light text-secondary border font-monospace small py-0">
                      {{ j.kode_jadwal }}
                    </span>
                    <span v-else-if="j.kode_mapel" class="text-muted small font-monospace">
                      {{ j.kode_mapel }}
                    </span>
                  </div>
                </div>
              </td>
              <td>
                <span class="badge bg-secondary-subtle text-secondary px-2 py-1">
                  <i class="fas fa-school me-1"></i>
                  {{ j.nama_kelas || '-' }}
                </span>
                <span v-if="j.jurusan_kelas" class="text-muted small d-block mt-1">
                  {{ j.jurusan_kelas }}
                </span>
              </td>
              <td>
                <div v-if="j.nama_guru" class="small">
                  <span class="fw-semibold text-dark d-block">
                    <i class="fas fa-chalkboard-teacher me-1 text-success"></i>
                    {{ j.nama_guru }}
                  </span>
                  <span v-if="j.nip_guru" class="text-muted font-monospace" style="font-size: 0.72rem;">
                    NIP: {{ j.nip_guru }}
                  </span>
                </div>
                <span v-else class="text-muted small fst-italic">Belum Ditentukan</span>
              </td>
              <td>
                <div class="small">
                  <span class="badge bg-light text-dark border d-inline-block mb-1">
                    <i class="fas fa-map-marker-alt text-danger me-1"></i>
                    {{ j.ruang || 'Umum / Kelas' }}
                  </span>
                  <div class="text-muted" style="font-size: 0.75rem;">
                    {{ j.jenis_kegiatan || 'Tatap Muka' }}
                  </div>
                </div>
              </td>
              <td>
                <div class="small text-muted">
                  <span class="text-dark fw-semibold d-block">{{ j.tahun_ajaran || '2026/2027' }}</span>
                  <span>Sem. {{ j.semester || 'Ganjil' }}</span>
                </div>
              </td>
              <td>
                <span v-if="j.aktif" class="badge bg-success-subtle text-success">
                  <i class="fas fa-check-circle me-1"></i> Aktif
                </span>
                <span v-else class="badge bg-secondary-subtle text-secondary">
                  Nonaktif
                </span>
              </td>
              <td class="text-end">
                <div class="btn-group btn-group-sm">
                  <button @click="openDetailModal(j)" class="btn btn-outline-info" title="Lihat Rincian Jadwal">
                    <i class="fas fa-eye"></i>
                  </button>
                  <button @click="openEditModal(j)" class="btn btn-outline-primary" title="Ubah Jadwal">
                    <i class="fas fa-edit"></i>
                  </button>
                  <button @click="openDeleteModal(j)" class="btn btn-outline-danger" title="Hapus Jadwal">
                    <i class="fas fa-trash-alt"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Tampilan Matriks Kalender Mingguan (View Mode: Calendar) -->
    <div v-else class="card card-custom p-4 border-0 shadow-sm">
      <div class="d-flex justify-content-between align-items-center mb-3">
        <h5 class="fw-bold text-dark mb-0">
          <i class="fas fa-calendar-week text-primary me-2"></i>Matriks Jadwal Mingguan Bebas Bentrok
        </h5>
        <span class="badge bg-info-subtle text-info small">Visualisasi Slot Waktu KBM</span>
      </div>

      <div class="row g-3">
        <div v-for="(dayName, dayIndex) in listHari.slice(0, 6)" :key="dayIndex" class="col-lg-4 col-md-6">
          <div class="border rounded-3 p-3 bg-light h-100 shadow-xs">
            <div class="d-flex justify-content-between align-items-center pb-2 border-bottom mb-2">
              <span class="fw-bold text-primary">{{ dayName }}</span>
              <span class="badge bg-primary text-white" style="font-size: 0.7rem;">
                {{ getJadwalByDay(dayIndex + 1).length }} Sesi
              </span>
            </div>

            <div v-if="getJadwalByDay(dayIndex + 1).length === 0" class="text-center py-4 text-muted small">
              <i class="far fa-calendar-times d-block mb-1 text-secondary"></i>
              Tidak ada sesi KBM
            </div>
            <div v-else class="d-flex flex-column gap-2">
              <div 
                v-for="j in getJadwalByDay(dayIndex + 1)" 
                :key="j.id"
                class="card card-sm border p-2 bg-white rounded-2 shadow-xs cursor-pointer hover-shadow"
                @click="openDetailModal(j)"
              >
                <div class="d-flex justify-content-between align-items-start mb-1">
                  <span class="fw-bold text-dark small">{{ j.nama_pelajaran }}</span>
                  <span class="badge bg-light text-secondary border font-monospace" style="font-size: 0.65rem;">
                    {{ j.jam_mulai }} - {{ j.jam_selesai }}
                  </span>
                </div>
                <div class="d-flex justify-content-between align-items-center text-muted" style="font-size: 0.75rem;">
                  <span><i class="fas fa-school text-primary me-1"></i>{{ j.nama_kelas || '-' }}</span>
                  <span><i class="fas fa-chalkboard-teacher text-success me-1"></i>{{ j.nama_guru || '-' }}</span>
                </div>
                <div class="mt-1 d-flex justify-content-between align-items-center" style="font-size: 0.7rem;">
                  <span class="badge bg-info-subtle text-info py-0 px-1">{{ j.ruang || 'Umum' }}</span>
                  <span v-if="j.aktif" class="text-success fw-semibold"><i class="fas fa-check-circle me-1"></i>Aktif</span>
                  <span v-else class="text-secondary fst-italic">Nonaktif</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Form (Create / Edit) -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showModal" class="app-modal-overlay" tabindex="-1" @click.self="closeModal">
          <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable modal-lg">
            <div class="modal-content border-0 shadow-lg bg-white d-flex flex-column" style="max-height: 90vh;">
              <div class="modal-header bg-primary text-white flex-shrink-0">
                <h5 class="modal-title fw-bold d-flex align-items-center gap-2">
                  <i :class="isEditing ? 'fas fa-edit' : 'fas fa-plus'"></i>
                  <span>{{ isEditing ? 'Ubah Jadwal Pelajaran' : 'Tambah Jadwal Pelajaran Baru' }}</span>
                </h5>
                <button type="button" class="btn-close btn-close-white" @click="closeModal"></button>
              </div>
              <form @submit.prevent="handleSubmit" class="d-flex flex-column overflow-hidden flex-grow-1 mb-0">
                <div class="modal-body p-4 overflow-y-auto" style="max-height: calc(90vh - 130px);">
                  <!-- Alert bentrok jika ada -->
                  <div v-if="conflictAlert" class="alert alert-danger small py-2 d-flex align-items-center gap-2 mb-3">
                    <i class="fas fa-exclamation-triangle fs-5"></i>
                    <div>{{ conflictAlert }}</div>
                  </div>

                  <!-- Section 1: Identitas Jadwal & Mapel -->
                  <div class="mb-4">
                    <h6 class="fw-bold text-primary border-bottom pb-2 mb-3">
                      <i class="fas fa-book me-2"></i>1. Identitas Sesi &amp; Mata Pelajaran
                    </h6>
                    <div class="row g-3">
                      <div class="col-md-7">
                        <label class="form-label small fw-semibold text-secondary">Mata Pelajaran <span class="text-danger">*</span></label>
                        <select v-model="form.mata_pelajaran_id" class="form-select" required @change="onMapelChange">
                          <option value="">-- Pilih Mata Pelajaran --</option>
                          <option v-for="m in listMapel" :key="m.id" :value="m.id">
                            {{ m.mata_pelajaran }} {{ m.kode_mapel ? '[' + m.kode_mapel + ']' : '' }}
                          </option>
                        </select>
                      </div>
                      <div class="col-md-5">
                        <label class="form-label small fw-semibold text-secondary">Kode Jadwal (Opsional)</label>
                        <input
                          v-model="form.kode_jadwal"
                          type="text"
                          class="form-control font-monospace"
                          placeholder="Contoh: JDW-10MIPA1-MAT-01"
                        />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Jenis Kegiatan Belajar</label>
                        <select v-model="form.jenis_kegiatan" class="form-select">
                          <option value="Tatap Muka">Tatap Muka Teori</option>
                          <option value="Praktikum">Praktikum / Eksperimen Lab</option>
                          <option value="Projek (P5)">Projek Penguatan Pelajar Pancasila (P5)</option>
                          <option value="Remedial / Pengayaan">Remedial / Bimbingan Belajar</option>
                        </select>
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Ruang Belajar / Laboratorium</label>
                        <input
                          v-model="form.ruang"
                          type="text"
                          class="form-control"
                          placeholder="Contoh: Ruang 101, Lab Fisika, Lab Komputer"
                        />
                      </div>
                    </div>
                  </div>

                  <!-- Section 2: Penugasan Akademik -->
                  <div class="mb-4">
                    <h6 class="fw-bold text-primary border-bottom pb-2 mb-3">
                      <i class="fas fa-chalkboard-teacher me-2"></i>2. Penugasan Guru &amp; Rombongan Belajar
                    </h6>
                    <div class="row g-3">
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Guru Pengampu</label>
                        <select v-model="form.id_guru" class="form-select">
                          <option value="">-- Sesuai Guru Mapel Default --</option>
                          <option v-for="g in listGuru" :key="g.id" :value="g.id">
                            {{ g.nama }} {{ g.nip ? '(' + g.nip + ')' : '' }}
                          </option>
                        </select>
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Rombongan Belajar (Kelas)</label>
                        <select v-model="form.kelas_id" class="form-select">
                          <option value="">-- Sesuai Kelas Mapel Default --</option>
                          <option v-for="k in listKelas" :key="k.id" :value="k.id">
                            {{ k.kelas }} {{ k.jurusan ? '(' + k.jurusan + ')' : '' }}
                          </option>
                        </select>
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Tahun Ajaran</label>
                        <input v-model="form.tahun_ajaran" type="text" class="form-control" placeholder="2026/2027" />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Semester</label>
                        <select v-model="form.semester" class="form-select">
                          <option value="Ganjil">Semester Ganjil</option>
                          <option value="Genap">Semester Genap</option>
                        </select>
                      </div>
                    </div>
                  </div>

                  <!-- Section 3: Alokasi Waktu Belajar -->
                  <div class="mb-4">
                    <h6 class="fw-bold text-primary border-bottom pb-2 mb-3">
                      <i class="far fa-clock me-2"></i>3. Waktu &amp; Jam Pelajaran (JP)
                    </h6>
                    <div class="row g-3">
                      <div class="col-md-4">
                        <label class="form-label small fw-semibold text-secondary">Hari KBM <span class="text-danger">*</span></label>
                        <select v-model="form.hari" class="form-select" required>
                          <option value="1">Senin</option>
                          <option value="2">Selasa</option>
                          <option value="3">Rabu</option>
                          <option value="4">Kamis</option>
                          <option value="5">Jumat</option>
                          <option value="6">Sabtu</option>
                          <option value="7">Minggu</option>
                        </select>
                      </div>
                      <div class="col-md-4">
                        <label class="form-label small fw-semibold text-secondary">Jam Pelajaran Mulai (JP)</label>
                        <input
                          v-model.number="form.jam_ke_mulai"
                          type="number"
                          min="1"
                          max="15"
                          class="form-control"
                          placeholder="Jam ke-1"
                        />
                      </div>
                      <div class="col-md-4">
                        <label class="form-label small fw-semibold text-secondary">Jam Pelajaran Selesai (JP)</label>
                        <input
                          v-model.number="form.jam_ke_selesai"
                          type="number"
                          min="1"
                          max="15"
                          class="form-control"
                          placeholder="Jam ke-2"
                        />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Waktu Mulai <span class="text-danger">*</span></label>
                        <input v-model="form.jam_mulai" type="time" class="form-control" required />
                      </div>
                      <div class="col-md-6">
                        <label class="form-label small fw-semibold text-secondary">Waktu Selesai <span class="text-danger">*</span></label>
                        <input v-model="form.jam_selesai" type="time" class="form-control" required />
                      </div>
                    </div>
                  </div>

                  <!-- Section 4: Status & Keterangan -->
                  <div>
                    <h6 class="fw-bold text-primary border-bottom pb-2 mb-3">
                      <i class="fas fa-sticky-note me-2"></i>4. Status &amp; Keterangan Tambahan
                    </h6>
                    <div class="row g-3">
                      <div class="col-12">
                        <label class="form-label small fw-semibold text-secondary">Keterangan / Catatan Sesi</label>
                        <textarea
                          v-model="form.keterangan"
                          rows="2"
                          class="form-control"
                          placeholder="Contoh: Menggunakan modul ajar bab 3, harap bawa jas praktikum"
                        ></textarea>
                      </div>
                      <div class="col-12">
                        <div class="form-check form-switch">
                          <input class="form-check-input" type="checkbox" id="aktifSwitch" v-model="form.aktif">
                          <label class="form-check-label small fw-semibold text-dark" for="aktifSwitch">
                            Status Jadwal Aktif (Terjadwal dalam KBM berjalan)
                          </label>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="modal-footer bg-light flex-shrink-0">
                  <button type="button" class="btn btn-light" @click="closeModal">Batal</button>
                  <button type="submit" class="btn btn-primary text-white fw-semibold d-flex align-items-center gap-2" :disabled="submitting">
                    <span v-if="submitting" class="spinner-border spinner-border-sm"></span>
                    <span>{{ isEditing ? 'Simpan Perubahan' : 'Tambah Jadwal' }}</span>
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Modal Detail Lengkap Jadwal -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showDetailModal" class="app-modal-overlay" tabindex="-1" @click.self="closeDetailModal">
          <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable modal-lg">
            <div class="modal-content border-0 shadow-lg bg-white d-flex flex-column" style="max-height: 90vh;">
              <div class="modal-header bg-primary text-white flex-shrink-0">
                <h5 class="modal-title fw-bold d-flex align-items-center gap-2">
                  <i class="fas fa-calendar-check"></i>
                  <span>Rincian Lengkap Jadwal Pembelajaran</span>
                </h5>
                <button type="button" class="btn-close btn-close-white" @click="closeDetailModal"></button>
              </div>
              <div v-if="selectedJadwal" class="modal-body p-4 overflow-y-auto" style="max-height: calc(90vh - 130px);">
                <!-- Header Profil Jadwal -->
                <div class="d-flex flex-column flex-sm-row align-items-sm-center gap-3 p-3 bg-light rounded-3 mb-4">
                  <div class="avatar-detail rounded-circle d-flex align-items-center justify-content-center text-primary fs-2 bg-white shadow-sm" style="width: 64px; height: 64px;">
                    <i class="fas fa-clock"></i>
                  </div>
                  <div class="flex-grow-1">
                    <div class="d-flex flex-wrap align-items-center gap-2 mb-1">
                      <h4 class="fw-bold text-dark mb-0">{{ selectedJadwal.nama_pelajaran }}</h4>
                      <span v-if="selectedJadwal.kode_jadwal" class="badge bg-secondary font-monospace">
                        {{ selectedJadwal.kode_jadwal }}
                      </span>
                    </div>
                    <div class="d-flex flex-wrap gap-2 align-items-center">
                      <span class="badge" :class="selectedJadwal.aktif ? 'bg-success' : 'bg-secondary'">
                        {{ selectedJadwal.aktif ? 'Sesi Aktif' : 'Nonaktif' }}
                      </span>
                      <span class="badge bg-primary-subtle text-primary">
                        {{ selectedJadwal.nama_hari }}, {{ selectedJadwal.jam_mulai }} - {{ selectedJadwal.jam_selesai }}
                      </span>
                      <span class="badge bg-warning-subtle text-dark border">
                        {{ selectedJadwal.jenis_kegiatan || 'Tatap Muka' }}
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Informasi Terperinci -->
                <div class="row g-3">
                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <h6 class="fw-bold text-primary mb-3 small text-uppercase">
                        <i class="fas fa-info-circle me-2"></i>Waktu &amp; Alokasi Belajar
                      </h6>
                      <dl class="row mb-0 small">
                        <dt class="col-sm-5 text-secondary">Hari</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedJadwal.nama_hari }}</dd>

                        <dt class="col-sm-5 text-secondary">Waktu KBM</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedJadwal.jam_mulai }} - {{ selectedJadwal.jam_selesai }}</dd>

                        <dt class="col-sm-5 text-secondary">Jam Pelajaran</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">
                          {{ selectedJadwal.jam_ke_mulai ? 'Jam ke-' + selectedJadwal.jam_ke_mulai + (selectedJadwal.jam_ke_selesai ? ' s/d ' + selectedJadwal.jam_ke_selesai : '') : '-' }}
                        </dd>

                        <dt class="col-sm-5 text-secondary">Ruang Belajar</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedJadwal.ruang || 'Umum / Kelas' }}</dd>

                        <dt class="col-sm-5 text-secondary">Jenis Sesi</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedJadwal.jenis_kegiatan || 'Tatap Muka' }}</dd>
                      </dl>
                    </div>
                  </div>

                  <div class="col-md-6">
                    <div class="border rounded-3 p-3 h-100 bg-white">
                      <h6 class="fw-bold text-primary mb-3 small text-uppercase">
                        <i class="fas fa-school me-2"></i>Rombel &amp; Pengajar
                      </h6>
                      <dl class="row mb-0 small">
                        <dt class="col-sm-5 text-secondary">Rombel (Kelas)</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedJadwal.nama_kelas || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Jurusan / Tk</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">
                          {{ selectedJadwal.tingkat_kelas ? 'Kelas ' + selectedJadwal.tingkat_kelas : '' }} {{ selectedJadwal.jurusan_kelas ? '(' + selectedJadwal.jurusan_kelas + ')' : '' }}
                        </dd>

                        <dt class="col-sm-5 text-secondary">Guru Pengampu</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedJadwal.nama_guru || 'Belum Ditentukan' }}</dd>

                        <dt class="col-sm-5 text-secondary">NIP Guru</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedJadwal.nip_guru || '-' }}</dd>

                        <dt class="col-sm-5 text-secondary">Tahun / Sem</dt>
                        <dd class="col-sm-7 fw-semibold text-dark">{{ selectedJadwal.tahun_ajaran }} (Sem. {{ selectedJadwal.semester }})</dd>
                      </dl>
                    </div>
                  </div>

                  <!-- Catatan / Keterangan Sesi -->
                  <div class="col-12" v-if="selectedJadwal.keterangan">
                    <div class="border rounded-3 p-3 bg-light">
                      <h6 class="fw-bold text-secondary mb-1 small text-uppercase">
                        <i class="fas fa-sticky-note me-2"></i>Keterangan Sesi Pembelajaran
                      </h6>
                      <p class="mb-0 text-dark small">{{ selectedJadwal.keterangan }}</p>
                    </div>
                  </div>
                </div>
              </div>
              <div class="modal-footer bg-light flex-shrink-0">
                <button type="button" class="btn btn-secondary" @click="closeDetailModal">Tutup</button>
                <button type="button" class="btn btn-primary" @click="openEditModal(selectedJadwal); closeDetailModal()">
                  <i class="fas fa-edit me-1"></i> Ubah Sesi Ini
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Modal Konfirmasi Hapus Jadwal -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showDeleteModal" class="app-modal-overlay" tabindex="-1" @click.self="closeDeleteModal">
          <div class="modal-dialog modal-dialog-centered modal-sm">
            <div class="modal-content border-0 shadow-lg text-center p-4 bg-white">
              <div class="bg-danger-subtle text-danger rounded-circle p-3 d-inline-flex align-items-center justify-content-center mx-auto mb-3" style="width: 58px; height: 58px;">
                <i class="fas fa-exclamation-triangle fs-4"></i>
              </div>
              <h5 class="fw-bold text-dark mb-1">Hapus Jadwal?</h5>
              <p class="text-secondary small mb-3">
                Apakah Anda yakin ingin menghapus jadwal <strong class="text-dark">{{ jadwalToDelete?.nama_pelajaran }}</strong> hari {{ jadwalToDelete?.nama_hari }} ({{ jadwalToDelete?.jam_mulai }} - {{ jadwalToDelete?.jam_selesai }})?
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
import { useAuthStore } from '@/stores/auth'
import { getJadwal, createJadwal, updateJadwal, deleteJadwal } from '@/api/jadwalApi'
import { mapelApi } from '@/api/mapelApi'
import { kelasApi } from '@/api/kelasApi'
import { guruApi } from '@/api/guruApi'

const authStore = useAuthStore()

const listHari = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu']
const jadwals = ref([])
const listMapel = ref([])
const listKelas = ref([])
const listGuru = ref([])

const viewMode = ref('table') // 'table' | 'calendar'

const getJadwalByDay = (dayInt) => {
  return filteredJadwal.value
    .filter((j) => parseInt(j.hari) === dayInt)
    .sort((a, b) => (a.jam_mulai || '').localeCompare(b.jam_mulai || ''))
}

const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)

const searchQuery = ref('')
const filterHari = ref('')
const filterKelas = ref('')
const filterKegiatan = ref('')
const filterStatus = ref('')

const alertMessage = ref('')
const alertType = ref('success')
const conflictAlert = ref('')

const showModal = ref(false)
const isEditing = ref(false)
const currentEditId = ref(null)

const showDetailModal = ref(false)
const selectedJadwal = ref(null)

const showDeleteModal = ref(false)
const jadwalToDelete = ref(null)

const form = ref({
  kode_jadwal: '',
  mata_pelajaran_id: '',
  id_guru: '',
  kelas_id: '',
  hari: '1',
  jam_ke_mulai: 1,
  jam_ke_selesai: 2,
  jam_mulai: '07:30',
  jam_selesai: '09:00',
  ruang: '',
  jenis_kegiatan: 'Tatap Muka',
  tahun_ajaran: '2026/2027',
  semester: 'Ganjil',
  keterangan: '',
  aktif: true
})

const alertClass = computed(() => {
  return alertType.value === 'success' ? 'alert-success' : 'alert-danger'
})

const alertIcon = computed(() => {
  return alertType.value === 'success' ? 'fas fa-check-circle' : 'fas fa-exclamation-circle'
})

const showAlert = (msg, type = 'success') => {
  alertMessage.value = msg
  alertType.value = type
  window.scrollTo({ top: 0, behavior: 'smooth' })
  setTimeout(() => {
    alertMessage.value = ''
  }, 5000)
}

const fetchData = async () => {
  loading.value = true
  try {
    const [resJadwal, resMapel, resKelas, resGuru] = await Promise.all([
      getJadwal(),
      mapelApi.getAll(),
      kelasApi.getAll(),
      guruApi.getAll()
    ])

    jadwals.value = resJadwal.data?.data || []
    listMapel.value = resMapel.data?.data || []
    listKelas.value = resKelas.data?.data || []
    listGuru.value = resGuru.data?.data || []
  } catch (err) {
    showAlert(err.response?.data?.message || 'Gagal memuat data jadwal pelajaran', 'danger')
  } finally {
    loading.value = false
  }
}

const onMapelChange = () => {
  const selected = listMapel.value.find((m) => m.id === form.value.mata_pelajaran_id)
  if (selected) {
    if (selected.id_guru && !form.value.id_guru) {
      form.value.id_guru = selected.id_guru
    }
    if (selected.kelas_id && !form.value.kelas_id) {
      form.value.kelas_id = selected.kelas_id
    }
  }
}

const filteredJadwal = computed(() => {
  return jadwals.value.filter((item) => {
    if (filterHari.value && String(item.hari) !== String(filterHari.value)) {
      return false
    }

    if (filterKelas.value && String(item.kelas_id) !== String(filterKelas.value)) {
      return false
    }

    if (filterKegiatan.value && item.jenis_kegiatan !== filterKegiatan.value) {
      return false
    }

    if (filterStatus.value) {
      const isAktif = filterStatus.value === 'Aktif'
      if (item.aktif !== isAktif) {
        return false
      }
    }

    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase().trim()
      const matchMapel = (item.nama_pelajaran || '').toLowerCase().includes(q)
      const matchKodeMapel = (item.kode_mapel || '').toLowerCase().includes(q)
      const matchKodeJadwal = (item.kode_jadwal || '').toLowerCase().includes(q)
      const matchGuru = (item.nama_guru || '').toLowerCase().includes(q)
      const matchKelas = (item.nama_kelas || '').toLowerCase().includes(q)
      const matchRuang = (item.ruang || '').toLowerCase().includes(q)
      const matchHari = (item.nama_hari || '').toLowerCase().includes(q)
      const matchKegiatan = (item.jenis_kegiatan || '').toLowerCase().includes(q)

      if (
        !matchMapel &&
        !matchKodeMapel &&
        !matchKodeJadwal &&
        !matchGuru &&
        !matchKelas &&
        !matchRuang &&
        !matchHari &&
        !matchKegiatan
      ) {
        return false
      }
    }

    return true
  })
})

const openCreateModal = () => {
  isEditing.value = false
  currentEditId.value = null
  conflictAlert.value = ''
  form.value = {
    kode_jadwal: '',
    mata_pelajaran_id: listMapel.value[0]?.id || '',
    id_guru: listMapel.value[0]?.id_guru || '',
    kelas_id: listMapel.value[0]?.kelas_id || '',
    hari: '1',
    jam_ke_mulai: 1,
    jam_ke_selesai: 2,
    jam_mulai: '07:30',
    jam_selesai: '09:00',
    ruang: '',
    jenis_kegiatan: 'Tatap Muka',
    tahun_ajaran: '2026/2027',
    semester: 'Ganjil',
    keterangan: '',
    aktif: true
  }
  showModal.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  currentEditId.value = item.id
  conflictAlert.value = ''
  form.value = {
    kode_jadwal: item.kode_jadwal || '',
    mata_pelajaran_id: item.mata_pelajaran_id || '',
    id_guru: item.id_guru || '',
    kelas_id: item.kelas_id || '',
    hari: String(item.hari),
    jam_ke_mulai: item.jam_ke_mulai || 1,
    jam_ke_selesai: item.jam_ke_selesai || 2,
    jam_mulai: item.jam_mulai || '',
    jam_selesai: item.jam_selesai || '',
    ruang: item.ruang || '',
    jenis_kegiatan: item.jenis_kegiatan || 'Tatap Muka',
    tahun_ajaran: item.tahun_ajaran || '2026/2027',
    semester: item.semester || 'Ganjil',
    keterangan: item.keterangan || '',
    aktif: Boolean(item.aktif)
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  conflictAlert.value = ''
}

const openDetailModal = (item) => {
  selectedJadwal.value = item
  showDetailModal.value = true
}

const closeDetailModal = () => {
  showDetailModal.value = false
  selectedJadwal.value = null
}

const handleSubmit = async () => {
  conflictAlert.value = ''
  submitting.value = true

  const payload = {
    kode_jadwal: form.value.kode_jadwal,
    mata_pelajaran_id: form.value.mata_pelajaran_id,
    id_guru: form.value.id_guru,
    kelas_id: form.value.kelas_id,
    hari: parseInt(form.value.hari),
    jam_ke_mulai: parseInt(form.value.jam_ke_mulai) || 0,
    jam_ke_selesai: parseInt(form.value.jam_ke_selesai) || 0,
    jam_mulai: form.value.jam_mulai,
    jam_selesai: form.value.jam_selesai,
    ruang: form.value.ruang,
    jenis_kegiatan: form.value.jenis_kegiatan,
    tahun_ajaran: form.value.tahun_ajaran,
    semester: form.value.semester,
    keterangan: form.value.keterangan,
    aktif: form.value.aktif
  }

  try {
    if (isEditing.value) {
      await updateJadwal(currentEditId.value, payload)
      showAlert('Jadwal pelajaran berhasil diperbarui')
    } else {
      await createJadwal(payload)
      showAlert('Jadwal pelajaran berhasil ditambahkan')
    }
    closeModal()
    fetchData()
  } catch (err) {
    const errMsg = err.response?.data?.message || err.response?.data || err.message || 'Terjadi kesalahan'
    conflictAlert.value = errMsg
  } finally {
    submitting.value = false
  }
}

const openDeleteModal = (item) => {
  jadwalToDelete.value = item
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  jadwalToDelete.value = null
}

const confirmDelete = async () => {
  if (!jadwalToDelete.value) return
  deleting.value = true
  try {
    await deleteJadwal(jadwalToDelete.value.id)
    showAlert('Jadwal pelajaran berhasil dihapus')
    closeDeleteModal()
    fetchData()
  } catch (err) {
    showAlert(err.response?.data?.message || 'Gagal menghapus jadwal pelajaran', 'danger')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.avatar-jadwal {
  width: 36px;
  height: 36px;
  background-color: rgba(13, 110, 253, 0.1);
}

.table-custom th {
  background-color: #f8fafc;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 0.85rem 1rem;
}

.table-custom td {
  padding: 0.85rem 1rem;
  vertical-align: middle;
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
  padding: 1.5rem 1rem;
  overflow-y: auto;
}

.modal-dialog {
  max-height: 90vh;
  margin: auto;
}

.modal-content {
  background-color: #ffffff !important;
  color: #1e293b !important;
  border-radius: 12px;
  overflow: hidden;
  max-height: 90vh;
}
</style>
