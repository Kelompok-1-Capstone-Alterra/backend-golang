package controller

import (
	"html/template"
	"log"
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func Show_all_DB(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	type TemplateData struct {
		Admins   []model.Admin
		Users    []model.User
		Plants   []model.Plant
		MyPlants []model.MyPlant
		Pictures []model.Picture
		// Articles []model.Article
		// Articles []model.Article
		// Articles []model.Article
		// Articles []model.Article
		// Articles []model.Article
		// Articles []model.Article
		// Articles []model.Article
		// Articles []model.Article
	}

	// Get Admins
	var admins []model.Admin
	if err_find_admins := config.DB.Find(&admins).Error; err_find_admins != nil {
		log.Print(color.RedString(err_find_admins.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get Users
	var users []model.User
	if err_find_users := config.DB.Find(&users).Error; err_find_users != nil {
		log.Print(color.RedString(err_find_users.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get Plants
	var plants []model.Plant
	if err_find_plants := config.DB.Find(&plants).Error; err_find_plants != nil {
		log.Print(color.RedString(err_find_plants.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get MyPlants
	var myPlants []model.MyPlant
	if err_find_myplants := config.DB.Find(&myPlants).Error; err_find_myplants != nil {
		log.Print(color.RedString(err_find_myplants.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get Pictures
	var pictures []model.Picture
	if err_find_pictures := config.DB.Find(&pictures).Error; err_find_pictures != nil {
		log.Print(color.RedString(err_find_pictures.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	data := TemplateData{
		Admins:   admins,
		Users:    users,
		Plants:   plants,
		MyPlants: myPlants,
		Pictures: pictures,
	}

	err := tmpl.Execute(c.Response().Writer, data)
	if err != nil {
		return err
	}

	return nil
}

func Show_all_DB_Plants(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles("templates/index_plants.html"))
	type TemplateData struct {
		WateringInfos    []model.WateringInfo
		TemperatureInfos []model.TemperatureInfo
		PlantingInfos    []model.PlantingInfo
		FertilizingInfos []model.FertilizingInfo
		ContainerInfos   []model.ContainerInfo
		GroundInfos      []model.GroundInfo
	}

	// Get WateringInfos
	var wateringInfos []model.WateringInfo
	if err_find_wateringInfos := config.DB.Find(&wateringInfos).Error; err_find_wateringInfos != nil {
		log.Print(color.RedString(err_find_wateringInfos.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get TemperatureInfos
	var temperatureInfos []model.TemperatureInfo
	if err_find_temperatureInfos := config.DB.Find(&temperatureInfos).Error; err_find_temperatureInfos != nil {
		log.Print(color.RedString(err_find_temperatureInfos.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get PlantingInfos
	var plantingInfos []model.PlantingInfo
	if err_find_plantingInfos := config.DB.Find(&plantingInfos).Error; err_find_plantingInfos != nil {
		log.Print(color.RedString(err_find_plantingInfos.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get FertilizingInfos
	var fertilizingInfos []model.FertilizingInfo
	if err_find_fertilizingInfos := config.DB.Find(&fertilizingInfos).Error; err_find_fertilizingInfos != nil {
		log.Print(color.RedString(err_find_fertilizingInfos.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get ContainerInfos
	var containerInfos []model.ContainerInfo
	if err_find_containerInfos := config.DB.Find(&containerInfos).Error; err_find_containerInfos != nil {
		log.Print(color.RedString(err_find_containerInfos.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get GroundInfos
	var groundInfos []model.GroundInfo
	if err_find_groundInfos := config.DB.Find(&groundInfos).Error; err_find_groundInfos != nil {
		log.Print(color.RedString(err_find_groundInfos.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	data := TemplateData{
		WateringInfos:    wateringInfos,
		TemperatureInfos: temperatureInfos,
		PlantingInfos:    plantingInfos,
		FertilizingInfos: fertilizingInfos,
		ContainerInfos:   containerInfos,
		GroundInfos:      groundInfos,
	}

	err := tmpl.Execute(c.Response().Writer, data)
	if err != nil {
		return err
	}

	return nil
}

func Show_all_DB_MyPlants(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles("templates/index_myplants.html"))
	type TemplateData struct {
		Waterings      []model.Watering
		Fertilizings   []model.Fertilizing
		WeeklyProgress []model.WeeklyProgress
	}

	// Get Waterings
	var waterings []model.Watering
	if err_find_waterings := config.DB.Find(&waterings).Error; err_find_waterings != nil {
		log.Print(color.RedString(err_find_waterings.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get Fertilizings
	var fertilizings []model.Fertilizing
	if err_find_fertilizings := config.DB.Find(&fertilizings).Error; err_find_fertilizings != nil {
		log.Print(color.RedString(err_find_fertilizings.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get WeeklyProgress
	var WeeklyProgress []model.WeeklyProgress
	if err_find_WeeklyProgress := config.DB.Find(&WeeklyProgress).Error; err_find_WeeklyProgress != nil {
		log.Print(color.RedString(err_find_WeeklyProgress.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	data := TemplateData{
		Waterings:      waterings,
		Fertilizings:   fertilizings,
		WeeklyProgress: WeeklyProgress,
	}

	err := tmpl.Execute(c.Response().Writer, data)
	if err != nil {
		return err
	}

	return nil
}

func Show_all_DB_Admins(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles("templates/index_admins.html"))
	type TemplateData struct {
		Weathers      []model.Weather
		Products      []model.Product
		Articles      []model.Article
		LikedArticles []model.LikedArticles
	}

	// Get Weathers
	var weathers []model.Weather
	if err_find_weathers := config.DB.Find(&weathers).Error; err_find_weathers != nil {
		log.Print(color.RedString(err_find_weathers.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get Products
	var products []model.Product
	if err_find_products := config.DB.Find(&products).Error; err_find_products != nil {
		log.Print(color.RedString(err_find_products.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get Articles
	var articles []model.Article
	if err_find_articles := config.DB.Find(&articles).Error; err_find_articles != nil {
		log.Print(color.RedString(err_find_articles.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get LikedArticles
	var likedArticles []model.LikedArticles
	if err_find_likedArticles := config.DB.Find(&likedArticles).Error; err_find_likedArticles != nil {
		log.Print(color.RedString(err_find_likedArticles.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	data := TemplateData{
		Weathers:      weathers,
		Products:      products,
		Articles:      articles,
		LikedArticles: likedArticles,
	}

	err := tmpl.Execute(c.Response().Writer, data)
	if err != nil {
		return err
	}

	return nil
}

func Show_all_DB_Users(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles("templates/index_users.html"))
	type TemplateData struct {
		InfoWeathers []model.InfoWeather
		Suggestions  []model.Suggestions
		Complaints   []model.Complaints
	}

	// Get Weathers
	var infoWeathers []model.InfoWeather
	if err_find_infoWeathers := config.DB.Find(&infoWeathers).Error; err_find_infoWeathers != nil {
		log.Print(color.RedString(err_find_infoWeathers.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get Suggestions
	var suggestions []model.Suggestions
	if err_find_suggestions := config.DB.Find(&suggestions).Error; err_find_suggestions != nil {
		log.Print(color.RedString(err_find_suggestions.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get Suggestions
	var complaints []model.Complaints
	if err_find_complaints := config.DB.Find(&complaints).Error; err_find_complaints != nil {
		log.Print(color.RedString(err_find_complaints.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	data := TemplateData{
		InfoWeathers: infoWeathers,
		Suggestions:  suggestions,
		Complaints:   complaints,
	}

	err := tmpl.Execute(c.Response().Writer, data)
	if err != nil {
		return err
	}

	return nil
}
