package admin

import (
	"log"
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/agriplant/utils"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func CreateAdmin(c echo.Context) error {
	admin := model.Admin{}

	c.Bind(&admin)

	admin.BeforeCreateAdmin(config.DB)

	if err := config.DB.Save(&admin).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(500, map[string]interface{}{
			"message": "Failed to create admin",
		})
	}

	return c.JSON(200, map[string]interface{}{
		"message": "Success create admin",
		"data":    admin,
	})
}

func GetAdmins(c echo.Context) error {
	admins := []model.Admin{}

	if err := config.DB.Find(&admins).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(500, map[string]interface{}{
			"message": "Failed to get admins",
		})
	}

	return c.JSON(200, map[string]interface{}{
		"message": "Success get admins",
		"data":    admins,
	})
}

func LoginAdmin(c echo.Context) error {
	admin := model.Admin{}

	var loginData struct {
		Email    string `json:"admin_email" validate:"required"`
		Password string `json:"admin_password" validate:"required"`
	}

	err := c.Bind(&loginData)
	if err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "Failed to bind data",
		})
	}

	// check if inputed email is empty
	if loginData.Email == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  401,
			"message": "Failed to login, Email cannot be empty",
		})
	}

	// check if inputed password is empty
	if loginData.Password == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  401,
			"message": "Failed to login, Password cannot be empty",
		})
	}

	if err := config.DB.Where("email = ?", loginData.Email).First(&admin).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  401,
			"message": "Failed to login, invalid email",
		})
	}

	// Verify the password
	if !utils.ComparePassword(admin.Password, loginData.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": "Invalid password",
		})
	}

	// Create token
	token, err := utils.CreateTokenAdmin(admin.ID, admin.Name)
	if err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(500, map[string]interface{}{
			"message": "Failed to create token",
		})
	}

	return c.JSON(200, map[string]interface{}{
		"message": "Success login",
		"data": map[string]interface{}{
			"admin": admin,
			"token": token,
		},
	})
}

func GetOverview(c echo.Context) error {
	// Count how many users are registered
	var countUser int64
	if err := config.DB.Model(&model.User{}).Count(&countUser).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Count how many plants are created
	var countPlant int64
	if err := config.DB.Model(&model.Plant{}).Count(&countPlant).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Count how many articles are created
	var countArticle int64
	if err := config.DB.Model(&model.Article{}).Count(&countArticle).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Count how many products are created
	var countProduct int64
	if err := config.DB.Model(&model.Product{}).Count(&countProduct).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Retrieve plant summaries
	var plantSummaries []struct {
		PlantName  string
		TotalUsers int64
	}
	if err := config.DB.
		Table("my_plants").
		Select("plants.name AS plant_name, COUNT(DISTINCT my_plants.user_id) AS total_users").
		Joins("JOIN plants ON my_plants.plant_id = plants.id").
		Where("plants.deleted_at IS NULL"). // Filter out deleted plants
		Group("my_plants.plant_id").
		Order("total_users DESC").
		Limit(10). // Limit the number of plant summaries to 10
		Find(&plantSummaries).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Prepare the plant summaries
	plantData := make([]map[string]interface{}, 0)
	for _, plantSummary := range plantSummaries {
		plantData = append(plantData, map[string]interface{}{
			"plant_name":  plantSummary.PlantName,
			"total_users": plantSummary.TotalUsers,
		})
	}

	// Retrieve weather data from InfoWeather table, sort from the newest and group by location
	var weatherData []model.InfoWeather
	if err := config.DB.
		Raw("SELECT * FROM info_weathers iw WHERE iw.created_at = (SELECT MAX(created_at) FROM info_weathers WHERE location = iw.location)").
		Limit(10). // Limit the number of weather entries to 10
		Find(&weatherData).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Prepare the weather data response with selected fields
	var weatherResponse []map[string]interface{}
	for _, weather := range weatherData {
		weatherResponse = append(weatherResponse, map[string]interface{}{
			"location":    weather.Location,
			"temperature": weather.Temperature,
			"label":       weather.Label,
		})
	}

	// Prepare the response data
	response := map[string]interface{}{
		"metrics_summary": map[string]interface{}{
			"total_users":    countUser,
			"total_plants":   countPlant,
			"total_articles": countArticle,
			"total_products": countProduct,
		},
		"weather_summary": weatherResponse,
		"plant_summary": map[string]interface{}{
			"plant": plantData,
		},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "Success get overview",
		"data":    response,
	})
}
