package d2m

import (
	"log"
	"os"
)

func Execute() {
	app := AppCmd()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
