package models

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model
	Id      uint `gorm:"primaryKey"`
	Expired time.Time
	UserId  uint
}
