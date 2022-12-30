package store

import (
	"encoding/json"
	"fmt"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"net/http"
)

func CreateProduct(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var product db.Product
	if err := json.NewDecoder(request.Body).Decode(&product); err != nil {
		sendMessage(writer, "Server error", http.StatusInternalServerError, product)
		fmt.Println(err)
		return
	}
	if err := product.Create(); err != nil {
		sendMessage(writer, "Cannot create product", http.StatusNotAcceptable, product)
		fmt.Println(err)
		return
	}
	sendMessage(writer, "Successfully created", http.StatusOK, product)

}
