package route

import (
	"github.com/agriplant/constant"
	admin "github.com/agriplant/controller/admin"
	user "github.com/agriplant/controller/user"
	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	e.POST("/admins", admin.CreateAdmin)
	e.GET("/admins", admin.GetAdmins)
	e.POST("/admins/login", admin.LoginAdmin)

	// Protected route
	eAuth := e.Group("/auth")
	eAuth.Use(JWTMiddleware())
	// Article
	eAuth.POST("/admins/articles/add", admin.CreateArticle)
	eAuth.GET("/admins/articles", admin.GetArticles)
	eAuth.GET("/admins/articles/search", admin.GetArticlesByKeyword)

	e.GET("/users/hello", user.User_Hello_World)

	return e
}

func JWTMiddleware() echo.MiddlewareFunc {
	config := mid.JWTConfig{
		SigningKey: []byte(constant.SECRET_JWT),
	}
	return mid.JWTWithConfig(config)
}
