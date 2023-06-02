package route

import (
	"github.com/agriplant/constant"
	"github.com/agriplant/controller"
	admin "github.com/agriplant/controller/admin"
	user "github.com/agriplant/controller/user"
	"github.com/agriplant/middleware"
	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {

	e := echo.New()

	e.Use(middleware.MiddlewareLogging)
	e.HTTPErrorHandler = middleware.ErrorHandler

	// ENDPOINT GLOBAL (no token)
	e.GET("/hello", controller.Hello_World)
	e.POST("/pictures", controller.Upload_pictures)
	e.GET("/pictures/:url", controller.Get_picture)

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

	// Product
	eAuth.POST("/admins/products/add", admin.CreateProduct)
	eAuth.GET("/admins/products", admin.GetProducts)
	eAuth.GET("/admins/products/display", admin.GetProductsDisplay)
	eAuth.GET("/admins/products/archive", admin.GetProductsArchive)
	eAuth.GET("/admins/products/:id/detail", admin.GetProductByID)
	eAuth.DELETE("/admins/products/:id", admin.DeleteProductByID)
	eAuth.PUT("/admins/products/:id", admin.UpdateProductByID)
	eAuth.GET("/admins/products/search", admin.GetProductsByKeyword)

	// ENDPOINT MOBILE (with token)
	// Recomendation
	eAuth.GET("/users/products", user.GetProducts)
	eAuth.GET("/users/products/:category", user.GetProductsByCategory)
	eAuth.GET("/users/products/:category/search", user.GetProductsByCategoryAndName)
	eAuth.GET("/users/products/:id/detail", user.GetProductByID)

	return e
}

func JWTMiddleware() echo.MiddlewareFunc {
	config := mid.JWTConfig{
		SigningKey: []byte(constant.SECRET_JWT),
	}
	return mid.JWTWithConfig(config)
}
