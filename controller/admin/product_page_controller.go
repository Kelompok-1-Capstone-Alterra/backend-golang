package admin

import (
	"log"
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func CreateProduct(c echo.Context) error {
	product := model.Product{}

	c.Bind(&product)

	admin := model.Admin{}

	// Get user by id
	// If user not found, return error
	if err := config.DB.First(&admin, product.AdminID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// set admin id to article
	product.AdminID = admin.ID

	// save article to database
	if err := config.DB.Save(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    product,
	})
}
