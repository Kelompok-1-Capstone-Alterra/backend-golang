package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Hello_World(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "Hello World. OK",
		"no_test": 4,
	})
}
