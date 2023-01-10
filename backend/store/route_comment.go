package store

import (
	"encoding/json"
	"github.com/TutorialEdge/realtime-chat-go-react/db/db_store"
	"github.com/TutorialEdge/realtime-chat-go-react/utils"
	"net/http"
	"strconv"
)

func CommentProduct(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	var comment db_store.ProductComment
	if err := json.NewDecoder(request.Body).Decode(&comment); err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	comment.UserID = session.User_ID
	if err = comment.Create(); err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.SuccessfullyCreatedMessage, http.StatusCreated, comment)
}

func RemoveCommentProduct(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	session, err := utils.Session(writer, request)
	if err != nil {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}
	commentID, err := strconv.Atoi(request.URL.Query().Get("comment_id"))
	if err != nil {
		utils.SendMessage(writer, utils.BadRequestMessage, http.StatusBadRequest, nil)
		return
	}
	comment, err := db_store.CommentByID(commentID)
	if err != nil {
		utils.SendMessage(writer, utils.BadRequestMessage, http.StatusBadRequest, nil)
		return
	}
	if comment.UserID != session.User_ID {
		utils.SendMessage(writer, utils.AuthorizationRequestMessage, http.StatusUnauthorized, nil)
		return
	}

	err = comment.Delete()
	if err != nil {
		utils.SendAndPrintErrorMessage(writer, err)
		return
	}
	utils.SendMessage(writer, utils.SuccessfullyUpdatedMessage, http.StatusOK, nil)
}
