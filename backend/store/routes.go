package store

import (
	"encoding/json"
	"fmt"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"github.com/TutorialEdge/realtime-chat-go-react/utils"
	"net/http"
)

const (
	AmountTopProducts = 10
)

func CreateProductHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendMessage(writer, "Authorization request", http.StatusUnauthorized, nil)
		return
	}
	var product db.Product
	product.UserID = session.User_ID
	if err := json.NewDecoder(request.Body).Decode(&product); err != nil {
		utils.SendMessage(writer, "Server error", http.StatusInternalServerError, product)
		fmt.Println(err)
		return
	}
	if err := product.Create(); err != nil {
		utils.SendMessage(writer, "Cannot create product", http.StatusNotAcceptable, product)
		fmt.Println(err)
		return
	}
	utils.SendMessage(writer, "Successfully created", http.StatusOK, product)
}

func RecentProductsHander(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	products, err := db.RecentProducts(AmountTopProducts)
	if err != nil {
		utils.SendMessage(writer, "Server error", http.StatusInternalServerError, products)
		fmt.Println(err)
		return
	}
	utils.SendMessage(writer, "Requested successfully", http.StatusOK, products)
}
