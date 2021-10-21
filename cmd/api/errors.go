package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) serverErrorResponse(err error, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	if err != nil {
		app.logger.PrintError(err, nil)
	}
}

func (app *application) BadRequest(err error, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	if err != nil {
		app.logger.PrintError(err, nil)
	}
}

func (app *application) NotFoundResponse(err error, c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	if err != nil {
		app.logger.PrintError(err, nil)
	}
}

func (app *application) NotAllowedResponse(err error, c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	if err != nil {
		app.logger.PrintError(err, nil)
	}
}

func (app *application) InvalidCredentials(err error, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	if err != nil {
		app.logger.PrintError(err, nil)
	}

}
