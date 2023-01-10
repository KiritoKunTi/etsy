package db_store

import (
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"time"
)

type ProductLike struct {
	ID        int    `json:"-"`
	ProductID int    `json:"product_id"`
	UserID    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func (like *ProductLike) Create() (err error) {
	if !like.isUserLiked() {
		st, err := db.DB.Prepare("INSERT INTO PRODUCT_LIKES(PRODUCT_ID, USER_ID, CREATED_AT) VALUES ($1, $2, $3) RETURNING ID, CREATED_AT")
		if err != nil {
			return err
		}
		defer st.Close()
		err = st.QueryRow(like.ProductID, like.UserID, time.Now()).Scan(&like.ID, &like.CreatedAt)
	}
	return
}

func (like *ProductLike) Delete() (err error) {
	_, err = db.DB.Exec("DELETE FROM PRODUCT_LIKES WHERE PRODUCT_ID=$1 AND USER_ID=$2", like.ProductID, like.UserID)
	return
}

func (like *ProductLike) isUserLiked() bool {
	var isLiked bool
	db.DB.QueryRow(
		"SELECT EXISTS (SELECT * FROM PRODUCT_LIKES WHERE USER_ID=$1 AND PRODUCT_ID=$2)", like.UserID, like.ProductID,
	).Scan(&isLiked)
	return isLiked
}
