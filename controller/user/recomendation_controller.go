package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func GetProducts(c echo.Context) error {
	product := []model.Product{}

	if err := config.DB.Find(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		},
		)
	}

	// Populate Pictures field for each product
	// limit the pictures to 1
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
		product[i].Pictures = product[i].Pictures[:1]
	}

	// split product by its category
	var (
		seeds       []model.Product
		pesticides  []model.Product
		fertilizers []model.Product
		tools       []model.Product
	)

	for i := 0; i < len(product); i++ {
		if product[i].Category == "Bibit" {
			seeds = append(seeds, product[i])
		} else if product[i].Category == "Pestisida" {
			pesticides = append(pesticides, product[i])
		} else if product[i].Category == "Alat tani" {
			tools = append(tools, product[i])
		} else if product[i].Category == "Pupuk" {
			fertilizers = append(fertilizers, product[i])
		}
	}

	// use product response struct
	var (
		seedsResponse       []model.ProductResponse
		pesticidesResponse  []model.ProductResponse
		fertilizersResponse []model.ProductResponse
		toolsResponse       []model.ProductResponse
	)

	for i := 0; i < len(seeds); i++ {
		seedPictures := make([]string, len(seeds[i].Pictures))
		for j, pic := range seeds[i].Pictures {
			seedPictures[j] = pic.URL
		}
		seedsResponse = append(seedsResponse, model.ProductResponse{
			ID:       seeds[i].ID,
			Pictures: seedPictures,
			Name:     seeds[i].Name,
			Price:    seeds[i].Price,
			Seen:     seeds[i].Seen,
		})
	}

	for i := 0; i < len(pesticides); i++ {
		pesticidePictures := make([]string, len(pesticides[i].Pictures))
		for j, pic := range pesticides[i].Pictures {
			pesticidePictures[j] = pic.URL
		}
		pesticidesResponse = append(pesticidesResponse, model.ProductResponse{
			ID:       pesticides[i].ID,
			Pictures: pesticidePictures,
			Name:     pesticides[i].Name,
			Price:    pesticides[i].Price,
			Seen:     pesticides[i].Seen,
		})
	}

	for i := 0; i < len(fertilizers); i++ {
		fertilizerPictures := make([]string, len(fertilizers[i].Pictures))
		for j, pic := range fertilizers[i].Pictures {
			fertilizerPictures[j] = pic.URL
		}
		fertilizersResponse = append(fertilizersResponse, model.ProductResponse{
			ID:       fertilizers[i].ID,
			Pictures: fertilizerPictures,
			Name:     fertilizers[i].Name,
			Price:    fertilizers[i].Price,
			Seen:     fertilizers[i].Seen,
		})
	}

	for i := 0; i < len(tools); i++ {
		toolPictures := make([]string, len(tools[i].Pictures))
		for j, pic := range tools[i].Pictures {
			toolPictures[j] = pic.URL
		}
		toolsResponse = append(toolsResponse, model.ProductResponse{
			ID:       tools[i].ID,
			Pictures: toolPictures,
			Name:     tools[i].Name,
			Price:    tools[i].Price,
			Seen:     tools[i].Seen,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success to retrieve products data",
		"data": map[string]interface{}{
			"seeds":       seedsResponse,
			"pesticides":  pesticidesResponse,
			"fertilizers": fertilizersResponse,
			"tools":       toolsResponse,
		},
	})

}

func GetProductsByName(c echo.Context) error {
	name := c.QueryParam("name")

	product := []model.Product{}

	if err := config.DB.Where("name LIKE ?", "%"+name+"%").Find(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		},
		)
	}

	// Populate Pictures field for each product
	// limit the pictures to 1
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
		product[i].Pictures = product[i].Pictures[:1]
	}

	// use product response struct
	var productResponse []model.ProductResponse

	for i := 0; i < len(product); i++ {
		pictureURLs := make([]string, len(product[i].Pictures))
		for j, pic := range product[i].Pictures {
			pictureURLs[j] = pic.URL
		}
		productResponse = append(productResponse, model.ProductResponse{
			ID:       product[i].ID,
			Pictures: pictureURLs,
			Name:     product[i].Name,
			Price:    product[i].Price,
			Seen:     product[i].Seen,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get products by name",
		"data":    productResponse,
	})
}

func GetProductsByCategory(c echo.Context) error {
	categoryParam := c.Param("category")
	categoryNumber, err := strconv.Atoi(categoryParam)
	if err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  400,
			"message": "invalid category parameter",
		})
	}

	category := getCategory(categoryNumber)

	product := []model.Product{}

	if err := config.DB.Where("category = ?", category).Find(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Populate Pictures field for each product
	// limit the pictures to 1
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
		product[i].Pictures = product[i].Pictures[:1]
	}

	// use product response struct
	var productResponse []model.ProductResponse

	for i := 0; i < len(product); i++ {
		pictureURLs := make([]string, len(product[i].Pictures))
		for j, pic := range product[i].Pictures {
			pictureURLs[j] = pic.URL
		}
		productResponse = append(productResponse, model.ProductResponse{
			ID:       product[i].ID,
			Pictures: pictureURLs,
			Name:     product[i].Name,
			Price:    product[i].Price,
			Seen:     product[i].Seen,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get products by category",
		"data": map[string]interface{}{
			"category": category,
			"products": productResponse,
		},
	})
}

func GetProductsByCategoryAndName(c echo.Context) error {
	category := c.Param("category")
	name := c.QueryParam("name")

	product := []model.Product{}

	if err := config.DB.Where("category = ? AND name LIKE ?", category, "%"+name+"%").Find(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		},
		)
	}

	// Populate Pictures field for each product
	// limit the pictures to 1
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
		product[i].Pictures = product[i].Pictures[:1]
	}

	// use product response struct
	var productResponse []model.ProductResponse

	for i := 0; i < len(product); i++ {
		pictureURLs := make([]string, len(product[i].Pictures))
		for j, pic := range product[i].Pictures {
			pictureURLs[j] = pic.URL
		}
		productResponse = append(productResponse, model.ProductResponse{
			ID:       product[i].ID,
			Pictures: pictureURLs,
			Name:     product[i].Name,
			Price:    product[i].Price,
			Seen:     product[i].Seen,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get products by category and name",
		"data": map[string]interface{}{
			"category": category,
			"products": productResponse,
		},
	})
}

func GetProductByID(c echo.Context) error {
	id := c.Param("id")

	product := model.Product{}

	if err := config.DB.Where("id = ?", id).Find(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Increment the "seen" field by 1
	product.Seen++
	if err := config.DB.Model(&model.Product{}).Where("id = ?", id).Update("seen", product.Seen).Error; err != nil {
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

	// Use product response struct
	productResponse := struct {
		ID          uint     `json:"id"`
		Pictures    []string `json:"product_pictures"`
		Name        string   `json:"product_name"`
		Category    string   `json:"product_category"`
		Description string   `json:"product_description"`
		Price       int      `json:"product_price"`
		Seen        int      `json:"product_seen"`
		Status      bool     `json:"product_status"`
		Brand       string   `json:"product_brand"`
		Condition   string   `json:"product_condition"`
		Unit        int      `json:"product_unit"`
		Weight      int      `json:"product_weight"`
		Form        string   `json:"product_form"`
	}{
		ID:          product.ID,
		Pictures:    pictureURLs,
		Name:        product.Name,
		Category:    product.Category,
		Description: product.Description,
		Price:       product.Price,
		Seen:        product.Seen,
		Status:      product.Status,
		Brand:       product.Brand,
		Condition:   product.Condition,
		Unit:        product.Unit,
		Weight:      product.Weight,
		Form:        product.Form,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get product by id",
		"data": map[string]interface{}{
			"product":          productResponse,
			"related-products": GetRelatedProducts(product.Category),
		},
	})
}

// GetRelatedProducts is a function to get related products by category
// used to get related products in GetProductByID function
func GetRelatedProducts(category string) []model.ProductResponse {
	product := []model.Product{}

	if err := config.DB.Where("category = ?", category).Find(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
	}

	// Populate Pictures field for each product
	// limit the pictures to 1
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
		product[i].Pictures = product[i].Pictures[:1]
	}

	// remove the first product (the product that is being viewed)
	product = append(product[:0], product[1:]...)

	// use product response struct
	var productResponse []model.ProductResponse

	for i := 0; i < len(product); i++ {
		pictureURLs := make([]string, len(product[i].Pictures))
		for j, pic := range product[i].Pictures {
			pictureURLs[j] = pic.URL
		}
		productResponse = append(productResponse, model.ProductResponse{
			ID:       product[i].ID,
			Pictures: pictureURLs,
			Name:     product[i].Name,
			Price:    product[i].Price,
			Seen:     product[i].Seen,
		})
	}

	return productResponse
}

// getCategory is a function to get category name by its number
// used in GetProductsByCategory function & GetProductsByCategoryAndName function
func getCategory(number int) string {
	switch number {
	case 1:
		return "Alat tani"
	case 2:
		return "Bibit"
	case 3:
		return "Pestisida"
	case 4:
		return "Pupuk"
	default:
		return ""
	}
}
