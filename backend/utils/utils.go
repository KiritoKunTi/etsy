package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"net/http"
)

type ErrorMessage struct {
	Message string      `json:"message"`
	Object  interface{} `json:"object"`
}

func SendMessage(res http.ResponseWriter, message string, status int, obj interface{}) {
	res.WriteHeader(status)
	errMessage := ErrorMessage{Message: message, Object: obj}
	jsonResp, err := json.Marshal(errMessage)
	if err != nil {
		fmt.Println("error while marshalling", jsonResp)
	}
	res.Write(jsonResp)
}

func SendErrorMessage(res http.ResponseWriter, message string, status int, obj interface{}, err error) {
	fmt.Println(err)
	res.WriteHeader(status)
	errMessage := ErrorMessage{Message: message, Object: obj}
	jsonResp, err := json.Marshal(errMessage)
	if err != nil {
		fmt.Println("error while marshalling", jsonResp)
	}
	res.Write(jsonResp)
}

func Session(writer http.ResponseWriter, request *http.Request) (session db.Session, err error) {
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		session = db.Session{UUID: cookie.Value}
		if ok, _ := session.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}
