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
	Articles []Article `json:"-" gorm:"foreignKey:AdminID"`
	Products []Product `json:"-" gorm:"foreignKey:AdminID"`
	Weathers []Weather `json:"-" gorm:"foreignKey:AdminID"`
	Plants   []Plant   `json:"-" gorm:"foreignKey:AdminID"`
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
	Label       string    `json:"weather_label"`
	Pictures    []Picture `json:"weather_pictures" gorm:"foreignKey:WeatherID"`
	Description string    `json:"weather_description"`
	AdminID     uint      `json:"-"`
}

type Plant struct {
	gorm.Model      `json:"-"`
	Name            string          `json:"plant_name"`
	Latin           string          `json:"plant_latin"`
	Description     string          `json:"plant_description"`
	Pictures        []Picture       `json:"plant_pictures" gorm:"foreignKey:PlantID"`
	WateringInfo    WateringInfo    `json:"watering_info" gorm:"foreignKey:PlantID"`
	TemperatureInfo TemperatureInfo `json:"temperature_info" gorm:"foreignKey:PlantID"`
	FertilizingInfo FertilizingInfo `json:"fertilizing_info" gorm:"foreignKey:PlantID"`
	PlantingInfo    PlantingInfo    `json:"planting_info" gorm:"foreignKey:PlantID"`
	MyPlant         MyPlant         `json:"my_plant" gorm:"foreignKey:PlantID"`
	AdminID         uint            `json:"-"`
}

type WateringInfo struct {
	gorm.Model  `json:"-"`
	Period      int       `json:"watering_period"`
	Pictures    []Picture `json:"watering_pictures" gorm:"foreignKey:WateringInfoID"`
	Description string    `json:"watering_description"`
	PlantID     uint      `json:"-"`
}

type TemperatureInfo struct {
	gorm.Model  `json:"-"`
	Min         int       `json:"temperature_min"`
	Max         int       `json:"temperature_max"`
	Description string    `json:"temperature_description"`
	Pictures    []Picture `json:"temperature_pictures" gorm:"foreignKey:TemperatureInfoID"`
	PlantID     uint      `json:"-"`
}

type FertilizingInfo struct {
	gorm.Model  `json:"-"`
	Limit       int       `json:"fertilizing_limit"`
	Period      int       `json:"fertilizing_period"`
	Pictures    []Picture `json:"fertilizing_pictures" gorm:"foreignKey:FertilizingInfoID"`
	Description string    `json:"fertilizing_description"`
	PlantID     uint      `json:"-"`
}

type PlantingInfo struct {
	gorm.Model    `json:"-"`
	Container     bool          `json:"planting_container"`
	Ground        bool          `json:"planting_ground"`
	ContainerInfo ContainerInfo `json:"container_info" gorm:"foreignKey:PlantingInfoID"`
	GroundInfo    GroundInfo    `json:"ground_info" gorm:"foreignKey:PlantingInfoID"`
	PlantID       uint          `json:"-"`
}

type ContainerInfo struct {
	gorm.Model     `json:"-"`
	Instructions   string    `json:"container_instruction"`
	Materials      string    `json:"container_materials"`
	Video          string    `json:"container_video"`
	Pictures       []Picture `json:"container_pictures" gorm:"foreignKey:ContainerInfoID"`
	PlantingInfoID uint      `json:"-"`
}

type GroundInfo struct {
	gorm.Model     `json:"-"`
	Instructions   string    `json:"ground_instruction"`
	Materials      string    `json:"ground_materials"`
	Video          string    `json:"ground_video"`
	Pictures       []Picture `json:"ground_pictures" gorm:"foreignKey:GroundInfoID"`
	PlantingInfoID uint      `json:"-"`
}

type Picture struct {
	gorm.Model        `json:"-"`
	URL               string `json:"url"`
	ArticleID         *uint  `json:"-"`
	ProductID         *uint  `json:"-"`
	WeatherID         *uint  `json:"-"`
	PlantID           *uint  `json:"-"`
	WateringInfoID    *uint  `json:"-"`
	TemperatureInfoID *uint  `json:"-"`
	FertilizingInfoID *uint  `json:"-"`
	ContainerInfoID   *uint  `json:"-"`
	GroundInfoID      *uint  `json:"-"`
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
