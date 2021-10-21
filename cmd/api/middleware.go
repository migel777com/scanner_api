package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func (app *application) requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		bearerToken := c.GetHeader("Authorization")
		bearerToken = strings.Trim(bearerToken, "Bearer")
		bearerToken = strings.Trim(bearerToken, " ")
		fmt.Println(bearerToken)

		_, err := app.models.Token.GetByToken(bearerToken)
		if err != nil {
			if err.Error() == "not found" {
				app.NotFoundResponse(err, c)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
