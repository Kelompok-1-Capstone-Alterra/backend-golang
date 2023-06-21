package controller

import (
	"context"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Create a global variable to hold the Google Cloud Storage client.
var client *storage.Client

func init() {
	// Initialize the Google Cloud Storage client.
	ctx := context.Background()
	// Replace "path/to/service-account-key.json" with the path to your service account key JSON file.
	// You can download the key file from the Google Cloud Console.
	var err error
	client, err = storage.NewClient(ctx, option.WithCredentialsFile("capstonealterra-0457bfb5b315.json"))
	if err != nil {
		log.Fatalf("Failed to create Google Cloud Storage client: %v", err)
	}
}

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
		return c.JSON(http.StatusBadRequest, "Failed to read pictures files")
	}

	pictures := form.File["pictures"] // pictures adalah nama field untuk file-file gambar

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
	encodedImage := Encode_base64(pictureFile)
	url, err := UploadToCloudStorage(encodedImage)
	if err != nil {
		return "", err
	}
	return url, nil
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

func UploadToCloudStorage(encodedImage string) (string, error) {
	ctx := context.Background()
	bucketName := "agriplant-image-bucket" // Replace with your actual bucket name.
	dummyURL := "dummy.png"
	encodedImage = strings.TrimPrefix(encodedImage, "data:image/png;base64,")

	decodedData, errDecode := base64.StdEncoding.DecodeString(encodedImage)
	if errDecode != nil {
		return dummyURL, errDecode
	}

	imageURL := uuid.New().String() + ".png"
	file, errCreate := os.Create(imageURL)
	if errCreate != nil {
		return dummyURL, errCreate
	}
	defer file.Close()

	_, errWrite := file.Write(decodedData)
	if errWrite != nil {
		return dummyURL, errWrite
	}

	// Open the bucket.
	bucket := client.Bucket(bucketName)

	// Open the file.
	obj := bucket.Object(imageURL)
	wc := obj.NewWriter(ctx)

	// Set the content type.
	wc.ContentType = "image/png"

	// Write the file to Cloud Storage.
	if _, err := wc.Write(decodedData); err != nil {
		return dummyURL, err
	}

	// Close the writer.
	if err := wc.Close(); err != nil {
		return dummyURL, err
	}

	return imageURL, nil
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
