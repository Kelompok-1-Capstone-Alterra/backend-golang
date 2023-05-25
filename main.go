package main

import (
	"github.com/agriplant/route"
)

func main() {
	e := route.New()
	e.Start(":8080")
}
