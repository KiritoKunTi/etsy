package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorMessage struct {
	Message string      `json:"message"`
	Object  interface{} `json:"object"`
}

func sendErrorMessage(res http.ResponseWriter, message string, status int) {
	res.WriteHeader(status)
	errMessage := ErrorMessage{Message: message}
	jsonResp, err := json.Marshal(errMessage)
	if err != nil {
		fmt.Println("error while marshalling", jsonResp)
	}
	res.Write(jsonResp)
}
