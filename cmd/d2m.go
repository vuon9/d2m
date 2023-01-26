package cmd

import (
	"context"

	"github.com/vuon9/d2m/d2m"
)

func Execute() {
	prog := d2m.NewApp()
	prog.Run(context.Background())
}
