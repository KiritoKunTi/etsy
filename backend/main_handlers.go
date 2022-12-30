package main

func MainHandlers() {
	Router.HandleFunc("/authorization/signup", signUp)
	Router.HandleFunc("/authorization/login", login)
	Router.HandleFunc("/authorization/logout", logout)
}
