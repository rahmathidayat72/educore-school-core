package model

import (
	"context"
	"fmt"
	"go_rest_native_sekolah/features/dashboard"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardModel struct {
	db *pgxpool.Pool
}

func NewDashboardModel(db *pgxpool.Pool) dashboard.DataDashboardInterface {
	if db == nil {
		panic("dashboard model: nil database")
	}
	return &DashboardModel{db: db}
}

func (dm *DashboardModel) GetSummary(role string, userID string) (*dashboard.DashboardSummary, error) {
	ctx := context.Background()
	summary := &dashboard.DashboardSummary{
		Distribution:     make([]dashboard.ClassStudentCount, 0),
		UserRoles:        make([]dashboard.UserRoleCount, 0),
		RecentActivities: make([]dashboard.RecentActivity, 0),
		AttentionItems:   make([]dashboard.AttentionItem, 0),
	}

	// 1. Hitung Totals dasar
	_ = dm.db.QueryRow(ctx, "SELECT COUNT(*) FROM siswa WHERE delete_at IS NULL").Scan(&summary.Totals.TotalSiswa)
	_ = dm.db.QueryRow(ctx, "SELECT COUNT(*) FROM guru WHERE delete_at IS NULL").Scan(&summary.Totals.TotalGuru)
	_ = dm.db.QueryRow(ctx, "SELECT COUNT(*) FROM kelas WHERE delete_at IS NULL").Scan(&summary.Totals.TotalKelas)
	_ = dm.db.QueryRow(ctx, "SELECT COUNT(*) FROM mata_pelajaran WHERE delete_at IS NULL").Scan(&summary.Totals.TotalMapel)

	// Hitung data master pendukung (Unit Sekolah, Wali Murid, Tahun Ajaran Aktif)
	_ = dm.db.QueryRow(ctx, "SELECT COUNT(*) FROM school_units WHERE delete_at IS NULL").Scan(&summary.Totals.TotalUnit)
	_ = dm.db.QueryRow(ctx, "SELECT COUNT(*) FROM guardians WHERE delete_at IS NULL").Scan(&summary.Totals.TotalWali)
	_ = dm.db.QueryRow(ctx, "SELECT COUNT(*) FROM academic_years WHERE delete_at IS NULL AND is_active = true").Scan(&summary.Totals.TotalTahunAktif)

	if role == "admin" || role == "super_admin" {
		_ = dm.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE delete_at IS NULL").Scan(&summary.Totals.TotalUsers)
	}

	// 2. Metrik Kelengkapan Data (Siswa)
	_ = dm.db.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM siswa 
		WHERE delete_at IS NULL 
		  AND kelas_id IS NOT NULL 
		  AND TRIM(kelas_id) != ''
	`).Scan(&summary.Completeness.SiswaDenganKelas)

	summary.Completeness.SiswaTanpaKelas = summary.Totals.TotalSiswa - summary.Completeness.SiswaDenganKelas
	if summary.Totals.TotalSiswa > 0 {
		pct := (float64(summary.Completeness.SiswaDenganKelas) / float64(summary.Totals.TotalSiswa)) * 100.0
		summary.Completeness.PersenSiswaKelas = math.Round(pct*10) / 10
	} else {
		summary.Completeness.PersenSiswaKelas = 0
	}

	// 3. Metrik Kelengkapan Data (Kelas)
	_ = dm.db.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM kelas 
		WHERE delete_at IS NULL 
		  AND id_guru IS NOT NULL 
		  AND TRIM(id_guru) != ''
	`).Scan(&summary.Completeness.KelasDenganWali)

	summary.Completeness.KelasTanpaWali = summary.Totals.TotalKelas - summary.Completeness.KelasDenganWali

	if summary.Totals.TotalKelas > 0 {
		avg := float64(summary.Totals.TotalSiswa) / float64(summary.Totals.TotalKelas)
		summary.Completeness.RataSiswaKelas = math.Round(avg*10) / 10
	} else {
		summary.Completeness.RataSiswaKelas = 0
	}

	// 4. Metrik Kelengkapan Data (Guru)
	_ = dm.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT id_guru) 
		FROM kelas 
		WHERE delete_at IS NULL 
		  AND id_guru IS NOT NULL 
		  AND TRIM(id_guru) != ''
	`).Scan(&summary.Completeness.GuruWaliKelas)

	_ = dm.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT id_guru) 
		FROM mata_pelajaran 
		WHERE delete_at IS NULL 
		  AND id_guru IS NOT NULL 
		  AND TRIM(id_guru) != ''
	`).Scan(&summary.Completeness.GuruPengampu)

	var guruBertugas int
	_ = dm.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT g.id)
		FROM guru g
		WHERE g.delete_at IS NULL
		  AND (
		    EXISTS (SELECT 1 FROM kelas k WHERE k.id_guru::text = g.id::text AND k.delete_at IS NULL)
		    OR EXISTS (SELECT 1 FROM mata_pelajaran mp WHERE mp.id_guru::text = g.id::text AND mp.delete_at IS NULL)
		  )
	`).Scan(&guruBertugas)

	summary.Completeness.GuruTanpaTugas = summary.Totals.TotalGuru - guruBertugas
	if summary.Completeness.GuruTanpaTugas < 0 {
		summary.Completeness.GuruTanpaTugas = 0
	}

	_ = dm.db.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM guru 
		WHERE delete_at IS NULL 
		  AND (id_user IS NULL OR TRIM(id_user) = '')
	`).Scan(&summary.Completeness.GuruTanpaAkun)

	// 5. Metrik Kelengkapan Data (Mata Pelajaran)
	_ = dm.db.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM mata_pelajaran 
		WHERE delete_at IS NULL 
		  AND id_guru IS NOT NULL 
		  AND TRIM(id_guru) != ''
	`).Scan(&summary.Completeness.MapelDenganGuru)

	summary.Completeness.MapelTanpaGuru = summary.Totals.TotalMapel - summary.Completeness.MapelDenganGuru

	_ = dm.db.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM mata_pelajaran 
		WHERE delete_at IS NULL 
		  AND (kelas_id IS NULL OR TRIM(kelas_id) = '')
	`).Scan(&summary.Completeness.MapelTanpaKelas)

	// 6. Distribusi Siswa per Rombel Kelas (Urutkan dari yang terbanyak)
	distRows, err := dm.db.Query(ctx, `
		SELECT 
			k.id::text, 
			k.kelas, 
			COALESCE(g.nama, '') AS wali_guru, 
			COUNT(s.id) AS total_siswa
		FROM kelas k
		LEFT JOIN guru g ON g.id::text = k.id_guru::text AND g.delete_at IS NULL
		LEFT JOIN siswa s ON s.kelas_id::text = k.id::text AND s.delete_at IS NULL
		WHERE k.delete_at IS NULL
		GROUP BY k.id, k.kelas, g.nama
		ORDER BY total_siswa DESC, k.kelas ASC
		LIMIT 10
	`)
	if err == nil {
		defer distRows.Close()
		for distRows.Next() {
			var item dashboard.ClassStudentCount
			if err := distRows.Scan(&item.ID, &item.Kelas, &item.WaliGuru, &item.TotalSiswa); err == nil {
				summary.Distribution = append(summary.Distribution, item)
			}
		}
	}

	// 7. Khusus Admin: Breakdown User Roles
	if role == "admin" {
		roleRows, err := dm.db.Query(ctx, `
			SELECT role, COUNT(*) 
			FROM users 
			WHERE delete_at IS NULL 
			GROUP BY role 
			ORDER BY count DESC
		`)
		if err == nil {
			defer roleRows.Close()
			for roleRows.Next() {
				var item dashboard.UserRoleCount
				if err := roleRows.Scan(&item.Role, &item.Total); err == nil {
					summary.UserRoles = append(summary.UserRoles, item)
				}
			}
		}

		// 8. Khusus Admin: Ringkasan Aktivitas Sistem Hari Ini (WIB)
		_ = dm.db.QueryRow(ctx, `
			SELECT COUNT(*) 
			FROM transaction_logs 
			WHERE timestamp >= (CURRENT_DATE AT TIME ZONE 'Asia/Jakarta')
		`).Scan(&summary.TodayActivity.TotalHariIni)

		_ = dm.db.QueryRow(ctx, `
			SELECT COUNT(*) 
			FROM transaction_logs 
			WHERE timestamp >= (CURRENT_DATE AT TIME ZONE 'Asia/Jakarta')
			  AND result = 'Success'
		`).Scan(&summary.TodayActivity.BerhasilHariIni)

		_ = dm.db.QueryRow(ctx, `
			SELECT COUNT(*) 
			FROM transaction_logs 
			WHERE timestamp >= (CURRENT_DATE AT TIME ZONE 'Asia/Jakarta')
			  AND result = 'Failed'
		`).Scan(&summary.TodayActivity.GagalHariIni)

		_ = dm.db.QueryRow(ctx, `
			SELECT COUNT(DISTINCT user_id) 
			FROM transaction_logs 
			WHERE timestamp >= (CURRENT_DATE AT TIME ZONE 'Asia/Jakarta')
			  AND user_id IS NOT NULL 
			  AND TRIM(user_id) != ''
		`).Scan(&summary.TodayActivity.UserAktifHariIni)

		// 9. Khusus Admin: 5 Aktivitas Terbaru
		recRows, err := dm.db.Query(ctx, `
			SELECT 
				l.id, 
				l.timestamp, 
				COALESCE(u.username, 'Anonim') AS username, 
				COALESCE(u.role, 'user') AS role, 
				COALESCE(l.activity_type, '') AS activity_type, 
				COALESCE(l.activity_description, '') AS activity_description, 
				COALESCE(l.module, 'Sistem') AS module, 
				COALESCE(l.result, 'Success') AS result
			FROM transaction_logs l
			LEFT JOIN users u ON u.id = l.user_id
			ORDER BY l.timestamp DESC, l.id DESC
			LIMIT 5
		`)
		if err == nil {
			defer recRows.Close()
			for recRows.Next() {
				var rItem dashboard.RecentActivity
				if err := recRows.Scan(
					&rItem.ID,
					&rItem.Timestamp,
					&rItem.Username,
					&rItem.Role,
					&rItem.ActivityType,
					&rItem.ActivityDescription,
					&rItem.Module,
					&rItem.Result,
				); err == nil {
					summary.RecentActivities = append(summary.RecentActivities, rItem)
				}
			}
		}
	}

	// 10. Bangun Item "Perlu Perhatian" (Actionable Insights)
	if summary.Completeness.SiswaTanpaKelas > 0 {
		summary.AttentionItems = append(summary.AttentionItems, dashboard.AttentionItem{
			ID:      "siswa-tanpa-kelas",
			Title:   "Siswa Belum Memiliki Kelas",
			Count:   summary.Completeness.SiswaTanpaKelas,
			Message: fmt.Sprintf("%d siswa terdaftar belum dimasukkan ke rombel kelas.", summary.Completeness.SiswaTanpaKelas),
			Level:   "warning",
			Route:   "/admin/siswa",
			Action:  "Tempatkan Siswa",
		})
	}

	if summary.Completeness.KelasTanpaWali > 0 {
		summary.AttentionItems = append(summary.AttentionItems, dashboard.AttentionItem{
			ID:      "kelas-tanpa-wali",
			Title:   "Kelas Belum Memiliki Wali",
			Count:   summary.Completeness.KelasTanpaWali,
			Message: fmt.Sprintf("%d rombongan belajar belum memiliki guru pembimbing/wali kelas.", summary.Completeness.KelasTanpaWali),
			Level:   "danger",
			Route:   "/admin/kelas",
			Action:  "Tugaskan Wali Kelas",
		})
	}

	if summary.Completeness.MapelTanpaGuru > 0 {
		summary.AttentionItems = append(summary.AttentionItems, dashboard.AttentionItem{
			ID:      "mapel-tanpa-guru",
			Title:   "Mata Pelajaran Tanpa Guru",
			Count:   summary.Completeness.MapelTanpaGuru,
			Message: fmt.Sprintf("%d mata pelajaran belum ditentukan guru pengampunya.", summary.Completeness.MapelTanpaGuru),
			Level:   "warning",
			Route:   "/admin/mapel",
			Action:  "Pilih Guru Pengampu",
		})
	}

	if summary.Completeness.MapelTanpaKelas > 0 {
		summary.AttentionItems = append(summary.AttentionItems, dashboard.AttentionItem{
			ID:      "mapel-tanpa-kelas",
			Title:   "Mata Pelajaran Belum Dialokasikan",
			Count:   summary.Completeness.MapelTanpaKelas,
			Message: fmt.Sprintf("%d mata pelajaran belum terikat ke rombel kelas manapun.", summary.Completeness.MapelTanpaKelas),
			Level:   "info",
			Route:   "/admin/mapel",
			Action:  "Alokasikan Kelas",
		})
	}

	if role == "admin" && summary.Completeness.GuruTanpaAkun > 0 {
		summary.AttentionItems = append(summary.AttentionItems, dashboard.AttentionItem{
			ID:      "guru-tanpa-akun",
			Title:   "Guru Belum Terhubung Akun",
			Count:   summary.Completeness.GuruTanpaAkun,
			Message: fmt.Sprintf("%d guru belum dihubungkan dengan akun login pengguna sistem.", summary.Completeness.GuruTanpaAkun),
			Level:   "info",
			Route:   "/admin/guru",
			Action:  "Hubungkan Akun",
		})
	}

	summary.Role = role

	// 11. Khusus Role GURU: Buat TeacherDashboardData secara mendalam & personal
	if role == "guru" && userID != "" {
		teacherData, err := dm.getTeacherDashboardData(ctx, userID)
		if err == nil && teacherData != nil {
			summary.TeacherData = teacherData
		}
	}

	return summary, nil
}

func (dm *DashboardModel) getTeacherDashboardData(ctx context.Context, userID string) (*dashboard.TeacherDashboardData, error) {
	// 1. Ambil Profil Guru berdasarkan id_user
	var tData dashboard.TeacherDashboardData
	tData.WeeklyScheduleSummary = make(map[string]int)
	tData.Classes = make([]dashboard.TeacherClassDetail, 0)
	tData.TodaySchedules = make([]dashboard.TeacherScheduleItem, 0)
	tData.NamaKelasAjar = make([]string, 0)
	tData.NamaMapelDiampu = make([]string, 0)

	err := dm.db.QueryRow(ctx, `
		SELECT id, nama, COALESCE(email, ''), COALESCE(alamat, '')
		FROM guru
		WHERE id_user = $1 AND delete_at IS NULL
		LIMIT 1
	`, userID).Scan(&tData.GuruID, &tData.NamaGuru, &tData.Email, &tData.Alamat)
	if err != nil {
		return nil, err
	}

	// 2. Ambil Kelas Wali milik guru ini
	var kelasWaliID string
	_ = dm.db.QueryRow(ctx, `
		SELECT id, kelas
		FROM kelas
		WHERE id_guru = $1 AND delete_at IS NULL
		LIMIT 1
	`, tData.GuruID).Scan(&kelasWaliID, &tData.NamaKelasWali)
	if kelasWaliID != "" {
		tData.TotalKelasWali = 1
	}

	// 3. Ambil seluruh Mapel yang diampu guru ini beserta alokasi kelasnya
	mapelRows, err := dm.db.Query(ctx, `
		SELECT mp.id, mp.nama_pelajaran, COALESCE(mp.kelas_id, ''), COALESCE(k.kelas, '')
		FROM mata_pelajaran mp
		LEFT JOIN kelas k ON k.id = mp.kelas_id AND k.delete_at IS NULL
		WHERE mp.id_guru = $1 AND mp.delete_at IS NULL
		ORDER BY mp.nama_pelajaran ASC
	`, tData.GuruID)
	
	kelasAjarMap := make(map[string]string)       // id -> nama
	kelasMapelMap := make(map[string][]string)     // kelas_id -> list nama mapel
	mapelUnikSet := make(map[string]bool)

	if err == nil {
		defer mapelRows.Close()
		for mapelRows.Next() {
			var mID, mNama, kID, kNama string
			if err := mapelRows.Scan(&mID, &mNama, &kID, &kNama); err == nil {
				if !mapelUnikSet[mNama] {
					mapelUnikSet[mNama] = true
					tData.NamaMapelDiampu = append(tData.NamaMapelDiampu, mNama)
				}
				if kID != "" {
					kelasAjarMap[kID] = kNama
					kelasMapelMap[kID] = append(kelasMapelMap[kID], mNama)
				}
			}
		}
	}
	tData.TotalMapelDiampu = len(tData.NamaMapelDiampu)

	for _, kNama := range kelasAjarMap {
		tData.NamaKelasAjar = append(tData.NamaKelasAjar, kNama)
	}
	tData.TotalKelasAjar = len(kelasAjarMap)

	// 4. Kumpulkan seluruh kelas terkait guru (Kelas Wali + Kelas yang Diajar)
	semuaKelasIDs := make(map[string]bool)
	if kelasWaliID != "" {
		semuaKelasIDs[kelasWaliID] = true
	}
	for kID := range kelasAjarMap {
		semuaKelasIDs[kID] = true
	}

	totalSiswaAjarCount := 0
	for kID := range semuaKelasIDs {
		var detail dashboard.TeacherClassDetail
		detail.KelasID = kID
		detail.SampleSiswa = make([]dashboard.TeacherStudentBrief, 0)
		detail.AllSiswa = make([]dashboard.TeacherStudentBrief, 0)
		detail.MapelAjar = kelasMapelMap[kID]
		if detail.MapelAjar == nil {
			detail.MapelAjar = []string{}
		}

		if kID == kelasWaliID {
			detail.IsWali = true
			detail.NamaKelas = tData.NamaKelasWali
		} else {
			detail.NamaKelas = kelasAjarMap[kID]
		}
		if _, isAjar := kelasAjarMap[kID]; isAjar {
			detail.IsMengajar = true
		}

		// Ambil daftar siswa untuk kelas ini
		siswaRows, err := dm.db.Query(ctx, `
			SELECT id, nama, COALESCE(email, '')
			FROM siswa
			WHERE kelas_id = $1 AND delete_at IS NULL
			ORDER BY nama ASC
		`, kID)
		if err == nil {
			for siswaRows.Next() {
				var sBrief dashboard.TeacherStudentBrief
				if err := siswaRows.Scan(&sBrief.ID, &sBrief.Nama, &sBrief.Email); err == nil {
					detail.AllSiswa = append(detail.AllSiswa, sBrief)
					if len(detail.SampleSiswa) < 5 {
						detail.SampleSiswa = append(detail.SampleSiswa, sBrief)
					}
				}
			}
			siswaRows.Close()
		}
		detail.TotalSiswa = len(detail.AllSiswa)
		if detail.IsWali {
			tData.TotalSiswaWali = detail.TotalSiswa
		}
		totalSiswaAjarCount += detail.TotalSiswa
		tData.Classes = append(tData.Classes, detail)
	}
	tData.TotalSiswaAjar = totalSiswaAjarCount

	// 5. Ambil Jadwal Mengajar Guru Hari Ini (WIB) & Mingguan
	wibLoc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		wibLoc = time.FixedZone("WIB", 7*3600)
	}
	nowWIB := time.Now().In(wibLoc)
	todayWeekday := int(nowWIB.Weekday()) // Sunday = 0
	if todayWeekday == 0 {
		todayWeekday = 7
	} // 1: Senin, ..., 7: Minggu
	nowTimeStr := nowWIB.Format("15:04")

	// Summary mingguan
	weekRows, err := dm.db.Query(ctx, `
		SELECT j.hari, COUNT(*)
		FROM jadwal_pelajaran j
		JOIN mata_pelajaran mp ON mp.id = j.mata_pelajaran_id AND mp.delete_at IS NULL
		WHERE mp.id_guru = $1 AND j.delete_at IS NULL AND j.aktif = TRUE
		GROUP BY j.hari
	`, tData.GuruID)
	if err == nil {
		defer weekRows.Close()
		for weekRows.Next() {
			var h, count int
			if err := weekRows.Scan(&h, &count); err == nil {
				namaH := namaHariIndo(h)
				tData.WeeklyScheduleSummary[namaH] = count
			}
		}
	}

	// Jadwal hari ini
	todayRows, err := dm.db.Query(ctx, `
		SELECT 
			j.id,
			mp.nama_pelajaran,
			COALESCE(k.kelas, 'Umum') AS nama_kelas,
			j.hari,
			to_char(j.jam_mulai, 'HH24:MI') AS jam_mulai,
			to_char(j.jam_selesai, 'HH24:MI') AS jam_selesai,
			COALESCE(j.ruang, '') AS ruang
		FROM jadwal_pelajaran j
		JOIN mata_pelajaran mp ON mp.id = j.mata_pelajaran_id AND mp.delete_at IS NULL
		LEFT JOIN kelas k ON k.id = mp.kelas_id AND k.delete_at IS NULL
		WHERE mp.id_guru = $1 
		  AND j.hari = $2
		  AND j.delete_at IS NULL
		  AND j.aktif = TRUE
		ORDER BY j.jam_mulai ASC
	`, tData.GuruID, todayWeekday)

	var nextSlot *dashboard.TeacherScheduleItem

	if err == nil {
		defer todayRows.Close()
		for todayRows.Next() {
			var sch dashboard.TeacherScheduleItem
			if err := todayRows.Scan(
				&sch.ID,
				&sch.MataPelajaran,
				&sch.NamaKelas,
				&sch.Hari,
				&sch.JamMulai,
				&sch.JamSelesai,
				&sch.Ruang,
			); err == nil {
				sch.NamaHari = namaHariIndo(sch.Hari)

				// Hitung status jadwal berdasarkan waktu WIB
				if nowTimeStr >= sch.JamMulai && nowTimeStr <= sch.JamSelesai {
					sch.Status = "BERLANGSUNG"
					if nextSlot == nil {
						temp := sch
						nextSlot = &temp
					}
				} else if nowTimeStr < sch.JamMulai {
					sch.Status = "MENDATANG"
					if nextSlot == nil {
						sch.Status = "SELANJUTNYA"
						temp := sch
						nextSlot = &temp
					}
				} else {
					sch.Status = "SELESAI"
				}

				tData.TodaySchedules = append(tData.TodaySchedules, sch)
			}
		}
	}
	tData.TotalJadwalHariIni = len(tData.TodaySchedules)
	tData.NextSchedule = nextSlot

	return &tData, nil
}

func namaHariIndo(h int) string {
	switch h {
	case 1:
		return "Senin"
	case 2:
		return "Selasa"
	case 3:
		return "Rabu"
	case 4:
		return "Kamis"
	case 5:
		return "Jumat"
	case 6:
		return "Sabtu"
	case 7:
		return "Minggu"
	default:
		return "-"
	}
}
