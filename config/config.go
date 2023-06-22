package config

import (
	"fmt"

	"github.com/agriplant/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	godotenv.Load(".env")
	InitDB()
	InitialMigration()
}

type Config struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_Name     string
}

func InitDB() {
	config := Config{
		DB_Username: "developergolang",
		DB_Password: "plantagridb123",
		DB_Port:     "3306",
		DB_Host:     "mysql",
		DB_Name:     "agriplant_db",
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB_Username,
		config.DB_Password,
		config.DB_Host,
		config.DB_Port,
		config.DB_Name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func InitialMigration() {
	DB.AutoMigrate(
		&model.Admin{},
		&model.User{},
		&model.Article{},
		&model.Product{},
		&model.Weather{},
		&model.InfoWeather{},
		&model.LikedArticles{},
		&model.Plant{},
		&model.WateringInfo{},
		&model.TemperatureInfo{},
		&model.FertilizingInfo{},
		&model.PlantingInfo{},
		&model.ContainerInfo{},
		&model.GroundInfo{},
		&model.WeeklyProgress{},
		&model.Picture{},
		&model.MyPlant{},
		&model.Watering{},
		&model.Fertilizing{},
		&model.Complaints{},
		&model.Suggestions{},
		&model.Notification{},
	)
}
