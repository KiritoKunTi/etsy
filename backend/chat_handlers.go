package main

import (
	"github.com/TutorialEdge/realtime-chat-go-react/chat"
)

func ChatHandlers() {
	Router.HandleFunc("/chat", chat.Main)
	Router.HandleFunc("/chat/ws", chat.ServeWs)
}
