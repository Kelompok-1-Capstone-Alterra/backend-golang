package route

import (
	"github.com/agriplant/constant"
	controller "github.com/agriplant/controller"
	admin "github.com/agriplant/controller/admin"
	user "github.com/agriplant/controller/user"
	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	// ENDPOINT GLOBAL
	e.GET("/hello", controller.Hello_World)

	// ENDPOINT WEB (no token)
	e.POST("/admins", admin.CreateAdmin)
	e.GET("/admins", admin.GetAdmins)
	e.POST("/admins/login", admin.LoginAdmin)

	// ENDPOINT MOBILE (no token)
	e.POST("/users/register", user.Register)
	e.POST("/users/login", user.Login)

	// Protected route
	eAuth := e.Group("/auth")
	eAuth.Use(JWTMiddleware())

	// ENDPOINT WEB (with token)
	// Article
	eAuth.POST("/admins/articles/add", admin.CreateArticle)
	eAuth.GET("/admins/articles", admin.GetArticles)
	eAuth.GET("/admins/articles/search", admin.GetArticlesByKeyword)

	// ENDPOINT MOBILE (with token)

	return e
}

func JWTMiddleware() echo.MiddlewareFunc {
	config := mid.JWTConfig{
		SigningKey: []byte(constant.SECRET_JWT),
	}
	return mid.JWTWithConfig(config)
}
