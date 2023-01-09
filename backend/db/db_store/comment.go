package db_store

import (
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"time"
)

type ProductComment struct {
	ID        int    `json:"-"`
	ProductID int    `json:"product_id"`
	UserID    int    `json:"user_id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

func (product *Product) Comments() (comments []ProductComment) {
	rows, err := db.DB.Query()
}

func (comment *ProductComment) Create() (err error) {
	st, err := db.DB.Prepare("INSERT INTO PRODUCT_COMMENTS(PRODUCT_ID, USER_ID, TEXT, CREATED_AT) VALUES ($1, $2, $3, $4) RETURNING ID, CREATED_AT")
	if err != nil {
		return
	}
	defer st.Close()
	err = st.QueryRow(comment.ProductID, comment.UserID, comment.Text, time.Now()).Scan(&comment.ID, &comment.CreatedAt)
	return
}

func (comment *ProductComment) Delete() (err error) {
	_, err = db.DB.Exec("DELETE FROM PRODUCT_COMMENTS WHERE ID=$1", comment.ID)
	return
}
