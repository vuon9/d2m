package cmd

import (
	"context"

	"github.com/vuon9/d2m/d2m"
)

func Execute() {
	d2m.RunProgram(context.Background())
}
