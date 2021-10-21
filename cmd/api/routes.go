package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := gin.Default()

	router.Use(
		//enabling AllowAllOrigins = true
		cors.Default(),
	)

	router.POST("/api/registration", app.Registration)
	router.POST("/api/auth", app.Auth)
	router.POST("/api/recover", app.Recover)
	router.POST("/api/update", app.requireAuth(), app.Update)

	router.POST("/getBasket", app.GetBasket)

	return router
}
