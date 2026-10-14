package helper

import (
	"net/http"
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
