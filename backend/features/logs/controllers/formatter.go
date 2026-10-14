package controllers

import (
	"encoding/json"
	"go_rest_native_sekolah/features/logs"
	"time"
)

type LogResponse struct {
	ID                  int         `json:"id"`
	Timestamp           time.Time   `json:"timestamp"`
	UserID              string      `json:"user_id"`
	Username            string      `json:"username"`
	Email               string      `json:"email"`
	Role                string      `json:"role"`
	ActivityType        string      `json:"activity_type"`
	ActivityDescription string      `json:"activity_description"`
	HTTPMethod          string      `json:"http_method"`
	Endpoint            string      `json:"endpoint"`
	Module              string      `json:"module"`
	TargetID            string      `json:"target_id"`
	StatusCode          int         `json:"status_code"`
	IPAddress           string      `json:"ip_address"`
	DurationMS          int64       `json:"duration_ms"`
	Perangkat           string      `json:"perangkat"`
	ServiceName         string      `json:"service_name"`
	RequestBody         interface{} `json:"request_body"`
	ResponseBody        interface{} `json:"response_body"`
	RequestParam        interface{} `json:"request_param"`
	Result              string      `json:"result"`
	Header              interface{} `json:"header"`
}

type LogPaginationResponse struct {
	Total int           `json:"total"`
	Limit int           `json:"limit"`
	Page  int           `json:"page"`
	Data  []LogResponse `json:"data"`
}

func parseJSONOrRaw(raw string) interface{} {
	if raw == "" {
		return map[string]interface{}{}
	}
	var parsed interface{}
	if err := json.Unmarshal([]byte(raw), &parsed); err == nil {
		return parsed
	}
	return raw
}

// ToLogResponse mengonversi domain LogCore ke LogResponse dengan JSON body yang diparse
func ToLogResponse(core logs.LogCore) LogResponse {
	return LogResponse{
		ID:                  core.ID,
		Timestamp:           core.Timestamp,
		UserID:              core.UserID,
		Username:            core.Username,
		Email:               core.Email,
		Role:                core.Role,
		ActivityType:        core.ActivityType,
		ActivityDescription: core.ActivityDescription,
		HTTPMethod:          core.HTTPMethod,
		Endpoint:            core.Endpoint,
		Module:              core.Module,
		TargetID:            core.TargetID,
		StatusCode:          core.StatusCode,
		IPAddress:           core.IPAddress,
		DurationMS:          core.DurationMS,
		Perangkat:           core.Perangkat,
		ServiceName:         core.ServiceName,
		RequestBody:         parseJSONOrRaw(core.RequestBody),
		ResponseBody:        parseJSONOrRaw(core.ResponseBody),
		RequestParam:        parseJSONOrRaw(core.RequestParam),
		Result:              core.Result,
		Header:              parseJSONOrRaw(core.Header),
	}
}

// ToLogResponses mengonversi slice domain LogCore ke slice LogResponse
func ToLogResponses(cores []logs.LogCore) []LogResponse {
	list := make([]LogResponse, 0, len(cores))
	for _, c := range cores {
		list = append(list, ToLogResponse(c))
	}
	return list
}
