package store

import (
	"github.com/TutorialEdge/realtime-chat-go-react/db/db_store"
	"github.com/TutorialEdge/realtime-chat-go-react/utils"
	"net/http"
	"strconv"
)

func ProductsByCategoryHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	category_id, _ := strconv.Atoi(request.URL.Query().Get("category_id"))
	products, err := db_store.ProductsByCategoryID(category_id)
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.RequestedSuccessfullyMessage, http.StatusOK, products)
}

func AllCategoriesHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	categories, err := db_store.AllCategories()
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.RequestedSuccessfullyMessage, http.StatusOK, categories)
}
