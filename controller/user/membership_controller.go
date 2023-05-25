package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/agriplant/config"
	_ "github.com/agriplant/model"
)

func User_Hello_World(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Hello World, welcome to the user endpoint",
	})
}
