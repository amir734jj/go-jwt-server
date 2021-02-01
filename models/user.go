package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id       uint `gorm:"primaryKey"`
	Name     string
	Username string
	Password string
	Realm    string
	Email    string
	Sessions []Session `gorm:"foreignKey:UserId"`
}
