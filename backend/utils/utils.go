package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"io"
	"net/http"
	"os"
)

type Message struct {
	Message string      `json:"message"`
	Object  interface{} `json:"object"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

const errorMessage = "ServerError"

func SendMessage(res http.ResponseWriter, message string, status int, obj interface{}) {
	res.WriteHeader(status)
	errMessage := Message{Message: message, Object: obj}
	jsonResp, err := json.Marshal(errMessage)
	if err != nil {
		fmt.Println("error while marshalling", jsonResp)
	}
	res.Write(jsonResp)
}

func SendErrorMessage(res http.ResponseWriter, err error) {
	fmt.Println(err)
	res.WriteHeader(http.StatusInternalServerError)
	errMessage := ErrorMessage{Message: errorMessage}
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

func PasteFile(request *http.Request, file string) (filename string, err error) {
	in, header, err := request.FormFile("photo")
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create("private/" + file + "/" + header.Filename)

	if err != nil {
		return
	}
	defer out.Close()
	io.Copy(out, in)
	return out.Name(), nil
}
