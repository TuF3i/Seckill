package dao

import "gorm.io/gorm"

type DaoReliance struct {
	Pgdb *gorm.DB
}

type Dao struct {
	*DaoReliance
}

func NewDao(m *DaoReliance) *Dao {
	return &Dao{m}
}
