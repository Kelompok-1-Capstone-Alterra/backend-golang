package controller

import (
	"log"
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/agriplant/utils"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	var user model.User

	// binding struct
	if err_bind := c.Bind(&user); err_bind != nil {
		// echo.NewHTTPError(http.StatusBadRequest)
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// hashing password
	user.BeforeCreateUser(config.DB)

	// register
	if err_insert := config.DB.Save(&user).Error; err_insert != nil {
		log.Print(color.RedString(err_insert.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "ok",
	})
}

func Login(c echo.Context) error {
	var user model.User

	var loginData struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	// binding struct
	if err_bind := c.Bind(&loginData); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// check email validity
	if err_select := config.DB.Where("email=?", loginData.Email).First(&user).Error; err_select != nil {
		log.Print(color.RedString(err_select.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// verify the password
	if !utils.ComparePassword(user.Password, loginData.Password) {
		log.Print(color.RedString("code=401, message=internal server error"))
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  401,
			"message": "unauthorized",
		})
	}

	// create token
	token, err_token := utils.CreateTokenUser(user.ID, user.Name)
	if err_token != nil {
		log.Print(color.RedString(err_token.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "successfully login",
		"token":   token,
	})
}

func Check_email_valid(c echo.Context) error {
	var user model.User
	if err_bind := c.Bind(&user); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	if err_select := config.DB.Where("email=?", user.Email).First(&user).Error; err_select != nil {
		// email not found
		log.Print(color.RedString(err_select.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// email found
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to check email",
		"user_id":      user.ID,
	})
}
