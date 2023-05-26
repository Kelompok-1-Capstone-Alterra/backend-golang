package controller

import (
	"fmt"
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	var user model.User

	// binding struct
	if err_bind := c.Bind(&user); err_bind != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// hashing password
	user.BeforeCreateUser(config.DB)

	// register
	if err_insert := config.DB.Save(&user).Error; err_insert != nil {
		fmt.Println(err_insert.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "ok",
	})
}
