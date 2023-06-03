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
	admin := model.Admin{}

	if err := c.Bind(&weather); err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Check if the label already exists
	existingWeather := model.Weather{}
	result := config.DB.Where("label = ?", weather.Label).First(&existingWeather)
	if result.Error == nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request, label already exists",
		})
	}

	// Get admin by ID
	// If admin not found, return error
	if err := config.DB.First(&admin, weather.AdminID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Set admin ID to weather
	weather.AdminID = admin.ID

	// Save weather to database
	if err := config.DB.Save(&weather).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Populate Pictures field for each weather
	config.DB.Model(&weather).Association("Pictures").Find(&weather.Pictures)

	// Extract picture URLs
	pictureURLs := make([]string, len(weather.Pictures))
	for i, pic := range weather.Pictures {
		pictureURLs[i] = pic.URL
	}

	response := struct {
		ID          uint     `json:"id"`
		Created_at  string   `json:"created_at"`
		Updated_at  string   `json:"updated_at"`
		Deleted_at  string   `json:"deleted_at"`
		Title       string   `json:"weather_title"`
		Label       string   `json:"weather_label"`
		Pictures    []string `json:"weather_pictures"`
		Description string   `json:"weather_description"`
	}{
		ID:          weather.ID,
		Created_at:  weather.CreatedAt.String(),
		Updated_at:  weather.UpdatedAt.String(),
		Deleted_at:  weather.DeletedAt.Time.String(),
		Title:       weather.Title,
		Label:       weather.Label,
		Pictures:    pictureURLs,
		Description: weather.Description,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
	})
}

func GetWeathers(c echo.Context) error {
	var weathers []model.Weather

	if err := config.DB.Find(&weathers).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Iterate over each weather record and generate custom response
	var responses []interface{}
	for _, weather := range weathers {
		// Populate Pictures field for each weather
		config.DB.Model(&weather).Association("Pictures").Find(&weather.Pictures)

		// Extract picture URLs
		pictureURLs := make([]string, len(weather.Pictures))
		for i, pic := range weather.Pictures {
			pictureURLs[i] = pic.URL
		}

		response := struct {
			ID          uint     `json:"id"`
			Created_at  string   `json:"created_at"`
			Updated_at  string   `json:"updated_at"`
			Deleted_at  string   `json:"deleted_at"`
			Title       string   `json:"weather_title"`
			Label       string   `json:"weather_label"`
			Pictures    []string `json:"weather_pictures"`
			Description string   `json:"weather_description"`
		}{
			ID:          weather.ID,
			Created_at:  weather.CreatedAt.String(),
			Updated_at:  weather.UpdatedAt.String(),
			Deleted_at:  weather.DeletedAt.Time.String(),
			Title:       weather.Title,
			Label:       weather.Label,
			Pictures:    pictureURLs,
			Description: weather.Description,
		}

		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    responses,
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
	// Extract picture URLs
	pictureURLs := make([]string, len(weather.Pictures))
	for i, pic := range weather.Pictures {
		pictureURLs[i] = pic.URL
	}

	response := struct {
		ID          uint     `json:"id"`
		Created_at  string   `json:"created_at"`
		Updated_at  string   `json:"updated_at"`
		Deleted_at  string   `json:"deleted_at"`
		Title       string   `json:"weather_title"`
		Label       string   `json:"weather_label"`
		Pictures    []string `json:"weather_pictures"`
		Description string   `json:"weather_description"`
	}{
		ID:          weather.ID,
		Created_at:  weather.CreatedAt.String(),
		Updated_at:  weather.UpdatedAt.String(),
		Deleted_at:  weather.DeletedAt.Time.String(),
		Title:       weather.Title,
		Label:       weather.Label,
		Pictures:    pictureURLs,
		Description: weather.Description,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
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

	if err := config.DB.Where("id = ?", weatherID).First(&weather).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	if err := config.DB.Delete(&weather).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
