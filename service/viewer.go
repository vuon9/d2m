package service

import "github.com/urfave/cli/v2"

type Viewer interface {
	Name() string
	Action() func(*cli.Context) error
}
