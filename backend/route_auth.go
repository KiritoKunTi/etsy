package main

import (
	"encoding/json"
	"errors"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"github.com/TutorialEdge/realtime-chat-go-react/utils"
	"net/http"
)

func signUp(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var user db.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		utils.SendAndPrintErrorMessage(res, err)
		return
	}
	if user.Password == user.Repassword {
		if err := user.Create(); err != nil {
			if errors.Is(err, db.ErrExistsUsernameOrEmail) {
				utils.SendMessage(res, utils.ExistsUserMessage, http.StatusNotAcceptable, user)
				return
			}
			if err != nil {
				utils.SendAndPrintErrorMessage(res, err)
				return
			}
		}
		user.HideInfo()
		utils.SendMessage(res, utils.SuccessfullyCreatedMessage, http.StatusCreated, user)
	} else {
		utils.SendMessage(res, utils.PasswordsMismatchMessage, http.StatusNotAcceptable, user)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendAndPrintErrorMessage(w, err)
		return
	}
	userFromDB, err := db.UserByEmailOrUsername(user.UsernameOrEmail)
	if err != nil {
		utils.SendMessage(w, utils.NotExistsUserMessage, http.StatusNotAcceptable, user)
		return
	}
	if userFromDB.Password == db.Encrypt(user.Password) {
		session, err := userFromDB.CreateSession()
		if err != nil {
			utils.SendAndPrintErrorMessage(w, err)
			return
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		userFromDB.HideInfo()
		utils.SendMessage(w, utils.SuccessfullyRequested, http.StatusOK, userFromDB)
		return
	}
	utils.SendMessage(w, "Password or email is not correct", http.StatusNotAcceptable, nil)
}

func logout(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		session := db.Session{UUID: cookie.Value}
		err = session.DeleteByUUID()
		if err != nil {
			utils.SendAndPrintErrorMessage(writer, err)
			return
		}
	}
	utils.SendMessage(writer, utils.SuccessfullyRequested, http.StatusOK, nil)
}
