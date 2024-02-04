package main

import "github.com/vuon9/d2m/internal"

func main() {
	err := internal.NewWebAPI().Run()
	if err != nil {
		panic(err)
	}
}
