package d2m

import (
	"context"

	"github.com/vuon9/d2m/d2m"
)

func Execute() {
	d2m.GetCLIMatches(context.Background())
}
