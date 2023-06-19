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

// SETTING - [Endpoint 1 : Get Profile]
func GetProfile(c echo.Context) error {
	user := model.User{}

	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	if err := config.DB.First(&user, user_id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to retrieve users profiles data",
		"data": map[string]interface{}{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"picture": user.URL,
		},
	})
}

// SETTING - [Endpoint 2 : Get Username]
func GetUsername(c echo.Context) error {
	user := model.User{}

	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	if err := config.DB.First(&user, user_id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to retrieve users profiles data",
		"data":    user.Name,
	})
}

// SETTING - [Endpoint 3 : Update Username]
func UpdateUsername(c echo.Context) error {
	var Request struct {
		Name string `json:"name" validate:"required"`
	}

	if err := c.Bind(&Request); err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if Request.Name == "" {
		log.Print(color.RedString("request body can't empty"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	user := model.User{}

	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	if err := config.DB.First(&user, user_id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	user.Name = Request.Name

	if err := config.DB.Model(&user).Updates(model.User{Name: user.Name}).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "successfully update user name",
		"data": map[string]interface{}{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"picture": user.URL,
		},
	})
}

// SETTING - [Endpoint 4 : Update User Password]
func UpdateUserPassword(c echo.Context) error {
	var Request struct {
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&Request); err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if Request.Password == "" {
		log.Print(color.RedString("request body can't empty"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if len(Request.Password) > 20 || len(Request.Password) < 8 {
		log.Print(color.RedString("password min 8 and max 20 character"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	user := model.User{}

	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	if err := config.DB.First(&user, user_id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	user.Password = Request.Password

	// hashing password
	user.BeforeCreateUser(config.DB)
	if err := config.DB.Model(&user).Updates(model.User{Password: user.Password}).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "successfully update user password",
		"data": map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"picture":  user.URL,
			"password": user.Password,
		},
	})
}

// SETTING - [Endpoint 5 : Get my plants stats]
func GetMyPlantsStats(c echo.Context) error {
	type Response struct {
		MyPlantID uint            `json:"myplant_id"`
		Pictures  []model.Picture `json:"pictures"`
		Name      string          `json:"name"`
		Latin     string          `json:"latin"`
		Status    string          `json:"status"`
	}

	status := c.QueryParam("status")
	MyPlants := []model.MyPlant{}
	Responses := []Response{}

	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	// Get all my plants
	if err := config.DB.Where("status = ? AND user_id = ?", status, user_id).Find(&MyPlants).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	for i := 0; i < len(MyPlants); i++ {
		var temp Response

		var plant model.Plant
		if err := config.DB.First(&plant, MyPlants[i].PlantID).Error; err != nil {
			log.Print(color.RedString(err.Error()))
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  400,
				"message": "bad request",
			})
		}
		temp.MyPlantID = MyPlants[i].ID

		picture := []model.Picture{}
		if err := config.DB.Where("plant_id = ?", MyPlants[i].PlantID).Find(&picture).Error; err != nil {
			log.Print(color.RedString(err.Error()))
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  400,
				"message": "bad request",
			})
		}
		temp.Pictures = picture

		temp.Name = MyPlants[i].Name
		temp.Latin = plant.Latin
		temp.Status = MyPlants[i].Status
		Responses = append(Responses, temp)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get plants stats",
		"data":    Responses,
	})
}

// SETTING - [Endpoint 6 : Send complaint email]
func SendComplaintEmail(c echo.Context) error {
	complaint := model.Complaints{}

	if err := c.Bind(&complaint); err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	complaint.UserID = user_id

	if err := config.DB.Save(&complaint).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to send complaint email",
	})
}

// SETTING - [Endpoint 7 : Send suggestion]
func SendSuggestion(c echo.Context) error {
	suggestion := model.Suggestions{}

	if err := c.Bind(&suggestion); err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	suggestion.UserID = user_id

	if err := config.DB.Save(&suggestion).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to send suggestion",
	})
}

// SETTING - [Endpoint 8 : Update profile picture]
func UpdateProfilePicture(c echo.Context) error {
	var Request struct {
		Picture string `json:"picture"`
	}

	if err := c.Bind(&Request); err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	user := model.User{}

	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	if err := config.DB.First(&user, user_id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	user.URL = Request.Picture

	if err := config.DB.Model(&user).Updates(model.User{URL: user.URL}).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to update profile picture",
	})
}
