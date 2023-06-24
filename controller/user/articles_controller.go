package controller

import (
	"log"
	"net/http"
	"strconv"
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

	data := []map[string]interface{}{}
	//Populate Pictures field for each article
	for i := 0; i < len(articles); i++ {
		config.DB.Model(&articles[i]).Association("Pictures").Find(&articles[i].Pictures)

		result := map[string]interface{}{
			"id":      articles[i].ID,
			"title":   articles[i].Title,
			"picture": articles[i].Pictures[0].URL,
			"post_at": articles[i].CreatedAt.UTC(),
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

		result := map[string]interface{}{
			"id":      articles[i].ID,
			"title":   articles[i].Title,
			"picture": articles[i].Pictures[0].URL,
			"post_at": articles[i].CreatedAt.UTC(),
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

	data := []map[string]interface{}{}
	for i := 0; i < len(articles); i++ {
		config.DB.Model(&articles[i]).Association("Pictures").Find(&articles[i].Pictures)

		result := map[string]interface{}{
			"id":      articles[i].ID,
			"title":   articles[i].Title,
			"picture": articles[i].Pictures[0].URL,
			"post_at": articles[i].CreatedAt.UTC(),
		}

		data = append(data, result)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to retrieve liked articles data",
		"data":    data,
	})
}

func GetArticlesByID(c echo.Context) error {
	articles := model.Article{}
	id, _ := strconv.Atoi(c.Param("id"))

	// Get product by id
	// If product not found, return error
	if err := config.DB.First(&articles, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Increment the "view" field by 1
	articles.View++
	if err := config.DB.Model(&model.Article{}).Where("id = ?", id).Update("view", articles.View).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	var data []map[string]interface{}
	// Populate Pictures field for each product
	config.DB.Model(&articles).Association("Pictures").Find(&articles.Pictures)

	// Check if the user has liked the article
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	userID, _ := utils.GetUserIDFromToken(token)
	var likedArticle model.LikedArticles
	if err := config.DB.Where("user_id = ? AND article_id = ?", userID, id).First(&likedArticle).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	result := map[string]interface{}{
		"id":          articles.ID,
		"title":       articles.Title,
		"picture":     articles.Pictures[0].URL,
		"description": articles.Description,
		"is_liked":    (likedArticle.ID != 0), // Check if the user has liked the article
	}
	data = append(data, result)

	// Remove article_id from articles_pictures
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
	article := model.Article{}

	articles_id, _ := StringToUintPointer(c.Param("article_id"))
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	user_id, _ := utils.GetUserIDFromToken(token)

	// increment the "like" field by 1
	article.Like++
	if err := config.DB.Model(&model.Article{}).Where("id = ?", articles_id).Update("like", article.Like).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to give like",
	})
}

func DeleteLikes(c echo.Context) error {
	articleID, _ := StringToUintPointer(c.Param("article_id"))

	// Fetch the article from the database
	article := model.Article{}
	if err := config.DB.Where("id = ?", articleID).First(&article).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// Extract user ID from the token
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	userID, _ := utils.GetUserIDFromToken(token)

	// Check if the user has liked the article
	var likedArticle model.LikedArticles
	if err := config.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&likedArticle).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// Delete the like entry
	if err := config.DB.Delete(&likedArticle).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Decrement the "like" field if it is greater than 0
	if article.Like > 0 {
		article.Like--
		if err := config.DB.Model(&model.Article{}).Where("id = ?", articleID).Update("like", article.Like).Error; err != nil {
			log.Print(color.RedString(err.Error()))
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status":  500,
				"message": "internal server error",
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to unlike",
	})
}
