package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *application) GetBasket(c *gin.Context) {
	var input struct {
		Url string `json:"url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		app.BadRequest(err, c)
		return
	}

	basket := app.models.Product.GetBasket(input.Url, app.browser)

	c.JSON(http.StatusOK, gin.H{"basket": basket})
	return
}
