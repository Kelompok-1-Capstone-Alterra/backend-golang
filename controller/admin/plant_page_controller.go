package admin

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

func CreatePlant(c echo.Context) error {
	plant := model.Plant{}

	if err := c.Bind(&plant); err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Get admin id from JWT token
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	adminId, _ := utils.GetAdminIDFromToken(token)

	plant.AdminID = adminId

	// save plant to database
	if err := config.DB.Save(&plant).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to add plant",
		"data":    plant,
	})
}

func GetPlants(c echo.Context) error {
	type Response struct {
		ID          uint
		Name        string
		Latin       string
		Description string
		Pictures    []model.Picture
		Watering    int
		Fertilizing int
		Min         int
		Max         int
	}
	plants := []model.Plant{}
	responses := []Response{}

	// Get all plants
	if err := config.DB.Find(&plants).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}
	//Populate all field for each plant to response
	for i := 0; i < len(plants); i++ {
		var temp Response
		temp.ID = plants[i].ID
		picture := []model.Picture{}
		config.DB.Where("plant_id = ?", plants[i].ID).Find(&picture)
		temp.Pictures = picture
		temp.Name = plants[i].Name
		temp.Latin = plants[i].Latin
		temp.Description = plants[i].Description
		watering := model.WateringInfo{}
		config.DB.Where("plant_id = ?", plants[i].ID).Find(&watering)
		temp.Watering = watering.Period
		fertilizing := model.FertilizingInfo{}
		config.DB.Where("plant_id = ?", plants[i].ID).Find(&fertilizing)
		temp.Fertilizing = fertilizing.Period
		temperature := model.TemperatureInfo{}
		config.DB.Where("plant_id = ?", plants[i].ID).Find(&temperature)
		temp.Min = temperature.Min
		temp.Max = temperature.Max
		responses = append(responses, temp)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get list of plants",
		"data":    responses,
	})
}

func GetPlantsByKeyword(c echo.Context) error {
	type Response struct {
		ID          uint
		Name        string
		Latin       string
		Description string
		Pictures    []model.Picture
		Watering    int
		Fertilizing int
		Min         int
		Max         int
	}
	keyword := c.QueryParam("keyword")

	plants := []model.Plant{}
	responses := []Response{}

	// Get all plants
	if err := config.DB.Where("`name` LIKE ?", "%"+keyword+"%").Find(&plants).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	//Populate all field for each plant to response
	for i := 0; i < len(plants); i++ {
		var temp Response
		temp.ID = plants[i].ID
		picture := []model.Picture{}
		config.DB.Where("plant_id = ?", plants[i].ID).Find(&picture)
		temp.Pictures = picture
		temp.Name = plants[i].Name
		temp.Latin = plants[i].Latin
		temp.Description = plants[i].Description
		watering := model.WateringInfo{}
		config.DB.Where("plant_id = ?", plants[i].ID).Find(&watering)
		temp.Watering = watering.Period
		fertilizing := model.FertilizingInfo{}
		config.DB.Where("plant_id = ?", plants[i].ID).Find(&fertilizing)
		temp.Fertilizing = fertilizing.Period
		temperature := model.TemperatureInfo{}
		config.DB.Where("plant_id = ?", plants[i].ID).Find(&temperature)
		temp.Min = temperature.Min
		temp.Max = temperature.Max
		responses = append(responses, temp)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get plant by name",
		"data":    responses,
	})
}

func GetPlantDetails(c echo.Context) error {
	plant := model.Plant{}

	plantID := c.Param("id")

	// Get plant by id
	if err := config.DB.First(&plant, plantID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	//Get all related data
	config.DB.Model(&plant).Association("Pictures").Find(&plant.Pictures)
	config.DB.Model(&plant).Association("WateringInfo").Find(&plant.WateringInfo)
	config.DB.Model(&plant.WateringInfo).Association("Pictures").Find(&plant.WateringInfo.Pictures)
	config.DB.Model(&plant).Association("FertilizingInfo").Find(&plant.FertilizingInfo)
	config.DB.Model(&plant.FertilizingInfo).Association("Pictures").Find(&plant.FertilizingInfo.Pictures)
	config.DB.Model(&plant).Association("TemperatureInfo").Find(&plant.TemperatureInfo)
	config.DB.Model(&plant.TemperatureInfo).Association("Pictures").Find(&plant.TemperatureInfo.Pictures)
	config.DB.Model(&plant).Association("PlantingInfo").Find(&plant.PlantingInfo)
	config.DB.Model(&plant.PlantingInfo).Association("ContainerInfo").Find(&plant.PlantingInfo.ContainerInfo)
	config.DB.Model(&plant.PlantingInfo.ContainerInfo).Association("Pictures").Find(&plant.PlantingInfo.ContainerInfo.Pictures)
	config.DB.Model(&plant.PlantingInfo).Association("GroundInfo").Find(&plant.PlantingInfo.GroundInfo)
	config.DB.Model(&plant.PlantingInfo.GroundInfo).Association("Pictures").Find(&plant.PlantingInfo.GroundInfo.Pictures)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to retrieve plant detailed information",
		"data":    plant,
	})
}

func UpdatePlantDetails(c echo.Context) error {
	plant := model.Plant{}

	plantID := c.Param("id")

	// Updating plant details
	if err := config.DB.First(&plant, plantID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if err := c.Bind(&plant); err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Updating pictures on plant details
	pictureDB := []model.Picture{}
	if err := config.DB.Where("plant_id = ?", plantID).Find(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if err := config.DB.Unscoped().Delete(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	pictures := plant.Pictures
	for i := 0; i < len(pictures); i++ {
		pictures[i].PlantID = &plant.ID
	}

	if err := config.DB.Save(pictures).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Updating watering info on plant details
	watering := model.WateringInfo{}
	if err := config.DB.Where("plant_id = ?", plantID).First(&watering).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	watering.Period = plant.WateringInfo.Period
	watering.Pictures = plant.WateringInfo.Pictures
	watering.Description = plant.WateringInfo.Description

	if err := config.DB.Model(&watering).Updates(model.WateringInfo{Period: watering.Period, Pictures: watering.Pictures, Description: watering.Description}).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Updating pictures on watering info
	pictureDB = []model.Picture{}
	if err := config.DB.Where("watering_info_id = ?", watering.ID).Find(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if err := config.DB.Unscoped().Delete(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	pictures = plant.WateringInfo.Pictures
	for i := 0; i < len(pictures); i++ {
		pictures[i].WateringInfoID = &watering.ID
	}

	if err := config.DB.Save(pictures).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Updating temperature info on plant details
	temperature := model.TemperatureInfo{}
	if err := config.DB.Where("plant_id = ?", plantID).First(&temperature).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	temperature.Min = plant.TemperatureInfo.Min
	temperature.Max = plant.TemperatureInfo.Max
	temperature.Description = plant.TemperatureInfo.Description

	if err := config.DB.Model(&temperature).Updates(model.TemperatureInfo{Min: temperature.Min, Max: temperature.Max, Description: temperature.Description}).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Updating pictures on temperature info
	pictureDB = []model.Picture{}
	if err := config.DB.Where("temperature_info_id = ?", temperature.ID).Find(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if err := config.DB.Unscoped().Delete(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	pictures = plant.TemperatureInfo.Pictures
	for i := 0; i < len(pictures); i++ {
		pictures[i].TemperatureInfoID = &temperature.ID
	}

	if err := config.DB.Save(pictures).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Updating fertilizing info on plant details
	fertilizing := model.FertilizingInfo{}
	if err := config.DB.Where("plant_id = ?", plantID).First(&fertilizing).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	fertilizing.Limit = plant.FertilizingInfo.Limit
	fertilizing.Period = plant.FertilizingInfo.Period
	fertilizing.Description = plant.FertilizingInfo.Description

	if err := config.DB.Model(&fertilizing).Updates(model.FertilizingInfo{Limit: fertilizing.Limit, Period: fertilizing.Period, Description: fertilizing.Description}).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Updating pictures on fertilizing info
	pictureDB = []model.Picture{}
	if err := config.DB.Where("fertilizing_info_id = ?", fertilizing.ID).Find(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if err := config.DB.Unscoped().Delete(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	pictures = plant.FertilizingInfo.Pictures
	for i := 0; i < len(pictures); i++ {
		pictures[i].FertilizingInfoID = &fertilizing.ID
	}

	if err := config.DB.Save(pictures).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Updating planting info on plant details
	planting := model.PlantingInfo{}
	if err := config.DB.Where("plant_id = ?", plantID).First(&planting).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	planting.Container = plant.PlantingInfo.Container
	planting.Ground = plant.PlantingInfo.Ground

	if err := config.DB.Model(&planting).Updates(model.PlantingInfo{Container: planting.Container, Ground: planting.Ground}).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Updating container info on planting info
	container := model.ContainerInfo{}
	if err := config.DB.Where("planting_info_id = ?", planting.ID).First(&container).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	container.Instructions = plant.PlantingInfo.ContainerInfo.Instructions
	container.Materials = plant.PlantingInfo.ContainerInfo.Materials
	container.Video = plant.PlantingInfo.ContainerInfo.Video

	if err := config.DB.Model(&container).Updates(model.ContainerInfo{Instructions: container.Instructions, Materials: container.Materials, Video: container.Video}).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Updating pictures on container info
	pictureDB = []model.Picture{}
	if err := config.DB.Where("container_info_id = ?", container.ID).Find(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if err := config.DB.Unscoped().Delete(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	pictures = plant.PlantingInfo.ContainerInfo.Pictures
	for i := 0; i < len(pictures); i++ {
		pictures[i].ContainerInfoID = &container.ID
	}

	if err := config.DB.Save(pictures).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Updating ground info on planting info
	ground := model.GroundInfo{}
	if err := config.DB.Where("planting_info_id = ?", planting.ID).First(&ground).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	ground.Instructions = plant.PlantingInfo.GroundInfo.Instructions
	ground.Materials = plant.PlantingInfo.GroundInfo.Materials
	ground.Video = plant.PlantingInfo.GroundInfo.Video

	if err := config.DB.Model(&ground).Updates(model.GroundInfo{Instructions: ground.Instructions, Materials: ground.Materials, Video: ground.Video}).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Updating pictures on ground info
	pictureDB = []model.Picture{}
	if err := config.DB.Where("ground_info_id = ?", ground.ID).Find(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if err := config.DB.Unscoped().Delete(&pictureDB).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	pictures = plant.PlantingInfo.GroundInfo.Pictures
	for i := 0; i < len(pictures); i++ {
		pictures[i].GroundInfoID = &ground.ID
	}

	if err := config.DB.Save(pictures).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Save plant details
	if err := config.DB.Model(&plant).Select("name", "latin", "description").Updates(model.Plant{Name: plant.Name, Latin: plant.Latin, Description: plant.Description}).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to update plant detailed information",
		"data":    plant,
	})
}

func DeletePlantDetails(c echo.Context) error {
	plant := model.Plant{}

	plantID := c.Param("id")

	// Get plant by id
	if err := config.DB.First(&plant, plantID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Delete plant from database
	if err := config.DB.Delete(&plant).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Delete all related data
	planting_info := model.PlantingInfo{}
	config.DB.Where("plant_id = ?", plantID).First(&planting_info)
	config.DB.Where("planting_info_id = ?", planting_info.ID).Delete(&model.ContainerInfo{})
	config.DB.Where("planting_info_id = ?", planting_info.ID).Delete(&model.GroundInfo{})
	config.DB.Where("plant_id = ?", plantID).Delete(&model.Picture{})
	config.DB.Where("plant_id = ?", plantID).Delete(&model.WateringInfo{})
	config.DB.Where("plant_id = ?", plantID).Delete(&model.FertilizingInfo{})
	config.DB.Where("plant_id = ?", plantID).Delete(&model.TemperatureInfo{})
	config.DB.Where("plant_id = ?", plantID).Delete(&model.PlantingInfo{})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to delete plant",
	})
}
