package main

import (
	"log"
	"net/http"
	"os"

	"github.com/agriplant/config"
	"github.com/agriplant/route"
)

func main() {
	config.InitDB()

	keyFilePath := "./capstonealterra-2426d3155b94.json"
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", keyFilePath)

	e := route.New()
	if err := e.StartTLS(":8080", "certificate.crt", "private.key"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
