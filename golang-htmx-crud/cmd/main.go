package main

import (
	_ "github.com/lib/pq"

	"github.com/deeprecession/golang-htmx-crud/pkg/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		return
	}

	app.Run()
}
