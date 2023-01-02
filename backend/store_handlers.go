package main

import "github.com/TutorialEdge/realtime-chat-go-react/store"

func Store_Handlers() {
	Router.HandleFunc("/store/createproduct", store.CreateProduct)
}
