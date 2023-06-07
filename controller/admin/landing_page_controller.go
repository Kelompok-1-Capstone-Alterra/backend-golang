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
		return c.JSON(500, map[string]interface{}{
			"message": "Failed to bind data",
		})
	}

	if err := config.DB.Where("email = ?", loginData.Email).First(&admin).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(500, map[string]interface{}{
			"message": "Failed to find admin",
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
		PlantID    uint
		PlantName  string
		TotalUsers int64
	}
	if err := config.DB.
		Model(&model.MyPlant{}).
		Select("plant_id, name AS plant_name, COUNT(DISTINCT user_id) AS total_users").
		Group("plant_id, name").
		Order("total_users DESC").
		Find(&plantSummaries).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Check if plantSummaries is empty
	if len(plantSummaries) == 0 {
		// Handle the case where there are no plant summaries
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Success get overview",
			"data": map[string]interface{}{
				"metrics_summary": map[string]interface{}{
					"total_users":    countUser,
					"total_plants":   countPlant,
					"total_articles": countArticle,
					"total_products": countProduct,
				},
				"weather_summary": map[string]interface{}{
					"location":    "Jakarta",
					"temperature": "30",
					"weather":     "Rainy",
				},
				"plant_summary": map[string]interface{}{
					"plant": map[string]interface{}{
						"plant_name":  "",
						"total_users": 0,
					},
				},
			},
		})
	}

	// Prepare the response data
	response := map[string]interface{}{
		"message": "Success get overview",
		"data": map[string]interface{}{
			"metrics_summary": map[string]interface{}{
				"total_users":    countUser,
				"total_plants":   countPlant,
				"total_articles": countArticle,
				"total_products": countProduct,
			},
			"weather_summary": map[string]interface{}{
				"location":    "Jakarta",
				"temperature": "30",
				"weather":     "Rainy",
			},
			"plant_summary": map[string]interface{}{
				"plant": map[string]interface{}{
					"plant_name":  plantSummaries[0].PlantName,
					"total_users": plantSummaries[0].TotalUsers,
				},
			},
		},
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"message": "Success get overview",
		"data":    response,
	})
}
