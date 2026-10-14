package helper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShouldLogRequest(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		statusCode int
		expected   bool
	}{
		{
			name:       "GET berhasil 200 - tidak perlu dicatat",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			expected:   false,
		},
		{
			name:       "GET berhasil 201 - tidak perlu dicatat",
			method:     http.MethodGet,
			statusCode: http.StatusCreated,
			expected:   false,
		},
		{
			name:       "HEAD berhasil 200 - tidak perlu dicatat",
			method:     http.MethodHead,
			statusCode: http.StatusOK,
			expected:   false,
		},
		{
			name:       "GET gagal 401 Unauthorized - wajib dicatat untuk audit keamanan",
			method:     http.MethodGet,
			statusCode: http.StatusUnauthorized,
			expected:   true,
		},
		{
			name:       "GET gagal 403 Forbidden - wajib dicatat untuk audit keamanan",
			method:     http.MethodGet,
			statusCode: http.StatusForbidden,
			expected:   true,
		},
		{
			name:       "GET gagal 404 Not Found - wajib dicatat untuk audit keamanan",
			method:     http.MethodGet,
			statusCode: http.StatusNotFound,
			expected:   true,
		},
		{
			name:       "GET error 500 Internal Server Error - wajib dicatat",
			method:     http.MethodGet,
			statusCode: http.StatusInternalServerError,
			expected:   true,
		},
		{
			name:       "OPTIONS preflight - tidak pernah dicatat",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			expected:   false,
		},
		{
			name:       "POST tambah data 201 - selalu dicatat",
			method:     http.MethodPost,
			statusCode: http.StatusCreated,
			expected:   true,
		},
		{
			name:       "POST tambah data gagal 400 - selalu dicatat",
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			expected:   true,
		},
		{
			name:       "PUT update data 200 - selalu dicatat",
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			expected:   true,
		},
		{
			name:       "DELETE hapus data 200 - selalu dicatat",
			method:     http.MethodDelete,
			statusCode: http.StatusOK,
			expected:   true,
		},
		{
			name:       "PATCH ubah data 200 - selalu dicatat",
			method:     http.MethodPatch,
			statusCode: http.StatusOK,
			expected:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ShouldLogRequest(tc.method, tc.statusCode)
			if got != tc.expected {
				t.Errorf("ShouldLogRequest(%s, %d) = %v; want %v", tc.method, tc.statusCode, got, tc.expected)
			}
		})
	}
}

func TestMaskSensitiveData(t *testing.T) {
	rawJSON := `{"email":"admin@sekolah.id","password":"SuperSecretPassword123","token":"jwt-token-val","nama":"Admin Sekolah"}`
	sanitized := maskSensitiveData(rawJSON)

	if strings.Contains(sanitized, "SuperSecretPassword123") {
		t.Errorf("Sanitasi gagal: password masih bocor di string output: %s", sanitized)
	}
	if strings.Contains(sanitized, "jwt-token-val") {
		t.Errorf("Sanitasi gagal: token masih bocor di string output: %s", sanitized)
	}
	if !strings.Contains(sanitized, "[REDACTED]") {
		t.Errorf("Sanitasi gagal: label [REDACTED] tidak ditemukan: %s", sanitized)
	}
	if !strings.Contains(sanitized, "admin@sekolah.id") {
		t.Errorf("Sanitasi keliru: email non-kritis terhapus: %s", sanitized)
	}
}

func TestDetermineActivityDetail_NewModules(t *testing.T) {
	// School Units
	mod, act, _, desc := DetermineActivityDetail("Super Admin", "POST", "/school-units/tambah", nil, nil, 201)
	if mod != "Unit Sekolah" || act != "CREATE" {
		t.Errorf("DetermineActivityDetail school-units failed: mod=%s, act=%s", mod, act)
	}
	if !strings.Contains(desc, "menambahkan unit sekolah") {
		t.Errorf("DetermineActivityDetail desc mismatch: %s", desc)
	}

	// Academic Years
	mod, act, _, _ = DetermineActivityDetail("Admin", "POST", "/academic-years/activate", nil, nil, 200)
	if mod != "Tahun Ajaran" || act != "UPDATE" {
		t.Errorf("DetermineActivityDetail academic-years failed: mod=%s, act=%s", mod, act)
	}

	// Guardians
	mod, act, _, _ = DetermineActivityDetail("Admin", "POST", "/guardians/link-student", nil, nil, 201)
	if mod != "Wali Murid" || act != "CREATE" {
		t.Errorf("DetermineActivityDetail guardians failed: mod=%s, act=%s", mod, act)
	}
}

func TestCORSMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	handler := CORSMiddleware(nextHandler)

	req, _ := http.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 for OPTIONS, got %d", rec.Code)
	}

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "http://localhost:5173" {
		t.Errorf("Expected Allow-Origin to reflect request origin, got %s", allowOrigin)
	}

	allowCreds := rec.Header().Get("Access-Control-Allow-Credentials")
	if allowCreds != "true" {
		t.Errorf("Expected Allow-Credentials true, got %s", allowCreds)
	}
}
