package route

import (
	"github.com/labstack/echo/v4"
	admin "github.com/agriplant/controller/admin"
	user "github.com/agriplant/controller/user"
)

func New() *echo.Echo{
	e:= echo.New()

	e.GET("/users/hello", user.User_Hello_World)
	e.GET("/admins/hello", admin.Admin_Hello_World)

	return e
}