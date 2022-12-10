package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
)

func signUp(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "application/json")
	var user db.User
	json.NewDecoder(req.Body).Decode(&user)
	if user.Password == user.Repassword{
		if err := user.Create(); err != nil{
			if errors.Is(err, db.ErrExistsUsernameOrEmail) {
				res.WriteHeader(http.StatusUnprocessableEntity)
				errMessage := ErrorMessage{Message: err.Error(), Object: user}
				jsonResp, err := json.Marshal(errMessage)
				if err != nil {
					fmt.Println("error while marshaling", jsonResp)
					return
				}
				res.Write(jsonResp)
				return
			}
			if err != nil{
				fmt.Println("error while creating user", err)
				return
			}
		}
		res.WriteHeader(http.StatusCreated)
	}else{
		res.WriteHeader(http.StatusUnprocessableEntity)
		errMessage := ErrorMessage{Message: "Passwords don't match", Object: user}
		jsonResp, err := json.Marshal(errMessage)
		if err != nil{
			fmt.Println("error while marshaling", jsonResp)
			return
		}
		res.Write(jsonResp)
	}
}