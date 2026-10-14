package logs

import "time"

type (
	// LogCore merepresentasikan entitas log transaksi sistem
	LogCore struct {
		ID                  int       `json:"id"`
		Timestamp           time.Time `json:"timestamp"`
		UserID              string    `json:"user_id"`
		Username            string    `json:"username"`
		Email               string    `json:"email"`
		Role                string    `json:"role"`
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

	// FilterLog merepresentasikan filter dan paginasi pencarian log
	FilterLog struct {
		Search       string
		Result       string
		Module       string
		ActivityType string
		Limit        int
		Offset       int
		NoCount      bool
	}

	// DataLogInterface merepresentasikan interface data repository untuk transaction_logs
	DataLogInterface interface {
		SelectAllLogs(filter FilterLog) ([]LogCore, int, error)
		SelectLogByID(id int) (*LogCore, error)
	}

	// ServiceLogInterface merepresentasikan interface business logic untuk logs
	ServiceLogInterface interface {
		GetAllLogs(filter FilterLog) ([]LogCore, int, error)
		GetLogByID(id int) (*LogCore, error)
	}
)
