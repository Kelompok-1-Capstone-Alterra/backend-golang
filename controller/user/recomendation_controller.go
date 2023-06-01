package controller

import (
	"log"
	"net/http"

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
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
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
		seedsResponse = append(seedsResponse, model.ProductResponse{
			ID:       seeds[i].ID,
			Pictures: seeds[i].Pictures,
			Name:     seeds[i].Name,
			Price:    seeds[i].Price,
			Seen:     seeds[i].Seen,
		})
	}

	for i := 0; i < len(pesticides); i++ {
		pesticidesResponse = append(pesticidesResponse, model.ProductResponse{
			ID:       pesticides[i].ID,
			Pictures: pesticides[i].Pictures,
			Name:     pesticides[i].Name,
			Price:    pesticides[i].Price,
			Seen:     pesticides[i].Seen,
		})
	}

	for i := 0; i < len(fertilizers); i++ {
		fertilizersResponse = append(fertilizersResponse, model.ProductResponse{
			ID:       fertilizers[i].ID,
			Pictures: fertilizers[i].Pictures,
			Name:     fertilizers[i].Name,
			Price:    fertilizers[i].Price,
			Seen:     fertilizers[i].Seen,
		})
	}

	for i := 0; i < len(tools); i++ {
		toolsResponse = append(toolsResponse, model.ProductResponse{
			ID:       tools[i].ID,
			Pictures: tools[i].Pictures,
			Name:     tools[i].Name,
			Price:    tools[i].Price,
			Seen:     tools[i].Seen,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get products",
		"data": map[string]interface{}{
			"seeds":       seedsResponse,
			"pesticides":  pesticidesResponse,
			"fertilizers": fertilizersResponse,
			"tools":       toolsResponse,
		},
	})

}

func GetProductsByCategory(c echo.Context) error {
	category := c.Param("category")

	product := []model.Product{}

	if err := config.DB.Where("category = ?", category).Find(&product).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		},
		)
	}

	// Populate Pictures field for each product
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
	}

	// use product response struct
	var productResponse []model.ProductResponse

	for i := 0; i < len(product); i++ {
		productResponse = append(productResponse, model.ProductResponse{
			ID:       product[i].ID,
			Pictures: product[i].Pictures,
			Name:     product[i].Name,
			Price:    product[i].Price,
			Seen:     product[i].Seen,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get products",
		"data":    productResponse,
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
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
	}

	// use product response struct
	var productResponse []model.ProductResponse

	for i := 0; i < len(product); i++ {
		productResponse = append(productResponse, model.ProductResponse{
			ID:       product[i].ID,
			Pictures: product[i].Pictures,
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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get product by id",
		"data": map[string]interface{}{
			"product":          product,
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
	for i := 0; i < len(product); i++ {
		config.DB.Model(&product[i]).Association("Pictures").Find(&product[i].Pictures)
	}

	// remove the first product (the product that is being viewed)
	product = append(product[:0], product[1:]...)

	// use product response struct
	var productResponse []model.ProductResponse

	for i := 0; i < len(product); i++ {
		productResponse = append(productResponse, model.ProductResponse{
			ID:       product[i].ID,
			Pictures: product[i].Pictures,
			Name:     product[i].Name,
			Price:    product[i].Price,
			Seen:     product[i].Seen,
		})
	}

	return productResponse
}
