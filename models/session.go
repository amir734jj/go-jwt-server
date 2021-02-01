package models

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model
	Id      uint `gorm:"primaryKey"`
	Token   string
	Expires time.Time
	UserId  uint
}
