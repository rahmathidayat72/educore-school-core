package router

import (
	authcontroller "go_rest_native_sekolah/features/auth/controllers"
	authmodels "go_rest_native_sekolah/features/auth/model"
	serviceauth "go_rest_native_sekolah/features/auth/service"
	gurucontroller "go_rest_native_sekolah/features/guru/controllers"
	gurumodels "go_rest_native_sekolah/features/guru/model"
	"go_rest_native_sekolah/features/guru/service"
	kelascontroller "go_rest_native_sekolah/features/kelas/controllers"
	kelasmodels "go_rest_native_sekolah/features/kelas/model"
	servicekelas "go_rest_native_sekolah/features/kelas/service"
	mapelcontroller "go_rest_native_sekolah/features/mata_pelajaran/controllers"
	mapelsmodels "go_rest_native_sekolah/features/mata_pelajaran/model"
	servicemapel "go_rest_native_sekolah/features/mata_pelajaran/service"
	siswacontroller "go_rest_native_sekolah/features/siswa/controllers"
	siswamodels "go_rest_native_sekolah/features/siswa/model"
	servicesiswa "go_rest_native_sekolah/features/siswa/service"
	userscontroller "go_rest_native_sekolah/features/users/controllers"
	usersmodels "go_rest_native_sekolah/features/users/model"
	serviceuser "go_rest_native_sekolah/features/users/service"

	logscontroller "go_rest_native_sekolah/features/logs/controllers"
	logsmodels "go_rest_native_sekolah/features/logs/model"
	servicelog "go_rest_native_sekolah/features/logs/service"

	dashboardcontroller "go_rest_native_sekolah/features/dashboard/controllers"
	dashboardmodels "go_rest_native_sekolah/features/dashboard/model"
	servicedashboard "go_rest_native_sekolah/features/dashboard/service"

	jadwalcontroller "go_rest_native_sekolah/features/jadwal/controllers"
	jadwalmodels "go_rest_native_sekolah/features/jadwal/model"
	servicejadwal "go_rest_native_sekolah/features/jadwal/service"

	schoolunitscontroller "go_rest_native_sekolah/features/school_units/controllers"
	schoolunitsmodels "go_rest_native_sekolah/features/school_units/model"
	serviceschoolunits "go_rest_native_sekolah/features/school_units/service"

	academicyearscontroller "go_rest_native_sekolah/features/academic_years/controllers"
	academicyearsmodels "go_rest_native_sekolah/features/academic_years/model"
	serviceacademicyears "go_rest_native_sekolah/features/academic_years/service"

	guardianscontroller "go_rest_native_sekolah/features/guardians/controllers"
	guardiansmodels "go_rest_native_sekolah/features/guardians/model"
	serviceguardians "go_rest_native_sekolah/features/guardians/service"

	"go_rest_native_sekolah/helper"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitRouter digunakan untuk menginisialisasi router.
// Fungsi ini akan menginisialisasi router untuk fitur auth, guru, users, dan kelas.
func InitRouter(db *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	// Pasang semua route
	// Endpoint /login digunakan untuk mengotentikasi user
	loginRouter(mux, db)
	// Endpoint /guru digunakan untuk mengelola data guru
	guruRouter(mux, db)
	// Endpoint /users digunakan untuk mengelola data user
	usersRouter(mux, db)
	// Endpoint /kelas digunakan untuk mengelola data kelas
	kelasRouter(mux, db)
	// Endpoint /siswa digunakan untuk mengelola data siswa
	siswaRouter(mux, db)
	// Endpoint /mapel digunakan untuk mengelola data mata pelajaran
	mataPelajaranRouter(mux, db)
	// Endpoint /logs digunakan untuk audit log sistem (hanya admin)
	logsRouter(mux, db)
	// Endpoint /dashboard digunakan untuk statistik & ringkasan komprehensif sistem
	dashboardRouter(mux, db)
	// Endpoint /jadwal digunakan untuk jadwal pelajaran sekolah & mengajar guru
	jadwalRouter(mux, db)
	// Endpoint /school-units digunakan untuk manajemen unit sekolah yayasan (Super Admin Yayasan)
	schoolUnitsRouter(mux, db)
	// Endpoint /academic-years dan /class-enrollments untuk pengelolaan tahun ajaran & riwayat kelas siswa
	academicYearsRouter(mux, db)
	// Endpoint /guardians untuk master data orang tua/wali murid & relasi student_guardians
	guardiansRouter(mux, db)

	// Bungkus mux dengan middleware logging dan cors
	return helper.CORSMiddleware(helper.LoggingMiddleware(mux, db))
}

// loginRouter digunakan untuk menginisialisasi router untuk fitur auth.
// Fungsi ini akan menginisialisasi router untuk endpoint /login yang digunakan untuk mengotentikasi user.
// Endpoint /login akan menerima request dengan method POST dan mengembalikan response JSON.
func loginRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	// Inisialisasi repository
	authRepo := authmodels.NewAuthData(db)
	// Inisialisasi service
	authService := serviceauth.NewServiceAuth(authRepo)
	// Inisialisasi controller
	authController := authcontroller.NewAutController(authService)
	// Membuat handler untuk endpoint /login
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// Cek apakah request menggunakan method POST
		if r.Method == http.MethodPost {
			err := authController.Auth(w, r)
			if err != nil {
				// Jika terjadi error maka akan mengembalikan response JSON dengan status Internal Server Error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			// Jika request tidak menggunakan method POST maka akan mengembalikan response JSON dengan status Method Not Allowed
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})
}

func guruRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	// Inisialisasi repository
	guruRepo := gurumodels.NewDataGuru(db)

	// Inisialisasi service
	guruService := service.NewServiceGuru(guruRepo, db)

	// Inisialisasi controller
	guruController := gurucontroller.NewGuruController(guruService)

	// Endpoint GET untuk menampilkan data guru (Admin dan Guru untuk kebutuhan relasi kelas/mapel)
	mux.HandleFunc("/guru", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, []string{"admin", "guru"}, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := guruController.Guru(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// Endpoint POST untuk menambah data guru (Khusus Admin)
	mux.HandleFunc("/guru/tambah", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := guruController.InsertGuru(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/guru/update", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			// Langsung jalankan fungsi UpdateGuru milik controller
			err := guruController.UpdateGuru(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/guru/gurubyid", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodPut {
			// Langsung jalankan fungsi GuruById milik controller
			err := guruController.GetGuruById(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/guru/deleted", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete || r.Method == http.MethodPut {
			// Langsung jalankan fungsi DeletedById milik controller
			err := guruController.DeleteGuru(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))
}

func usersRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	usersRepo := usersmodels.NewUserData(db)
	usersService := serviceuser.NewServiceUser(usersRepo, db)
	usersController := userscontroller.NewUsesController(usersService)

	mux.HandleFunc("/users", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := usersController.Users(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/users/tambah", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := usersController.InsertUser(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/users/userbyid", helper.AuthMiddleware(helper.UserSelfOrAdminMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodPut {
			// Langsung jalankan fungsi GetUserById milik controller
			err := usersController.GetUserById(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/users/update", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			// Langsung jalankan fungsi UpdateUser milik controller
			err := usersController.UpdateUser(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/users/deleted", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete || r.Method == http.MethodPut {
			// Langsung jalankan fungsi DeletedById milik controller
			err := usersController.DeleteUser(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))
}

func kelasRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	kelasRepo := kelasmodels.NewDataKelas(db)
	kelasService := servicekelas.NewServiceKelas(kelasRepo)
	kelasController := kelascontroller.NewKelasController(kelasService)

	// Admin dan Guru diizinkan mengelola data kelas
	allowedRoles := []string{"admin", "guru"}

	mux.HandleFunc("/kelas/tambah", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := kelasController.Insert(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/kelas", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := kelasController.Kelas(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/kelas/kelasbyid", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Langsung jalankan fungsi GetKelasById milik controller
			err := kelasController.GetKelasById(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/kelas/update", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			// Langsung jalankan fungsi UpdateKelas milik controller
			err := kelasController.UpdateKelas(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/kelas/deleted", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete || r.Method == http.MethodPut {
			// Langsung jalankan fungsi DeleteKelas milik controller
			err := kelasController.DeleteKelas(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))
}

func siswaRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	siswaRepo := siswamodels.NewSiswaData(db)
	siswaService := servicesiswa.NewServiceSiswa(siswaRepo)
	siswaController := siswacontroller.NewSiswaController(siswaService)

	// Data siswa hanya boleh dikelola oleh admin
	mux.HandleFunc("/siswa/tambah", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := siswaController.InsertSiswa(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/siswa", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := siswaController.Siswa(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/siswa/siswabyid", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Langsung jalankan fungsi GetSiswaById milik controller
			err := siswaController.GetSiswaById(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/siswa/update", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			// Langsung jalankan fungsi UpdateSiswa milik controller
			err := siswaController.UpdateSiswa(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/siswa/deleted", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete || r.Method == http.MethodPut {
			// Langsung jalankan fungsi DeleteSiswa milik controller
			err := siswaController.DeleteSiswa(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))
}

func mataPelajaranRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	mataPelajaranRepo := mapelsmodels.NewDataMataPelajaran(db)
	mataPelajaranService := servicemapel.NewMataPelajaranService(mataPelajaranRepo)
	mataPelajaranController := mapelcontroller.NewMataPelajaranController(mataPelajaranService)

	// Admin dan Guru diizinkan mengelola data mata pelajaran
	allowedRoles := []string{"admin", "guru"}

	mux.HandleFunc("/mapel/tambah", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := mataPelajaranController.InsertMapel(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/mapel", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := mataPelajaranController.Mapel(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/mapel/mapelbyid", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Langsung jalankan fungsi GetMapelById milik controller
			err := mataPelajaranController.GetMapelById(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/mapel/update", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			// Langsung jalankan fungsi UpdateMapel milik controller
			err := mataPelajaranController.UpdateMapel(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/mapel/deleted", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete || r.Method == http.MethodPut {
			// Langsung jalankan fungsi DeleteMapel milik controller
			err := mataPelajaranController.DeleteMapel(w, r)
			if err != nil {
				// Tampilkan response JSON error
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))
}

func logsRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	logRepo := logsmodels.NewLogModel(db)
	logService := servicelog.NewServiceLog(logRepo)
	logController := logscontroller.NewLogController(logService)

	// GET /logs - Mengambil daftar seluruh log transaksi sistem (Admin Only)
	mux.HandleFunc("/logs", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := logController.GetLogs(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// GET /logs/detail?id=123 - Mengambil detail log transaksi sistem berdasarkan ID (Admin Only)
	mux.HandleFunc("/logs/detail", helper.AuthMiddleware(helper.AdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			idStr := r.URL.Query().Get("id")
			id, err := strconv.Atoi(idStr)
			if err != nil || id <= 0 {
				helper.JSONResponse(w, http.StatusBadRequest, "Parameter id tidak valid")
				return
			}
			err = logController.GetLogByID(w, r, id)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))
}

func dashboardRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	dashboardRepo := dashboardmodels.NewDashboardModel(db)
	dashboardService := servicedashboard.NewDashboardService(dashboardRepo)
	dashboardController := dashboardcontroller.NewDashboardController(dashboardService)

	// GET /dashboard/summary - Mengambil seluruh metrik, distribusi, dan insight kualitas data sekolah
	mux.HandleFunc("/dashboard/summary", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, []string{"admin", "guru", "user"}, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := dashboardController.GetSummary(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))
}

func jadwalRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	jadwalRepo := jadwalmodels.NewJadwalModel(db)
	jadwalService := servicejadwal.NewJadwalService(jadwalRepo)
	jadwalController := jadwalcontroller.NewJadwalController(jadwalService)

	allowedRoles := []string{"admin", "guru"}

	// GET /jadwal - Mengambil daftar jadwal pelajaran (Admin: semua / filter, Guru: jadwal miliknya)
	mux.HandleFunc("/jadwal", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := jadwalController.GetAll(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// GET /jadwal/detail?id=... - Mengambil rincian jadwal
	mux.HandleFunc("/jadwal/detail", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := jadwalController.GetByID(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// POST /jadwal/tambah - Menambah jadwal baru
	mux.HandleFunc("/jadwal/tambah", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := jadwalController.Create(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// PUT & POST /jadwal/update?id=... - Mengubah jadwal
	mux.HandleFunc("/jadwal/update", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			err := jadwalController.Update(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// DELETE & PUT /jadwal/deleted?id=... - Menghapus jadwal
	mux.HandleFunc("/jadwal/deleted", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete || r.Method == http.MethodPut {
			err := jadwalController.Delete(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))
}

func schoolUnitsRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	unitRepo := schoolunitsmodels.NewSchoolUnitModel(db)
	unitService := serviceschoolunits.NewSchoolUnitService(unitRepo)
	unitController := schoolunitscontroller.NewSchoolUnitController(unitService)

	allowedReadRoles := []string{"super_admin", "admin"}

	// GET & POST /school-units - Mengambil daftar seluruh unit sekolah (GET) atau menambah unit baru (POST)
	mux.HandleFunc("/school-units", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			helper.RolesAllowedMiddleware(db, allowedReadRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := unitController.GetAll(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodPost:
			helper.SuperAdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
				if err := unitController.Create(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodPut:
			helper.SuperAdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
				if err := unitController.Update(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodDelete:
			helper.SuperAdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
				if err := unitController.Delete(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		default:
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	// GET /school-units/detail?id=... - Mengambil rincian unit sekolah
	mux.HandleFunc("/school-units/detail", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedReadRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			err := unitController.GetByID(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// POST /school-units/tambah - Menambah unit sekolah baru (Khusus Super Admin Yayasan)
	mux.HandleFunc("/school-units/tambah", helper.AuthMiddleware(helper.SuperAdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := unitController.Create(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// PUT & POST /school-units/update?id=... - Mengubah profil unit sekolah (Khusus Super Admin Yayasan)
	mux.HandleFunc("/school-units/update", helper.AuthMiddleware(helper.SuperAdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			err := unitController.Update(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// DELETE & PUT /school-units/deleted?id=... - Menghapus unit sekolah secara soft delete (Khusus Super Admin Yayasan)
	mux.HandleFunc("/school-units/deleted", helper.AuthMiddleware(helper.SuperAdminOnlyMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete || r.Method == http.MethodPut {
			err := unitController.Delete(w, r)
			if err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))
}

func academicYearsRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	ayRepo := academicyearsmodels.NewAcademicYearModel(db)
	ayService := serviceacademicyears.NewAcademicYearService(ayRepo)
	ayController := academicyearscontroller.NewAcademicYearController(ayService)

	allowedAdminRoles := []string{"super_admin", "admin"}
	allowedReadRoles := []string{"super_admin", "admin", "guru"}

	// --- Academic Years Endpoints ---
	// GET, POST, PUT, DELETE /academic-years (RESTful multiplexing)
	mux.HandleFunc("/academic-years", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			helper.RolesAllowedMiddleware(db, allowedReadRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := ayController.GetAllYears(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodPost:
			helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := ayController.CreateYear(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodPut:
			helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := ayController.UpdateYear(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodDelete:
			helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := ayController.DeleteYear(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		default:
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/academic-years/detail", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedReadRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := ayController.GetYearByID(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/academic-years/active", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedReadRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := ayController.GetActiveYear(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/academic-years/tambah", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := ayController.CreateYear(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/academic-years/update", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if err := ayController.UpdateYear(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/academic-years/set-active", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if err := ayController.SetActiveYear(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// Alias endpoint /academic-years/activate
	mux.HandleFunc("/academic-years/activate", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if err := ayController.SetActiveYear(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/academic-years/deleted", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodDelete {
			if err := ayController.DeleteYear(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// --- Class Enrollments Endpoints ---
	// GET, POST, PUT, DELETE /class-enrollments (RESTful multiplexing)
	mux.HandleFunc("/class-enrollments", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			helper.RolesAllowedMiddleware(db, allowedReadRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := ayController.GetEnrollments(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodPost:
			helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := ayController.EnrollStudent(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodPut:
			helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := ayController.UpdateEnrollment(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodDelete:
			helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := ayController.DeleteEnrollment(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		default:
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/class-enrollments/detail", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedReadRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := ayController.GetEnrollmentByID(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/class-enrollments/tambah", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := ayController.EnrollStudent(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/class-enrollments/update", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if err := ayController.UpdateEnrollment(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/class-enrollments/deleted", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodDelete {
			if err := ayController.DeleteEnrollment(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/class-enrollments/student-history", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedReadRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := ayController.GetStudentEnrollmentHistory(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))
}

func guardiansRouter(mux *http.ServeMux, db *pgxpool.Pool) {
	gRepo := guardiansmodels.NewGuardianModel(db)
	gService := serviceguardians.NewGuardianService(gRepo)
	gController := guardianscontroller.NewGuardianController(gService)

	allowedAdminRoles := []string{"super_admin", "admin"}
	allowedStaffRoles := []string{"super_admin", "admin", "guru"}
	allowedPortalRoles := []string{"super_admin", "admin", "guru", "wali"}

	// --- Guardians Master Endpoints ---
	// GET, POST, PUT, DELETE /guardians (RESTful multiplexing)
	mux.HandleFunc("/guardians", helper.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			helper.RolesAllowedMiddleware(db, allowedStaffRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := gController.GetAll(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodPost:
			helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := gController.Create(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodPut:
			helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := gController.Update(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		case http.MethodDelete:
			helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
				if err := gController.Delete(w, r); err != nil {
					helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
				}
			})(w, r)
		default:
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}))

	mux.HandleFunc("/guardians/detail", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedStaffRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := gController.GetByID(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/guardians/my-profile", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedPortalRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := gController.GetMyProfile(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/guardians/tambah", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := gController.Create(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/guardians/update", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if err := gController.Update(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/guardians/deleted", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodDelete {
			if err := gController.Delete(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	// --- Student Guardians Relation Endpoints ---
	mux.HandleFunc("/guardians/link-student", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := gController.LinkStudent(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/guardians/unlink-student", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedAdminRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodDelete {
			if err := gController.UnlinkStudent(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/guardians/students", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedPortalRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := gController.GetStudentsByGuardianID(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))

	mux.HandleFunc("/guardians/by-student", helper.AuthMiddleware(helper.RolesAllowedMiddleware(db, allowedStaffRoles, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := gController.GetGuardiansBySiswaID(w, r); err != nil {
				helper.JSONResponse(w, http.StatusInternalServerError, err.Error())
			}
		} else {
			helper.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})))
}
