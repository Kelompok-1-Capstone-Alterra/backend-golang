package config

import (
	"fmt"
	"os"

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
		DB_Username: os.Getenv("DB_USERNAME"),
		DB_Password: os.Getenv("DB_PASSWORD"),
		DB_Port:     os.Getenv("DB_PORT"),
		DB_Host:     os.Getenv("DB_HOST"),
		DB_Name:     os.Getenv("DB_DB"),
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
		&model.Picture{},
		&model.Weather{},
		&model.InfoWeather{},
		&model.Plant{},
		&model.WateringInfo{},
		&model.TemperatureInfo{},
		&model.FertilizingInfo{},
		&model.PlantingInfo{},
		&model.ContainerInfo{},
		&model.GroundInfo{},
		&model.MyPlant{},
	)
}
