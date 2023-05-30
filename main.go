package main

import (
	"github.com/agriplant/config"
	"github.com/agriplant/route"
)

func main() {
	config.InitDB()
	
	e := route.New()
	e.Start(":8080")
}
