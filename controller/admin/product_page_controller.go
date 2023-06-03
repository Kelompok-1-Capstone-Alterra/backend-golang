package admin

import (
	"log"
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func GetProducts(c echo.Context) error {
	products := []model.Product{}

	if err := config.DB.Find(&products).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Iterate over each product record and generate custom response
	var responses []interface{}
	for _, product := range products {
		// Populate Pictures field for each product
		config.DB.Model(&product).Association("Pictures").Find(&product.Pictures)

		// Extract the first picture URL
		pictureURL := ""
		if len(product.Pictures) > 0 {
			pictureURL = product.Pictures[0].URL
		}

		response := struct {
			ID         uint   `json:"id"`
			Created_at string `json:"created_at"`
			Updated_at string `json:"updated_at"`
			Deleted_at string `json:"deleted_at"`
			Picture    string `json:"product_picture"`
			Name       string `json:"product_name"`
			SellerName string `json:"product_seller_name"`
			Price      int    `json:"product_price"`
			Category   string `json:"product_category"`
			Status     bool   `json:"product_status"`
		}{
			ID:         product.ID,
			Created_at: product.CreatedAt.String(),
			Updated_at: product.UpdatedAt.String(),
			Deleted_at: product.DeletedAt.Time.String(),
			Picture:    pictureURL,
			Name:       product.Name,
			SellerName: product.SellerName,
			Price:      product.Price,
			Category:   product.Category,
			Status:     product.Status,
		}
		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    responses,
	})
}

func CreateProduct(c echo.Context) error {
	product := model.Product{}

	c.Bind(&product)

	admin := model.Admin{}

	// Get user by id
	// If user not found, return error
	if err := config.DB.First(&admin, product.AdminID).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// set admin id to article
	product.AdminID = admin.ID

	// save article to database
	if err := config.DB.Save(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Extract picture URLs
	pictureURLs := make([]string, len(product.Pictures))
	for i, pic := range product.Pictures {
		pictureURLs[i] = pic.URL
	}

	response := struct {
		ID          uint     `json:"id"`
		Pictures    []string `json:"product_pictures"`
		Name        string   `json:"product_name"`
		Category    string   `json:"product_category"`
		Description string   `json:"product_description"`
		Price       int      `json:"product_price"`
		Status      bool     `json:"product_status"`
		Brand       string   `json:"product_brand"`
		Condition   string   `json:"product_condition"`
		Unit        int      `json:"product_unit"`
		Weight      int      `json:"product_weight"`
		Form        string   `json:"product_form"`
		SellerName  string   `json:"product_seller_name"`
		SellerPhone string   `json:"product_seller_phone"`
		AdminID     uint     `json:"admin_id"`
	}{
		ID:          product.ID,
		Pictures:    pictureURLs,
		Name:        product.Name,
		Category:    product.Category,
		Description: product.Description,
		Price:       product.Price,
		Status:      product.Status,
		Brand:       product.Brand,
		Condition:   product.Condition,
		Unit:        product.Unit,
		Weight:      product.Weight,
		Form:        product.Form,
		SellerName:  product.SellerName,
		SellerPhone: product.SellerPhone,
		AdminID:     product.AdminID,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
	})
}

func GetProductByID(c echo.Context) error {
	product := model.Product{}

	id := c.Param("id")

	// Get product by id
	// If product not found, return error
	if err := config.DB.First(&product, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Populate Pictures field for each product
	config.DB.Model(&product).Association("Pictures").Find(&product.Pictures)

	// Extract picture URLs
	pictureURLs := make([]string, len(product.Pictures))
	for i, pic := range product.Pictures {
		pictureURLs[i] = pic.URL
	}

	response := struct {
		ID          uint     `json:"id"`
		Pictures    []string `json:"product_pictures"`
		Name        string   `json:"product_name"`
		Category    string   `json:"product_category"`
		Description string   `json:"product_description"`
		Price       int      `json:"product_price"`
		Status      bool     `json:"product_status"`
		Brand       string   `json:"product_brand"`
		Condition   string   `json:"product_condition"`
		Unit        int      `json:"product_unit"`
		Weight      int      `json:"product_weight"`
		Form        string   `json:"product_form"`
		SellerName  string   `json:"product_seller_name"`
		SellerPhone string   `json:"product_seller_phone"`
		AdminID     uint     `json:"admin_id"`
	}{
		ID:          product.ID,
		Pictures:    pictureURLs,
		Name:        product.Name,
		Category:    product.Category,
		Description: product.Description,
		Price:       product.Price,
		Status:      product.Status,
		Brand:       product.Brand,
		Condition:   product.Condition,
		Unit:        product.Unit,
		Weight:      product.Weight,
		Form:        product.Form,
		SellerName:  product.SellerName,
		SellerPhone: product.SellerPhone,
		AdminID:     product.AdminID,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    response,
	})
}

func DeleteProductByID(c echo.Context) error {
	product := model.Product{}

	id := c.Param("id")

	// Get product by id
	// If product not found, return error
	if err := config.DB.First(&product, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Delete product
	if err := config.DB.Delete(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}

func UpdateProductByID(c echo.Context) error {
	product := model.Product{}

	id := c.Param("id")

	// Get product by id
	// If product not found, return error
	if err := config.DB.First(&product, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	c.Bind(&product)

	// save product to database
	if err := config.DB.Save(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// remove article_id from product_pictures
	for i := 0; i < len(product.Pictures); i++ {
		product.Pictures[i].ArticleID = nil
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    product,
	})
}

func GetProductsByKeyword(c echo.Context) error {
	product := []model.Product{}

	keyword := c.QueryParam("keyword")

	// Get product by keyword
	// If product not found, return error
	if err := config.DB.Where("name LIKE ?", "%"+keyword+"%").Find(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Populate Pictures field for each product
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
	}

	// remove article_id from product_pictures
	for i := 0; i < len(product); i++ {
		for j := 0; j < len(product[i].Pictures); j++ {
			product[i].Pictures[j].ArticleID = nil
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    product,
	})
}

func GetProductsDisplay(c echo.Context) error {
	product := []model.Product{}

	// Get product by keyword
	// If product not found, return error
	if err := config.DB.Where("status = ?", true).Find(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Populate Pictures field for each product
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
	}

	// remove article_id from product_pictures
	for i := 0; i < len(product); i++ {
		for j := 0; j < len(product[i].Pictures); j++ {
			product[i].Pictures[j].ArticleID = nil
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    product,
	})
}

func GetProductsArchive(c echo.Context) error {
	product := []model.Product{}

	// Get product by keyword
	// If product not found, return error
	if err := config.DB.Where("status = ?", false).Find(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Populate Pictures field for each product
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
	}

	// remove article_id field from product_pictures
	for i := 0; i < len(product); i++ {
		for j := 0; j < len(product[i].Pictures); j++ {
			product[i].Pictures[j].ArticleID = nil
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    product,
	})
}
