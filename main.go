package main

import (
	"log"
	"net/http"

	"github.com/agriplant/config"
	"github.com/agriplant/route"
)

func main() {
	config.InitDB()

	e := route.New()
	if err := e.StartTLS(":8080", "certificate.crt", "private.key"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
