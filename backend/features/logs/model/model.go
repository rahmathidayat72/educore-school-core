package model

import (
	"go_rest_native_sekolah/features/logs"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LogModel struct {
	db *pgxpool.Pool
}

// NewLogModel membuat instance DataLogInterface dengan pool koneksi PostgreSQL
func NewLogModel(db *pgxpool.Pool) logs.DataLogInterface {
	return &LogModel{
		db: db,
	}
}
