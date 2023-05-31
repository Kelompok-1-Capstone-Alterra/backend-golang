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
