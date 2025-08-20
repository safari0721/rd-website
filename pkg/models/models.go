package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	AgentID  string
	Password string
}

type Accounts struct {
	gorm.Model
	AccountNo       int
	FullName        string
	OpeningDate     time.Time
	Amount          int
	TotalAmount     int
	MonthsPaid      int
	MonthsUnpaid    int
	NextInstallment time.Time
	PrevInstallment time.Time
}
