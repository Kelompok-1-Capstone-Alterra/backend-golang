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
	e.GET("/alldb", controller.Show_all_DB)
	e.GET("/alldb/plants", controller.Show_all_DB_Plants)
	e.GET("/alldb/myplants", controller.Show_all_DB_MyPlants)
	e.GET("/alldb/admins", controller.Show_all_DB_Admins)
	e.GET("/alldb/users", controller.Show_all_DB_Users)
	e.POST("/pictures", controller.Upload_pictures)
	e.GET("/pictures/:url", controller.Get_picture)
	e.DELETE("/pictures/:url", controller.Delete_picture_from_local)

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
	// Landing Page
	eAuth.GET("/admins/overview", admin.GetOverview)

	// Article
	eAuth.POST("/admins/articles/add", admin.CreateArticle)
	eAuth.GET("/admins/articles", admin.GetArticles)
	eAuth.GET("/admins/articles/search", admin.GetArticlesByTitle)
	eAuth.GET("/admins/articles/:id/detail", admin.GetArticleByID)
	eAuth.PUT("/admins/articles/:id", admin.UpdateArticleByID)
	eAuth.DELETE("/admins/articles/:id", admin.DeleteArticleByID)

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

	// Suggestions
	eAuth.GET("/admins/suggestions", admin.GetAllSuggestions)
	eAuth.GET("/admins/suggestions/:suggestion_id", admin.GetSuggestionByID)
	eAuth.DELETE("/admins/suggestions/:suggestion_id", admin.DeleteSuggestionByID)

	// ENDPOINT MOBILE (with token)
	// Recomendation
	eAuth.GET("/users/products", user.GetProducts)
	eAuth.GET("/users/products/search", user.GetProductsByName)
	eAuth.GET("/users/products/:category", user.GetProductsByCategory)
	eAuth.GET("/users/products/:category/search", user.GetProductsByCategoryAndName)
	eAuth.GET("/users/products/:id/detail", user.GetProductByID)
	eAuth.GET("/users/products/:id/contact", user.GetProductContactByID)

	// Explore & Monitoring
	eAuth.GET("/users/weather/:latitude/:longitude", user.Get_weather)
	eAuth.GET("/users/weather/:label_id", user.Get_weather_article)
	eAuth.POST("/users/plants/notifications", user.Generate_notifications)
	eAuth.GET("/plants", user.Get_available_plants)
	eAuth.GET("/plants/search", user.Search_available_plants)
	eAuth.GET("/plants/:plant_id", user.Get_plant_detail)
	eAuth.GET("/plants/:plant_id/location", user.Get_plant_location)
	eAuth.POST("/plants/:plant_id", user.Add_my_plant)
	eAuth.GET("/users/plants/:myplant_id/name", user.Get_myplant_name)
	eAuth.PUT("/users/plants/:myplant_id/name", user.Update_myplant_name)
	eAuth.POST("/users/plants/:myplant_id/start", user.Start_planting)
	eAuth.GET("/users/plants/:myplant_id/overview", user.Get_myplant_overview)
	eAuth.POST("/users/plants/:myplant_id/watering", user.Add_watering)
	eAuth.POST("/users/plants/:myplant_id/fertilizing", user.Add_fertilizing)
	eAuth.POST("/users/plants/:myplant_id/progress", user.Add_weekly_progress)
	eAuth.GET("/users/plants/:myplant_id/progress", user.Get_all_myplant_weekly_progress)
	eAuth.GET("/users/plants/:myplant_id/progress/:weekly_progress_id", user.Get_my_plant_weekly_progress_by_id)
	eAuth.PUT("/users/plants/progress/:weekly_progress_id", user.Update_weekly_progress)
	eAuth.POST("/users/plants/:myplant_id/progress/dead", user.Add_dead_plant_progress)
	eAuth.POST("/users/plants/:myplant_id/progress/harvest", user.Add_harvest_plant_progress)
	eAuth.GET("/articles/planting/:plant_id/:location", user.GetPlantingArticle)
	eAuth.GET("/articles/fertilizing/:plant_id", user.GetFertilizingArticle)
	eAuth.GET("/articles/watering/:plant_id", user.GetWateringArticle)
	eAuth.GET("/articles/temperature/:plant_id", user.GetTemperatureArticle)

	// MyPlants
	eAuth.GET("/users/plants", user.GetMyPlantList)
	eAuth.GET("/users/plants/search", user.GetMyPlantListBYKeyword)
	eAuth.DELETE("/users/plants", user.DeleteMyPlants)

	// Articles (with token)
	eAuth.GET("/users/articles/trending", user.GetArticlesTrending)
	eAuth.GET("/users/articles/latest", user.GetArticlesLatest)
	eAuth.GET("/users/articles/:id", user.GetArticlesbyID)
	eAuth.GET("/users/articles/liked", user.GetArticlesLiked)
	eAuth.POST("/users/articles/:article_id/liked", user.AddLikes)
	eAuth.DELETE("/users/articles/:article_id/liked", user.DeleteLikes)

	// Settings
	eAuth.GET("/users/profiles", user.GetProfile)
	eAuth.GET("/users/profiles/name", user.GetUsername)
	eAuth.PUT("/users/profiles/name", user.UpdateUsername)
	eAuth.PUT("/users/profiles/password", user.UpdateUserPassword)
	eAuth.GET("/users/plants/stats", user.GetMyPlantsStats)
	eAuth.POST("/users/helps", user.SendComplaintEmail)
	eAuth.POST("/users/suggestions", user.SendSuggestion)
	eAuth.PUT("/users/profiles/pictures", user.UpdateProfilePicture)

	return e
}

func JWTMiddleware() echo.MiddlewareFunc {
	config := mid.JWTConfig{
		SigningKey: []byte(constant.SECRET_JWT),
	}
	return mid.JWTWithConfig(config)
}
