package main

import (
	"backend/internal/routes"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	routes.ConfigureRoutes(app)

	err := app.Start(":3000")
	if err != nil {
		log.Fatal("Port already used")
	}
}
