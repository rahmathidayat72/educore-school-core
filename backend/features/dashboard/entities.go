package dashboard

import "time"

type (
	// DashboardTotals menyimpan angka akumulasi total data sistem
	DashboardTotals struct {
		TotalSiswa int `json:"total_siswa"`
		TotalGuru  int `json:"total_guru"`
		TotalKelas int `json:"total_kelas"`
		TotalMapel int `json:"total_mapel"`
		TotalUsers int `json:"total_users"`
	}

	// DashboardCompleteness menyimpan indikator kualitas & kelengkapan data
	DashboardCompleteness struct {
		SiswaDenganKelas int     `json:"siswa_dengan_kelas"`
		SiswaTanpaKelas  int     `json:"siswa_tanpa_kelas"`
		PersenSiswaKelas float64 `json:"persen_siswa_kelas"`

		GuruWaliKelas  int `json:"guru_wali_kelas"`
		GuruPengampu   int `json:"guru_pengampu"`
		GuruTanpaTugas int `json:"guru_tanpa_tugas"`
		GuruTanpaAkun  int `json:"guru_tanpa_akun"`

		KelasDenganWali int     `json:"kelas_dengan_wali"`
		KelasTanpaWali  int     `json:"kelas_tanpa_wali"`
		RataSiswaKelas  float64 `json:"rata_siswa_kelas"`

		MapelDenganGuru int `json:"mapel_dengan_guru"`
		MapelTanpaGuru  int `json:"mapel_tanpa_guru"`
		MapelTanpaKelas int `json:"mapel_tanpa_kelas"`
	}

	// ClassStudentCount merepresentasikan data distribusi jumlah siswa per rombel
	ClassStudentCount struct {
		ID         string `json:"id"`
		Kelas      string `json:"kelas"`
		WaliGuru   string `json:"wali_guru"`
		TotalSiswa int    `json:"total_siswa"`
	}

	// UserRoleCount merepresentasikan jumlah akun berdasarkan role
	UserRoleCount struct {
		Role  string `json:"role"`
		Total int    `json:"total"`
	}

	// TodayActivity merangkum aktivitas operasional hari ini
	TodayActivity struct {
		TotalHariIni     int `json:"total_hari_ini"`
		BerhasilHariIni   int `json:"berhasil_hari_ini"`
		GagalHariIni      int `json:"gagal_hari_ini"`
		UserAktifHariIni int `json:"user_aktif_hari_ini"`
	}

	// RecentActivity menampilkan 5 aktivitas terakhir untuk ringkasan cepat
	RecentActivity struct {
		ID                  int       `json:"id"`
		Timestamp           time.Time `json:"timestamp"`
		Username            string    `json:"username"`
		Role                string    `json:"role"`
		ActivityType        string    `json:"activity_type"`
		ActivityDescription string    `json:"activity_description"`
		Module              string    `json:"module"`
		Result              string    `json:"result"`
	}

	// AttentionItem merepresentasikan data yang butuh perhatian admin
	AttentionItem struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Count   int    `json:"count"`
		Message string `json:"message"`
		Level   string `json:"level"` // "danger", "warning", "info"
		Route   string `json:"route"`
		Action  string `json:"action"`
	}

	// TeacherClassDetail memuat detail kelas wali atau kelas ajar beserta daftar siswanya
	TeacherStudentBrief struct {
		ID    string `json:"id"`
		Nama  string `json:"nama"`
		Email string `json:"email"`
	}

	TeacherClassDetail struct {
		KelasID      string                `json:"kelas_id"`
		NamaKelas    string                `json:"nama_kelas"`
		IsWali       bool                  `json:"is_wali"`
		IsMengajar   bool                  `json:"is_mengajar"`
		TotalSiswa   int                   `json:"total_siswa"`
		MapelAjar    []string              `json:"mapel_ajar"` // Mapel yang diajar guru ini di kelas tersebut
		SampleSiswa  []TeacherStudentBrief `json:"sample_siswa"`
		AllSiswa     []TeacherStudentBrief `json:"all_siswa"`
	}

	// TeacherScheduleItem merepresentasikan slot jadwal mengajar guru
	TeacherScheduleItem struct {
		ID            string `json:"id"`
		MataPelajaran string `json:"mata_pelajaran"`
		NamaKelas     string `json:"nama_kelas"`
		Hari          int    `json:"hari"`
		NamaHari      string `json:"nama_hari"`
		JamMulai      string `json:"jam_mulai"`
		JamSelesai    string `json:"jam_selesai"`
		Ruang         string `json:"ruang"`
		Status        string `json:"status"` // "BERLANGSUNG", "SELANJUTNYA", "SELESAI", "MENDATANG"
	}

	// TeacherDashboardData menyimpan profil dan data personal khusus role Guru
	TeacherDashboardData struct {
		GuruID             string                `json:"guru_id"`
		NamaGuru           string                `json:"nama_guru"`
		Email              string                `json:"email"`
		Alamat             string                `json:"alamat"`
		TotalKelasWali     int                   `json:"total_kelas_wali"`
		TotalKelasAjar     int                   `json:"total_kelas_ajar"`
		TotalMapelDiampu   int                   `json:"total_mapel_diampu"`
		TotalSiswaWali     int                   `json:"total_siswa_wali"`
		TotalSiswaAjar     int                   `json:"total_siswa_ajar"`
		TotalJadwalHariIni int                   `json:"total_jadwal_hari_ini"`
		NamaKelasWali      string                `json:"nama_kelas_wali"`
		NamaKelasAjar      []string              `json:"nama_kelas_ajar"`
		NamaMapelDiampu    []string              `json:"nama_mapel_diampu"`
		Classes            []TeacherClassDetail  `json:"classes"`
		TodaySchedules     []TeacherScheduleItem `json:"today_schedules"`
		NextSchedule       *TeacherScheduleItem  `json:"next_schedule"`
		WeeklyScheduleSummary map[string]int     `json:"weekly_schedule_summary"`
	}

	// DashboardSummary adalah paket lengkap ringkasan dasbor
	DashboardSummary struct {
		Role             string                `json:"role"`
		Totals           DashboardTotals       `json:"totals"`
		Completeness     DashboardCompleteness `json:"completeness"`
		Distribution     []ClassStudentCount   `json:"distribution"`
		UserRoles        []UserRoleCount       `json:"user_roles"`
		TodayActivity    TodayActivity         `json:"today_activity"`
		RecentActivities []RecentActivity      `json:"recent_activities"`
		AttentionItems   []AttentionItem       `json:"attention_items"`
		TeacherData      *TeacherDashboardData `json:"teacher_data,omitempty"`
	}

	DataDashboardInterface interface {
		GetSummary(role string, userID string) (*DashboardSummary, error)
	}

	ServiceDashboardInterface interface {
		GetSummary(role string, userID string) (*DashboardSummary, error)
	}
)
