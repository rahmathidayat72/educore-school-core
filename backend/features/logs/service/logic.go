package service

import (
	"errors"
	"go_rest_native_sekolah/features/logs"
)

type logService struct {
	logData logs.DataLogInterface
}

// NewServiceLog menginisialisasi service untuk logs
func NewServiceLog(logData logs.DataLogInterface) logs.ServiceLogInterface {
	if logData == nil {
		panic("NewServiceLog: logData is nil")
	}
	return &logService{
		logData: logData,
	}
}

// GetAllLogs mengambil semua log transaksi dengan filter
func (s *logService) GetAllLogs(filter logs.FilterLog) ([]logs.LogCore, int, error) {
	return s.logData.SelectAllLogs(filter)
}

// GetLogByID mengambil detail 1 log transaksi
func (s *logService) GetLogByID(id int) (*logs.LogCore, error) {
	if id <= 0 {
		return nil, errors.New("invalid log id")
	}
	return s.logData.SelectLogByID(id)
}
