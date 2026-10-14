package model

import (
	"context"
	"fmt"
	"go_rest_native_sekolah/features/logs"
	"strings"
)

// SelectAllLogs mengambil data log aktivitas sistem beserta data user (username, email, role)
func (lm *LogModel) SelectAllLogs(filter logs.FilterLog) ([]logs.LogCore, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(u.username ILIKE $%d OR u.email ILIKE $%d OR l.service_name ILIKE $%d OR l.activity_description ILIKE $%d OR l.module ILIKE $%d OR l.ip_address ILIKE $%d OR l.target_id ILIKE $%d)", argIdx, argIdx, argIdx, argIdx, argIdx, argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	if filter.Result != "" {
		conditions = append(conditions, fmt.Sprintf("l.result = $%d", argIdx))
		args = append(args, filter.Result)
		argIdx++
	}

	if filter.Module != "" {
		conditions = append(conditions, fmt.Sprintf("l.module = $%d", argIdx))
		args = append(args, filter.Module)
		argIdx++
	}

	if filter.ActivityType != "" {
		conditions = append(conditions, fmt.Sprintf("l.activity_type = $%d", argIdx))
		args = append(args, filter.ActivityType)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	// Jika no_count diaktifkan, lewati query COUNT(*) agar menghemat beban CPU & I/O database
	if !filter.NoCount {
		countQuery := fmt.Sprintf(`
			SELECT COUNT(*) 
			FROM transaction_logs l
			LEFT JOIN users u ON u.id = l.user_id
			%s
		`, whereClause)

		err := lm.db.QueryRow(context.Background(), countQuery, args...).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("error counting logs: %v", err)
		}
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	} else if limit > 200 {
		limit = 200
	}

	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	dataQuery := fmt.Sprintf(`
		SELECT 
			l.id, 
			l.timestamp, 
			l.user_id, 
			COALESCE(u.username, '') as username,
			COALESCE(u.email, '') as email,
			COALESCE(u.role, '') as role,
			COALESCE(l.activity_type, '') as activity_type,
			COALESCE(l.activity_description, '') as activity_description,
			COALESCE(l.http_method, '') as http_method,
			COALESCE(l.endpoint, '') as endpoint,
			COALESCE(l.module, '') as module,
			COALESCE(l.target_id, '') as target_id,
			COALESCE(l.status_code, 200) as status_code,
			COALESCE(l.ip_address, '') as ip_address,
			COALESCE(l.duration_ms, 0) as duration_ms,
			l.perangkat, 
			l.service_name, 
			COALESCE(l.request_body::text, '{}') as request_body,
			COALESCE(l.response_body::text, '{}') as response_body,
			COALESCE(l.request_param::text, '{}') as request_param,
			l.result, 
			COALESCE(l.header::text, '{}') as header
		FROM transaction_logs l
		LEFT JOIN users u ON u.id = l.user_id
		%s
		ORDER BY l.timestamp DESC, l.id DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := lm.db.Query(context.Background(), dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying logs: %v", err)
	}
	defer rows.Close()

	var result []logs.LogCore
	for rows.Next() {
		var logItem logs.LogCore
		err := rows.Scan(
			&logItem.ID,
			&logItem.Timestamp,
			&logItem.UserID,
			&logItem.Username,
			&logItem.Email,
			&logItem.Role,
			&logItem.ActivityType,
			&logItem.ActivityDescription,
			&logItem.HTTPMethod,
			&logItem.Endpoint,
			&logItem.Module,
			&logItem.TargetID,
			&logItem.StatusCode,
			&logItem.IPAddress,
			&logItem.DurationMS,
			&logItem.Perangkat,
			&logItem.ServiceName,
			&logItem.RequestBody,
			&logItem.ResponseBody,
			&logItem.RequestParam,
			&logItem.Result,
			&logItem.Header,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning log row: %v", err)
		}
		result = append(result, logItem)
	}

	if filter.NoCount {
		total = len(result)
	}

	return result, total, nil
}

// SelectLogByID mengambil 1 detail data log berdasarkan ID
func (lm *LogModel) SelectLogByID(id int) (*logs.LogCore, error) {
	query := `
		SELECT 
			l.id, 
			l.timestamp, 
			l.user_id, 
			COALESCE(u.username, '') as username,
			COALESCE(u.email, '') as email,
			COALESCE(u.role, '') as role,
			COALESCE(l.activity_type, '') as activity_type,
			COALESCE(l.activity_description, '') as activity_description,
			COALESCE(l.http_method, '') as http_method,
			COALESCE(l.endpoint, '') as endpoint,
			COALESCE(l.module, '') as module,
			COALESCE(l.target_id, '') as target_id,
			COALESCE(l.status_code, 200) as status_code,
			COALESCE(l.ip_address, '') as ip_address,
			COALESCE(l.duration_ms, 0) as duration_ms,
			l.perangkat, 
			l.service_name, 
			COALESCE(l.request_body::text, '{}') as request_body,
			COALESCE(l.response_body::text, '{}') as response_body,
			COALESCE(l.request_param::text, '{}') as request_param,
			l.result, 
			COALESCE(l.header::text, '{}') as header
		FROM transaction_logs l
		LEFT JOIN users u ON u.id = l.user_id
		WHERE l.id = $1
	`

	var logItem logs.LogCore
	err := lm.db.QueryRow(context.Background(), query, id).Scan(
		&logItem.ID,
		&logItem.Timestamp,
		&logItem.UserID,
		&logItem.Username,
		&logItem.Email,
		&logItem.Role,
		&logItem.ActivityType,
		&logItem.ActivityDescription,
		&logItem.HTTPMethod,
		&logItem.Endpoint,
		&logItem.Module,
		&logItem.TargetID,
		&logItem.StatusCode,
		&logItem.IPAddress,
		&logItem.DurationMS,
		&logItem.Perangkat,
		&logItem.ServiceName,
		&logItem.RequestBody,
		&logItem.ResponseBody,
		&logItem.RequestParam,
		&logItem.Result,
		&logItem.Header,
	)
	if err != nil {
		return nil, fmt.Errorf("log not found: %v", err)
	}

	return &logItem, nil
}
