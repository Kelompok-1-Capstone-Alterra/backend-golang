package route

import (
	admin "github.com/agriplant/controller/admin"
	user "github.com/agriplant/controller/user"
	"github.com/agriplant/middleware"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {

	e := echo.New()

	e.Use(middleware.MiddlewareLogging)
	e.HTTPErrorHandler = middleware.ErrorHandler

	// ENDPOINT WEB
	e.POST("/admins", admin.CreateAdmin)
	e.GET("/admins", admin.GetAdmins)
	e.POST("/admins/login", admin.LoginAdmin)

	// ENDPOINT MOBILE
	e.POST("/users/register", user.Register)
	e.POST("/users/login", user.Login)

	return e
}
