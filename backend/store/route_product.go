package store

import (
	"encoding/json"
	"github.com/TutorialEdge/realtime-chat-go-react/db/db_store"
	"github.com/TutorialEdge/realtime-chat-go-react/utils"
	"net/http"
	"strconv"
)

const (
	productPhotoDir = "product"
	productPhotoKey = "photo"
)

func CreateProductHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	var product db_store.Product
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
	var product db_store.Product
	if err := json.NewDecoder(request.Body).Decode(&product); err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	if prod, err := db_store.ProductByIDDetail(product.ID); err != nil || prod.UserID != session.User_ID {
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
	request.ParseMultipartForm(10 << 20)
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	productID, err := strconv.Atoi(request.URL.Query().Get("product_id"))
	if err != nil {
		utils.SendMessage(writer, utils.BadRequestMessage, http.StatusBadRequest, nil)
		return
	}
	product, err := db_store.ProductByIDDetail(productID)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	if product.UserID != session.User_ID {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	filename, err := utils.PasteFile(request, productPhotoDir, productPhotoKey, product.ID)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	err = product.UpdatePhoto(filename)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.SuccessfullyUpdatedMessage, http.StatusOK, product)
}

func RecentProductsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	products, err := db_store.RecentProducts(utils.AmountTopProducts)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.RequestedSuccessfullyMessage, http.StatusOK, products)
}

func UpgradeProductPhotosHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	productID, err := strconv.Atoi(request.URL.Query().Get("product_id"))
	if err != nil {
		utils.SendMessage(writer, utils.BadRequestMessage, http.StatusBadRequest, nil)
		return
	}
	product, err := db_store.ProductByIDDetail(productID)
	if err != nil {
		utils.SendMessage(writer, utils.BadRequestMessage, http.StatusBadRequest, nil)
		return
	}
	if session, err := utils.Session(writer, request); session.User_ID != product.UserID || err != nil {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	product.ProductPhotos, _ = utils.PasteProductPhoto(request, product)
	product.CreatePhotos()
	utils.SendMessage(writer, utils.SuccessfullyCreatedMessage, http.StatusOK, product)
}

func ProductHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	productID, err := strconv.Atoi(request.URL.Query().Get("product_id"))
	if err != nil {
		utils.SendMessage(writer, utils.BadRequestMessage, http.StatusBadRequest, nil)
		return
	}
	product, err := db_store.ProductByIDDetail(productID)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.SuccessfullyRequested, http.StatusOK, product)
}

func ProductDeactivateHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	productID, err := strconv.Atoi(request.URL.Query().Get("product_id"))
	if err != nil {
		utils.SendMessage(writer, utils.BadRequestMessage, http.StatusBadRequest, nil)
		return
	}
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	product, err := db_store.ProductByID(productID)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	if product.UserID != session.User_ID {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	product.Deactivate()
	utils.SendMessage(writer, utils.SuccessfullyUpdatedMessage, http.StatusOK, nil)
}
