package main

import (
	"github.com/gorilla/mux"
)

var Router = mux.NewRouter().StrictSlash(true)

func main() {
	ChatHandlers()

}
