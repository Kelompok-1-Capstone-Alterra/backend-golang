package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/agriplant/utils"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

// EXPLORE & MONITORING (Menu Home) - [Endpoint 1 : Get weather]
func Get_weather(c echo.Context) error {
	var coordinate model.Coordinate

	if err_bind := c.Bind(&coordinate); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	apikey := "869a5f0aa562d21ec64ff37c7c6c157f"
	baseURL := "https://api.openweathermap.org/data/2.5/weather"

	u, err_parseURL := url.Parse(baseURL)
	if err_parseURL != nil {
		log.Print(color.RedString(err_parseURL.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Add query parameters to the URL
	q := u.Query()
	q.Set("lat", coordinate.Latitude)
	q.Set("lon", coordinate.Longitude)
	q.Set("appid", apikey)
	q.Set("units", "metric")
	u.RawQuery = q.Encode()

	// Make the GET request
	response, err_consume := http.Get(u.String())
	if err_consume != nil {
		log.Print(color.RedString(err_consume.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}
	defer response.Body.Close()

	// Read the response body
	body, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		log.Print(color.RedString(err_read.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Create a map to store the JSON data
	var data map[string]interface{}

	// Unmarshal the response body into the map
	err_unmarshal := json.Unmarshal(body, &data)
	if err_unmarshal != nil {
		log.Print(color.RedString(err_unmarshal.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	label_id := 4
	var label string

	weather, ok := data["weather"].([]interface{})
	if ok && len(weather) > 0 {
		weatherData, ok := weather[0].(map[string]interface{})
		if ok {
			labelValue, ok := weatherData["main"].(string)
			if ok {
				switch labelValue {
				case "Thunderstorm", "Drizzle", "Rain", "Snow":
					label = "Hujan"
					label_id = 1
				case "Atmosphere":
					label = "Mendung"
					label_id = 2
				case "Clear":
					label = "Cerah"
					label_id = 3
				case "Clouds":
					label = "Berawan"
					label_id = 4
				default:
					label = "Cerah"
				}
			}
		}
	}

	city := fmt.Sprintf("%v", data["name"])
	tempereture := fmt.Sprintf("%v", data["main"].(map[string]interface{})["temp"])

	// Save weather info
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)
	save_weather_info(city, tempereture, label, user_id)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "successfully obtained weather information",
		"data": map[string]interface{}{
			"label_id":    label_id,
			"label":       label,
			"city":        city,
			"temperature": tempereture,
		},
	})
}

// HERLPER FUNCTION
func save_weather_info(location, temperature, label string, user_id uint) bool {
	var infoWeather model.InfoWeather
	err_select := config.DB.Where("user_id=?", user_id).First(&infoWeather).Error
	if err_select == nil {
		// Query update
		infoWeather.User_id = user_id
		infoWeather.Location = location
		infoWeather.Temperature = temperature
		infoWeather.Label = label

		if err_update := config.DB.Save(&infoWeather).Error; err_update != nil {
			log.Print(color.RedString(err_update.Error()))
		}
		return true
	}
	// Record not found
	// Query insert
	var infoWeather2 model.InfoWeather

	infoWeather2.User_id = user_id
	infoWeather2.Location = location
	infoWeather2.Temperature = temperature
	infoWeather2.Label = label

	if err_insert := config.DB.Save(&infoWeather2).Error; err_insert != nil {
		log.Print(color.RedString(err_insert.Error()))
	}

	return true
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 2 : Get weather article]
func Get_weather_article(c echo.Context) error {
	var weatherArticle model.Weather
	id, _ := strconv.Atoi(c.Param("label_id"))

	label := get_label_by_id(id)
	if err_first := config.DB.Where("label=?", label).First(&weatherArticle).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var picture model.Picture
	if err_first2 := config.DB.Where("weather_id=?", weatherArticle.ID).First(&picture).Error; err_first2 != nil {
		log.Print(color.RedString(err_first2.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	result := map[string]interface{}{
		"label_id":   id,
		"label":      label,
		"article_id": weatherArticle.ID,
		"picture":    picture.URL,
		"title":      weatherArticle.Title,
		"desc":       weatherArticle.Description,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to retrieve current weather information detail",
		"data":    result,
	})
}

// HELPER FUNCTION
func get_label_by_id(id int) string {
	switch id {
	case 1:
		return "Hujan"
	case 2:
		return "Mendung"
	case 3:
		return "Cerah"
	case 4:
		return "Berawan"
	default:
		return "Berawan"
	}
}

// HELPER FUNCTION
func StringToUintPointer(value string) (*uint, error) {
	intValue, err := strconv.ParseUint(value, 10, 0)
	if err != nil {
		return nil, err
	}

	uintValue := uint(intValue)
	return &uintValue, nil
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 5 : Get available plants]
func Get_available_plants(c echo.Context) error {
	var plants []model.Plant

	if err_find := config.DB.Find(&plants).Error; err_find != nil {
		log.Print(color.RedString(err_find.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var responses []map[string]interface{}
	for _, plant := range plants {
		config.DB.Model(&plant).Association("Pictures").Find(&plant.Pictures)

		var url string
		for _, picture := range plant.Pictures {
			url = picture.URL
			break
		}

		response := map[string]interface{}{
			"plant_id": plant.ID,
			"pictures": url,
			"name":     plant.Name,
			"latin":    plant.Latin,
		}
		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get list of plants",
		"data":    responses,
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 6 : Search available plants]
func Search_available_plants(c echo.Context) error {
	name := c.FormValue("name")
	name = "%" + name + "%"
	var plants []model.Plant

	if err_find := config.DB.Where("name LIKE ?", name).Find(&plants).Error; err_find != nil {
		log.Print(color.RedString(err_find.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var responses []map[string]interface{}
	for _, plant := range plants {
		config.DB.Model(&plant).Association("Pictures").Find(&plant.Pictures)

		var url string
		for _, picture := range plant.Pictures {
			url = picture.URL
			break
		}

		response := map[string]interface{}{
			"plant_id": plant.ID,
			"picture":  url,
			"name":     plant.Name,
			"latin":    plant.Latin,
		}
		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to search available plants by name",
		"data":    responses,
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 7 : Get plant detail]
func Get_plant_detail(c echo.Context) error {
	plant_id, _ := strconv.Atoi(c.Param("plant_id"))
	var plant model.Plant

	if err_first := config.DB.First(&plant, plant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	config.DB.Model(&plant).Association("Pictures").Find(&plant.Pictures)
	var url string
	for _, picture := range plant.Pictures {
		url = picture.URL
		break
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get plant detail",
		"data": map[string]interface{}{
			"plant_id":    plant.ID,
			"picture":     url,
			"name":        plant.Name,
			"latin":       plant.Latin,
			"description": plant.Description,
		},
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 8 : Get plant location]
func Get_plant_location(c echo.Context) error {
	plant_id, _ := strconv.Atoi(c.Param("plant_id"))

	var planting_info model.PlantingInfo
	if err_first := config.DB.Where("plant_id=?", plant_id).First(&planting_info).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var container_info model.ContainerInfo
	config.DB.Model(&container_info).Association("Pictures").Find(&container_info.Pictures)

	var url_container string
	for _, picture := range container_info.Pictures {
		url_container = picture.URL
		break
	}

	if err_container := config.DB.Where("planting_info_id=?", planting_info.ID).First(&container_info).Error; err_container != nil {
		log.Print(color.RedString(err_container.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var ground_info model.GroundInfo
	config.DB.Model(&ground_info).Association("Pictures").Find(&ground_info.Pictures)

	var url_ground string
	for _, picture := range container_info.Pictures {
		url_ground = picture.URL
		break
	}

	if err_ground := config.DB.Where("planting_info_id=?", planting_info.ID).First(&ground_info).Error; err_ground != nil {
		log.Print(color.RedString(err_ground.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	data := []map[string]interface{}{
		{
			"container":           planting_info.Container,
			"planting_article_id": container_info.ID,
			"picture":             url_container,
		},
		{
			"ground":              planting_info.Ground,
			"planting_article_id": ground_info.ID,
			"picture":             url_ground,
		},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get planting location",
		"data":    data,
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 8 : Get plant location]
func Add_my_plant(c echo.Context) error {
	plant_id, _ := strconv.Atoi(c.Param("plant_id"))
	var myplant model.MyPlant

	if err_bind := c.Bind(&myplant); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	myplant.UserID = user_id
	myplant.PlantID = uint(plant_id)
	myplant.IsStartPlanting = false

	if err_save := config.DB.Save(&myplant).Error; err_save != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to add user plant",
		"data": map[string]interface{}{
			"myplant_id": myplant.ID,
			"plant_id":   plant_id,
		},
	})

}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 16 : Start planting]
func Start_planting(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var myplant model.MyPlant

	// VALIDATION1
	var watering_check model.Watering
	if err_first := config.DB.Where("my_plant_id=? AND week=?", myplant_id, 1).First(&watering_check).Error; err_first == nil {
		log.Print(color.RedString(echo.ErrBadRequest.Error()), " is already start planting")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if err_first := config.DB.First(&myplant, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var myplant_binding model.MyPlant
	if err_bind := c.Bind(&myplant_binding); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// START SET1 - myplant table : longitude(current), latitude(current), is_start_planting(true), is_start_planting(current date)
	myplant.Longitude = myplant_binding.Longitude
	myplant.Latitude = myplant_binding.Latitude
	myplant.IsStartPlanting = true
	myplant.StartPlantingDate = time.Now()
	myplant.Status = "planting"

	if err_save1 := config.DB.Save(&myplant).Error; err_save1 != nil {
		log.Print(color.RedString(err_save1.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}
	// END SET1

	// START SET2 - watering table : all columns
	watering := model.Watering{
		MyPlantID: uint(myplant_id),
		Week:      1,
		Day1:      0,
		Day2:      0,
		Day3:      0,
		Day4:      0,
		Day5:      0,
		Day6:      0,
		Day7:      0,
	}

	watering.Week = 1
	if err_save2 := config.DB.Save(&watering).Error; err_save2 != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}
	// END SET2

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to start planting",
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 17 : Get my plant overview]
func Get_myplant_overview(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var myplant model.MyPlant

	if err_first := config.DB.First(&myplant, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	diff := time.Now().Sub(myplant.StartPlantingDate)

	day := int(diff.Hours()/24) + 1
	week := int(diff.Hours()/(24*7)) + 1
	if day > 7 {
		day = day % 7
		var wateringCheck model.Watering
		if err_first_check_watering := config.DB.Where("my_plant_id=? AND week=?", myplant_id, week).First(&wateringCheck).Error; err_first_check_watering != nil {
			wateringCheck.MyPlantID = uint(myplant_id)
			wateringCheck.Week = week
			wateringCheck.Day1 = 0
			wateringCheck.Day2 = 0
			wateringCheck.Day3 = 0
			wateringCheck.Day4 = 0
			wateringCheck.Day5 = 0
			wateringCheck.Day6 = 0
			wateringCheck.Day7 = 0

			if err_save2 := config.DB.Save(&wateringCheck).Error; err_save2 != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"status":  500,
					"message": "internal server error",
				})
			}
		}
	}

	// START GET WATERING ------------------------------------------------------------------------------
	// get watering period
	var wateringInfo model.WateringInfo
	if err_first2 := config.DB.Where("plant_id=?", myplant.PlantID).First(&wateringInfo).Error; err_first2 != nil {
		log.Print(color.RedString(err_first2.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	var watering model.Watering
	if err_first4 := config.DB.Where("my_plant_id=? AND week=?", myplant_id, week).First(&watering).Error; err_first4 != nil {
		log.Print(color.RedString(err_first4.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}
	watering_history := []int{watering.Day1, watering.Day2, watering.Day3, watering.Day4, watering.Day5, watering.Day6, watering.Day7}

	// get watering is_active
	is_active_watering := true
	if watering_history[day-1] >= 2 {
		is_active_watering = false
	}

	response_watering := map[string]interface{}{
		"week":      watering.Week,
		"day":       day,
		"period":    wateringInfo.Period,
		"is_active": is_active_watering,
		"history":   watering_history,
	}
	// END GET WATERING --------------------------------------------------------------------------------------

	// START GET FERTILIZING ---------------------------------------------------------------------------------
	// get fertilizing period
	var fertilizingInfo model.FertilizingInfo
	if err_first3 := config.DB.Where("plant_id=?", myplant.PlantID).First(&fertilizingInfo).Error; err_first3 != nil {
		log.Print(color.RedString(err_first3.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	day_fertilizing := int(diff.Hours()/24) + 1
	is_active_fertilizing := false
	if week == 1 && day == 1 {
		is_active_fertilizing = true
	} else if day_fertilizing%fertilizingInfo.Period == 0 {
		is_active_fertilizing = true
	}

	is_enabled_fertilizing := false
	if is_active_fertilizing {
		var fertilizing model.Fertilizing
		if err_first_fertilizing := config.DB.Where("my_plant_id=? AND week=?", myplant_id, week).First(&fertilizing).Error; err_first_fertilizing != nil {
			is_enabled_fertilizing = true
		}
	}

	response_fertilizing := map[string]interface{}{
		"is_active":  is_active_fertilizing,
		"is_enabled": is_enabled_fertilizing,
	}
	// END GET FERTILIZING ----------------------------------------------------------------------------------

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get my plant overview",
		"data": map[string]interface{}{
			"watering":    response_watering,
			"fertilizing": response_fertilizing,
		},
	})
}
