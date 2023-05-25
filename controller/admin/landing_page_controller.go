package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Admin_Hello_World(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Hello World, welcome to the admin endpoint",
	})
}
