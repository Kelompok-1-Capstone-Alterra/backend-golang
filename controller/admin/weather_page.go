package admin

import (
	"log"
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func CreateWeather(c echo.Context) error {
	weather := model.Weather{}

	c.Bind(&weather)

	admin := model.Admin{}

	// Get admin by id
	// If admin not found, return error
	if err := config.DB.First(&admin, weather.AdminID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// set admin id to weather
	weather.AdminID = admin.ID

	// save weather to database
	if err := config.DB.Save(&weather).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    weather,
	})
}

func GetWeathers(c echo.Context) error {
	weather := []model.Weather{}

	if err := config.DB.Find(&weather).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Populate Pictures field for each weather
	for i := 0; i < len(weather); i++ {
		config.DB.Model(&weather[i]).Association("Pictures").Find(&weather[i].Pictures)
	}

	// remove article_id from weather_pictures
	for i := 0; i < len(weather); i++ {
		for j := 0; j < len(weather[i].Pictures); j++ {
			weather[i].Pictures[j].ArticleID = nil
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    weather,
	})
}

func GetWeatherByID(c echo.Context) error {
	weather := model.Weather{}

	weatherID := c.Param("id")

	if err := config.DB.First(&weather, weatherID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Populate Pictures field for weather
	config.DB.Model(&weather).Association("Pictures").Find(&weather.Pictures)

	// remove article_id from weather_pictures
	for i := 0; i < len(weather.Pictures); i++ {
		weather.Pictures[i].ArticleID = nil
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    weather,
	})
}

func UpdateWeatherByID(c echo.Context) error {
	weather := model.Weather{}

	weatherID := c.Param("id")

	if err := config.DB.First(&weather, weatherID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	c.Bind(&weather)

	if err := config.DB.Save(&weather).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    weather,
	})
}

func DeleteWeatherByID(c echo.Context) error {
	weather := model.Weather{}

	weatherID := c.Param("id")

	if err := config.DB.First(&weather, weatherID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := config.DB.Delete(&weather).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
