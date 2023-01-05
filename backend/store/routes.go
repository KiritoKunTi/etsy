package store

import (
	"encoding/json"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"github.com/TutorialEdge/realtime-chat-go-react/utils"
	"net/http"
	"strconv"
)

func CreateProductHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	var product db.Product
	product.SetUser(session.User_ID)
	if err := json.NewDecoder(request.Body).Decode(&product); err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	if !product.User.IsShop {
		utils.SendMessage(writer, "Incorrect account type", http.StatusNotAcceptable, nil)
		return
	}

	if err := product.Create(); err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.SuccessfullyCreatedMessage, http.StatusOK, product)
}

func UpdateProductHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	var product db.Product
	if err := json.NewDecoder(request.Body).Decode(&product); err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	if prod, err := db.ProductByID(product.ID); err != nil || prod.UserID != session.User_ID {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	product.SetUser(session.User_ID)
	err = product.Update()
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.SuccessfullyUpdatedMessage, http.StatusOK, product)
}

func UpdateProductPhotoHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	productID, err := strconv.Atoi(request.URL.Query().Get("product_id"))
	if err != nil {
		utils.SendMessage(writer, utils.BadRequestMessage, http.StatusBadRequest, nil)
		return
	}
	product, err := db.ProductByID(productID)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
	}
	if product.UserID != session.User_ID {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	filename, err := utils.PasteFile(request, "product")
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	product.MainPhoto = filename
	err = product.UpdatePhoto()
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
	}
	utils.SendMessage(writer, utils.SuccessfullyUpdatedMessage, http.StatusOK, product)
}

func RecentProductsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	products, err := db.RecentProducts(utils.AmountTopProducts)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.RequestedSuccessfullyMessage, http.StatusOK, products)
}

func ProductsByCategoryHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	category_id, _ := strconv.Atoi(request.URL.Query().Get("category_id"))
	products, err := db.ProductsByCategoryID(category_id, utils.PaginationSize)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.RequestedSuccessfullyMessage, http.StatusOK, products)
}

func AllCategoriesHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	products, err := db.AllCategoriesWithGivenAmountProducts(utils.PaginationSize)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.RequestedSuccessfullyMessage, http.StatusOK, products)
}

func UserProductsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	userID, err := strconv.Atoi(request.URL.Query().Get("user_id"))
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	products, err := db.ProductsByUserID(userID)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.SuccessfullyRequested, http.StatusOK, products)
}
