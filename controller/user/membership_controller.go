package controller

import (
	"fmt"
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/agriplant/utils"
	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	var user model.User

	// binding struct
	if err_bind := c.Bind(&user); err_bind != nil {
		fmt.Println(err_bind.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// hashing password
	user.BeforeCreateUser(config.DB)

	// register
	if err_insert := config.DB.Save(&user).Error; err_insert != nil {
		fmt.Println(err_insert.Error())
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
		fmt.Println(err_bind.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// check email validity
	if err_select := config.DB.Where("email=?", loginData.Email).First(&user).Error; err_select != nil {
		fmt.Println(err_select.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// verify the password
	if !utils.ComparePassword(user.Password, loginData.Password) {
		fmt.Println("Invalid password")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  401,
			"message": "unauthorized",
		})
	}

	// create token
	token, err_token := utils.CreateTokenUser(user.ID, user.Name)
	if err_token != nil {
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
