package admin

import (
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/agriplant/utils"
	"github.com/labstack/echo/v4"
)

func CreateAdmin(c echo.Context) error {
	admin := model.Admin{}

	c.Bind(&admin)

	admin.BeforeCreateAdmin(config.DB)

	if err := config.DB.Save(&admin).Error; err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "Failed to create admin",
			"error":   err.Error(),
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
		return c.JSON(500, map[string]interface{}{
			"message": "Failed to get admins",
			"error":   err.Error(),
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
		return c.JSON(500, map[string]interface{}{
			"message": "Failed to bind data",
			"error":   err.Error(),
		})
	}

	if err := config.DB.Where("email = ?", loginData.Email).First(&admin).Error; err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "Failed to find admin",
			"error":   err.Error(),
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
		return c.JSON(500, map[string]interface{}{
			"message": "Failed to create token",
			"error":   err.Error(),
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
