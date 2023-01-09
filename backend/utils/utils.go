package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"github.com/TutorialEdge/realtime-chat-go-react/db/db_store"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Message struct {
	Message string      `json:"message"`
	Object  interface{} `json:"object"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

const (
	AmountTopProducts = 10
	PaginationSize    = 10
)

func SendMessage(res http.ResponseWriter, message string, status int, obj interface{}) {
	res.WriteHeader(status)
	errMessage := Message{Message: message, Object: obj}
	jsonResp, err := json.Marshal(errMessage)
	if err != nil {
		fmt.Println("error while marshalling", jsonResp)
	}
	res.Write(jsonResp)
}

func SendAndPrintErrorMessage(res http.ResponseWriter, err error) {
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

func PasteFile(request *http.Request, file string, key string, userID int) (filename string, err error) {
	in, header, err := request.FormFile(key)
	if err != nil {
		return
	}
	defer in.Close()
	filename = "private/" + file + "/" + strconv.Itoa(userID) + "/" + header.Filename
	os.Mkdir("private/"+file+"/"+strconv.Itoa(userID), os.ModePerm)
	out, err := os.Create(filename)
	if err != nil {
		return
	}
	defer out.Close()
	io.Copy(out, in)
	return filename, nil
}

func PasteProductPhoto(request *http.Request, product db_store.Product) (photos []db_store.ProductPhoto, err error) {
	for i := 1; i > 0; i++ {
		filename, err := PasteFile(request, "product_photos", "photo"+strconv.Itoa(i), product.ID)
		if err != nil {
			return photos, err
		}
		photo := db_store.ProductPhoto{ProductID: product.ID, Photo: filename}
		photos = append(photos, photo)
	}
	return
}
