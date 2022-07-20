package routes

import (
	"backend/server"
	"backend/server/handlers"
	m "backend/server/middlewares"
	"backend/services/auth"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func WrapHandler(h http.HandlerFunc) echo.HandlerFunc {
	return echo.WrapHandler(http.HandlerFunc(h))
}

func WrapMiddleware(m func(http.Handler) http.Handler) echo.MiddlewareFunc {
	return echo.WrapMiddleware(m)
}

func ConfigureRoutes(server *server.Server) {
	redisService, err := auth.NewRedisDB("localhost", "6379", "")
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()

	publicHandler := handlers.NewPublicHandler(server)
	authHandler := handlers.NewAuthHandler(redisService.Auth, tk)
	privateHandler := handlers.NewPrivateHandler(server)

	server.Echo.GET("/", WrapHandler(publicHandler.Hello))
	server.Echo.POST("/", WrapHandler(publicHandler.Post))

	server.Echo.POST("/auth/signin", WrapHandler(authHandler.SignIn))
	server.Echo.POST("/auth/signout", WrapHandler(authHandler.SignOut))
	server.Echo.GET("/auth/refresh", WrapHandler(authHandler.Refresh))

	server.Echo.GET("/private", WrapHandler(privateHandler.Tes), WrapMiddleware(m.IsLoggedIn))
	server.Echo.GET("/private/admin", WrapHandler(privateHandler.Tes), WrapMiddleware(m.IsLoggedIn), WrapMiddleware(m.IsAdmin))
}
