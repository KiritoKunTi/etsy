package main

import (
	"encoding/json"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"github.com/TutorialEdge/realtime-chat-go-react/utils"
	"net/http"
	"os"
)

const (
	userPhotoDir = "avatar"
	photoKey     = "photo"
)

func UpdateUserPhotoHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	request.ParseMultipartForm(10 << 20)
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	user, err := db.UserByID(session.User_ID)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	os.Remove(user.Photo)
	photo, err := utils.PasteFile(request, userPhotoDir, photoKey, user.ID)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	user.Photo = photo
	user.UpdatePhoto()
	user.HideInfo()
	utils.SendMessage(writer, utils.SuccessfullyUpdatedMessage, http.StatusOK, user)
}

func UpdateUserHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var user db.User
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
	}
	if user.Password != user.Repassword {
		utils.SendMessage(writer, utils.PasswordsMismatchMessage, http.StatusNotAcceptable, user)
		return
	}
	userFromDB, err := db.UserByID(user.ID)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	if userFromDB.Password != db.Encrypt(user.Password) {
		utils.SendMessage(writer, utils.PasswordsMismatchMessage, http.StatusNotAcceptable, user)
		return
	}
	if err = user.Update(); err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.SuccessfullyUpdatedMessage, http.StatusOK, user)
}

func DeactivateUserHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "json/application")

}
