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
		return echo.NewHTTPError(http.StatusInternalServerError)
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
