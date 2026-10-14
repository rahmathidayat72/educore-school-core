package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionLog merepresentasikan data log transaksi yang disimpan di database.
type TransactionLog struct {
	ID                  int       `json:"id"`
	Timestamp           time.Time `json:"timestamp"`
	UserID              string    `json:"user_id"`
	ActivityType        string    `json:"activity_type"`
	ActivityDescription string    `json:"activity_description"`
	HTTPMethod          string    `json:"http_method"`
	Endpoint            string    `json:"endpoint"`
	Module              string    `json:"module"`
	TargetID            string    `json:"target_id"`
	StatusCode          int       `json:"status_code"`
	IPAddress           string    `json:"ip_address"`
	DurationMS          int64     `json:"duration_ms"`
	Perangkat           string    `json:"perangkat"`
	ServiceName         string    `json:"service_name"`
	RequestBody         string    `json:"request_body"`
	ResponseBody        string    `json:"response_body"`
	RequestParam        string    `json:"request_param"`
	Result              string    `json:"result"`
	Header              string    `json:"header"`
}

// responseCategory adalah ResponseWriter custom untuk menangkap response body & status
type responseCategory struct {
	http.ResponseWriter
	body       *[]byte
	statusCode int
	wrote      bool // untuk mencegah WriteHeader dipanggil 2x
}

func (rw *responseCategory) WriteHeader(status int) {
	if rw.wrote {
		return
	}
	rw.statusCode = status
	rw.ResponseWriter.WriteHeader(status)
	rw.wrote = true
}

// ShouldLogRequest menentukan apakah sebuah request harus dicatat ke audit log database.
// Aturan penghematan penyimpanan:
// 1. Method OPTIONS (CORS preflight) diabaikan.
// 2. Request GET & HEAD yang berhasil (statusCode < 400 / sekedar melihat data) TIDAK dicatat.
// 3. Request GET & HEAD yang gagal/ditolak (statusCode >= 400) TETAP dicatat sebagai audit keamanan.
// 4. Request mutasi data (POST, PUT, PATCH, DELETE) SELALU dicatat (login, create, update, delete).
func ShouldLogRequest(method string, statusCode int) bool {
	switch strings.ToUpper(method) {
	case http.MethodOptions:
		return false
	case http.MethodGet, http.MethodHead:
		return statusCode >= 400
	default:
		return true
	}
}

const maxResponseBodyLog = 64 * 1024 // 64 KB limit untuk response body log

func (rw *responseCategory) Write(b []byte) (int, error) {
	if rw.body != nil && len(*rw.body) < maxResponseBodyLog {
		remaining := maxResponseBodyLog - len(*rw.body)
		if len(b) > remaining {
			*rw.body = append(*rw.body, b[:remaining]...)
		} else {
			*rw.body = append(*rw.body, b...)
		}
	}
	if !rw.wrote {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// GetServiceNameFromEndpoint digunakan untuk memberi nama service berdasarkan endpoint
func GetServiceNameFromEndpoint(endpoint string) string {
	endpoint = strings.Split(endpoint, "?")[0]
	parts := strings.Split(strings.Trim(endpoint, "/"), "/")

	if len(parts) >= 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	} else if len(parts) == 1 {
		return parts[0]
	}
	return ""
}

// LoggingMiddleware mengumpulkan data transaksi dan detail aktivitas setiap request
func LoggingMiddleware(next http.Handler, db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		var requestBody []byte
		var responseBody []byte

		if db == nil {
			http.Error(w, "DB pool belum diinisialisasi", http.StatusInternalServerError)
			return
		}

		// Baca request body jika ada
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				requestBody = bodyBytes
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Ambil userID dari token JWT
		authHeader := r.Header.Get("Authorization")
		accessToken := GetTokenFromAuthorizationHeader(authHeader)

		var userID string
		if accessToken != "" {
			metaToken, err := VerifyTokenHeader(accessToken)
			if err == nil {
				userID = metaToken.ID
			}
		}

		// Fallback ambil userID berdasarkan email login pada body jika token belum ada
		if userID == "" && len(requestBody) > 0 {
			var bodyMap map[string]interface{}
			if err := json.Unmarshal(requestBody, &bodyMap); err == nil {
				if email, ok := bodyMap["email"].(string); ok {
					uid, err := GetUserIDByEmail(db, email)
					if err == nil {
						userID = uid
					}
				}
			}
		}

		// ResponseWriter custom
		responseWriter := &responseCategory{
			ResponseWriter: w,
			body:           &responseBody,
			statusCode:     http.StatusOK,
		}

		// Jalankan handler berikutnya
		next.ServeHTTP(responseWriter, r)

		// Filter efisiensi: jika aktivitas sekedar melihat data yang berhasil atau preflight OPTIONS,
		// tidak perlu dicatat ke database agar tabel transaction_logs tidak membengkak.
		if !ShouldLogRequest(r.Method, responseWriter.statusCode) {
			return
		}

		// Hitung durasi eksekusi request
		durationMS := time.Since(startTime).Milliseconds()

		// Simpan log secara asynchronous di background
		go func() {
			perangkat := r.UserAgent()
			if perangkat == "" {
				perangkat, _ = os.Hostname()
			}

			// Ambil IP address klien
			ipAddress := r.Header.Get("X-Forwarded-For")
			if ipAddress == "" {
				ipAddress = r.Header.Get("X-Real-IP")
			}
			if ipAddress == "" {
				ipAddress = r.RemoteAddr
			}
			if idx := strings.Index(ipAddress, ":"); idx != -1 && !strings.Contains(ipAddress, ",") && !strings.Contains(ipAddress, "[") {
				ipAddress = strings.Split(ipAddress, ":")[0]
			}

			paramJSON, _ := json.Marshal(r.URL.Query())
			headerJSON := maskHeaders(r.Header)

			// Masking request body sensitif
			requestStr := "{}"
			if len(requestBody) > 0 {
				requestStr = maskSensitiveData(string(requestBody))
			}

			// Response body JSON valid
			responseStr := "{}"
			if json.Valid(responseBody) {
				responseStr = string(responseBody)
			} else if len(responseBody) > 0 {
				tmp, _ := json.Marshal(string(responseBody))
				responseStr = string(tmp)
			}

			resultStatus := "Success"
			if responseWriter.statusCode >= 400 {
				resultStatus = "Failed"
			}

			// Dapatkan nama/username pelaku jika ada untuk deskripsi yang ramah
			actorLabel := "Pengguna"
			if userID != "" {
				var actorName, actorRole string
				err := db.QueryRow(context.Background(), "SELECT username, role FROM users WHERE id = $1", userID).Scan(&actorName, &actorRole)
				if err == nil && actorName != "" {
					switch strings.ToLower(actorRole) {
					case "admin":
						actorLabel = fmt.Sprintf("Admin %s", actorName)
					case "guru":
						actorLabel = fmt.Sprintf("Guru %s", actorName)
					default:
						actorLabel = fmt.Sprintf("Pengguna %s", actorName)
					}
				}
			}

			// Tentukan detail aktivitas user dengan bahasa Indonesia yang ramah
			module, activityType, targetID, activityDesc := DetermineActivityDetail(
				actorLabel,
				r.Method,
				r.URL.Path,
				r.URL.Query(),
				requestBody,
				responseWriter.statusCode,
			)

			// Jika aktivitas READ gagal/ditolak (karena hanya READ yang gagal yang dicatat untuk audit keamanan)
			if activityType == "READ" && responseWriter.statusCode >= 400 {
				activityDesc = fmt.Sprintf("%s (Gagal diakses - HTTP %d)", activityDesc, responseWriter.statusCode)
			}

			_, err := db.Exec(
				context.Background(),
				`INSERT INTO transaction_logs 
					(timestamp, user_id, activity_type, activity_description, http_method, endpoint, module, target_id, status_code, ip_address, duration_ms, perangkat, service_name, request_body, response_body, request_param, result, header)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14::jsonb, $15::jsonb, $16::jsonb, $17, $18::jsonb)`,
				time.Now(),
				userID,
				activityType,
				activityDesc,
				r.Method,
				r.RequestURI,
				module,
				targetID,
				responseWriter.statusCode,
				ipAddress,
				durationMS,
				perangkat,
				GetServiceNameFromEndpoint(r.RequestURI),
				requestStr,
				responseStr,
				string(paramJSON),
				resultStatus,
				headerJSON,
			)
			if err != nil {
				log.Println("[ERROR] Gagal simpan log transaksi:", err)
			}
		}()
	})
}

// DetermineActivityDetail menganalisis request untuk menghasilkan informasi aktivitas yang ramah dan mudah dipahami admin sekolah
func DetermineActivityDetail(actorLabel, method, path string, query map[string][]string, requestBody []byte, statusCode int) (module, activityType, targetID, description string) {
	path = strings.TrimSuffix(path, "/")
	if actorLabel == "" {
		actorLabel = "Pengguna"
	}

	var bodyMap map[string]interface{}
	if len(requestBody) > 0 {
		_ = json.Unmarshal(requestBody, &bodyMap)
	}

	targetID = ""
	if ids, ok := query["id"]; ok && len(ids) > 0 {
		targetID = ids[0]
	} else if idVal, ok := bodyMap["id"].(string); ok {
		targetID = idVal
	}

	// 1. Modul Autentikasi
	if path == "/login" {
		module = "Autentikasi"
		activityType = "LOGIN"
		email := ""
		if em, ok := bodyMap["email"].(string); ok {
			email = em
		}
		if statusCode < 400 {
			if email != "" {
				description = fmt.Sprintf("%s berhasil masuk ke sistem (%s)", actorLabel, email)
			} else {
				description = fmt.Sprintf("%s berhasil masuk ke sistem", actorLabel)
			}
		} else {
			if email != "" {
				description = fmt.Sprintf("%s gagal masuk ke sistem (Email: %s)", actorLabel, email)
			} else {
				description = fmt.Sprintf("%s gagal masuk ke sistem (Email atau kata sandi tidak cocok)", actorLabel)
			}
		}
		return
	}

	// 2. Modul Pengguna (Users)
	if strings.HasPrefix(path, "/users") {
		module = "Pengguna"
		switch path {
		case "/users/tambah":
			activityType = "CREATE"
			username, _ := bodyMap["username"].(string)
			email, _ := bodyMap["email"].(string)
			role, _ := bodyMap["role"].(string)
			if role == "" {
				role = "user"
			}
			if username != "" {
				description = fmt.Sprintf("%s menambahkan akun pengguna baru: %s (Peran: %s)", actorLabel, username, role)
			} else {
				description = fmt.Sprintf("%s menambahkan akun pengguna baru (Email: %s)", actorLabel, email)
			}
		case "/users/update":
			activityType = "UPDATE"
			username, _ := bodyMap["username"].(string)
			if username != "" {
				description = fmt.Sprintf("%s memperbarui profil akun pengguna %s", actorLabel, username)
			} else {
				description = fmt.Sprintf("%s memperbarui data akun pengguna", actorLabel)
			}
		case "/users/deleted":
			activityType = "DELETE"
			description = fmt.Sprintf("%s menghapus akun pengguna dari sistem", actorLabel)
		case "/users/userbyid":
			activityType = "READ"
			description = fmt.Sprintf("%s melihat rincian akun pengguna", actorLabel)
		default:
			activityType = "READ"
			description = fmt.Sprintf("%s melihat daftar seluruh akun pengguna", actorLabel)
		}
		return
	}

	// 3. Modul Guru
	if strings.HasPrefix(path, "/guru") {
		module = "Guru"
		switch path {
		case "/guru/tambah":
			activityType = "CREATE"
			nama, _ := bodyMap["nama"].(string)
			if nama != "" {
				description = fmt.Sprintf("%s menambahkan data guru baru: %s", actorLabel, nama)
			} else {
				description = fmt.Sprintf("%s menambahkan data guru baru", actorLabel)
			}
		case "/guru/update":
			activityType = "UPDATE"
			nama, _ := bodyMap["nama"].(string)
			if nama != "" {
				description = fmt.Sprintf("%s memperbarui data guru: %s", actorLabel, nama)
			} else {
				description = fmt.Sprintf("%s memperbarui data guru", actorLabel)
			}
		case "/guru/deleted":
			activityType = "DELETE"
			description = fmt.Sprintf("%s menghapus data guru dari sistem", actorLabel)
		case "/guru/gurubyid":
			activityType = "READ"
			description = fmt.Sprintf("%s melihat rincian data guru", actorLabel)
		default:
			activityType = "READ"
			description = fmt.Sprintf("%s melihat daftar seluruh guru", actorLabel)
		}
		return
	}

	// 4. Modul Siswa
	if strings.HasPrefix(path, "/siswa") {
		module = "Siswa"
		switch path {
		case "/siswa/tambah":
			activityType = "CREATE"
			nama, _ := bodyMap["nama"].(string)
			if nama != "" {
				description = fmt.Sprintf("%s menambahkan siswa baru bernama %s", actorLabel, nama)
			} else {
				description = fmt.Sprintf("%s menambahkan data siswa baru", actorLabel)
			}
		case "/siswa/update":
			activityType = "UPDATE"
			nama, _ := bodyMap["nama"].(string)
			if nama != "" {
				description = fmt.Sprintf("%s memperbarui data siswa bernama %s", actorLabel, nama)
			} else {
				description = fmt.Sprintf("%s memperbarui data siswa", actorLabel)
			}
		case "/siswa/deleted":
			activityType = "DELETE"
			description = fmt.Sprintf("%s menghapus data siswa dari sistem", actorLabel)
		case "/siswa/siswabyid":
			activityType = "READ"
			description = fmt.Sprintf("%s melihat rincian data siswa", actorLabel)
		default:
			activityType = "READ"
			description = fmt.Sprintf("%s melihat daftar seluruh siswa", actorLabel)
		}
		return
	}

	// 5. Modul Kelas
	if strings.HasPrefix(path, "/kelas") {
		module = "Kelas"
		switch path {
		case "/kelas/tambah":
			activityType = "CREATE"
			namaKelas, _ := bodyMap["kelas"].(string)
			if namaKelas != "" {
				description = fmt.Sprintf("%s membuat rombel kelas baru: %s", actorLabel, namaKelas)
			} else {
				description = fmt.Sprintf("%s membuat rombel kelas baru", actorLabel)
			}
		case "/kelas/update":
			activityType = "UPDATE"
			namaKelas, _ := bodyMap["kelas"].(string)
			if namaKelas != "" {
				description = fmt.Sprintf("%s memperbarui data kelas %s", actorLabel, namaKelas)
			} else {
				description = fmt.Sprintf("%s memperbarui data kelas", actorLabel)
			}
		case "/kelas/deleted":
			activityType = "DELETE"
			description = fmt.Sprintf("%s menghapus rombel kelas dari sistem", actorLabel)
		case "/kelas/kelasbyid":
			activityType = "READ"
			description = fmt.Sprintf("%s melihat rincian informasi kelas", actorLabel)
		default:
			activityType = "READ"
			description = fmt.Sprintf("%s melihat daftar seluruh kelas", actorLabel)
		}
		return
	}

	// 6. Modul Mata Pelajaran
	if strings.HasPrefix(path, "/mapel") {
		module = "Mata Pelajaran"
		switch path {
		case "/mapel/tambah":
			activityType = "CREATE"
			namaMapel, _ := bodyMap["nama_pelajaran"].(string)
			if namaMapel != "" {
				description = fmt.Sprintf("%s menambahkan mata pelajaran baru: %s", actorLabel, namaMapel)
			} else {
				description = fmt.Sprintf("%s menambahkan mata pelajaran baru", actorLabel)
			}
		case "/mapel/update":
			activityType = "UPDATE"
			namaMapel, _ := bodyMap["nama_pelajaran"].(string)
			if namaMapel != "" {
				description = fmt.Sprintf("%s memperbarui mata pelajaran %s", actorLabel, namaMapel)
			} else {
				description = fmt.Sprintf("%s memperbarui mata pelajaran", actorLabel)
			}
		case "/mapel/deleted":
			activityType = "DELETE"
			description = fmt.Sprintf("%s menghapus mata pelajaran dari sistem", actorLabel)
		case "/mapel/mapelbyid":
			activityType = "READ"
			description = fmt.Sprintf("%s melihat rincian mata pelajaran", actorLabel)
		default:
			activityType = "READ"
			description = fmt.Sprintf("%s melihat daftar seluruh mata pelajaran", actorLabel)
		}
		return
	}

	// 7. Modul Log Sistem
	if strings.HasPrefix(path, "/logs") {
		module = "Log Sistem"
		activityType = "READ"
		if path == "/logs/detail" {
			description = fmt.Sprintf("%s melihat rincian detail aktivitas log sistem", actorLabel)
		} else {
			description = fmt.Sprintf("%s membuka riwayat audit aktivitas sistem", actorLabel)
		}
		return
	}

	// 8. Modul Jadwal Pelajaran
	if strings.HasPrefix(path, "/jadwal") {
		module = "Jadwal Pelajaran"
		switch path {
		case "/jadwal/tambah":
			activityType = "CREATE"
			description = fmt.Sprintf("%s menambahkan jadwal pelajaran baru", actorLabel)
		case "/jadwal/update":
			activityType = "UPDATE"
			description = fmt.Sprintf("%s memperbarui jadwal pelajaran", actorLabel)
		case "/jadwal/deleted":
			activityType = "DELETE"
			description = fmt.Sprintf("%s menghapus jadwal pelajaran", actorLabel)
		case "/jadwal/detail":
			activityType = "READ"
			description = fmt.Sprintf("%s melihat rincian jadwal pelajaran", actorLabel)
		default:
			activityType = "READ"
			description = fmt.Sprintf("%s melihat daftar jadwal pelajaran", actorLabel)
		}
		return
	}

	// 9. Modul Unit Sekolah (School Units)
	if strings.HasPrefix(path, "/school-units") {
		module = "Unit Sekolah"
		switch path {
		case "/school-units/tambah":
			activityType = "CREATE"
			description = fmt.Sprintf("%s menambahkan unit sekolah baru", actorLabel)
		case "/school-units/update":
			activityType = "UPDATE"
			description = fmt.Sprintf("%s memperbarui profil unit sekolah", actorLabel)
		case "/school-units/deleted":
			activityType = "DELETE"
			description = fmt.Sprintf("%s menghapus unit sekolah dari yayasan", actorLabel)
		default:
			activityType = "READ"
			description = fmt.Sprintf("%s melihat data unit sekolah", actorLabel)
		}
		return
	}

	// 10. Modul Tahun Ajaran (Academic Years)
	if strings.HasPrefix(path, "/academic-years") {
		module = "Tahun Ajaran"
		switch path {
		case "/academic-years/tambah":
			activityType = "CREATE"
			description = fmt.Sprintf("%s membuat periode tahun ajaran baru", actorLabel)
		case "/academic-years/update":
			activityType = "UPDATE"
			description = fmt.Sprintf("%s memperbarui periode tahun ajaran", actorLabel)
		case "/academic-years/activate":
			activityType = "UPDATE"
			description = fmt.Sprintf("%s mengaktifkan tahun ajaran berjalan", actorLabel)
		case "/academic-years/deleted":
			activityType = "DELETE"
			description = fmt.Sprintf("%s menghapus data tahun ajaran", actorLabel)
		default:
			activityType = "READ"
			description = fmt.Sprintf("%s melihat daftar tahun ajaran", actorLabel)
		}
		return
	}

	// 11. Modul Penempatan Kelas (Class Enrollments)
	if strings.HasPrefix(path, "/class-enrollments") {
		module = "Riwayat Rombel Siswa"
		switch path {
		case "/class-enrollments/enroll":
			activityType = "CREATE"
			description = fmt.Sprintf("%s mendaftarkan penempatan siswa ke rombel", actorLabel)
		case "/class-enrollments/update":
			activityType = "UPDATE"
			description = fmt.Sprintf("%s memperbarui status penempatan rombel siswa", actorLabel)
		case "/class-enrollments/deleted":
			activityType = "DELETE"
			description = fmt.Sprintf("%s membatalkan penempatan rombel siswa", actorLabel)
		default:
			activityType = "READ"
			description = fmt.Sprintf("%s melihat histori penempatan kelas siswa", actorLabel)
		}
		return
	}

	// 12. Modul Wali Murid (Guardians)
	if strings.HasPrefix(path, "/guardians") {
		module = "Wali Murid"
		switch path {
		case "/guardians/tambah":
			activityType = "CREATE"
			description = fmt.Sprintf("%s mendaftarkan profil wali murid baru", actorLabel)
		case "/guardians/update":
			activityType = "UPDATE"
			description = fmt.Sprintf("%s memperbarui data wali murid", actorLabel)
		case "/guardians/deleted":
			activityType = "DELETE"
			description = fmt.Sprintf("%s menghapus data wali murid", actorLabel)
		case "/guardians/link-student":
			activityType = "CREATE"
			description = fmt.Sprintf("%s menautkan wali murid dengan anak/siswa", actorLabel)
		case "/guardians/unlink-student":
			activityType = "DELETE"
			description = fmt.Sprintf("%s memutuskan tautan wali murid dari siswa", actorLabel)
		default:
			activityType = "READ"
			description = fmt.Sprintf("%s melihat rincian informasi wali murid", actorLabel)
		}
		return
	}

	// 13. Modul Dashboard
	if strings.HasPrefix(path, "/dashboard") {
		module = "Dashboard"
		activityType = "READ"
		description = fmt.Sprintf("%s membuka ringkasan infografis dashboard sistem", actorLabel)
		return
	}

	// Default fallback ramah
	module = "Sistem"
	switch method {
	case http.MethodGet:
		activityType = "READ"
		description = fmt.Sprintf("%s melihat data pada menu sistem", actorLabel)
	case http.MethodPost:
		activityType = "CREATE"
		description = fmt.Sprintf("%s mengirimkan data baru ke sistem", actorLabel)
	case http.MethodPut, http.MethodPatch:
		activityType = "UPDATE"
		description = fmt.Sprintf("%s memperbarui informasi di sistem", actorLabel)
	case http.MethodDelete:
		activityType = "DELETE"
		description = fmt.Sprintf("%s menghapus data di sistem", actorLabel)
	default:
		activityType = "OTHER"
		description = fmt.Sprintf("%s melakukan aktivitas di sistem", actorLabel)
	}

	return
}

// GetUserIDByEmail ambil ID user berdasarkan email
func GetUserIDByEmail(db *pgxpool.Pool, email string) (string, error) {
	var userID string
	query := "SELECT id FROM users WHERE email = $1"
	err := db.QueryRow(context.Background(), query, email).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// maskHeaders menyembunyikan data sensitif seperti token Authorization dan cookie
func maskHeaders(headers http.Header) string {
	cloned := make(http.Header)
	for k, v := range headers {
		lowerK := strings.ToLower(k)
		if lowerK == "authorization" {
			cloned[k] = []string{"Bearer [REDACTED_TOKEN]"}
		} else if lowerK == "cookie" {
			cloned[k] = []string{"[REDACTED_COOKIE]"}
		} else {
			cloned[k] = v
		}
	}
	bytes, _ := json.Marshal(cloned)
	return string(bytes)
}

// maskSensitiveData menyembunyikan field password, access_token, refresh_token, secret, dsb (AGENTS.md Section 18.5)
func maskSensitiveData(data string) string {
	var bodyMap map[string]interface{}
	if err := json.Unmarshal([]byte(data), &bodyMap); err != nil {
		return data
	}

	sensitiveFields := []string{
		"password", "kata_sandi", "old_password", "new_password",
		"access_token", "refresh_token", "token", "secret", "private_key",
		"authorization", "api_key",
	}

	for _, field := range sensitiveFields {
		for k := range bodyMap {
			if strings.EqualFold(k, field) {
				bodyMap[k] = "[REDACTED]"
			}
		}
	}

	hashedData, _ := json.Marshal(bodyMap)
	return string(hashedData)
}
