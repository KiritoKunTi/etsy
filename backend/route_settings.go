package main

import (
	"encoding/json"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"github.com/TutorialEdge/realtime-chat-go-react/utils"
	"net/http"
)

const photo = "avatar"

func UpdateUserPhotoHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	request.ParseMultipartForm(10 << 20)
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendMessage(writer, "Authorization request", http.StatusUnauthorized, nil)
		return
	}
	user, err := db.UserByID(session.User_ID)
	if err != nil {
		utils.SendErrorMessage(writer, err)
		return
	}
	photo, err := utils.PasteFile(request, photo)
	if err != nil {
		utils.SendErrorMessage(writer, err)
		return
	}
	user.Photo = photo
	user.UpdatePhoto()
	user.HideInfo()
	utils.SendMessage(writer, "Successfully updated", http.StatusOK, user)
}

func UpdateUserHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var user db.User
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		utils.SendErrorMessage(writer, err)
	}
	if user.Password != user.Repassword {
		utils.SendMessage(writer, "Passwords' don't match", http.StatusNotAcceptable, user)
		return
	}
	userFromDB, err := db.UserByID(user.ID)
	if err != nil {
		utils.SendErrorMessage(writer, err)
		return
	}
	if userFromDB.Password != db.Encrypt(user.Password) {
		utils.SendMessage(writer, "Incorrect password", http.StatusNotAcceptable, user)
		return
	}
	if err = user.Update(); err != nil {
		utils.SendErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, "Succesfully updated", http.StatusOK, user)
}