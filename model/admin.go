package model

import (
	"github.com/agriplant/utils"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Name     string    `json:"admin_name"`
	Email    string    `json:"admin_email" gorm:"unique"`
	Password string    `json:"admin_password"`
	Articles []Article `gorm:"foreignKey:AdminID"`
	Products []Product `gorm:"foreignKey:AdminID"`
	Weathers []Weather `gorm:"foreignKey:AdminID"`
}

type Product struct {
	gorm.Model
	Pictures    []Picture `json:"product_pictures" gorm:"foreignKey:ProductID"`
	Name        string    `json:"product_name"`
	Category    string    `json:"product_category"`
	Description string    `json:"product_description"`
	Price       int       `json:"product_price"`
	Seen        int       `json:"product_seen"`
	Status      bool      `json:"product_status"`
	Brand       string    `json:"product_brand"`
	Condition   string    `json:"product_condition"`
	Unit        int       `json:"product_unit"`
	Weight      int       `json:"product_weight"`
	Form        string    `json:"product_form"`
	SellerName  string    `json:"product_seller_name"`
	SellerPhone string    `json:"product_seller_phone"`
	AdminID     uint      `json:"admin_id"`
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

// Struct for save weather article made by admin
type Weather struct {
	gorm.Model
	Title       string    `json:"weather_title"`
	Label       string    `json:"weather_label" gorm:"unique"`
	Pictures    []Picture `json:"weather_pictures" gorm:"foreignKey:WeatherID"`
	Description string    `json:"weather_description"`
	AdminID     uint      `json:"admin_id"`
}

type Picture struct {
	gorm.Model
	URL       string `json:"url"`
	ArticleID *uint  `json:"article_id" `
	ProductID *uint  `json:"product_id"`
	WeatherID *uint  `json:"weather_id"`
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
