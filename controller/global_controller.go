package controller

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/agriplant/config"
	"github.com/agriplant/model"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func Hello_World(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "Hello World. OK",
		"no_test": 8,
	})
}

func Show_all_DB(c echo.Context) error {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	type TemplateData struct {
		Admins []model.Admin
		Users  []model.User
	}

	// Get Admins
	var admins []model.Admin
	if err_find_admins := config.DB.Find(&admins).Error; err_find_admins != nil {
		log.Print(color.RedString(err_find_admins.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	// Get Users
	var users []model.User
	if err_find_users := config.DB.Find(&users).Error; err_find_users != nil {
		log.Print(color.RedString(err_find_users.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	data := TemplateData{
		Admins: admins,
		Users: users,
	}

	err := tmpl.Execute(c.Response().Writer, data)
	if err != nil {
		return err
	}

	return nil
}

func Upload_pictures(c echo.Context) error {
	// Menerima file-file gambar dari request body
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, "Failed to read pictures files")
	}

	pictures := form.File["pictures"] // pictures adalah nama field untuk file-file gambar

	var urls []string
	for _, pictureFile := range pictures {
		url := Save_picture(pictureFile)
		urls = append(urls, url)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"urls": urls,
	})
}

func Save_picture(pictureFile *multipart.FileHeader) string {
	encodedImage := Encode_base64(pictureFile)
	url := Decode_and_write_base64(encodedImage)
	return url
}

func Encode_base64(pictureFile *multipart.FileHeader) string {

	// Buka file gambar
	src, err_open := pictureFile.Open()
	if err_open != nil {
		log.Println(err_open)
		return "Failed to open image file"
	}
	defer src.Close()

	// Baca isi file gambar
	imageData, err_read := ioutil.ReadAll(src)
	if err_read != nil {
		log.Println(err_read)
		return "Failed to read image file"
	}

	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	return encodedImage
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
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  404,
			"message": "not found",
		})
	}

	return c.File(url)
}
