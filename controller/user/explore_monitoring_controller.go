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
	"github.com/pariz/gountries"
	"github.com/zsefvlol/timezonemapper"
)

// EXPLORE & MONITORING (Menu Home) - [Endpoint 1 : Get weather]
func Get_weather(c echo.Context) error {
	latitude := c.Param("latitude")
	longitude := c.Param("longitude")
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
	q.Set("lat", latitude)
	q.Set("lon", longitude)
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
	countryCode, ok := data["sys"].(map[string]interface{})["country"].(string)
	if !ok {
		fmt.Println("Invalid country code")
		return nil
	}

	country, err := gountries.New().FindCountryByAlpha(countryCode)
	if err != nil {
		fmt.Println("Error finding country:", err)
		return err
	}

	countryName := country.Name.Common
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
			"country":     countryName,
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
		infoWeather.UserID = user_id
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

	infoWeather2.UserID = user_id
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
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var picture model.Picture
	if err_first2 := config.DB.Where("weather_id=?", weatherArticle.ID).First(&picture).Error; err_first2 != nil {
		log.Print(color.RedString(err_first2.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
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
		return "NoWeather"
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

// EXPLORE & MONITORING (Menu Home) - [Endpoint 3.1 : Get notifictions]
func Get_notifications(c echo.Context) error {
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)
	latitude := c.Param("latitude")
	longitude := c.Param("longitude")

	// get user userCurrentDate
	var user model.User
	if err_first_user := config.DB.First(&user, user_id).Error; err_first_user != nil {
		log.Print(color.RedString(err_first_user.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "bad request",
		})
	}

	userCurrentDate := get_current_time_from_latlong(latitude, longitude).Format("2006-01-02")

	var notificationCheck model.Notification
	if err_first_check := config.DB.Where("user_id=? AND date=?", user_id, userCurrentDate).First(&notificationCheck).Error; err_first_check != nil {
		// delete yesterday notification
		var deleteNotification model.Notification
		if err_deleteNotifications := config.DB.Where("user_id=?", user_id).Delete(&deleteNotification).Error; err_deleteNotifications != nil {
			log.Print(color.RedString(err_deleteNotifications.Error()))
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status":  500,
				"message": "internal server error",
			})
		}

		// insert notifications
		err_generateNotifications := generate_notifications(user_id, userCurrentDate)
		if !err_generateNotifications {
			log.Print(color.RedString("error while generate notifications 1"))
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status":  500,
				"message": "internal server error",
			})
		}

		err_generateNotifications2 := generate_notifications(user_id, userCurrentDate)
		if !err_generateNotifications2 {
			log.Print(color.RedString("error while generate notifications 2"))
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status":  500,
				"message": "internal server error",
			})
		}
	}

	// retrieve all notification datas
	var notifications []model.Notification
	if err_find := config.DB.Order("created_at DESC").Where("user_id=? AND date=?", user_id, userCurrentDate).Find(&notifications).Error; err_find != nil {
		log.Print(color.RedString(err_find.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	responses := []map[string]interface{}{}
	for _, notification := range notifications {
		// get my_plants data
		var myplant model.MyPlant
		if err_first := config.DB.First(&myplant, notification.MyPlantID).Error; err_first != nil {
			log.Print(color.RedString(err_first.Error()))
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status":  500,
				"message": "internal server error",
			})
		}

		// generate description
		description := ""
		switch notification.Activity {
		case "watering":
			description = "Kamu belum menyiram"
		case "fertilizing":
			description = "Kamu belum memupuk"
		case "weekly progress":
			description = "Kamu belum menambahkan weekly progress"
		}

		response := map[string]interface{}{
			"id_notif":       notification.ID,
			"nama_tanamanku": myplant.Name,
			"description":    description,
			"myplant_id":     notification.MyPlantID,
			"plant_id":       myplant.PlantID,
			"location":       myplant.Location,
			"read_status":    notification.ReadStatus,
		}

		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get all notifications",
		"data":    responses,
	})
}

// HELPER FUNCTION
func generate_notifications(user_id uint, userCurrentDate string) bool {
	var myPlants []model.MyPlant
	if err_find := config.DB.Where("user_id=? AND status='planting'", user_id).Find(&myPlants).Error; err_find != nil {
		log.Print(color.RedString(err_find.Error()))
		return false
	}

	// Loop each myPlants data to store data
	for _, myPlant := range myPlants {
		currentTime := get_current_time_from_latlong(myPlant.Latitude, myPlant.Longitude)
		diff := currentTime.Sub(myPlant.StartPlantingDate)
		day := int(diff.Hours()/24) + 1
		week := int(diff.Hours()/(24*7)) + 1
		if day > 6 {
			day = day % 7
			if day == 0 {
				day = 7
			}
		}

		// Check watering
		// get watering periode
		var wateringInfo model.WateringInfo
		if err_first_watering := config.DB.Where("plant_id=?", myPlant.PlantID).First(&wateringInfo).Error; err_first_watering != nil {
			log.Print(color.RedString(err_first_watering.Error()))
			return false
		}

		// get watering history
		var watering model.Watering
		if err_first := config.DB.Where("my_plant_id=? AND week=?", myPlant.ID, week).First(&watering).Error; err_first != nil {
			log.Print(color.RedString(err_first.Error(), "please hit endpoint Get myplant overview first"))

			var wateringCheck model.Watering
			wateringCheck.MyPlantID = myPlant.ID
			wateringCheck.Week = week
			wateringCheck.Day1 = 0
			wateringCheck.Day2 = 0
			wateringCheck.Day3 = 0
			wateringCheck.Day4 = 0
			wateringCheck.Day5 = 0
			wateringCheck.Day6 = 0
			wateringCheck.Day7 = 0

			if err_save2 := config.DB.Save(&wateringCheck).Error; err_save2 != nil {
				return false
			}
			fmt.Println("aman", myPlant.ID)
			return true
		}

		watering_history := []int{watering.Day1, watering.Day2, watering.Day3, watering.Day4, watering.Day5, watering.Day6, watering.Day7}
		if watering_history[day-1] < wateringInfo.Period {
			var notification model.Notification
			notification.UserID = user_id
			notification.MyPlantID = myPlant.ID
			notification.Date = userCurrentDate
			notification.Activity = "watering"
			notification.ReadStatus = false

			if err_save := config.DB.Save(&notification).Error; err_save != nil {
				log.Print(color.RedString(err_save.Error()))
				return false
			}
		}

		// Check fetilizing
		var fertilizingInfo model.FertilizingInfo
		if err_first := config.DB.Where("plant_id=?", myPlant.PlantID).First(&fertilizingInfo).Error; err_first != nil {
			log.Print(color.RedString(err_first.Error()))
			return false
		}

		day_fertilizing := int(diff.Hours()/24) + 1
		if (week == 1 && day == 1) || (day_fertilizing%fertilizingInfo.Period == 0) {
			var fertilizing model.Fertilizing
			if err_first_fertilizing := config.DB.Where("my_plant_id=? AND week=?", myPlant.ID, week).First(&fertilizing).Error; err_first_fertilizing != nil {
				var notification model.Notification
				notification.UserID = user_id
				notification.MyPlantID = myPlant.ID
				notification.Date = userCurrentDate
				notification.Activity = "fertilizing"
				notification.ReadStatus = false

				if err_save := config.DB.Save(&notification).Error; err_save != nil {
					log.Print(color.RedString(err_save.Error()))
					return false
				}
			}
		}

		// Check weekly progress
		var weeklyProgress model.WeeklyProgress
		fmt.Println(day)
		if day == 7 {
			if err_first := config.DB.Where("my_plant_id=? AND week=? status='planting'", myPlant.ID, week).First(&weeklyProgress).Error; err_first != nil {
				var notification model.Notification
				notification.UserID = user_id
				notification.MyPlantID = myPlant.ID
				notification.Date = userCurrentDate
				notification.Activity = "weekly progress"
				notification.ReadStatus = false

				if err_save := config.DB.Save(&notification).Error; err_save != nil {
					log.Print(color.RedString(err_save.Error()))
					return false
				}
			}
		}
	}

	return true
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 3.2 : Read notifiction]
func Read_notification(c echo.Context) error {
	notificationId := c.Param("notification_id")

	var notification model.Notification
	if err_first := config.DB.First(&notification, notificationId).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	notification.ReadStatus = true
	if err_save := config.DB.Save(&notification).Error; err_save != nil {
		log.Print(color.RedString(err_save.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to read notification",
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 5 : Get available plants]
func Get_available_plants(c echo.Context) error {
	var plants []model.Plant

	if err_find := config.DB.Find(&plants).Error; err_find != nil {
		log.Print(color.RedString(err_find.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "bad request",
		})
	}

	responses := []map[string]interface{}{}
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
			"status":  400,
			"message": "bad request",
		})
	}

	responses := []map[string]interface{}{}
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
		return c.JSON(http.StatusNotFound, map[string]interface{}{
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

	// Get fertilizing article
	var fertilizngInfo model.FertilizingInfo
	if err_first := config.DB.Where("plant_id=?", plant_id).First(&fertilizngInfo).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error(), "fertilizing article not found"))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// Get watering article
	var wateringInfo model.WateringInfo
	if err_first := config.DB.Where("plant_id=?", plant_id).First(&wateringInfo).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error(), "watering article not found"))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// Get temperature article
	var temperatureInfo model.TemperatureInfo
	if err_first := config.DB.Where("plant_id=?", plant_id).First(&temperatureInfo).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error(), "temperature article not found"))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get plant detail",
		"data": map[string]interface{}{
			"plant_id":               plant.ID,
			"picture":                url,
			"name":                   plant.Name,
			"latin":                  plant.Latin,
			"description":            plant.Description,
			"fertilizing_article_id": fertilizngInfo.ID,
			"watering_article_id":    wateringInfo.ID,
			"temperature_article_id": temperatureInfo.ID,
		},
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 8 : Get plant location]
func Get_plant_location(c echo.Context) error {
	plant_id, _ := strconv.Atoi(c.Param("plant_id"))

	var planting_info model.PlantingInfo
	if err_first := config.DB.Where("plant_id=?", plant_id).First(&planting_info).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	responses := []map[string]interface{}{}

	if planting_info.Container {
		var container_info model.ContainerInfo

		if err_container := config.DB.Where("planting_info_id=?", planting_info.ID).First(&container_info).Error; err_container != nil {
			log.Print(color.RedString(err_container.Error()))
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  404,
				"message": "not found",
			})
		}

		config.DB.Model(&container_info).Association("Pictures").Find(&container_info.Pictures)

		var url_container string
		for _, picture := range container_info.Pictures {
			url_container = picture.URL
			break
		}

		responses = append(responses, map[string]interface{}{
			"location_plant": 1,
			"picture":        url_container,
		})
	}

	if planting_info.Ground {
		var ground_info model.GroundInfo

		if err_ground := config.DB.Where("planting_info_id=?", planting_info.ID).First(&ground_info).Error; err_ground != nil {
			log.Print(color.RedString(err_ground.Error()))
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  404,
				"message": "not found",
			})
		}

		config.DB.Model(&ground_info).Association("Pictures").Find(&ground_info.Pictures)

		var url_ground string
		for _, picture := range ground_info.Pictures {
			url_ground = picture.URL
			break
		}

		responses = append(responses, map[string]interface{}{
			"location_plant": 2,
			"picture":        url_ground,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get planting location",
		"data":    responses,
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 9 : Get planting article]
func GetPlantingArticle(c echo.Context) error {
	location := c.Param("location")
	plantID := c.Param("plant_id")

	// make location to lowercase
	location = strings.ToLower(location)

	// check if plant id is valid
	plantIDUint, err := StringToUintPointer(plantID)
	if err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	var plant model.Plant
	if err_first := config.DB.First(&plant, plantIDUint).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var plantingInfo model.PlantingInfo
	if err_first := config.DB.Where("plant_id=?", plantIDUint).First(&plantingInfo).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var plantingArticle model.ContainerInfo
	if location == "container" {
		if err_first := config.DB.Where("planting_info_id=?", plantingInfo.ID).First(&plantingArticle).Error; err_first != nil {
			log.Print(color.RedString(err_first.Error()))
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  404,
				"message": "not found",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  200,
			"message": "success to get planting article",
			"data": map[string]interface{}{
				"plant_id":   plant.ID,
				"location":   location,
				"link_video": plantingArticle.Video,
				"description": map[string]interface{}{
					"material":    plantingArticle.Materials,
					"instruction": plantingArticle.Instructions,
				},
			},
		})
	} else if location == "ground" {
		if err_first := config.DB.Where("planting_info_id=?", plantingInfo.ID).First(&plantingArticle).Error; err_first != nil {
			log.Print(color.RedString(err_first.Error()))
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  404,
				"message": "not found",
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  200,
			"message": "success to get planting article",
			"data": map[string]interface{}{
				"plant_id":   plant.ID,
				"location":   location,
				"link_video": plantingArticle.Video,
				"description": map[string]interface{}{
					"material":    plantingArticle.Materials,
					"instruction": plantingArticle.Instructions,
				},
			},
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"status":  400,
		"message": "bad request",
	})
}

func GetFertilizingArticle(c echo.Context) error {
	plantID := c.Param("plant_id")

	// check if plant id is valid
	plantIDUint, err := StringToUintPointer(plantID)
	if err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	var plant model.Plant
	if err_first := config.DB.First(&plant, plantIDUint).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var fertilizingInfo model.FertilizingInfo
	if err_first := config.DB.Where("plant_id=?", plantIDUint).First(&fertilizingInfo).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// get the picture from picture by fertilizing info id
	var picture model.Picture
	if err_first := config.DB.Where("fertilizing_info_id=?", fertilizingInfo.ID).First(&picture).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get fertilizing article",
		"data": map[string]interface{}{
			"plant_id":               plant.ID,
			"name":                   plant.Name,
			"picture":                picture.URL,
			"description":            fertilizingInfo.Description,
			"products_recomendation": GetRelatedProducts("pupuk"),
		},
	})
}

func GetWateringArticle(c echo.Context) error {
	plantID := c.Param("plant_id")

	// check if plant id is valid
	plantIDUint, err := StringToUintPointer(plantID)
	if err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	var plant model.Plant
	if err_first := config.DB.First(&plant, plantIDUint).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var wateringInfo model.WateringInfo
	if err_first := config.DB.Where("plant_id=?", plantIDUint).First(&wateringInfo).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// get the picture from picture by watering info id
	var picture model.Picture
	if err_first := config.DB.Where("watering_info_id=?", wateringInfo.ID).First(&picture).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get watering article",
		"data": map[string]interface{}{
			"plant_id":    plant.ID,
			"name":        plant.Name,
			"picture":     picture.URL,
			"description": wateringInfo.Description,
		},
	})
}

func GetTemperatureArticle(c echo.Context) error {
	plantID := c.Param("plant_id")

	// check if plant id is valid
	plantIDUint, err := StringToUintPointer(plantID)
	if err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	var plant model.Plant
	if err_first := config.DB.First(&plant, plantIDUint).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var temperatureInfo model.TemperatureInfo
	if err_first := config.DB.Where("plant_id=?", plantIDUint).First(&temperatureInfo).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// get the picture from picture by temperature info id
	var picture model.Picture
	if err_first := config.DB.Where("temperature_info_id=?", temperatureInfo.ID).First(&picture).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get temperature article",
		"data": map[string]interface{}{
			"plant_id":    plant.ID,
			"name":        plant.Name,
			"picture":     picture.URL,
			"description": temperatureInfo.Description,
		},
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 13 : Add my plant]
func Add_my_plant(c echo.Context) error {
	plant_id, _ := strconv.Atoi(c.Param("plant_id"))
	var myplant model.MyPlant

	// Validation1 : is plant_id exist
	var plant_check model.Plant
	if err_val1 := config.DB.First(&plant_check, plant_id).Error; err_val1 != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

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

	// Set the current time as the start_planting_date
	myplant.StartPlantingDate = time.Now().UTC()

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

// EXPLORE & MONITORING (Menu Home) - [Endpoint 14 : Get my plant]
func Get_myplant_name(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var myplant model.MyPlant

	if err_first := config.DB.First(&myplant, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var plant model.Plant
	if err_first2 := config.DB.First(&plant, myplant.PlantID).Error; err_first2 != nil {
		log.Print(color.RedString(err_first2.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get my plant name",
		"data": map[string]interface{}{
			"myplant_id":        myplant.ID,
			"name":              myplant.Name,
			"latin":             plant.Latin,
			"is_start_planting": myplant.IsStartPlanting,
		},
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 15 : Update my plant name]

func Update_myplant_name(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var myplant model.MyPlant

	if err_first := config.DB.First(&myplant, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
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
	myplant.Name = myplant_binding.Name

	if err_save := config.DB.Save(&myplant).Error; err_save != nil {
		log.Print(color.RedString(err_save.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to update my plant name",
		"data": map[string]interface{}{
			"myplant_id":           myplant.ID,
			"updated_myplant_name": myplant.Name,
		},
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 16 : Start planting]
func Start_planting(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var myplant model.MyPlant

	// Validation1 : is already start planting
	var watering_check model.Watering
	if err_first := config.DB.Where("my_plant_id=? AND week=?", myplant_id, 1).First(&watering_check).Error; err_first == nil {
		log.Print(color.RedString("is already start planting"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Validation2 : is match beetwen myplant_id and user_id
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	var myplant_check model.MyPlant
	if err_val2 := config.DB.Where("user_id=?", int(user_id)).First(&myplant_check, myplant_id).Error; err_val2 != nil {
		log.Print(color.RedString(err_val2.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
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

	// Get current timestamp according to longitude and latitude
	currentTime := get_current_time_from_latlong(myplant_binding.Latitude, myplant_binding.Longitude)
	newCurrentTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 1, 0, currentTime.Location())

	// START SET1 - myplant table : longitude(current), latitude(current), is_start_planting(true), is_start_planting(current date)
	myplant.Longitude = myplant_binding.Longitude
	myplant.Latitude = myplant_binding.Latitude
	myplant.IsStartPlanting = true
	myplant.StartPlantingDate = newCurrentTime
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

	// START APPEND NOTIFICATION
	userCurrentDate := get_current_time_from_latlong(myplant_binding.Latitude, myplant_binding.Longitude).Format("2006-01-02")
	activities := []string{"watering", "fertilizing"}

	for _, activity := range activities {
		var notification model.Notification
		notification.UserID = user_id
		notification.MyPlantID = uint(myplant_id)
		notification.Date = userCurrentDate
		notification.Activity = activity
		notification.ReadStatus = false

		if err_save := config.DB.Save(&notification).Error; err_save != nil {
			log.Print(color.RedString(err_save.Error()))
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status":  500,
				"message": "internal server error",
			})
		}
	}
	// END APPEND NOTIFICATION

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to start planting",
	})
}

// HELPER FUNCTION
func get_current_time_from_latlong(lat, long string) time.Time {
	longitude, _ := strconv.ParseFloat(long, 32)
	latitude, _ := strconv.ParseFloat(lat, 32)

	// Get timezone string from lat/long
	timezone := timezonemapper.LatLngToTimezoneString(latitude, longitude)

	// Should print "Timezone: Asia/Shanghai"
	fmt.Printf("Timezone: %s\n", timezone)

	// Load location from timezone
	loc, _ := time.LoadLocation(timezone)

	// Get current time in the specified location
	currentTime := time.Now().In(loc)

	return currentTime
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 17 : Get my plant overview]
func Get_myplant_overview(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var myplant model.MyPlant

	// Validation1 : make sure other user cannot access data from other user
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	if err_first := config.DB.Where("user_id=?", user_id).First(&myplant, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// Get current timestamp according to longitude and latitude
	currentTime := get_current_time_from_latlong(myplant.Latitude, myplant.Longitude)
	// Print log for developergolang
	fmt.Println(map[string]interface{}{
		"startPlantingDate": myplant.StartPlantingDate,
		"currentTime":       currentTime,
		"difference":        currentTime.Sub(myplant.StartPlantingDate),
	})

	diff := currentTime.Sub(myplant.StartPlantingDate)
	day := int(diff.Hours()/24) + 1
	week := int(diff.Hours()/(24*7)) + 1

	if day > 6 {
		day = day % 7
		if day == 0 {
			day = 7
		}
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
	if watering_history[day-1] >= wateringInfo.Period {
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
		"period":     fertilizingInfo.Period,
	}
	// END GET FERTILIZING ----------------------------------------------------------------------------------

	// START GET WEEKLY PROGRESS ---------------------------------------------------------------------------------
	isActiveWeeklyProgress := false
	isEnabledWeeklyProgress := false
	var weeklyProgress model.WeeklyProgress

	if day == 7 {
		isActiveWeeklyProgress = true
		if err_first4 := config.DB.Where("my_plant_id=? AND week=?", myplant_id, week).First(&weeklyProgress).Error; err_first4 != nil {
			isEnabledWeeklyProgress = true
		}
	}

	response_weekly_progress := map[string]interface{}{
		"is_active":  isActiveWeeklyProgress,
		"from":       myplant.StartPlantingDate.Format("02-01-2006"),
		"to":         myplant.StartPlantingDate.Add(168 * time.Hour).Format("02-01-2006"),
		"is_enabled": isEnabledWeeklyProgress,
	}
	// END GET WEEKLY PROGRESS ---------------------------------------------------------------------------------

	// START GET BUTTON HARVEST & DEAD ---------------------------------------------------------------------------------
	isActiveButtonHarvest := false
	isActiveButtonDead := false
	day_button := int(diff.Hours()/24) + 1
	if day_button > 6 {
		isActiveButtonHarvest = true
		isActiveButtonDead = true
		if isActiveWeeklyProgress && !isEnabledWeeklyProgress {
			isActiveButtonHarvest = false
			isActiveButtonDead = false
		}
	}
	// END GET BUTTON HARVEST & DEAD ---------------------------------------------------------------------------------

	// START GET TEMPERATURE ALERT ---------------------------------------------------------------------------------
	currentTemperature, _ := strconv.ParseFloat(get_current_temperature(myplant.Longitude, myplant.Latitude), 32)
	var temperatureInfo model.TemperatureInfo
	if err_first_temperature := config.DB.Where("plant_id=?", myplant.PlantID).First(&temperatureInfo).Error; err_first_temperature != nil {
		log.Print(color.RedString(err_first_temperature.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	isActiveTemperatureAlert := false
	if float64(temperatureInfo.Min) > currentTemperature || float64(temperatureInfo.Max) < currentTemperature {
		isActiveTemperatureAlert = true
	}

	roundedNumber := float64(int(currentTemperature*100)) / 100
	responseTemperatureAlert := map[string]interface{}{
		"is_active": isActiveTemperatureAlert,
		"name":      myplant.Name,
		"current":   roundedNumber,
		"min":       temperatureInfo.Min,
		"max":       temperatureInfo.Max,
	}

	// END GET TEMPERATURE ALERT ---------------------------------------------------------------------------------

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get my plant overview",
		"data": map[string]interface{}{
			"watering":          response_watering,
			"fertilizing":       response_fertilizing,
			"weekly_progress":   response_weekly_progress,
			"temperature_alert": responseTemperatureAlert,
			"button_harvest":    isActiveButtonHarvest,
			"button_dead":       isActiveButtonDead,
		},
	})
}

// HELPER FUNCTION
func get_current_temperature(longitude, latitude string) string {
	apikey := "869a5f0aa562d21ec64ff37c7c6c157f"
	baseURL := "https://api.openweathermap.org/data/2.5/weather"

	u, err_parseURL := url.Parse(baseURL)
	if err_parseURL != nil {
		log.Print(color.RedString(err_parseURL.Error()))
		return "0"
	}

	// Add query parameters to the URL
	q := u.Query()
	q.Set("lat", latitude)
	q.Set("lon", longitude)
	q.Set("appid", apikey)
	q.Set("units", "metric")
	u.RawQuery = q.Encode()

	// Make the GET request
	response, err_consume := http.Get(u.String())
	if err_consume != nil {
		log.Print(color.RedString(err_consume.Error()))
		return "0"
	}
	defer response.Body.Close()

	// Read the response body
	body, err_read := ioutil.ReadAll(response.Body)
	if err_read != nil {
		log.Print(color.RedString(err_read.Error()))
		return "0"
	}

	// Create a map to store the JSON data
	var data map[string]interface{}

	// Unmarshal the response body into the map
	err_unmarshal := json.Unmarshal(body, &data)
	if err_unmarshal != nil {
		log.Print(color.RedString(err_unmarshal.Error()))
		return "0"
	}
	temperature := fmt.Sprintf("%v", data["main"].(map[string]interface{})["temp"])

	return temperature
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 18 : Add watering]
func Add_watering(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var myplant model.MyPlant

	// Validation1 : make sure other user cannot modify other user data
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	if err_first := config.DB.Where("user_id=?", user_id).First(&myplant, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// Get current timestamp according to longitude and latitude
	currentTime := get_current_time_from_latlong(myplant.Latitude, myplant.Longitude)
	diff := currentTime.Sub(myplant.StartPlantingDate)

	day := int(diff.Hours()/24) + 1
	if day > 6 {
		day = day % 7
		if day == 0 {
			day = 7
		}
	}
	week := int(diff.Hours()/(24*7)) + 1

	var watering model.Watering
	if err_first2 := config.DB.Where("my_plant_id=? AND week=?", myplant_id, week).First(&watering).Error; err_first2 != nil {
		log.Print(color.RedString(err_first2.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	watering_history := []int{watering.Day1, watering.Day2, watering.Day3, watering.Day4, watering.Day5, watering.Day6, watering.Day7}

	var wateringInfo model.WateringInfo
	if err_first2 := config.DB.Where("plant_id=?", myplant.PlantID).First(&wateringInfo).Error; err_first2 != nil {
		log.Print(color.RedString(err_first2.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	if watering_history[day-1] < wateringInfo.Period {

		switch day {
		case 1:
			watering.Day1 = watering.Day1 + 1
		case 2:
			watering.Day2 = watering.Day2 + 1
		case 3:
			watering.Day3 = watering.Day3 + 1
		case 4:
			watering.Day4 = watering.Day4 + 1
		case 5:
			watering.Day5 = watering.Day5 + 1
		case 6:
			watering.Day6 = watering.Day6 + 1
		case 7:
			watering.Day7 = watering.Day7 + 1
		}

		if err_update := config.DB.Save(&watering).Error; err_update != nil {
			log.Print(color.RedString(err_update.Error()))
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status":  500,
				"message": "internal server error",
			})
		}
	} else {
		log.Print(color.RedString("already do the watering according to the period"))
		return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
			"status":  429,
			"message": "too many request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to add my plant watering",
		"data": map[string]interface{}{
			"week": week,
			"day":  day,
		},
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 19 : Add fertilizing]
func Add_fertilizing(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var fertilizing model.Fertilizing
	var myplant model.MyPlant

	// Validation1 : make sure other user cannot modify other user data
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	if err_first := config.DB.Where("user_id=?", user_id).First(&myplant, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// get fertilizing period
	var fertilizingInfo model.FertilizingInfo
	if err_first3 := config.DB.Where("plant_id=?", myplant.PlantID).First(&fertilizingInfo).Error; err_first3 != nil {
		log.Print(color.RedString(err_first3.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get current timestamp according to longitude and latitude
	currentTime := get_current_time_from_latlong(myplant.Latitude, myplant.Longitude)
	diff := currentTime.Sub(myplant.StartPlantingDate)
	day := int(diff.Hours()/24) + 1

	week := int(diff.Hours()/(24*7)) + 1
	is_active_fertilizing := false
	if week == 1 && day == 1 {
		is_active_fertilizing = true
	} else if day%fertilizingInfo.Period == 0 {
		is_active_fertilizing = true
	}

	if is_active_fertilizing {
		if err_first2 := config.DB.Where("my_plant_id=? AND week=?", myplant_id, week).First(&fertilizing).Error; err_first2 == nil {
			log.Print(color.RedString("already fertilizing in this week"))
			return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
				"status":  429,
				"message": "too many request",
			})
		} else {
			fertilizing.MyPlantID = uint(myplant_id)
			fertilizing.Week = week
			fertilizing.Status = true
			if err_insert := config.DB.Save(&fertilizing).Error; err_insert != nil {
				log.Print(color.RedString(err_insert.Error()))
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"status":  500,
					"message": "internal server error",
				})
			}
		}
	} else {
		log.Print(color.RedString("today is not fertilizing period"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to add my plant fertilizing",
		"data": map[string]interface{}{
			"week": week,
			"day":  day,
		},
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 20 : Add weekly progress]
func Add_weekly_progress(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var myplant model.MyPlant

	// Validation1 : make sure other user cannot modify other user data
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	if err_first := config.DB.Where("user_id=?", user_id).First(&myplant, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// Get current timestamp according to longitude and latitude
	currentTime := get_current_time_from_latlong(myplant.Latitude, myplant.Longitude)
	diff := currentTime.Sub(myplant.StartPlantingDate)

	day := int(diff.Hours()/24) + 1
	if day > 6 {
		day = day % 7
		if day == 0 {
			day = 7
		}
	}
	week := int(diff.Hours()/(24*7)) + 1

	if day == 7 {
		var weeklyProgress model.WeeklyProgress
		if err_first2 := config.DB.Where("my_plant_id=? AND week=?", myplant_id, week).First(&weeklyProgress).Error; err_first2 == nil {
			log.Print(color.RedString("already report weekly progress for this week"))
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  400,
				"message": "bad request",
			})

		} else {
			if err_bind := c.Bind(&weeklyProgress); err_bind != nil {
				log.Print(color.RedString(err_bind.Error()))
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"status":  400,
					"message": "bad request",
				})
			}

			weeklyProgress.MyPlantID = uint(myplant_id)
			weeklyProgress.Week = week
			weeklyProgress.From = myplant.StartPlantingDate
			weeklyProgress.To = myplant.StartPlantingDate.Add(168 * time.Hour)
			weeklyProgress.Status = "planting"

			if err_insert := config.DB.Save(&weeklyProgress).Error; err_insert != nil {
				log.Print(color.RedString(err_insert.Error()))
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"status":  500,
					"message": "internal server error",
				})
			}

		}
	} else {
		log.Print(color.RedString("already report weekly progress for this week"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to add my plant weekly progress",
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 21 : Get all myplant weekly progress]
func Get_all_myplant_weekly_progress(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var weeklyProgresses []model.WeeklyProgress

	// Validation1 : make sure other user cannot modify other user data
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	var myplant_val1 model.MyPlant
	if err_first := config.DB.Where("user_id=?", user_id).First(&myplant_val1, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	if err_find := config.DB.Where("my_plant_id=?", myplant_id).Find(&weeklyProgresses).Error; err_find != nil {
		log.Print(color.RedString(err_find.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	responses := []map[string]interface{}{}
	for _, weeklyProgress := range weeklyProgresses {
		config.DB.Model(&weeklyProgress).Association("Pictures").Find(&weeklyProgress.Pictures)

		var url_ground string
		for _, picture := range weeklyProgress.Pictures {
			url_ground = picture.URL
			break
		}

		response := map[string]interface{}{
			"weekly_progress_id": weeklyProgress.ID,
			"week":               weeklyProgress.Week,
			"picture":            url_ground,
			"from":               weeklyProgress.From.Format("02-01-2006"),
			"to":                 weeklyProgress.To.Format("02-01-2006"),
			"status":             weeklyProgress.Status,
		}

		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get all my plant weekly progress",
		"data":    responses,
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 22 : Get myplant weekly progress by id]
func Get_my_plant_weekly_progress_by_id(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	weekly_progress_id, _ := strconv.Atoi(c.Param("weekly_progress_id"))
	var myPlant model.MyPlant
	var weeklyProgress model.WeeklyProgress

	// Validation1 : make sure other user cannot modify other user data
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	// Get MyPlant
	if err_first := config.DB.Where("user_id=?", user_id).First(&myPlant, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// Get WeeklyProgress
	if err_first2 := config.DB.First(&weeklyProgress, weekly_progress_id).Error; err_first2 != nil {
		log.Print(color.RedString(err_first2.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	var urls []string
	config.DB.Model(&weeklyProgress).Association("Pictures").Find(&weeklyProgress.Pictures)
	for _, picture := range weeklyProgress.Pictures {
		urls = append(urls, picture.URL)
	}

	// Get Watering
	var watering model.Watering
	var response_watering map[string]interface{}
	if err_first3 := config.DB.Where("my_plant_id=? AND week=?", myplant_id, weeklyProgress.Week).First(&watering).Error; err_first3 != nil {
		log.Print(color.RedString(err_first3.Error(), "this is not planting weekly progress"))
		response_watering = nil

	} else {
		wateringHistory := []int{watering.Day1, watering.Day2, watering.Day3, watering.Day4, watering.Day5, watering.Day6, watering.Day7}
		response_watering = map[string]interface{}{
			"watering_id": watering.ID,
			"week":        watering.Week,
			"history":     wateringHistory,
		}
	}

	// Get fertilizing
	var fertilizing model.Fertilizing
	var response_fertilizing map[string]interface{}
	if err_first4 := config.DB.Where("my_plant_id=? AND week=?", myplant_id, weeklyProgress.Week).First(&fertilizing).Error; err_first4 != nil {
		log.Print(color.RedString(err_first4.Error(), "this is not planting weekly progress"))
		response_fertilizing = nil

	} else {
		response_fertilizing = map[string]interface{}{
			"fertilizing_id": fertilizing.ID,
			"week":           watering.Week,
			"history":        1,
		}
	}

	response := map[string]interface{}{
		"status": myPlant.Status,
		"progress": map[string]interface{}{
			"weekly":      weeklyProgress.ID,
			"week":        weeklyProgress.Week,
			"pictures":    urls,
			"from":        weeklyProgress.From.Format("02-01-2006"),
			"to":          weeklyProgress.To.Format("02-01-2006"),
			"condition":   weeklyProgress.Condition,
			"description": weeklyProgress.Description,
		},
		"watering":    response_watering,
		"fertilizing": response_fertilizing,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get my plant weekly progress by weekly_progress_id",
		"data":    response,
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 23 : Update weekly progress]
func Update_weekly_progress(c echo.Context) error {
	weekly_progress_id, _ := strconv.Atoi(c.Param("weekly_progress_id"))
	var weeklyProgress model.WeeklyProgress

	if err_first := config.DB.First(&weeklyProgress, weekly_progress_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// Validation1 : make sure other user cannot modify other user data
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	var myplant_val1 model.MyPlant
	if err_first := config.DB.Where("user_id=?", user_id).First(&myplant_val1, weeklyProgress.MyPlantID).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var weeklyProgress_bind model.WeeklyProgress
	if err_bind := c.Bind(&weeklyProgress_bind); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	config.DB.Model(&weeklyProgress).Association("Pictures").Clear()

	weeklyProgress.Condition = weeklyProgress_bind.Condition
	weeklyProgress.Description = weeklyProgress_bind.Description
	weeklyProgress.Pictures = weeklyProgress_bind.Pictures

	if err_update := config.DB.Save(&weeklyProgress).Error; err_update != nil {
		log.Print(color.RedString(err_update.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to update my plant weekly progress",
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 24 : Add dead plant progress]
func Add_dead_plant_progress(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var myplant model.MyPlant

	// Validation1 : make sure other user cannot modify other user data
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	if err_first := config.DB.Where("user_id=?", user_id).First(&myplant, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var weeklyProgress_bind model.WeeklyProgress
	if err_bind := c.Bind(&weeklyProgress_bind); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Get current timestamp according to longitude and latitude
	currentTime := get_current_time_from_latlong(myplant.Latitude, myplant.Longitude)
	diff := currentTime.Sub(myplant.StartPlantingDate)
	day := int(diff.Hours()/24) + 1
	week := int(diff.Hours()/(24*7)) + 1

	var weeklyProgress model.WeeklyProgress

	if err_first2 := config.DB.Where("my_plant_id=? AND (status=? OR status=?)", myplant_id, "dead", "harvest").First(&weeklyProgress).Error; err_first2 == nil {
		log.Print(color.RedString("already add dead plant progress"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Get current timestamp according to longitude and latitude
	currentTime2 := get_current_time_from_latlong(myplant.Latitude, myplant.Longitude)
	weeklyProgress.MyPlantID = uint(myplant_id)
	weeklyProgress.Week = week
	weeklyProgress.From = currentTime2
	weeklyProgress.Status = "dead"
	weeklyProgress.To = currentTime2

	weeklyProgress.Condition = weeklyProgress_bind.Condition
	weeklyProgress.Description = weeklyProgress_bind.Description
	weeklyProgress.Pictures = weeklyProgress_bind.Pictures

	if day > 6 {
		if err_insert := config.DB.Save(&weeklyProgress).Error; err_insert != nil {
			log.Print(color.RedString(err_insert.Error()))
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status":  500,
				"message": "internal server error",
			})
		}
	} else {
		log.Print(color.RedString("should not add to the progress of dead plants less than 1 week"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	myplant.Status = "dead"
	if err_update := config.DB.Save(&myplant).Error; err_update != nil {
		log.Print(color.RedString(err_update.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Delete notification that used deleted plant_id
	var notification model.Notification
	config.DB.Where("my_plant_id=?", myplant_id).Delete(&notification)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to add dead plant progress",
	})
}

// EXPLORE & MONITORING (Menu Home) - [Endpoint 25 : Add harvest plant progress]
func Add_harvest_plant_progress(c echo.Context) error {
	myplant_id, _ := strconv.Atoi(c.Param("myplant_id"))
	var myplant model.MyPlant

	// Validation1 : make sure other user cannot modify other user data
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	if err_first := config.DB.Where("user_id=?", user_id).First(&myplant, myplant_id).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	var weeklyProgress_bind model.WeeklyProgress
	if err_bind := c.Bind(&weeklyProgress_bind); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Get current timestamp according to longitude and latitude
	currentTime := get_current_time_from_latlong(myplant.Latitude, myplant.Longitude)
	diff := currentTime.Sub(myplant.StartPlantingDate)
	day := int(diff.Hours()/24) + 1
	week := int(diff.Hours()/(24*7)) + 1

	var weeklyProgress model.WeeklyProgress

	if day > 6 {
		if err_first2 := config.DB.Where("my_plant_id=? AND (status=? OR status=?)", myplant_id, "dead", "harvest").First(&weeklyProgress).Error; err_first2 == nil {
			log.Print(color.RedString("already add harvest plant progress"))
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  400,
				"message": "bad request",
			})
		}
	} else {
		log.Print(color.RedString("should not add to the progress of harvest plants less than 1 week"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Get current timestamp according to longitude and latitude
	currentTime2 := get_current_time_from_latlong(myplant.Latitude, myplant.Longitude)
	weeklyProgress.MyPlantID = uint(myplant_id)
	weeklyProgress.Week = week
	weeklyProgress.From = currentTime2
	weeklyProgress.Status = "harvest"
	weeklyProgress.To = currentTime2

	weeklyProgress.Condition = weeklyProgress_bind.Condition
	weeklyProgress.Description = weeklyProgress_bind.Description
	weeklyProgress.Pictures = weeklyProgress_bind.Pictures

	if err_insert := config.DB.Save(&weeklyProgress).Error; err_insert != nil {
		log.Print(color.RedString(err_insert.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	myplant.Status = "harvest"
	if err_update := config.DB.Save(&myplant).Error; err_update != nil {
		log.Print(color.RedString(err_update.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Delete notification that used deleted plant_id
	var notification model.Notification
	config.DB.Where("my_plant_id=?", myplant_id).Delete(&notification)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to add harvest plant progress",
	})
}
