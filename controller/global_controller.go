package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
)

func Hello_World(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "Hello World. OK",
		"no_test": 12,
	})
}

func Upload_pictures(c echo.Context) error {
	// Menerima file-file gambar dari request body
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, "Failed to read picture files")
	}

	pictures := form.File["pictures"]

	var urls []string
	for _, pictureFile := range pictures {
		url, err := Save_picture(pictureFile)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Failed to save pictures")
		}
		urls = append(urls, url)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"urls": urls,
	})
}

func Save_picture(pictureFile *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	// Create a GCS client
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return "", err
	}
	defer client.Close()

	bucketName := "agriplant-image-bucket"
	objectName := generateBase64FileName(pictureFile.Filename)

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	file, err := pictureFile.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
	return url, nil
}

func generateBase64FileName(filename string) string {
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	base64Name := base64.RawURLEncoding.EncodeToString([]byte(baseName))

	return "images/" + base64Name + ext
}

// GLOBAL - [Endpoint 3 : Get picture]
func Get_picture(c echo.Context) error {
	url := c.Param("url")

	ctx := context.Background()

	// Create a GCS client
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return err
	}
	defer client.Close()

	bucketName := "agriplant-image-bucket"
	objectName := url

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	reader, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	w := c.Response().Writer
	if _, err := io.Copy(w, reader); err != nil {
		return err
	}

	return nil
}

// GLOBAL - [Endpoint 3 : Delete picture]
func Delete_picture_from_local(c echo.Context) error {
	url := c.Param("url")

	ctx := context.Background()

	// Create a GCS client
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return err
	}
	defer client.Close()

	bucketName := "agriplant-image-bucket"
	objectName := url

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	if err := obj.Delete(ctx); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "success to delete picture",
	})
}
