package model

import (
	"database/sql"
	schoolunits "go_rest_native_sekolah/features/school_units"
	"time"
)

type SchoolUnit struct {
	ID          string
	NamaSekolah string
	Jenjang     string
	NPSN        sql.NullString
	Alamat      sql.NullString
	Telepon     sql.NullString
	Email       sql.NullString
	Website     sql.NullString
	Status      string
	CreatedAt   *time.Time
	UpdateAt    *time.Time
	DeleteAt    *time.Time
}

func (m *SchoolUnit) toCore() schoolunits.SchoolUnitCore {
	core := schoolunits.SchoolUnitCore{
		ID:          m.ID,
		NamaSekolah: m.NamaSekolah,
		Jenjang:     m.Jenjang,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdateAt:    m.UpdateAt,
		DeleteAt:    m.DeleteAt,
	}
	if m.NPSN.Valid {
		core.NPSN = m.NPSN.String
	}
	if m.Alamat.Valid {
		core.Alamat = m.Alamat.String
	}
	if m.Telepon.Valid {
		core.Telepon = m.Telepon.String
	}
	if m.Email.Valid {
		core.Email = m.Email.String
	}
	if m.Website.Valid {
		core.Website = m.Website.String
	}
	return core
}

func toCoreList(list []SchoolUnit) []schoolunits.SchoolUnitCore {
	cores := make([]schoolunits.SchoolUnitCore, 0, len(list))
	for _, item := range list {
		cores = append(cores, item.toCore())
	}
	return cores
}
