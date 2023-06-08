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

	e.Use(mid.CORS())

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
	e.PUT("/users/:user_id/password", user.Reset_password)
	e.GET("/users/emails/check", user.Check_email_valid)

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
	eAuth.GET("/admins/products/search", admin.GetProductsByName)

	// Weather Management
	eAuth.POST("/admins/weathers/add", admin.CreateWeather)
	eAuth.GET("/admins/weathers", admin.GetWeathers)
	eAuth.GET("/admins/weathers/:id/detail", admin.GetWeatherByID)
	eAuth.PUT("/admins/weathers/:id", admin.UpdateWeatherByID)
	eAuth.DELETE("/admins/weathers/:id", admin.DeleteWeatherByID)

	// Plant
	eAuth.GET("/admins/plants/search", admin.GetPlantsByKeyword)
	eAuth.GET("/admins/plants", admin.GetPlants)
	eAuth.GET("/admins/plants/:id/detail", admin.GetPlantDetails)
	eAuth.PUT("/admins/plants/:id/detail", admin.UpdatePlantDetails)
	eAuth.DELETE("/admins/plants/:id/detail", admin.DeletePlantDetails)
	eAuth.POST("/admins/plants/add", admin.CreatePlant)

	// ENDPOINT MOBILE (with token)
	// Recomendation
	eAuth.GET("/users/products", user.GetProducts)
	eAuth.GET("/users/products/search", user.GetProductsByName)
	eAuth.GET("/users/products/:category", user.GetProductsByCategory)
	eAuth.GET("/users/products/:category/search", user.GetProductsByCategoryAndName)
	eAuth.GET("/users/products/:id/detail", user.GetProductByID)

	// Explore & Monitoring
	eAuth.GET("/users/weather", user.Get_weather)
	eAuth.GET("/users/weather/:label_id", user.Get_weather_article)
	eAuth.GET("/plants", user.Get_available_plants)
	eAuth.GET("/plants/search", user.Search_available_plants)
	eAuth.GET("/plants/:plant_id", user.Get_plant_detail)
	eAuth.GET("/plants/:plant_id/location", user.Get_plant_location)
	eAuth.POST("/plants/:plant_id", user.Add_my_plant)

	eAuth.GET("/articles/planting/:plant_id/:location", user.GetPlantingArticle)
	eAuth.GET("/articles/fertilizing/:plant_id", user.GetFertilizingArticle)
	eAuth.GET("/articles/watering/:plant_id", user.GetWateringArticle)
	eAuth.GET("/articles/temperature/:plant_id", user.GetTemperatureArticle)

	return e
}

func JWTMiddleware() echo.MiddlewareFunc {
	config := mid.JWTConfig{
		SigningKey: []byte(constant.SECRET_JWT),
	}
	return mid.JWTWithConfig(config)
}
