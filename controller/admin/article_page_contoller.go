package admin

import (
	"log"
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func CreateArticle(c echo.Context) error {
	article := model.Article{}

	c.Bind(&article)

	admin := model.Admin{}
	// Get user by id
	// If user not found, return error
	if err := config.DB.First(&admin, article.AdminID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// set admin id to article
	article.AdminID = admin.ID

	// save article to database
	if err := config.DB.Save(&article).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    article,
	})
}

func GetArticles(c echo.Context) error {
	articles := []model.Article{}

	// Get all articles
	if err := config.DB.Find(&articles).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	//Populate Pictures field for each article
	for i := 0; i < len(articles); i++ {
		config.DB.Model(&articles[i]).Association("Pictures").Find(&articles[i].Pictures)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    articles,
	})
}

func GetArticlesByKeyword(c echo.Context) error {
	keyword := c.QueryParam("keyword")

	articles := []model.Article{}

	// Retrieve articles by keyword
	if err := config.DB.Where("title LIKE ?", "%"+keyword+"%").Find(&articles).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	//Populate Pictures field for each article
	for i := 0; i < len(articles); i++ {
		config.DB.Model(&articles[i]).Association("Pictures").Find(&articles[i].Pictures)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    articles,
	})
}
