package helper

import (
	"context"
	"encoding/json"
	"go_rest_native_sekolah/config"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type MetaToken struct {
	ID           string  `json:"id"`
	Role         string  `json:"role"`
	SchoolUnitID *string `json:"school_unit_id,omitempty"`
	Exp          string  `json:"exp"`
}

type AccessToken struct {
	Claims MetaToken
}

// SignToken digunakan untuk membuat token JWT baru berdasarkan data yang diberikan.
// Fungsi ini mengembalikan token yang ditandatangani, waktu kedaluwarsa, dan error jika ada.
func SignToken(data map[string]interface{}) (string, time.Time, error) {
	// Menetapkan waktu kedaluwarsa token secara hardcode
	expiryTime := time.Now().UTC().Add(time.Hour * 24) // Waktu kedaluwarsa 24 jam di UTC

	// Membuat klaim untuk token
	claims := jwt.MapClaims{}
	claims["exp"] = expiryTime.Unix() // Menetapkan waktu kedaluwarsa ke klaim

	// Menambahkan data tambahan ke klaim
	for key, value := range data {
		claims[key] = value
	}

	// Membuat token baru dengan metode penandatanganan HS256 dan klaim yang telah dibuat
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan secret key dari environment variable JWT_SECRET
	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", time.Time{}, err
	}

	// Mengembalikan token yang ditandatangani, waktu kedaluwarsa, dan nil untuk error
	return accessToken, expiryTime, nil
}

// Middleware untuk memverifikasi token JWT
// Middleware ini akan memverifikasi apakah token yang dikirimkan lewat header Authorization
// valid dan sesuai dengan secret key yang diatur di environment variable JWT_SECRET
// Jika token tidak valid maka akan dikembalikan error 401 Unauthorized
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mendapatkan header Authorization dari request
		authHeader := r.Header.Get("Authorization")

		// Mendapatkan token dari header Authorization
		tokenString := GetTokenFromAuthorizationHeader(authHeader)

		// Jika token tidak ditemukan maka akan dikembalikan error 401 Unauthorized
		if tokenString == "" {
			JSONResponse(w, http.StatusUnauthorized, "Authorization token required")
			return
		}

		// Memverifikasi token yang diterima
		// Jika token tidak valid maka akan dikembalikan error 401 Unauthorized
		_, err := VerifyTokenHeader(tokenString)
		if err != nil {
			JSONResponse(w, http.StatusUnauthorized, "Invalid token: "+err.Error())
			return
		}

		// Jika token valid maka akan dilanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	}
}

// RolesAllowedMiddleware memastikan request dilakukan oleh pengguna dengan salah satu role yang diizinkan
// Memeriksa token JWT sekaligus memverifikasi status role langsung dari database users
func RolesAllowedMiddleware(db *pgxpool.Pool, allowedRoles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenString := GetTokenFromAuthorizationHeader(authHeader)

		if tokenString == "" {
			JSONResponse(w, http.StatusUnauthorized, "Authorization token required")
			return
		}

		meta, err := VerifyTokenHeader(tokenString)
		if err != nil {
			JSONResponse(w, http.StatusUnauthorized, "Invalid token: "+err.Error())
			return
		}

		userRole := strings.ToLower(strings.TrimSpace(meta.Role))

		// Verifikasi role langsung dari database users untuk keamanan ketat
		if db != nil && meta.ID != "" {
			var dbRole string
			query := "SELECT role FROM users WHERE id = $1 AND delete_at IS NULL"
			err := db.QueryRow(context.Background(), query, meta.ID).Scan(&dbRole)
			if err != nil {
				JSONResponse(w, http.StatusForbidden, "Akses ditolak: Data pengguna tidak ditemukan atau dinonaktifkan")
				return
			}
			userRole = strings.ToLower(strings.TrimSpace(dbRole))
		}

		isAllowed := false
		if userRole == "super_admin" {
			isAllowed = true
		} else {
			for _, role := range allowedRoles {
				if userRole == strings.ToLower(strings.TrimSpace(role)) {
					isAllowed = true
					break
				}
			}
		}

		if !isAllowed {
			JSONResponse(w, http.StatusForbidden, "Akses ditolak: Anda tidak memiliki hak akses ke resource ini")
			return
		}

		next.ServeHTTP(w, r)
	}
}

// AdminOnlyMiddleware memastikan request dilakukan oleh pengguna dengan role admin atau super_admin
// Memeriksa token JWT sekaligus memverifikasi status role langsung dari database
func AdminOnlyMiddleware(db *pgxpool.Pool, next http.HandlerFunc) http.HandlerFunc {
	return RolesAllowedMiddleware(db, []string{"admin", "super_admin"}, next)
}

// SuperAdminOnlyMiddleware memastikan request khusus dilakukan oleh Super Admin Yayasan
func SuperAdminOnlyMiddleware(db *pgxpool.Pool, next http.HandlerFunc) http.HandlerFunc {
	return RolesAllowedMiddleware(db, []string{"super_admin"}, next)
}

// UserSelfOrAdminMiddleware mengizinkan admin untuk melihat/mengakses siapa saja,
// sedangkan role lain (guru/user) hanya boleh mengakses datanya sendiri berdasarkan ID
func UserSelfOrAdminMiddleware(db *pgxpool.Pool, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenString := GetTokenFromAuthorizationHeader(authHeader)

		if tokenString == "" {
			JSONResponse(w, http.StatusUnauthorized, "Authorization token required")
			return
		}

		meta, err := VerifyTokenHeader(tokenString)
		if err != nil {
			JSONResponse(w, http.StatusUnauthorized, "Invalid token: "+err.Error())
			return
		}

		userRole := strings.ToLower(strings.TrimSpace(meta.Role))

		if db != nil && meta.ID != "" {
			var dbRole string
			query := "SELECT role FROM users WHERE id = $1 AND delete_at IS NULL"
			err := db.QueryRow(context.Background(), query, meta.ID).Scan(&dbRole)
			if err != nil {
				JSONResponse(w, http.StatusForbidden, "Akses ditolak: Data pengguna tidak ditemukan atau dinonaktifkan")
				return
			}
			userRole = strings.ToLower(strings.TrimSpace(dbRole))
		}

		if userRole == "admin" || userRole == "super_admin" {
			next.ServeHTTP(w, r)
			return
		}

		// Jika bukan admin, cek apakah target ID sesuai dengan ID pengguna sendiri
		targetID := r.URL.Query().Get("id")
		if targetID != "" && targetID == meta.ID {
			next.ServeHTTP(w, r)
			return
		}

		JSONResponse(w, http.StatusForbidden, "Akses ditolak: Anda hanya dapat mengakses data akun Anda sendiri")
	}
}

// Verifikasi token JWT yang diterima dari header Authorization
// Token yang diterima harus sesuai dengan secret key yang diatur di environment variable JWT_SECRET
// Jika token tidak valid maka akan dikembalikan error
func VerifyTokenHeader(requestToken string) (MetaToken, error) {
	// Buat token JWT baru
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		// Gunakan secret key yang diatur di environment variable JWT_SECRET
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Jika terjadi error maka akan dikembalikan error
	if err != nil {
		log.Println(err)
		return MetaToken{}, err
	}

	// Jika token tidak valid maka akan dikembalikan error
	if !token.Valid {
		log.Println("Token tidak valid")
		return MetaToken{}, jwt.ErrSignatureInvalid
	}

	// Cek klaim jwt.MapClaims terlebih dahulu
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		var meta MetaToken
		if id, ok := claims["id"].(string); ok {
			meta.ID = id
		}
		if role, ok := claims["role"].(string); ok {
			meta.Role = role
		}
		if unitID, ok := claims["school_unit_id"].(string); ok && unitID != "" {
			meta.SchoolUnitID = &unitID
		}
		return meta, nil
	}

	// Fallback dekode token menjadi MetaToken
	claimToken := DecodeToken(token)
	return claimToken.Claims, nil
}

func VerifyToken(accessToken string) (*jwt.Token, error) {
	config.LoadEnv() // Pastikan environment variable sudah di-load
	// Ambil secret key dari environment variable
	jwtSecretKey := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	if !token.Valid {
		logrus.Error("Token is not valid")
		return nil, jwt.ErrSignatureInvalid
	}

	return token, nil
}

func DecodeToken(accessToken *jwt.Token) AccessToken {
	var token AccessToken
	stringify, err := json.Marshal(&accessToken)
	if err != nil {
		return token
	}
	err = json.Unmarshal(stringify, &token)
	if err != nil {
		return token
	}
	return token
}

// Fungsi untuk mengambil token dari header Authorization
func GetTokenFromAuthorizationHeader(authorizationHeader string) string {
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

// ExtractUserContext membaca klaim ID, Role, dan SchoolUnitID dari header Authorization request
func ExtractUserContext(r *http.Request) (id string, role string, unitID string, isSuperAdmin bool) {
	authHeader := r.Header.Get("Authorization")
	tokenString := GetTokenFromAuthorizationHeader(authHeader)
	if tokenString == "" {
		return "", "", "", false
	}

	meta, err := VerifyTokenHeader(tokenString)
	if err != nil {
		return "", "", "", false
	}

	currentRole := strings.ToLower(strings.TrimSpace(meta.Role))
	superAdmin := currentRole == "super_admin"

	currentUnit := ""
	if meta.SchoolUnitID != nil {
		currentUnit = *meta.SchoolUnitID
	}

	// Jika Super Admin mengirimkan query param ?unit_id=..., gunakan unit tersebut untuk filter
	if superAdmin {
		paramUnit := r.URL.Query().Get("unit_id")
		if paramUnit != "" {
			currentUnit = paramUnit
		}
	}

	return meta.ID, currentRole, currentUnit, superAdmin
}

// CORSMiddleware menangani Cross-Origin Resource Sharing (CORS) dengan pengetatan origin (CORS Hardening)
func CORSMiddleware(next http.Handler) http.Handler {
	allowedOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowedOrigins []string
	if strings.TrimSpace(allowedOriginsEnv) != "" {
		for _, o := range strings.Split(allowedOriginsEnv, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowOrigin := "*"

		// Jika ada konfigurasi spesifik CORS_ALLOWED_ORIGINS, validasi origin yang datang
		if len(allowedOrigins) > 0 {
			matched := false
			for _, allowed := range allowedOrigins {
				if allowed == "*" || strings.EqualFold(origin, allowed) {
					allowOrigin = origin
					matched = true
					break
				}
			}
			if !matched && origin != "" {
				// Origin tidak ada di daftar izinkan
				allowOrigin = allowedOrigins[0]
			}
		} else if origin != "" {
			// Default development: echo request origin jika aman
			allowOrigin = origin
		}

		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-School-Unit-ID")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

