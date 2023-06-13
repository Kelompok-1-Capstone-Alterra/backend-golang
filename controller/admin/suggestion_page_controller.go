package admin

import (
	"log"
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func GetAllSuggestions(c echo.Context) error {
	type Response struct {
		suggestion_id uint
		user_id       uint
		name          string
		picture       string
		email         string
		message       string
	}

	suggestions := []model.Suggestions{}
	responses := []Response{}

	if err := config.DB.Find(&suggestions).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	for i := 0; i < len(suggestions); i++ {
		var temp Response

		user := model.User{}
		if err := config.DB.First(&user, suggestions[i].User_id).Error; err != nil {
			log.Print(color.RedString(err.Error()))
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  400,
				"message": "bad request",
			})
		}

		temp.suggestion_id = suggestions[i].ID
		temp.user_id = suggestions[i].User_id
		temp.name = user.Name
		temp.picture = user.URL
		temp.email = user.Email
		temp.message = suggestions[i].Content
		responses = append(responses, temp)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get suggestions",
		"data":    responses,
	})
}

func GetSuggestionByID(c echo.Context) error {
	suggestion := model.Suggestions{}

	suggestionID := c.Param("id")

	if err := config.DB.First(&suggestion, suggestionID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	user := model.User{}
	if err := config.DB.First(&user, suggestion.User_id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to get suggestion by id",
		"data": map[string]interface{}{
			"suggestion_id": suggestion.ID,
			"user_id":       suggestion.User_id,
			"name":          user.Name,
			"picture":       user.URL,
			"email":         user.Email,
			"message":       suggestion.Content,
		},
	})
}

func DeleteSuggestionByID(c echo.Context) error {
	suggestion := model.Suggestions{}

	suggestionID := c.Param("id")

	if err := config.DB.First(&suggestion, suggestionID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if err := config.DB.Delete(&suggestion).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to delete suggestion",
	})
}
