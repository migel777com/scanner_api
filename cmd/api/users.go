package main

import (
	"demoapi/internal/data"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (app *application) Registration(c *gin.Context) {
	// Validate input
	var input data.User

	if err := c.ShouldBindJSON(&input); err != nil {
		app.BadRequest(err, c)
		return
	}
	_, err := app.models.User.GetByEmail(input.Email)

	if err != nil {
		if err.Error() == "not found" {
			//reg user
			layout := "2006-01-02"
			input.Birthdate, _ = time.Parse(layout, input.BirthdateString)

			// Create user input.Birthdate.Format(layout)
			user := &data.User{
				Email:     input.Email,
				Pass:      input.Pass,
				Name:      input.Name,
				Surname:   input.Surname,
				Birthdate: input.Birthdate,
			}
			//add validation here

			//hash password
			user.Pass, _ = HashPassword(user.Pass)

			id, err := app.models.User.Insert(user.Email, user.Pass, user.Name, user.Surname, user.Birthdate)
			if err != nil {
				app.serverErrorResponse(err, c)
				return
			}

			user.Id = int64(id)
			user.BirthdateString = user.Birthdate.Format(layout)
			c.JSON(http.StatusOK, gin.H{"user": user})
			return

		} else {
			app.serverErrorResponse(err, c)
			return
		}
	}

	//user already exists
	c.JSON(http.StatusOK, gin.H{"error": "User already exists"})
	return

}

func (app *application) Recover(c *gin.Context) {
	// Validate input
	var input struct {
		Email string `json:"Email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		app.BadRequest(err, c)
		return
	}

	user, err := app.models.User.GetByEmail(input.Email)

	if err != nil {
		app.serverErrorResponse(err, c)
		return
	}

	recoveredPassword := "123456!"

	defpass, _ := HashPassword(recoveredPassword)
	err = app.models.User.UpdatePassword(defpass, int(user.Id))
	if err != nil {
		app.serverErrorResponse(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"recovered password": recoveredPassword})
	return
}

func (app *application) Auth(c *gin.Context) {
	// Validate input
	var input struct {
		Email string `json:"Email"`
		Pass  string `json:"Pass"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		app.BadRequest(err, c)
		return
	}

	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		app.NotFoundResponse(err, c)
		return
	}

	possibleToken, err := app.models.Token.GetByUserId(int(user.Id))
	if err == nil {
		app.models.Token.Delete(int(possibleToken.Id))
	}

	token, err := uuid.NewRandom()
	if err != nil {
		app.serverErrorResponse(err, c)
		return
	}
	status, err := CompareHash(input.Pass, user.Pass)
	if !status {
		app.InvalidCredentials(err, c)
		return
	} else {

		_, err := app.models.Token.Insert(int(user.Id), token.String(), time.Now().Add(2*time.Minute))
		if err != nil {
			app.serverErrorResponse(err, c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"BEARER token": "Bearer " + token.String()})
	}

}

func (app *application) Update(c *gin.Context) {
	var input struct {
		Email           string `json:"Email"`
		Name            string `json:"Name"`
		Surname         string `json:"Surname"`
		BirthdateString string `json:"Birthdate"`
		Birthdate       time.Time
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		app.BadRequest(err, c)
		return
	}
	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		if err.Error() == "not found" {
			app.NotFoundResponse(err, c)
			return
		} else if err != nil {
			app.serverErrorResponse(err, c)
			return
		}

	} else {
		fmt.Println(user)
		layout := "2006-01-02"
		input.Birthdate, _ = time.Parse(layout, input.BirthdateString)
		err = app.models.User.Update(int(user.Id), user.Email, input.Name, input.Surname, input.Birthdate)
		if err != nil {
			app.serverErrorResponse(err, c)
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
		return
	}

}
