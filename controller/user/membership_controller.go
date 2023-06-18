package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/agriplant/utils"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

// MEMBERSHIP - [Endpoint 1 : Login]
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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  401,
			"message": "unauthorized",
		})
	}

	// verify the password
	fmt.Println(user.Password, loginData.Password)
	if !utils.ComparePassword(user.Password, loginData.Password) {
		log.Print(color.RedString("code=401, message=unauthorized"))
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
		"message": "success to login",
		"token":   token,
	})
}

// MEMBERSHIP - [Endpoint 2 : Register]
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

	// request body validation
	if user.Password == "" || user.Name == "" {
		log.Print(color.RedString("request body can't empty"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// email validation
	if !utils.Is_email_valid(user.Email) {
		log.Print(color.RedString("email not valid"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// password validation
	if len(user.Password) > 20 || len(user.Password) < 8 {
		log.Print(color.RedString("password min 8 and max 20 character"))
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to register",
	})
}

// MEMBERSHIP - [Endpoint 3 : Check email valid]
func Check_email_valid(c echo.Context) error {
	var user model.User
	if err_bind := c.Bind(&user); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
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
		"user_id": user.ID,
	})
}

// MEMBERSHIP - [Endpoint 4 : Reset password]
func Reset_password(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("user_id"))

	var user model.User
	if err_first := config.DB.First(&user, id).Error; err_first != nil {
		return c.JSON((http.StatusNotFound), map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	oldPassword := user.Password

	// binding struct
	if err_bind := c.Bind(&user); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// password validation
	if len(user.Password) > 20 || len(user.Password) < 8 {
		log.Print(color.RedString("password min 8 and max 20 character"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// verify the password
	fmt.Println(user.Password, oldPassword)
	if utils.ComparePassword(oldPassword, user.Password) {
		log.Print(color.RedString("new password can't same with old password"))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}
	// hashing password
	user.BeforeCreateUser(config.DB)

	if err_update := config.DB.Save(&user).Error; err_update != nil {
		log.Print(color.RedString(err_update.Error()))
		return c.JSON((http.StatusInternalServerError), map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to reset password",
	})
}
