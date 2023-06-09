package model

import (
	"time"

	"github.com/agriplant/utils"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
	URL      string `json:"-"`
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email" gorm:"unique; not null" validate:"required, email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type ProductResponse struct {
	ID       uint     `json:"id"`
	Pictures []string `json:"product_pictures"`
	Name     string   `json:"product_name"`
	Price    int      `json:"product_price"`
	Seen     int      `json:"product_seen"`
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

type Coordinate struct {
	Latitude  string `json:"latitude" form:"latitude"`
	Longitude string `json:"longitude" form:"longitude"`
}

// Struct for save weather history for each user
type InfoWeather struct {
	gorm.Model
	User_id     uint   `json:"user_id" form:"user_id" gorm:"unique"`
	Location    string `json:"location" form:"location"`
	Temperature string `json:"temperature" form:"temperature"`
	Label       string `json:"label" form:"label"`
}

// Struct for save my plants for each user
type MyPlant struct {
	gorm.Model
	User_id           uint      `json:"user_id" form:"user_id"`
	Plant_id          uint      `json:"plant_id" form:"plant_id"`
	Name              string    `json:"name" form:"name"`
	Location          string    `json:"location" form:"location"`
	IsStartPlanting   bool      `json:"is_start_planting" form:"is_start_planting"`
	StartPlantingDate time.Time `json:"start_planting_date" form:"start_planting_date"`
	Status            string    `json:"status" form:"status"`
	Latitude          string    `json:"latitude" form:"latitude"`
	Longitude         string    `json:"longitude" form:"longitude"`
}

type Complaints struct {
	gorm.Model
	User_id     uint   `json:"-"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Notes       string `json:"notes"`
}

type Suggestions struct {
	gorm.Model
	User_id uint   `json:"-"`
	Content string `json:"content"`
}
