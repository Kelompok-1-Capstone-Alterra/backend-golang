package route

import (
	admin "github.com/agriplant/controller/admin"
	user "github.com/agriplant/controller/user"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	e.POST("/admins", admin.CreateAdmin)
	e.GET("/admins", admin.GetAdmins)
	e.POST("/admins/login", admin.LoginAdmin)
	e.GET("/users/hello", user.User_Hello_World)

	return e
}
