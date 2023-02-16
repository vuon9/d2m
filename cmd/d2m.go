package cmd

import (
	"context"

	"github.com/vuon9/d2m/internal/app"
)

func Execute() {
	prog := app.NewApp()
	prog.Run(context.Background())
}
