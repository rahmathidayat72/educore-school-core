<template>
  <div>
    <!-- Header dengan Jam Realtime WIB & Info Pengguna -->
    <div class="card card-custom p-3 p-md-4 mb-4 border-0 shadow-sm bg-gradient-header">
      <div class="d-flex flex-column flex-lg-row justify-content-between align-items-lg-center gap-3">
        <div>
          <div class="d-flex align-items-center gap-2 mb-1 flex-wrap">
            <span class="badge bg-success-subtle text-success border border-success-subtle px-2 py-1 small fw-semibold">
              <i class="fas fa-circle-dot me-1 text-success pulse-beacon"></i>Sistem Berjalan Normal
            </span>
            <span class="badge bg-light text-dark border px-2 py-1 small">
              {{ authStore.isAdmin ? 'Hak Akses: Administrator' : (authStore.isGuru ? 'Hak Akses: Guru Pendidik' : 'Hak Akses: Pengguna') }}
            </span>
            <span class="badge bg-primary-subtle text-primary border border-primary-subtle px-2 py-1 small">
              Tahun Ajaran 2026/2027
            </span>
          </div>
          <h3 class="fw-bold text-dark mb-1">
            {{ greetingText }}, {{ authStore.isGuru && teacherData?.nama_guru ? teacherData.nama_guru : (authStore.currentUser?.username || 'Pengguna') }}!
          </h3>
          <p class="text-secondary small mb-0">
            {{ authStore.isAdmin ? 'Ringkasan operasional terpusat, audit aktivitas, dan kualitas integritas data sekolah.' : 'Informasi kelas bimbingan, daftar siswa, serta agenda jadwal mengajar Anda hari ini.' }}
          </p>
        </div>

        <!-- Widget Jam Digital Realtime WIB & Kalender Trigger -->
        <div class="d-flex align-items-center gap-3 bg-white p-2 p-md-3 rounded-3 shadow-sm border clock-card">
          <div class="text-center px-2 py-1 bg-primary-subtle text-primary rounded-3 clock-icon-box">
            <i class="fas fa-clock fs-3 animate-icon-float"></i>
          </div>
          <div class="border-start ps-3">
            <div class="d-flex align-items-baseline gap-2">
              <span class="font-monospace fw-bold fs-3 text-dark lh-1 clock-time">{{ currentTimeWIB }}</span>
              <span class="badge bg-danger text-white small pulse-beacon" style="font-size: 0.7rem; letter-spacing: 0.05em;">WIB</span>
            </div>
            <div
              class="text-secondary small mt-1 fw-medium calendar-trigger"
              style="font-size: 0.8rem; cursor: pointer;"
              @click="openCalendarModal"
              title="Klik untuk membuka Kalender Akademik & Hari Libur Indonesia"
            >
              <i class="far fa-calendar-alt me-1 text-primary"></i>
              <span class="calendar-link-text">{{ currentDateWIB }}</span>
              <i class="fas fa-external-link-alt ms-1 text-muted" style="font-size: 0.65rem;"></i>
            </div>
          </div>
          <button @click="loadDashboardData" class="btn btn-outline-secondary btn-sm ms-2 btn-living" :disabled="loading" title="Segarkan Data Dashboard">
            <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
          </button>
        </div>
      </div>
    </div>

    <!-- Alert jika error -->
    <div v-if="error" class="alert alert-danger small py-2 d-flex align-items-center gap-2 mb-4">
      <i class="fas fa-exclamation-triangle"></i>
      <div>{{ error }}</div>
    </div>

    <!-- ========================================================
         A. DASHBOARD KHUSUS ROLE GURU (PERSONAL & INFORMATIF)
         ======================================================== -->
    <div v-if="authStore.isGuru" class="animate-fade-up">
      <!-- Highlight Banner Jadwal Terdekat / Berlangsung -->
      <div v-if="teacherData?.next_schedule" class="card card-custom border-primary border-opacity-50 p-3 p-md-4 mb-4 bg-primary-subtle bg-opacity-10 shadow-sm">
        <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3">
          <div class="d-flex align-items-center gap-3">
            <div class="bg-primary text-white p-3 rounded-circle shadow-sm">
              <i class="fas fa-bell fs-4 animate-icon-float"></i>
            </div>
            <div>
              <div class="d-flex align-items-center gap-2 mb-1">
                <span class="badge" :class="teacherData.next_schedule.status === 'BERLANGSUNG' ? 'bg-success text-white pulse-beacon' : 'bg-primary text-white'">
                  {{ teacherData.next_schedule.status === 'BERLANGSUNG' ? 'SEDANG BERLANGSUNG' : 'SESI MENGAJAR BERIKUTNYA' }}
                </span>
                <span class="badge bg-light text-dark border">
                  Hari {{ teacherData.next_schedule.nama_hari }} &bull; {{ teacherData.next_schedule.jam_mulai }} - {{ teacherData.next_schedule.jam_selesai }} WIB
                </span>
              </div>
              <h5 class="fw-bold text-dark mb-0">
                {{ teacherData.next_schedule.mata_pelajaran }}
                <span class="text-secondary fw-normal">di kelas</span>
                <span class="text-primary">{{ teacherData.next_schedule.nama_kelas }}</span>
              </h5>
              <div class="text-muted small mt-1">
                <i class="fas fa-map-marker-alt me-1 text-danger"></i>Ruangan: {{ teacherData.next_schedule.ruang || 'Ruang Belajar Umum' }}
              </div>
            </div>
          </div>
          <router-link to="/admin/jadwal" class="btn btn-primary text-white fw-semibold px-3 py-2 d-flex align-items-center gap-2 shadow-sm align-self-md-center">
            <i class="fas fa-calendar-alt"></i>
            <span>Buka Jadwal Lengkap</span>
          </router-link>
        </div>
      </div>

      <!-- KPI Cards Guru (Lingkup Tugas Spesifik) -->
      <div class="row g-3 mb-4">
        <!-- Kelas Wali -->
        <div class="col-sm-6 col-xl-3">
          <div class="card card-custom card-living-shine p-3 border-start border-4 border-success h-100 stat-hover-card">
            <div class="d-flex align-items-start justify-content-between">
              <div>
                <div class="text-secondary small fw-semibold text-uppercase">Wali Kelas</div>
                <div class="fs-4 fw-bold text-dark mt-1">
                  {{ teacherData?.nama_kelas_wali || 'Bukan Wali Kelas' }}
                </div>
                <div class="mt-2">
                  <span v-if="teacherData?.total_kelas_wali" class="badge bg-success-subtle text-success border border-success-subtle small py-1 px-2">
                    <i class="fas fa-check-circle me-1"></i>{{ teacherData.total_siswa_wali }} Siswa Bimbingan
                  </span>
                  <span v-else class="badge bg-secondary-subtle text-secondary small py-1 px-2">
                    Hanya Pengampu Mapel
                  </span>
                </div>
              </div>
              <div class="bg-success-subtle text-success p-3 rounded-circle stat-icon-bubble">
                <i class="fas fa-user-shield fs-4"></i>
              </div>
            </div>
            <div class="mt-3 pt-2 border-top small d-flex justify-content-between align-items-center">
              <span class="text-muted small">Status Perwalian</span>
              <router-link to="/admin/kelas" class="text-decoration-none text-success fw-medium">
                Data Kelas <i class="fas fa-arrow-right ms-1 transition-icon"></i>
              </router-link>
            </div>
          </div>
        </div>

        <!-- Kelas yang Diajar -->
        <div class="col-sm-6 col-xl-3">
          <div class="card card-custom card-living-shine p-3 border-start border-4 border-primary h-100 stat-hover-card">
            <div class="d-flex align-items-start justify-content-between">
              <div>
                <div class="text-secondary small fw-semibold text-uppercase">Kelas Diajar</div>
                <div class="fs-4 fw-bold text-dark mt-1">
                  <AnimatedNumber :value="teacherData?.total_kelas_ajar || 0" /> Kelas
                </div>
                <div class="mt-2">
                  <span class="badge bg-primary-subtle text-primary border border-primary-subtle small py-1 px-2 text-truncate" style="max-width: 200px;">
                    {{ teacherData?.nama_kelas_ajar?.join(', ') || 'Belum ada kelas' }}
                  </span>
                </div>
              </div>
              <div class="bg-primary-subtle text-primary p-3 rounded-circle stat-icon-bubble">
                <i class="fas fa-chalkboard fs-4"></i>
              </div>
            </div>
            <div class="mt-3 pt-2 border-top small d-flex justify-content-between align-items-center">
              <span class="text-muted small">Cakupan Mengajar</span>
              <router-link to="/admin/kelas" class="text-decoration-none text-primary fw-medium">
                Rombel <i class="fas fa-arrow-right ms-1 transition-icon"></i>
              </router-link>
            </div>
          </div>
        </div>

        <!-- Mata Pelajaran Diampu -->
        <div class="col-sm-6 col-xl-3">
          <div class="card card-custom card-living-shine p-3 border-start border-4 border-warning h-100 stat-hover-card">
            <div class="d-flex align-items-start justify-content-between">
              <div>
                <div class="text-secondary small fw-semibold text-uppercase">Mata Pelajaran</div>
                <div class="fs-4 fw-bold text-dark mt-1">
                  <AnimatedNumber :value="teacherData?.total_mapel_diampu || 0" /> Mapel
                </div>
                <div class="mt-2">
                  <span class="badge bg-warning-subtle text-warning-emphasis border border-warning-subtle small py-1 px-2 text-truncate" style="max-width: 200px;">
                    {{ teacherData?.nama_mapel_diampu?.join(', ') || 'Belum ada mapel' }}
                  </span>
                </div>
              </div>
              <div class="bg-warning-subtle text-warning p-3 rounded-circle stat-icon-bubble">
                <i class="fas fa-book-open fs-4"></i>
              </div>
            </div>
            <div class="mt-3 pt-2 border-top small d-flex justify-content-between align-items-center">
              <span class="text-muted small">Kurikulum Sekolah</span>
              <router-link to="/admin/mapel" class="text-decoration-none text-warning fw-medium">
                Lihat Mapel <i class="fas fa-arrow-right ms-1 transition-icon"></i>
              </router-link>
            </div>
          </div>
        </div>

        <!-- Sesi Jadwal Hari Ini -->
        <div class="col-sm-6 col-xl-3">
          <div class="card card-custom card-living-shine p-3 border-start border-4 border-info h-100 stat-hover-card">
            <div class="d-flex align-items-start justify-content-between">
              <div>
                <div class="text-secondary small fw-semibold text-uppercase">Jadwal Hari Ini</div>
                <div class="fs-4 fw-bold text-dark mt-1">
                  <AnimatedNumber :value="teacherData?.total_jadwal_hari_ini || 0" /> Sesi
                </div>
                <div class="mt-2">
                  <span v-if="teacherData?.total_jadwal_hari_ini" class="badge bg-info-subtle text-info border border-info-subtle small py-1 px-2">
                    <i class="fas fa-clock me-1"></i>Hari {{ currentDateWIB.split(',')[0] }}
                  </span>
                  <span v-else class="badge bg-secondary-subtle text-secondary small py-1 px-2">
                    Tidak ada jadwal hari ini
                  </span>
                </div>
              </div>
              <div class="bg-info-subtle text-info p-3 rounded-circle stat-icon-bubble">
                <i class="fas fa-calendar-day fs-4"></i>
              </div>
            </div>
            <div class="mt-3 pt-2 border-top small d-flex justify-content-between align-items-center">
              <span class="text-muted small">Agenda Harian</span>
              <router-link to="/admin/jadwal" class="text-decoration-none text-info fw-medium">
                Atur Jadwal <i class="fas fa-arrow-right ms-1 transition-icon"></i>
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- Detail Kelas & Siswa (Kiri) + Highlight Jadwal Hari Ini (Kanan) -->
      <div class="row g-4 mb-4">
        <!-- Kolom Kiri: Infografis Kelas & Preview Siswa -->
        <div class="col-lg-7">
          <div class="card card-custom p-4 h-100">
            <div class="d-flex align-items-center justify-content-between mb-3 pb-2 border-bottom">
              <div>
                <h5 class="fw-bold text-dark mb-0">
                  <i class="fas fa-users-rectangle text-primary me-2"></i>Kelas & Siswa Terkait
                </h5>
                <div class="text-secondary small mt-1">
                  Rombongan belajar tempat Anda bertugas sebagai wali kelas maupun pengampu mata pelajaran
                </div>
              </div>
              <span class="badge bg-primary text-white">
                {{ teacherData?.classes?.length || 0 }} Kelas Terkait
              </span>
            </div>

            <!-- List Kartu Kelas -->
            <div v-if="teacherData?.classes && teacherData.classes.length > 0" class="d-flex flex-column gap-3">
              <div
                v-for="c in teacherData.classes"
                :key="c.kelas_id"
                class="card border p-3 rounded-3 shadow-none bg-light bg-opacity-25 class-detail-card"
              >
                <!-- Header Kelas -->
                <div class="d-flex flex-column flex-sm-row justify-content-between align-items-sm-center gap-2 mb-2 pb-2 border-bottom">
                  <div>
                    <div class="d-flex align-items-center gap-2">
                      <span class="fs-5 fw-bold text-dark">{{ c.nama_kelas }}</span>
                      <span v-if="c.is_wali" class="badge bg-success-subtle text-success border border-success-subtle small">
                        <i class="fas fa-user-shield me-1"></i>Wali Kelas
                      </span>
                      <span v-if="c.is_mengajar" class="badge bg-primary-subtle text-primary border border-primary-subtle small">
                        <i class="fas fa-chalkboard-teacher me-1"></i>Guru Pengampu
                      </span>
                    </div>
                    <div v-if="c.mapel_ajar && c.mapel_ajar.length > 0" class="text-secondary small mt-1">
                      <i class="fas fa-book me-1 text-muted"></i>Mapel Anda:
                      <strong class="text-dark">{{ c.mapel_ajar.join(', ') }}</strong>
                    </div>
                  </div>

                  <div class="d-flex align-items-center gap-2">
                    <span class="badge bg-white text-dark border px-2 py-1 small">
                      <i class="fas fa-user-graduate me-1 text-primary"></i>Total {{ c.total_siswa }} Siswa
                    </span>
                    <button
                      @click="openStudentListModal(c)"
                      class="btn btn-outline-primary btn-sm py-1 px-2 d-flex align-items-center gap-1"
                      title="Lihat seluruh daftar siswa kelas ini"
                    >
                      <i class="fas fa-eye"></i>
                      <span>Lihat Semua Siswa</span>
                    </button>
                  </div>
                </div>

                <!-- Preview Siswa (Maks 5 data) -->
                <div>
                  <div class="text-secondary small fw-semibold mb-2 d-flex justify-content-between align-items-center">
                    <span>Pratinjau Siswa (5 Pertama):</span>
                    <span v-if="c.total_siswa > 5" class="text-muted" style="font-size: 0.72rem;">
                      +{{ c.total_siswa - 5 }} siswa lainnya
                    </span>
                  </div>

                  <div v-if="c.sample_siswa && c.sample_siswa.length > 0" class="table-responsive">
                    <table class="table table-sm table-hover align-middle mb-0 bg-white rounded border small">
                      <thead class="table-light">
                        <tr>
                          <th style="width: 40px;" class="text-center">No</th>
                          <th>Nama Lengkap Siswa</th>
                          <th>Email / Kontak</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="(s, sIdx) in c.sample_siswa" :key="s.id">
                          <td class="text-center text-muted fw-semibold">{{ sIdx + 1 }}</td>
                          <td>
                            <div class="d-flex align-items-center gap-2">
                              <div class="rounded-circle bg-primary-subtle text-primary d-flex align-items-center justify-content-center fw-bold" style="width: 24px; height: 24px; font-size: 0.65rem;">
                                {{ s.nama?.substring(0, 1).toUpperCase() }}
                              </div>
                              <span class="fw-semibold text-dark">{{ s.nama }}</span>
                            </div>
                          </td>
                          <td class="text-secondary">{{ s.email || '-' }}</td>
                        </tr>
                      </tbody>
                    </table>
                  </div>

                  <div v-else class="text-muted small fst-italic py-2 text-center bg-white rounded border">
                    Belum ada siswa yang ditempatkan di rombel ini.
                  </div>
                </div>
              </div>
            </div>

            <!-- Empty State Kelas Guru -->
            <div v-else class="text-center py-5">
              <i class="fas fa-chalkboard-user text-muted fs-1 mb-2"></i>
              <h6 class="fw-bold text-dark">Belum Terhubung dengan Kelas</h6>
              <p class="text-secondary small mb-3">
                Anda belum ditugaskan sebagai wali kelas atau guru mata pelajaran di rombongan belajar manapun.
              </p>
              <span class="badge bg-light text-muted border">Silakan hubungi administrator sekolah</span>
            </div>
          </div>
        </div>

        <!-- Kolom Kanan: Highlight Jadwal Hari Ini & Timeline Mingguan -->
        <div class="col-lg-5">
          <div class="card card-custom p-4 h-100">
            <div class="d-flex align-items-center justify-content-between mb-3 pb-2 border-bottom">
              <div>
                <h5 class="fw-bold text-dark mb-0">
                  <i class="fas fa-calendar-check text-success me-2"></i>Jadwal Mengajar Hari Ini
                </h5>
                <div class="text-secondary small mt-1">Sesi tatap muka di kelas hari ini (WIB)</div>
              </div>
              <router-link to="/admin/jadwal" class="btn btn-sm btn-outline-primary" title="Kelola Jadwal">
                <i class="fas fa-cog"></i>
              </router-link>
            </div>

            <!-- List Sesi Hari Ini -->
            <div v-if="teacherData?.today_schedules && teacherData.today_schedules.length > 0" class="d-flex flex-column gap-2 mb-4">
              <div
                v-for="sch in teacherData.today_schedules"
                :key="sch.id"
                class="p-3 rounded-3 border d-flex justify-content-between align-items-center"
                :class="{
                  'bg-success-subtle border-success text-success-emphasis': sch.status === 'BERLANGSUNG',
                  'bg-primary-subtle border-primary text-primary-emphasis': sch.status === 'SELANJUTNYA',
                  'bg-white': sch.status === 'MENDATANG',
                  'bg-light text-muted': sch.status === 'SELESAI'
                }"
              >
                <div>
                  <div class="d-flex align-items-center gap-2 mb-1">
                    <span
                      class="badge"
                      :class="{
                        'bg-success text-white pulse-beacon': sch.status === 'BERLANGSUNG',
                        'bg-primary text-white': sch.status === 'SELANJUTNYA',
                        'bg-secondary text-white': sch.status === 'SELESAI',
                        'bg-light text-dark border': sch.status === 'MENDATANG'
                      }"
                    >
                      {{ sch.status }}
                    </span>
                    <span class="fw-bold font-monospace small">
                      <i class="far fa-clock me-1"></i>{{ sch.jam_mulai }} - {{ sch.jam_selesai }}
                    </span>
                  </div>
                  <div class="fw-bold text-dark">{{ sch.mata_pelajaran }}</div>
                  <div class="text-secondary small">
                    Kelas <strong class="text-dark">{{ sch.nama_kelas }}</strong> &bull; Ruang: {{ sch.ruang || '-' }}
                  </div>
                </div>
                <div class="fs-4 text-muted">
                  <i class="fas fa-chalkboard-teacher"></i>
                </div>
              </div>
            </div>

            <!-- Empty State Jadwal Hari Ini -->
            <div v-else class="text-center py-4 bg-light rounded-3 border mb-4">
              <i class="fas fa-coffee text-secondary fs-3 mb-2 d-block"></i>
              <div class="fw-bold text-dark small">Tidak Ada Jadwal Mengajar Hari Ini</div>
              <div class="text-muted small mt-1">Anda bebas dari sesi tatap muka hari ini.</div>
            </div>

            <!-- Ringkasan Jadwal Mingguan -->
            <div>
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span class="text-secondary small fw-bold text-uppercase">Beban Mengajar Mingguan</span>
                <router-link to="/admin/jadwal" class="small text-decoration-none">Lihat Jadwal &rarr;</router-link>
              </div>

              <div class="row g-2 text-center">
                <div v-for="h in ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']" :key="h" class="col-4 col-sm-2">
                  <div class="p-2 border rounded bg-light">
                    <div class="text-muted small" style="font-size: 0.7rem;">{{ h }}</div>
                    <div class="fw-bold fs-6 text-dark mt-1">
                      {{ teacherData?.weekly_schedule_summary?.[h] || 0 }}
                    </div>
                    <div class="text-secondary" style="font-size: 0.65rem;">Sesi</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ========================================================
         B. DASHBOARD KHUSUS ROLE ADMINISTRATOR
         ======================================================== -->
    <div v-else class="animate-fade-up">
      <!-- Executive KPI Cards Admin -->
      <div class="row g-3 mb-4">
        <!-- Total Siswa -->
        <div v-if="authStore.canManageSiswa" class="col-sm-6 col-xl-3">
          <div class="card card-custom card-living-shine p-3 border-start border-4 border-success h-100 stat-hover-card">
            <div class="d-flex align-items-start justify-content-between">
              <div>
                <div class="text-secondary small fw-semibold text-uppercase">Total Siswa</div>
                <div class="fs-3 fw-bold text-dark mt-1">
                  <AnimatedNumber :value="totals.total_siswa" />
                </div>
                <div class="mt-2">
                  <span v-if="completeness.siswa_tanpa_kelas > 0" class="badge bg-warning-subtle text-warning-emphasis border border-warning-subtle small py-1 px-2">
                    <i class="fas fa-exclamation-circle me-1"></i>{{ completeness.siswa_tanpa_kelas }} belum ada kelas
                  </span>
                  <span v-else class="badge bg-success-subtle text-success border border-success-subtle small py-1 px-2">
                    <i class="fas fa-check-circle me-1"></i>100% teralokasi di rombel
                  </span>
                </div>
              </div>
              <div class="bg-success-subtle text-success p-3 rounded-circle stat-icon-bubble">
                <i class="fas fa-user-graduate fs-4"></i>
              </div>
            </div>
            <div class="mt-3 pt-2 border-top small d-flex justify-content-between align-items-center">
              <span class="text-muted small">Rata-rata: {{ completeness.rata_siswa_kelas }} per kelas</span>
              <router-link to="/admin/siswa" class="text-decoration-none text-success fw-medium">
                Lihat Detail <i class="fas fa-arrow-right ms-1 transition-icon"></i>
              </router-link>
            </div>
          </div>
        </div>

        <!-- Total Guru -->
        <div v-if="authStore.canManageGuru" class="col-sm-6 col-xl-3">
          <div class="card card-custom card-living-shine p-3 border-start border-4 border-primary h-100 stat-hover-card">
            <div class="d-flex align-items-start justify-content-between">
              <div>
                <div class="text-secondary small fw-semibold text-uppercase">Tenaga Pendidik</div>
                <div class="fs-3 fw-bold text-dark mt-1">
                  <AnimatedNumber :value="totals.total_guru" />
                </div>
                <div class="mt-2">
                  <span class="badge bg-primary-subtle text-primary border border-primary-subtle small py-1 px-2">
                    <i class="fas fa-chalkboard me-1"></i>{{ completeness.guru_pengampu }} mengajar &bull; {{ completeness.guru_wali_kelas }} wali
                  </span>
                </div>
              </div>
              <div class="bg-primary-subtle text-primary p-3 rounded-circle stat-icon-bubble">
                <i class="fas fa-chalkboard-teacher fs-4"></i>
              </div>
            </div>
            <div class="mt-3 pt-2 border-top small d-flex justify-content-between align-items-center">
              <span class="text-muted small" v-if="completeness.guru_tanpa_akun > 0">
                <i class="fas fa-info-circle text-warning me-1"></i>{{ completeness.guru_tanpa_akun }} tanpa akun
              </span>
              <span class="text-muted small" v-else>
                <i class="fas fa-shield-alt text-success me-1"></i>Semua berakun
              </span>
              <router-link to="/admin/guru" class="text-decoration-none text-primary fw-medium">
                Lihat Detail <i class="fas fa-arrow-right ms-1 transition-icon"></i>
              </router-link>
            </div>
          </div>
        </div>

        <!-- Total Rombel Kelas -->
        <div class="col-sm-6 col-xl-3">
          <div class="card card-custom card-living-shine p-3 border-start border-4 border-info h-100 stat-hover-card">
            <div class="d-flex align-items-start justify-content-between">
              <div>
                <div class="text-secondary small fw-semibold text-uppercase">Rombongan Belajar</div>
                <div class="fs-3 fw-bold text-dark mt-1">
                  <AnimatedNumber :value="totals.total_kelas" />
                </div>
                <div class="mt-2">
                  <span v-if="completeness.kelas_tanpa_wali > 0" class="badge bg-danger-subtle text-danger border border-danger-subtle small py-1 px-2">
                    <i class="fas fa-user-times me-1"></i>{{ completeness.kelas_tanpa_wali }} rombel tanpa wali
                  </span>
                  <span v-else class="badge bg-info-subtle text-info border border-info-subtle small py-1 px-2">
                    <i class="fas fa-user-check me-1"></i>Semua berwali kelas
                  </span>
                </div>
              </div>
              <div class="bg-info-subtle text-info p-3 rounded-circle stat-icon-bubble">
                <i class="fas fa-school fs-4"></i>
              </div>
            </div>
            <div class="mt-3 pt-2 border-top small d-flex justify-content-between align-items-center">
              <span class="text-muted small">{{ completeness.siswa_dengan_kelas }} siswa aktif</span>
              <router-link to="/admin/kelas" class="text-decoration-none text-info fw-medium">
                Atur Rombel <i class="fas fa-arrow-right ms-1 transition-icon"></i>
              </router-link>
            </div>
          </div>
        </div>

        <!-- Total Mata Pelajaran -->
        <div class="col-sm-6 col-xl-3">
          <div class="card card-custom card-living-shine p-3 border-start border-4 border-warning h-100 stat-hover-card">
            <div class="d-flex align-items-start justify-content-between">
              <div>
                <div class="text-secondary small fw-semibold text-uppercase">Mata Pelajaran</div>
                <div class="fs-3 fw-bold text-dark mt-1">
                  <AnimatedNumber :value="totals.total_mapel" />
                </div>
                <div class="mt-2">
                  <span v-if="completeness.mapel_tanpa_guru > 0" class="badge bg-warning-subtle text-warning-emphasis border border-warning-subtle small py-1 px-2">
                    <i class="fas fa-exclamation-triangle me-1"></i>{{ completeness.mapel_tanpa_guru }} belum ada guru
                  </span>
                  <span v-else class="badge bg-success-subtle text-success border border-success-subtle small py-1 px-2">
                    <i class="fas fa-check me-1"></i>Tersedia pengampu
                  </span>
                </div>
              </div>
              <div class="bg-warning-subtle text-warning p-3 rounded-circle stat-icon-bubble">
                <i class="fas fa-book-open fs-4"></i>
              </div>
            </div>
            <div class="mt-3 pt-2 border-top small d-flex justify-content-between align-items-center">
              <span class="text-muted small">{{ completeness.mapel_dengan_guru }} teralokasi guru</span>
              <router-link to="/admin/mapel" class="text-decoration-none text-warning fw-medium">
                Kelola Mapel <i class="fas fa-arrow-right ms-1 transition-icon"></i>
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- Section: Data Perlu Perhatian Admin -->
      <div class="mb-4" v-if="attentionItems.length > 0">
        <div class="card card-custom border-warning border-opacity-50 p-3 p-md-4 bg-warning-subtle bg-opacity-10 shadow-sm">
          <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-3">
            <div class="d-flex align-items-center gap-2">
              <div class="bg-warning text-dark p-2 rounded-3">
                <i class="fas fa-bell"></i>
              </div>
              <div>
                <h6 class="fw-bold text-dark mb-0">Perlu Tindakan Lanjut Admin</h6>
                <span class="text-secondary small">Terdapat {{ attentionItems.length }} item data yang membutuhkan pembaruan atau alokasi</span>
              </div>
            </div>
          </div>

          <div class="row g-3">
            <div v-for="item in attentionItems" :key="item.id" class="col-md-6 col-lg-4">
              <div class="card h-100 border p-3 bg-white shadow-none attention-card">
                <div class="d-flex align-items-start justify-content-between gap-2 mb-2">
                  <span
                    class="badge small"
                    :class="{
                      'bg-danger text-white': item.level === 'danger',
                      'bg-warning text-dark': item.level === 'warning',
                      'bg-info text-dark': item.level === 'info',
                    }"
                  >
                    <i class="fas fa-circle-exclamation me-1"></i>{{ item.count }} Data
                  </span>
                  <router-link :to="item.route" class="btn btn-outline-primary btn-sm py-0 px-2 small">
                    {{ item.action }} &rarr;
                  </router-link>
                </div>
                <div class="fw-bold text-dark small mb-1">{{ item.title }}</div>
                <div class="text-secondary small">{{ item.message }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Alert Kualitas Data Prima Jika Tidak Ada Pending Item -->
      <div class="mb-4" v-else>
        <div class="card card-custom p-3 border-0 bg-success-subtle bg-opacity-25 shadow-sm">
          <div class="d-flex align-items-center justify-content-between flex-wrap gap-2">
            <div class="d-flex align-items-center gap-3">
              <div class="bg-success text-white p-2 px-3 rounded-circle">
                <i class="fas fa-check-double"></i>
              </div>
              <div>
                <div class="fw-bold text-dark small">Kualitas Data Sekolah Sempurna</div>
                <div class="text-secondary small">Seluruh siswa, guru, rombel kelas, dan mata pelajaran sudah teralokasi dengan rapi.</div>
              </div>
            </div>
            <span class="badge bg-success text-white px-3 py-2 fw-medium small">
              Integritas 100%
            </span>
          </div>
        </div>
      </div>

      <!-- Row: Distribusi Siswa & Beban Operasional Sistem -->
      <div class="row g-4 mb-4">
        <!-- Kolom Kiri: Distribusi Siswa per Rombel Kelas -->
        <div class="col-lg-7">
          <div class="card card-custom p-4 h-100">
            <div class="d-flex align-items-center justify-content-between mb-3 pb-2 border-bottom">
              <div>
                <h5 class="fw-bold text-dark mb-0">
                  <i class="fas fa-chart-bar text-primary me-2"></i>Distribusi Siswa per Rombel
                </h5>
                <div class="text-secondary small mt-1">Kepadatan rombongan belajar dan wali kelas</div>
              </div>
              <router-link to="/admin/kelas" class="btn btn-sm btn-outline-primary">
                Kelola Rombel
              </router-link>
            </div>

            <!-- List Kelas dengan Bar Progres Visual -->
            <div v-if="distribution.length > 0" class="d-flex flex-column gap-3">
              <div v-for="item in distribution" :key="item.id" class="p-2 rounded-2 bg-light bg-opacity-50 border">
                <div class="d-flex justify-content-between align-items-center mb-1">
                  <div class="d-flex align-items-center gap-2">
                    <span class="badge bg-dark-subtle text-dark fw-bold px-2 py-1">{{ item.kelas }}</span>
                    <span class="text-secondary small">
                      <i class="fas fa-user-tie me-1 text-muted"></i>
                      {{ item.wali_guru ? item.wali_guru : 'Belum ada wali kelas' }}
                    </span>
                  </div>
                  <div class="fw-bold text-dark small">
                    {{ item.total_siswa }} Siswa
                  </div>
                </div>
                <!-- Progress Bar -->
                <div class="progress" style="height: 7px;">
                  <div
                    class="progress-bar progress-bar-striped rounded"
                    :class="getBarColor(item.total_siswa, maxClassStudents)"
                    role="progressbar"
                    :style="{ width: getPercentage(item.total_siswa, maxClassStudents) + '%' }"
                  ></div>
                </div>
              </div>
            </div>

            <!-- Empty State Kelas -->
            <div v-else class="text-center py-5">
              <i class="fas fa-school-circle-xmark text-muted fs-1 mb-2"></i>
              <p class="text-secondary small mb-3">Belum ada data rombongan belajar kelas yang terdaftar.</p>
              <router-link to="/admin/kelas" class="btn btn-sm btn-primary">
                <i class="fas fa-plus me-1"></i>Buat Kelas Baru
              </router-link>
            </div>
          </div>
        </div>

        <!-- Kolom Kanan: Operasional Sistem Hari Ini -->
        <div class="col-lg-5">
          <div class="card card-custom p-4 h-100">
            <div class="d-flex align-items-center justify-content-between mb-3 pb-2 border-bottom">
              <div>
                <h5 class="fw-bold text-dark mb-0">
                  <i class="fas fa-server text-info me-2"></i>Operasional Sistem Hari Ini
                </h5>
                <div class="text-secondary small mt-1">Audit transaksi sistem & komposisi akun</div>
              </div>
              <router-link to="/admin/logs" class="btn btn-sm btn-outline-secondary" title="Audit Log">
                <i class="fas fa-clock-rotate-left"></i>
              </router-link>
            </div>

            <!-- Ringkasan Transaksi Hari Ini -->
            <div class="mb-4">
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span class="text-secondary small fw-medium">Total Aktivitas Hari Ini</span>
                <span class="badge bg-secondary text-white fw-bold">{{ todayActivity.total_hari_ini }} Event</span>
              </div>

              <!-- Bar Sukses vs Gagal -->
              <div class="progress mb-2" style="height: 10px;">
                <div
                  class="progress-bar bg-success"
                  role="progressbar"
                  :style="{ width: successRate + '%' }"
                  :title="'Berhasil: ' + todayActivity.berhasil_hari_ini"
                ></div>
                <div
                  class="progress-bar bg-danger"
                  role="progressbar"
                  :style="{ width: (100 - successRate) + '%' }"
                  :title="'Gagal: ' + todayActivity.gagal_hari_ini"
                ></div>
              </div>

              <div class="row g-2 text-center pt-2">
                <div class="col-4">
                  <div class="p-2 border rounded bg-light">
                    <div class="text-success fw-bold fs-6">{{ todayActivity.berhasil_hari_ini }}</div>
                    <div class="text-muted small" style="font-size: 0.75rem;">Berhasil</div>
                  </div>
                </div>
                <div class="col-4">
                  <div class="p-2 border rounded bg-light">
                    <div class="text-danger fw-bold fs-6">{{ todayActivity.gagal_hari_ini }}</div>
                    <div class="text-muted small" style="font-size: 0.75rem;">Gagal/Ditolak</div>
                  </div>
                </div>
                <div class="col-4">
                  <div class="p-2 border rounded bg-light">
                    <div class="text-primary fw-bold fs-6">{{ todayActivity.user_aktif_hari_ini }}</div>
                    <div class="text-muted small" style="font-size: 0.75rem;">Pengguna Aktif</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Komposisi Akun Pengguna -->
            <div>
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span class="text-secondary small fw-medium">Komposisi Akun Pengguna</span>
                <router-link to="/admin/users" class="small text-decoration-none">Kelola Akun &rarr;</router-link>
              </div>
              <div class="d-flex flex-wrap gap-2">
                <div
                  v-for="roleItem in userRoles"
                  :key="roleItem.role"
                  class="p-2 px-3 border rounded-3 bg-light d-flex align-items-center gap-2 flex-grow-1"
                >
                  <div
                    class="rounded-circle p-2 d-flex align-items-center justify-content-center"
                    :class="roleItem.role === 'admin' ? 'bg-danger text-white' : (roleItem.role === 'guru' ? 'bg-primary text-white' : 'bg-secondary text-white')"
                    style="width: 28px; height: 28px;"
                  >
                    <i :class="roleItem.role === 'admin' ? 'fas fa-shield-halved' : (roleItem.role === 'guru' ? 'fas fa-chalkboard-user' : 'fas fa-user')" style="font-size: 0.75rem;"></i>
                  </div>
                  <div>
                    <div class="text-uppercase fw-bold text-dark small" style="font-size: 0.75rem;">{{ roleItem.role }}</div>
                    <div class="fw-bold text-secondary small">{{ roleItem.total }} Akun</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Section: Riwayat Aktivitas Terbaru (Khusus Admin) -->
      <div class="card card-custom p-4 mb-4">
        <div class="d-flex flex-column flex-sm-row justify-content-between align-items-sm-center gap-2 mb-3 pb-2 border-bottom">
          <div>
            <h5 class="fw-bold text-dark mb-0">
              <i class="fas fa-list-check text-success me-2"></i>Aktivitas Sistem Terbaru
            </h5>
            <div class="text-secondary small mt-1">Audit jejak operasional realtime pada modul aplikasi</div>
          </div>
          <router-link to="/admin/logs" class="btn btn-outline-secondary btn-sm">
            Buka Log Selengkapnya &rarr;
          </router-link>
        </div>

        <div class="table-responsive" v-if="recentActivities.length > 0">
          <table class="table table-hover align-middle mb-0 small">
            <thead class="table-light">
              <tr>
                <th scope="col" style="width: 140px;">Waktu</th>
                <th scope="col" style="width: 160px;">Pengguna</th>
                <th scope="col" style="width: 100px;">Aksi</th>
                <th scope="col" style="width: 130px;">Modul</th>
                <th scope="col">Keterangan</th>
                <th scope="col" class="text-center" style="width: 90px;">Hasil</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="act in recentActivities" :key="act.id">
                <td class="text-nowrap text-secondary font-monospace" style="font-size: 0.8rem;">
                  {{ formatDateTime(act.timestamp) }}
                </td>
                <td>
                  <div class="d-flex align-items-center gap-2">
                    <div class="rounded-circle bg-light border d-flex align-items-center justify-content-center text-secondary small fw-bold" style="width: 26px; height: 26px; font-size: 0.7rem;">
                      {{ act.username?.substring(0, 2).toUpperCase() }}
                    </div>
                    <div>
                      <div class="fw-semibold text-dark">{{ act.username }}</div>
                      <span class="text-muted" style="font-size: 0.7rem;">{{ act.role }}</span>
                    </div>
                  </div>
                </td>
                <td>
                  <span
                    class="badge"
                    :class="{
                      'bg-success-subtle text-success border border-success-subtle': act.activity_type === 'CREATE',
                      'bg-info-subtle text-info border border-info-subtle': act.activity_type === 'READ',
                      'bg-warning-subtle text-warning-emphasis border border-warning-subtle': act.activity_type === 'UPDATE',
                      'bg-danger-subtle text-danger border border-danger-subtle': act.activity_type === 'DELETE',
                      'bg-primary-subtle text-primary border border-primary-subtle': act.activity_type === 'LOGIN',
                    }"
                  >
                    {{ act.activity_type }}
                  </span>
                </td>
                <td class="fw-medium text-dark">{{ act.module }}</td>
                <td class="text-secondary text-truncate" style="max-width: 320px;" :title="act.activity_description">
                  {{ act.activity_description }}
                </td>
                <td class="text-center">
                  <span
                    class="badge rounded-pill"
                    :class="act.result === 'Success' ? 'bg-success text-white' : 'bg-danger text-white'"
                  >
                    {{ act.result === 'Success' ? 'Sukses' : 'Gagal' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-else class="text-center py-4 text-muted small">
          Belum ada log transaksi tercatat.
        </div>
      </div>
    </div>

    <!-- ========================================================
         MODAL LIHAT SEMUA SISWA (SOLID WHITE BACKGROUND & TELEPORT)
         ======================================================== -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showStudentModal" class="app-modal-overlay" tabindex="-1" @click.self="closeStudentModal">
          <div class="modal-dialog modal-dialog-centered modal-lg">
            <div class="modal-content border-0 shadow-lg">
              <div class="modal-header bg-light">
                <div>
                  <h5 class="modal-title fw-bold text-dark d-flex align-items-center gap-2 mb-0">
                    <i class="fas fa-users text-primary"></i>
                    <span>Daftar Siswa Kelas {{ activeClassForModal?.nama_kelas }}</span>
                  </h5>
                  <span class="text-secondary small">
                    Total {{ activeClassForModal?.total_siswa || 0 }} siswa terdaftar di rombongan belajar ini
                  </span>
                </div>
                <button type="button" class="btn-close" @click="closeStudentModal"></button>
              </div>

              <div class="modal-body p-4">
                <!-- Search input dalam modal -->
                <div class="mb-3">
                  <div class="input-group">
                    <span class="input-group-text bg-light text-muted"><i class="fas fa-search"></i></span>
                    <input
                      v-model="modalStudentSearch"
                      type="text"
                      class="form-control"
                      placeholder="Cari siswa berdasarkan nama atau email..."
                    />
                    <button v-if="modalStudentSearch" class="btn btn-link text-secondary p-0 px-2" @click="modalStudentSearch = ''" title="Bersihkan">
                      <i class="fas fa-times-circle"></i>
                    </button>
                  </div>
                </div>

                <!-- Tabel Siswa Lengkap -->
                <div class="table-responsive" style="max-height: 400px; overflow-y: auto;">
                  <table class="table table-hover align-middle mb-0">
                    <thead class="table-light sticky-top">
                      <tr>
                        <th style="width: 50px;" class="text-center">No</th>
                        <th>Nama Siswa</th>
                        <th>Email / Kontak</th>
                        <th class="text-center" style="width: 100px;">Status</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-if="filteredModalStudents.length === 0">
                        <td colspan="4" class="text-center py-4 text-muted">
                          Tidak ditemukan siswa yang cocok dengan pencarian.
                        </td>
                      </tr>
                      <tr v-for="(s, idx) in filteredModalStudents" :key="s.id">
                        <td class="text-center text-muted fw-semibold">{{ idx + 1 }}</td>
                        <td>
                          <div class="d-flex align-items-center gap-2">
                            <div class="rounded-circle bg-primary-subtle text-primary d-flex align-items-center justify-content-center fw-bold" style="width: 32px; height: 32px; font-size: 0.8rem;">
                              {{ s.nama?.substring(0, 1).toUpperCase() }}
                            </div>
                            <span class="fw-bold text-dark">{{ s.nama }}</span>
                          </div>
                        </td>
                        <td class="text-secondary">{{ s.email || '-' }}</td>
                        <td class="text-center">
                          <span class="badge bg-success-subtle text-success">Aktif</span>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              <div class="modal-footer bg-light">
                <button type="button" class="btn btn-secondary btn-sm" @click="closeStudentModal">
                  Tutup
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Modal Kalender Indonesia -->
    <IndonesiaCalendarModal
      :is-open="showCalendarModal"
      @close="showCalendarModal = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { dashboardApi } from '@/api/dashboardApi'
import IndonesiaCalendarModal from '@/components/IndonesiaCalendarModal.vue'
import AnimatedNumber from '@/components/AnimatedNumber.vue'

const authStore = useAuthStore()

// State Modal Kalender Indonesia
const showCalendarModal = ref(false)
const openCalendarModal = () => {
  showCalendarModal.value = true
}

// State Modal Siswa Kelas (Guru)
const showStudentModal = ref(false)
const activeClassForModal = ref(null)
const modalStudentSearch = ref('')

const openStudentListModal = (classDetail) => {
  activeClassForModal.value = classDetail
  modalStudentSearch.value = ''
  showStudentModal.value = true
}

const closeStudentModal = () => {
  showStudentModal.value = false
  activeClassForModal.value = null
  modalStudentSearch.value = ''
}

const filteredModalStudents = computed(() => {
  if (!activeClassForModal.value?.all_siswa) return []
  if (!modalStudentSearch.value) return activeClassForModal.value.all_siswa
  const q = modalStudentSearch.value.toLowerCase()
  return activeClassForModal.value.all_siswa.filter(
    (s) => (s.nama || '').toLowerCase().includes(q) || (s.email || '').toLowerCase().includes(q)
  )
})

// Realtime WIB Clock
const currentTimeWIB = ref('--:--:--')
const currentDateWIB = ref('Memuat tanggal...')
let clockTimer = null

const updateClockWIB = () => {
  const now = new Date()
  currentTimeWIB.value = now.toLocaleTimeString('id-ID', {
    timeZone: 'Asia/Jakarta',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  })
  currentDateWIB.value = now.toLocaleDateString('id-ID', {
    timeZone: 'Asia/Jakarta',
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

// Ucapan greeting dinamis
const greetingText = computed(() => {
  const hour = new Date().getHours()
  if (hour >= 4 && hour < 11) return 'Selamat Pagi'
  if (hour >= 11 && hour < 15) return 'Selamat Siang'
  if (hour >= 15 && hour < 19) return 'Selamat Sore'
  return 'Selamat Malam'
})

// State Dashboard Data Global
const totals = ref({
  total_siswa: 0,
  total_guru: 0,
  total_kelas: 0,
  total_mapel: 0,
  total_users: 0,
})

const completeness = ref({
  siswa_dengan_kelas: 0,
  siswa_tanpa_kelas: 0,
  persen_siswa_kelas: 0,
  guru_wali_kelas: 0,
  guru_pengampu: 0,
  guru_tanpa_tugas: 0,
  guru_tanpa_akun: 0,
  kelas_dengan_wali: 0,
  kelas_tanpa_wali: 0,
  rata_siswa_kelas: 0,
  mapel_dengan_guru: 0,
  mapel_tanpa_guru: 0,
  mapel_tanpa_kelas: 0,
})

const distribution = ref([])
const userRoles = ref([])
const todayActivity = ref({
  total_hari_ini: 0,
  berhasil_hari_ini: 0,
  gagal_hari_ini: 0,
  user_aktif_hari_ini: 0,
})
const recentActivities = ref([])
const attentionItems = ref([])
const teacherData = ref(null)

const loading = ref(false)
const error = ref('')

// Hitung persentase keberhasilan aktivitas sistem
const successRate = computed(() => {
  const total = todayActivity.value.total_hari_ini
  if (!total || total === 0) return 100
  const success = todayActivity.value.berhasil_hari_ini || 0
  return Math.min(100, Math.max(0, Math.round((success / total) * 100)))
})

// Hitung kapasitas tertinggi kelas untuk skala visual bar
const maxClassStudents = computed(() => {
  if (!distribution.value || distribution.value.length === 0) return 30
  const max = Math.max(...distribution.value.map((d) => d.total_siswa || 0))
  return max > 0 ? max : 30
})

const getPercentage = (value, max) => {
  if (!max || max <= 0) return 10
  const pct = Math.round(((value || 0) / max) * 100)
  return Math.max(10, Math.min(100, pct))
}

const getBarColor = (count, max) => {
  const pct = (count / (max || 1)) * 100
  if (pct >= 80) return 'bg-success'
  if (pct >= 40) return 'bg-primary'
  return 'bg-info'
}

const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  try {
    const d = new Date(dateStr)
    return d.toLocaleString('id-ID', {
      timeZone: 'Asia/Jakarta',
      day: '2-digit',
      month: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    }).replace(',', '')
  } catch {
    return dateStr
  }
}

const loadDashboardData = async () => {
  loading.value = true
  error.value = ''

  try {
    const res = await dashboardApi.getSummary()
    const data = res?.data?.data || {}

    if (data.totals) totals.value = data.totals
    if (data.completeness) completeness.value = data.completeness
    distribution.value = data.distribution || []
    userRoles.value = data.user_roles || []
    if (data.today_activity) todayActivity.value = data.today_activity
    recentActivities.value = data.recent_activities || []
    attentionItems.value = data.attention_items || []
    if (data.teacher_data) teacherData.value = data.teacher_data
  } catch (err) {
    error.value = 'Gagal memuat ringkasan data dasbor. Silakan periksa koneksi server.'
    console.error('Error loading dashboard summary:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  updateClockWIB()
  clockTimer = setInterval(updateClockWIB, 1000)
  loadDashboardData()
})

onBeforeUnmount(() => {
  if (clockTimer) {
    clearInterval(clockTimer)
    clockTimer = null
  }
})
</script>

<style scoped>
.bg-gradient-header {
  background: linear-gradient(135deg, #f8fafc 0%, #edf2f7 100%);
  border: 1px solid #e2e8f0 !important;
}

.clock-card {
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.3s ease;
}

.clock-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.06) !important;
}

.clock-icon-box {
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.clock-card:hover .clock-icon-box {
  transform: scale(1.1);
}

.clock-time {
  letter-spacing: 0.04em;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.calendar-trigger {
  padding: 3px 6px;
  border-radius: 6px;
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  display: inline-flex;
  align-items: center;
}

.calendar-trigger:hover {
  background-color: #f1f5f9;
  color: #0d6efd !important;
  transform: translateX(2px);
}

.calendar-trigger:hover .calendar-link-text {
  text-decoration: underline;
  color: #0d6efd;
}

/* Stat Cards Living Dynamics */
.stat-hover-card {
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.3s ease;
}

.stat-hover-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 14px 28px rgba(0, 0, 0, 0.08) !important;
}

.stat-icon-bubble {
  transition: transform 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.stat-hover-card:hover .stat-icon-bubble {
  transform: scale(1.12) rotate(4deg);
}

.transition-icon {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  display: inline-block;
}

.stat-hover-card:hover .transition-icon {
  transform: translateX(4px);
}

.attention-card {
  transition: transform 0.25s ease, box-shadow 0.25s ease;
  border-left: 3px solid #ffc107 !important;
}

.attention-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.05) !important;
}

.class-detail-card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.class-detail-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.04) !important;
}
</style>
