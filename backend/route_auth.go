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
		sendErrorMessage(res, "Server error", http.StatusInternalServerError)
		return
	}
	if user.Password == user.Repassword {
		if err := user.Create(); err != nil {
			if errors.Is(err, db.ErrExistsUsernameOrEmail) {
				sendErrorMessage(res, err.Error(), http.StatusNotAcceptable)
				return
			}
			if err != nil {
				fmt.Println("error while creating user")
				fmt.Println(err)
				return
			}
		}
		res.WriteHeader(http.StatusCreated)
	} else {
		sendErrorMessage(res, "Passwords' doesn't match", http.StatusNotAcceptable)
	}
}
