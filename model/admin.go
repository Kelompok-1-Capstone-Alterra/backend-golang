package model

import (
	"github.com/agriplant/utils"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Name     string    `json:"admin_name"`
	Email    string    `json:"admin_email"`
	Password string    `json:"admin_password"`
	Articles []Article `gorm:"foreignKey:AdminID"`
}

type Article struct {
	gorm.Model
	Title       string    `json:"article_title"`
	Pictures    []Picture `json:"article_pictures" gorm:"foreignKey:ArticleID"`
	Description string    `json:"article_description"`
	View        int       `json:"article_view"`
	Like        int       `json:"article_like"`
	AdminID     uint      `json:"admin_id"`
}

type Picture struct {
	gorm.Model
	URL       string `json:"article_url"`
	ArticleID uint   `json:"article_id"`
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
