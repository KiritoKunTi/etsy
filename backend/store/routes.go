package store

import (
	"encoding/json"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"github.com/TutorialEdge/realtime-chat-go-react/utils"
	"net/http"
	"strconv"
)

const (
	AmountTopProducts    = 10
	AmountPaginationSize = 10
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
		utils.SendErrorMessage(writer, "Server error", http.StatusInternalServerError, product, err)
		return
	}
	if err := product.Create(); err != nil {
		utils.SendErrorMessage(writer, "Cannot create product", http.StatusNotAcceptable, product, err)
		return
	}
	utils.SendMessage(writer, "Successfully created", http.StatusOK, product)
}

func UpdateProductHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendMessage(writer, "Authorization request", http.StatusUnauthorized, nil)
		return
	}
	var product db.Product
	if err := json.NewDecoder(request.Body).Decode(&product); err != nil {
		utils.SendErrorMessage(writer, "Server error", http.StatusInternalServerError, product, err)
		return
	}
	if prod, err := db.ProductByID(product.ID); err != nil || prod.UserID != session.User_ID {
		utils.SendMessage(writer, "Bad Request", http.StatusBadRequest, nil)
		return
	}
	product.UserID = session.User_ID
	err = product.Update()
	if err != nil {
		utils.SendErrorMessage(writer, "Server error", http.StatusInternalServerError, nil, err)
		return
	}
	utils.SendMessage(writer, "Successfully updated", http.StatusOK, product)
}

func RecentProductsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	products, err := db.RecentProducts(AmountTopProducts)
	if err != nil {
		utils.SendErrorMessage(writer, "Server error", http.StatusInternalServerError, products, err)
		return
	}
	utils.SendMessage(writer, "Requested successfully", http.StatusOK, products)
}

func ProductsByCategoryHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	category_id, _ := strconv.Atoi(request.URL.Query().Get("category_id"))
	products, err := db.ProductsByCategoryID(category_id, AmountPaginationSize)
	if err != nil {
		utils.SendErrorMessage(writer, "Bad Request", http.StatusBadRequest, products, err)
		return
	}
	utils.SendMessage(writer, "Requested successfully", http.StatusOK, products)
}

func AllCategoriesHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	products, err := db.AllCategoriesWithGivenAmountProducts(AmountPaginationSize)
	if err != nil {
		utils.SendErrorMessage(writer, "Server error", http.StatusInternalServerError, products, err)
		return
	}
	utils.SendMessage(writer, "Requested successfully", http.StatusOK, products)
}
