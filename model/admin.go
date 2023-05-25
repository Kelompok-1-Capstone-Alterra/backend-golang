package model

import (
	"github.com/agriplant/utils"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Name     string `json:"admin_name"`
	Email    string `json:"admin_email"`
	Password string `json:"admin_password"`
}

func (a *Admin) BeforeCreateAdmin(tx *gorm.DB) (err error) {
	hashPassword, err := utils.HashPassword(a.Password)
	if err != nil {
		return err
	}
	a.Password = hashPassword

	return
}

func (a *Admin) ComparePassword(password string) string {
	utils.ComparePassword(a.Password, password)
	return a.Password
}
