package store

import (
	"github.com/TutorialEdge/realtime-chat-go-react/db/db_store"
	"github.com/TutorialEdge/realtime-chat-go-react/utils"
	"net/http"
	"strconv"
)

func UserProductsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	userID, err := strconv.Atoi(request.URL.Query().Get("user_id"))
	if err != nil {
		utils.SendMessage(writer, utils.BadRequestMessage, http.StatusBadRequest, nil)
		return
	}
	products, err := db_store.ProductsByUserID(userID)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.SuccessfullyRequested, http.StatusOK, products)
}
