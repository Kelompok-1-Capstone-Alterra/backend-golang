package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}
