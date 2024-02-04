package main

import (
	"log"
	"os"

	"github.com/vuon9/d2m/internal"
)

func main() {
	app := internal.NewCLI()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
