package admin

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

func CreateArticle(c echo.Context) error {
	article := model.Article{}

	if err := c.Bind(&article); err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Get admin id from JWT token
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	adminId, _ := utils.GetAdminIDFromToken(token)

	// Set admin ID to article
	article.AdminID = adminId

	// save article to database
	if err := config.DB.Save(&article).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	if err := config.DB.Save(&article).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Populate Pictures field for each article
	config.DB.Model(&article).Association("Pictures").Find(&article.Pictures)

	// extract picture urls
	pictureURLs := make([]string, len(article.Pictures))
	for i, pic := range article.Pictures {
		pictureURLs[i] = pic.URL
	}

	response := struct {
		ID          uint     `json:"id"`
		Created_at  string   `json:"created_at"`
		Updated_at  string   `json:"updated_at"`
		Deleted_at  string   `json:"deleted_at"`
		Title       string   `json:"article_title"`
		Pictures    []string `json:"article_pictures"`
		Description string   `json:"article_description"`
		View        int      `json:"article_view"`
		Like        int      `json:"article_like"`
	}{
		ID:          article.ID,
		Created_at:  article.CreatedAt.String(),
		Updated_at:  article.UpdatedAt.String(),
		Deleted_at:  article.DeletedAt.Time.String(),
		Title:       article.Title,
		Pictures:    pictureURLs,
		Description: article.Description,
		View:        article.View,
		Like:        article.Like,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
	})
}

func GetArticles(c echo.Context) error {
	articles := []model.Article{}

	// Get all articles
	if err := config.DB.Order("updated_at DESC").Find(&articles).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Iterate over each weather record and generate custom response
	var responses []interface{}
	for _, article := range articles {
		// Populate Pictures field for each weather
		config.DB.Model(&article).Association("Pictures").Find(&article.Pictures)

		// Extract picture URLs
		pictureURLs := make([]string, len(article.Pictures))
		for i, pic := range article.Pictures {
			pictureURLs[i] = pic.URL
		}

		response := struct {
			ID          uint     `json:"id"`
			Created_at  string   `json:"created_at"`
			Updated_at  string   `json:"updated_at"`
			Deleted_at  string   `json:"deleted_at"`
			Title       string   `json:"article_title"`
			Pictures    []string `json:"article_pictures"`
			Description string   `json:"article_description"`
			View        int      `json:"article_view"`
			Like        int      `json:"article_like"`
		}{
			ID:         article.ID,
			Created_at: article.CreatedAt.String(),
			Updated_at: article.UpdatedAt.String(),
			Deleted_at: article.DeletedAt.Time.String(),
			Title:      article.Title,
			Pictures:   pictureURLs,
			View:       article.View,
			Like:       article.Like,
		}

		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    responses,
	})
}

func GetArticlesByTitle(c echo.Context) error {
	title := c.QueryParam("title")

	articles := []model.Article{}

	// Retrieve articles by keyword
	if err := config.DB.Order("updated_at DESC").Where("title LIKE ?", "%"+title+"%").Find(&articles).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Iterate over each weather record and generate custom response
	var responses []interface{}
	for _, article := range articles {
		// Populate Pictures field for each weather
		config.DB.Model(&article).Association("Pictures").Find(&article.Pictures)

		// Extract picture URLs
		pictureURLs := make([]string, len(article.Pictures))
		for i, pic := range article.Pictures {
			pictureURLs[i] = pic.URL
		}

		response := struct {
			ID          uint     `json:"id"`
			Created_at  string   `json:"created_at"`
			Updated_at  string   `json:"updated_at"`
			Deleted_at  string   `json:"deleted_at"`
			Title       string   `json:"article_title"`
			Pictures    []string `json:"article_pictures"`
			Description string   `json:"article_description"`
			View        int      `json:"article_view"`
			Like        int      `json:"article_like"`
		}{
			ID:         article.ID,
			Created_at: article.CreatedAt.String(),
			Updated_at: article.UpdatedAt.String(),
			Deleted_at: article.DeletedAt.Time.String(),
			Title:      article.Title,
			Pictures:   pictureURLs,
			View:       article.View,
			Like:       article.Like,
		}

		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    responses,
	})
}

func GetArticleByID(c echo.Context) error {
	id := c.Param("id")

	article := model.Article{}

	// Get article by ID
	if err := config.DB.Where("id = ?", id).First(&article).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	// Populate Pictures field for each article
	config.DB.Model(&article).Association("Pictures").Find(&article.Pictures)

	// Extract picture URLs
	pictureURLs := make([]string, len(article.Pictures))
	for i, pic := range article.Pictures {
		pictureURLs[i] = pic.URL
	}

	response := struct {
		ID          uint     `json:"id"`
		Created_at  string   `json:"created_at"`
		Updated_at  string   `json:"updated_at"`
		Deleted_at  string   `json:"deleted_at"`
		Title       string   `json:"article_title"`
		Pictures    []string `json:"article_pictures"`
		Description string   `json:"article_description"`
		View        int      `json:"article_view"`
		Like        int      `json:"article_like"`
	}{
		ID:          article.ID,
		Created_at:  article.CreatedAt.String(),
		Updated_at:  article.UpdatedAt.String(),
		Deleted_at:  article.DeletedAt.Time.String(),
		Title:       article.Title,
		Pictures:    pictureURLs,
		Description: article.Description,
		View:        article.View,
		Like:        article.Like,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
	})
}

func UpdateArticleByID(c echo.Context) error {
	id := c.Param("id")

	article := model.Article{}

	if err := config.DB.First(&article, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	config.DB.Model(&article).Association("Pictures").Clear()

	// Bind new data to article
	if err := c.Bind(&article); err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// Get admin id from JWT token
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	adminId, _ := utils.GetAdminIDFromToken(token)

	// Set admin ID to article
	article.AdminID = adminId

	// Save article to database
	if err := config.DB.Save(&article).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Populate Pictures field for each article
	config.DB.Model(&article).Association("Pictures").Find(&article.Pictures)

	// Extract picture URLs
	pictureURLs := make([]string, len(article.Pictures))
	for i, pic := range article.Pictures {
		pictureURLs[i] = pic.URL
	}

	response := struct {
		ID          uint     `json:"id"`
		Created_at  string   `json:"created_at"`
		Updated_at  string   `json:"updated_at"`
		Deleted_at  string   `json:"deleted_at"`
		Title       string   `json:"article_title"`
		Pictures    []string `json:"article_pictures"`
		Description string   `json:"article_description"`
		View        int      `json:"article_view"`
		Like        int      `json:"article_like"`
	}{
		ID:          article.ID,
		Created_at:  article.CreatedAt.String(),
		Updated_at:  article.UpdatedAt.String(),
		Deleted_at:  article.DeletedAt.Time.String(),
		Title:       article.Title,
		Pictures:    pictureURLs,
		Description: article.Description,
		View:        article.View,
		Like:        article.Like,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
	})
}

func DeleteArticleByID(c echo.Context) error {
	id := c.Param("id")

	article := model.Article{}

	if err := config.DB.First(&article, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	config.DB.Model(&article).Association("Pictures").Find(&article.Pictures)
	for _, picture := range article.Pictures {
		if err_delete_picture := utils.Delete_picture(picture.URL); err_delete_picture != nil {
			log.Print(color.RedString(err_delete_picture.Error()))
		}
	}

	config.DB.Model(&article).Association("Pictures").Clear()

	// Delete article from database
	if err := config.DB.Delete(&article).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
