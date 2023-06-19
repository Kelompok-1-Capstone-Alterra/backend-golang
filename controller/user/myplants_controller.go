package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/agriplant/utils"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func GetMyPlantList(c echo.Context) error {
	var myPlants []model.MyPlant
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	//get data by trending
	if err := config.DB.Where("user_id=?", user_id).Find(&myPlants).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	var plants []model.Plant
	for _, idPlant := range myPlants {
		var plant model.Plant

		if err := config.DB.First(&plant, idPlant.PlantID).Error; err != nil {
			log.Print(color.RedString(err.Error()))
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  400,
				"message": "bad request",
			})
		}

		plants = append(plants, plant)
	}

	data := []map[string]interface{}{}
	//Populate Pictures field for each article
	for i := 0; i < len(plants); i++ {
		config.DB.Model(&plants[i]).Association("Pictures").Find(&plants[i].Pictures)
		result := map[string]interface{}{
			"myplant_id": myPlants[i].ID,
			"name":       myPlants[i].Name,
			"picture":    plants[i].Pictures[0].URL,
			"latin":      plants[i].Latin,
		}
		data = append(data, result)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to retrieve latest myPlants data",
		"data":    data,
	})
}

func GetMyPlantListBYKeyword(c echo.Context) error {
	myPlants := []model.MyPlant{}
	name := c.QueryParam("name")

	//get data by trending
	if err := config.DB.Where("name LIKE ?", "%"+name+"%").Find(&myPlants).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	var plants []model.Plant
	for _, idPlant := range myPlants {
		var plant model.Plant

		if err := config.DB.First(&plant, idPlant.PlantID).Error; err != nil {
			log.Print(color.RedString(err.Error()))
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  400,
				"message": "bad request",
			})
		}
		plants = append(plants, plant)
	}
	var data []map[string]interface{}
	//var expic []
	//Populate Pictures field for each article
	for i := 0; i < len(plants); i++ {
		config.DB.Model(&plants[i]).Association("Pictures").Find(&plants[i].Pictures)
		result := map[string]interface{}{
			"myplant_id": myPlants[i].ID,
			"name":       myPlants[i].Name,
			"picture":    plants[i].Pictures[0].URL,
			"latin":      plants[i].Latin,
		}
		data = append(data, result)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to retrieve latest myPlants data",
		"data":    data,
	})
}

func DeleteMyPlants(c echo.Context) error {
	type DeleteID struct {
		MyPlants_ID []int `json:"myplants_id"`
	}

	var deleteID DeleteID
	var myPlants model.MyPlant

	if err := c.Bind(&deleteID); err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if err := config.DB.Where("id IN ?", deleteID.MyPlants_ID).Delete(&myPlants).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to Delete My Plant",
	})
}
