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
	if err := config.DB.Order("created_at DESC").Where("user_id=? AND status NOT IN ?", user_id, []string{"harvest", "dead"}).Find(&myPlants).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	responses := []map[string]interface{}{}
	for _, myPlant := range myPlants {
		var plant model.Plant
		if err_first := config.DB.First(&plant, myPlant.PlantID).Error; err_first != nil {
			log.Print(color.RedString(err_first.Error()))
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status":  500,
				"message": "internal server error",
			})
		}
		config.DB.Model(&plant).Association("Pictures").Find(&plant.Pictures)

		response := map[string]interface{}{
			"plant_id":   myPlant.PlantID,
			"myplant_id": myPlant.ID,
			"name":       myPlant.Name,
			"location":   myPlant.Location,
			"picture":    plant.Pictures[0].URL,
			"latin":      plant.Latin,
		}

		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to retrieve latest myPlants data",
		"data":    responses,
	})
}

func GetMyPlantListBYKeyword(c echo.Context) error {
	myPlants := []model.MyPlant{}
	name := c.QueryParam("name")

	//get data by trending
	if err := config.DB.Order("created_at DESC").Where("name LIKE ?", "%"+name+"%").Find(&myPlants).Error; err != nil {
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

	// Validation1 : check if myplant_id valid
	var myPlants_check []model.MyPlant
	if err_find := config.DB.Where("id IN ?", deleteID.MyPlants_ID).Find(&myPlants_check).Error; err_find != nil {
		log.Print(color.RedString(err_find.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}
	if len(deleteID.MyPlants_ID) != len(myPlants_check) || len(deleteID.MyPlants_ID) == 0 {
		log.Print(color.RedString("there is myplant_id not valid"))
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

	// delete notification according to myplant_id
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	var notification model.Notification
	if err_del := config.DB.Where("user_id=? AND my_plant_id IN ?", user_id, deleteID.MyPlants_ID).Delete(&notification).Error; err_del != nil {
		log.Print(color.RedString(err_del.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to Delete My Plant",
	})
}
