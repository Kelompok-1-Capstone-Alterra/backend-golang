package admin

import (
	"github.com/agriplant/config"
	"github.com/agriplant/model"
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
