package main

import (
	"log"
	"os"

	"github.com/vuon9/d2m/service"
)

func main() {
	app := service.NewCLIApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
