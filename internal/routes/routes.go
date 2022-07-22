package routes

import (
	"backend/internal/handler"
	m "backend/internal/middlewares"
	"backend/pkg/auth"
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

func ConfigureRoutes(server *echo.Echo) {
	redisService, err := auth.NewRedisDB("localhost", "6379", "")
	if err != nil {
		log.Fatal(err)
	}

	tk := auth.NewToken()

	publicHandler := handler.PublicHandler{}
	authHandler := handler.NewAuthHandler(redisService.Auth, tk)
	privateHandler := handler.PrivateHandler{}
	adminHandler := handler.AdminHandler{}

	server.GET("/", WrapHandler(publicHandler.Hello))
	server.POST("/", WrapHandler(publicHandler.Post))

	auth := server.Group("/auth")
	auth.POST("/signin", WrapHandler(authHandler.SignIn))
	auth.POST("/signout", WrapHandler(authHandler.SignOut))
	auth.GET("/refresh", WrapHandler(authHandler.Refresh))

	private := server.Group("/private", WrapMiddleware(m.IsLoggedIn))
	private.GET("/", WrapHandler(privateHandler.Tes))

	admin := private.Group("/admin", WrapMiddleware(m.IsAdmin))
	admin.GET("/", WrapHandler(adminHandler.Tes))
}
