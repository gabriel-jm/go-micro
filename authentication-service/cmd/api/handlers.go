package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var reqPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &reqPayload)

	if err != nil {
		log.Print("Error here")
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Repo.GetByEmail(reqPayload.Email)

	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := app.Repo.PasswordMatches(reqPayload.Password, user)

	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) logRequest(name, data string) error {
	entry := struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}{
		Name: name,
		Data: data,
	}

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	loggerUrl := "http://logger-service:8000/log"

	request, err := http.NewRequest("POST", loggerUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	_, err = app.Client.Do(request)

	if err != nil {
		return err
	}

	return nil
}
