package application

import (
	"backend/server"
	"backend/server/routes"
	"log"
)

func Start() {
	app := server.NewServer()

	routes.ConfigureRoutes(app)

	err := app.Start(":3000")
	if err != nil {
		log.Fatal("Port already used")
	}
}
