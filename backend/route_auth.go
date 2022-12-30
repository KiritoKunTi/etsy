package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"net/http"
)

func signUp(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var user db.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		sendMessage(res, "Server error", http.StatusInternalServerError, user)
		return
	}
	if user.Password == user.Repassword {
		if err := user.Create(); err != nil {
			if errors.Is(err, db.ErrExistsUsernameOrEmail) {
				sendMessage(res, "Already exists username or email", http.StatusNotAcceptable, user)
				return
			}
			if err != nil {
				fmt.Println("error while creating user")
				fmt.Println(err)
				return
			}
		}
		user.Password = ""
		user.Repassword = ""
		sendMessage(res, "Successfully registered", http.StatusCreated, user)
	} else {
		sendMessage(res, "Passwords' doesn't match", http.StatusNotAcceptable, user)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendMessage(w, "Server error", http.StatusInternalServerError, user)
		return
	}
	userFromDB, err := db.UserByEmailOrUsername(user.UsernameOrEmail)
	if err != nil {
		sendMessage(w, "Username or email doesn't exist", http.StatusNotAcceptable, user)
		return
	}
	if userFromDB.Password == db.Encrypt(user.Password) {
		session, err := userFromDB.CreateSession()
		if err != nil {
			sendMessage(w, "Server error", http.StatusInternalServerError, user)
			return
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		userFromDB.Password = ""
		sendMessage(w, "Successfully login", http.StatusOK, userFromDB)
		return
	}
	sendMessage(w, "Password or email is not correct", http.StatusNotAcceptable, user)
}

func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		session := db.Session{UUID: cookie.Value}
		session.DeleteByUUID()
	}
	http.Redirect(writer, request, "/", 302)
}
