package model

import (
	"github.com/agriplant/utils"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email" gorm:"unique"`
	Password string `json:"password" form:"password"`
}

type ProductResponse struct {
	ID       uint      `json:"id"`
	Pictures []Picture `json:"product_pictures"`
	Name     string    `json:"product_name"`
	Price    int       `json:"product_price"`
	Seen     int       `json:"product_seen"`
}

func (u *User) BeforeCreateUser(tx *gorm.DB) (err error) {
	hashPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashPassword

	return
}

func (u *User) ComparePassword(password string) string {
	utils.ComparePassword(u.Password, password)
	return u.Password
}
