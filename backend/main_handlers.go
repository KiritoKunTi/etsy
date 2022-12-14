package main

func MainHandlers() {
	Router.HandleFunc("/authorization/signup", signUp)
}
