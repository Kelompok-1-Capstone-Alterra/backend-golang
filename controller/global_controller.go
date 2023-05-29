package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func Add_image(c echo.Context) error {
	type Image struct {
		StringBase64 string `form:"string_base64"`
	}

	var image Image
	c.Bind(&image)
	fmt.Println(image.StringBase64[:10])
	url := Decode_and_write_base64(image.StringBase64)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"url": url,
	})
}

func Decode_and_write_base64(stringBase64 string) string {
	dummyURL := "dummy.png"
	stringBase64 = strings.TrimPrefix(stringBase64, "data:image/png;base64,")

	decodedData, errDecode := base64.StdEncoding.DecodeString(stringBase64)
	if errDecode != nil {
		fmt.Println(errDecode)
		return dummyURL
	}

	imageURL := uuid.New().String() + ".png"
	file, errCreate := os.Create("assets/images/" + imageURL)
	if errCreate != nil {
		fmt.Println(errCreate)
		return dummyURL
	}
	defer file.Close()

	_, errWrite := file.Write(decodedData)
	if errWrite != nil {
		fmt.Println(errWrite)
		return dummyURL
	}

	return imageURL
}

func Get_picture(c echo.Context) error {
	url := c.Param("url")
	url = "assets/images/" + url

	// Memeriksa keberadaan file gambar
	_, err := os.Stat(url)

	if os.IsNotExist(err) {
		// Jika file tidak ditemukan, kirimkan pesan error
		return c.File("assets/images/dummy.png")
	}

	return c.File(url)
}
