package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"net/http"
	"net/url"

	// "github.com/agriplant/config"
	// "github.com/agriplant/model"
	// "github.com/agriplant/utils"
	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/agriplant/utils"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

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

func save_weather_info(location, temperature, label string, user_id uint) bool {
	var infoWeather model.InfoWeather
	err_select := config.DB.Where("user_id=?", user_id).First(&infoWeather).Error
	fmt.Println(infoWeather)
	if err_select == nil {
		fmt.Println("update")
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
	fmt.Println("insert")
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
