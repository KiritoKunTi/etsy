package store

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorMessage struct {
	Message string      `json:"message"`
	Object  interface{} `json:"object"`
}

func sendMessage(res http.ResponseWriter, message string, status int, obj interface{}) {
	res.WriteHeader(status)
	errMessage := ErrorMessage{Message: message, Object: obj}
	jsonResp, err := json.Marshal(errMessage)
	if err != nil {
		fmt.Println("error while marshalling", jsonResp)
	}
	res.Write(jsonResp)
}
