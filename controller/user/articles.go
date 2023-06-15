package controller

import (
	"log"
	"math"
	"net/http"
	"time"

	"strings"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/agriplant/utils"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func GetArticlesTrending(c echo.Context) error {
	var articles []model.Article

	//get data by trending
	if err := config.DB.Order("`like` DESC").Find(&articles).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	var data []map[string]interface{}
	//Populate Pictures field for each article
	for i := 0; i < len(articles); i++ {
		config.DB.Model(&articles[i]).Association("Pictures").Find(&articles[i].Pictures)
		now := time.Now()

		convertTime := now.Sub(articles[i].CreatedAt).Hours() / 24
		convertTimeInt := int(math.Round(convertTime))
		result := map[string]interface{}{
			"article_id": articles[i].ID,
			"title":      articles[i].Title,
			"picture":    articles[i].Pictures[0].URL,
			"hours":      convertTimeInt,
		}
		data = append(data, result)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to retrieve latest articles data",
		"data":    data,
	})
}

func GetArticlesLatest(c echo.Context) error {
	var articles []model.Article

	//get data by latest
	if err := config.DB.Order("`created_at` DESC").Find(&articles).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	//Populate Pictures field for each article
	var data []map[string]interface{}
	for i := 0; i < len(articles); i++ {
		config.DB.Model(&articles[i]).Association("Pictures").Find(&articles[i].Pictures)
		now := time.Now()

		convertTime := now.Sub(articles[i].CreatedAt).Minutes()
		convertTimeInt := int(math.Round(convertTime))
		result := map[string]interface{}{
			"article_id": articles[i].ID,
			"title":      articles[i].Title,
			"picture":    articles[i].Pictures[0].URL,
			"minutes":    convertTimeInt,
		}
		data = append(data, result)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to retrieve latest articles data",
		"data":    data,
	})
}
func GetArticlesLiked(c echo.Context) error {
	var likeds []model.LikedArticles
	var articles []model.Article
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	//get data by latest
	if err := config.DB.Where("user_id =?", user_id).Find(&likeds).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}
	var idArticles []int
	for _, liked := range likeds {
		idArticles = append(idArticles, int(liked.ArticleID))
	}

	if errer := config.DB.Where("id IN ?", idArticles).Find(&articles).Error; errer != nil {
		log.Print(color.RedString(errer.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}
	var data []map[string]interface{}
	//Populate Pictures field for each article
	for i := 0; i < len(articles); i++ {
		config.DB.Model(&articles[i]).Association("Pictures").Find(&articles[i].Pictures)
		now := time.Now()
		convertTime := now.Sub(articles[i].CreatedAt).Minutes()
		convertTimeInt := int(math.Round(convertTime))
		result := map[string]interface{}{
			"article_id": articles[i].ID,
			"title":      articles[i].Title,
			"picture":    articles[i].Pictures[0].URL,
			"time":       convertTimeInt,
		}
		data = append(data, result)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to retrieve liked articles data",
		"data":    data,
	})
}
func GetArticlesbyID(c echo.Context) error {
	articles := model.Article{}
	id := c.Param("id")

	// Get product by id
	// If product not found, return error
	if err := config.DB.First(&articles, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}
	var data []map[string]interface{}
	// Populate Pictures field for each product
	config.DB.Model(&articles).Association("Pictures").Find(&articles.Pictures)
	result := map[string]interface{}{
		"article_id":  articles.ID,
		"title":       articles.Title,
		"picture":     articles.Pictures[0].URL,
		"description": articles.Description,
		"isliked":     false,
	}
	data = append(data, result)
	// remove article_id from articles_pictures
	for i := 0; i < len(articles.Pictures); i++ {
		articles.Pictures[i].ArticleID = nil

	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to retrieve latest articles data",
		"data":    data,
	})
}

func AddLikes(c echo.Context) error {
	like := model.LikedArticles{}
	articles_id, _ := StringToUintPointer(c.Param("article_id"))
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)
	// Get product by id
	// If product not found, return error
	like.ArticleID = *articles_id
	like.UserID = user_id
	if err := config.DB.Save(&like).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// // Populate Pictures field for each product
	// config.DB.Model(&articles).Association("Pictures").Find(&articles.Pictures)

	// // remove article_id from articles_pictures
	// for i := 0; i < len(articles.Pictures); i++ {
	// 	articles.Pictures[i].ArticleID = nil
	// }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to give like",
	})
}

func DeleteLikes(c echo.Context) error {
	like := model.LikedArticles{}
	articles_id, _ := StringToUintPointer(c.Param("article_id"))
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)
	// Get product by id
	// If product not found, return error

	if err := config.DB.Where("articles_id =?", *articles_id).Where("user_id=?", user_id).Delete(&like).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// // Populate Pictures field for each product
	// config.DB.Model(&articles).Association("Pictures").Find(&articles.Pictures)

	// // remove article_id from articles_pictures
	// for i := 0; i < len(articles.Pictures); i++ {
	// 	articles.Pictures[i].ArticleID = nil
	// }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to unlike",
	})
}
