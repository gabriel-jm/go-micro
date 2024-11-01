package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var reqPayload *RequestPayload = &RequestPayload{}

	err := app.readJSON(w, r, reqPayload)

	if err != nil {
		app.errorJSON(w, fmt.Errorf("[broker] error: %v", err))
		return
	}

	switch reqPayload.Action {
	case "auth":
		app.authenticate(w, reqPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, payload AuthPayload) {
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	request, err := http.NewRequest(
		"POST",
		"http://authentication-service:8000/authenticate",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("[Auth] Response body: %v", jsonFromService)
		app.errorJSON(w, errors.New(jsonFromService.Message))
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    jsonFromService.Data,
	}

	app.writeJSON(w, http.StatusOK, responsePayload)
}
