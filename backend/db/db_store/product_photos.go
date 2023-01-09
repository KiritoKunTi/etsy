package db_store

import (
	"github.com/TutorialEdge/realtime-chat-go-react/db"
)

type ProductPhoto struct {
	ID        int
	ProductID int
	Photo     string
}

func (photo *ProductPhoto) Create() (err error) {
	stmt, err := db.DB.Prepare("INSERT INTO PRODUCT_PHOTO(PRODUCT_ID, PHOTO) VALUES ($1, $2) RETURNING ID")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(photo.ProductID, photo.Photo).Scan(&photo.ID)
	return
}

func (product *Product) DeletePhotos() (err error) {
	_, err = db.DB.Exec("DELETE FROM PRODUCT_PHOTO WHERE PRODUCT_ID=$1", product.ID)
	return
}

func (product *Product) CreatePhotos() (err error) {
	product.DeletePhotos()
	for _, photo := range product.ProductPhotos {
		photo.ProductID = product.ID
		if err = photo.Create(); err != nil {
			return
		}
	}
	return
}

func (product *Product) GetPhotos() (err error) {
	rows, err := db.DB.Query("SELECT * FROM PRODUCT_PHOTO WHERE PRODUCT_ID=$1", product.ID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var photo ProductPhoto
		err = rows.Scan(&photo.ID, &photo.ProductID, &photo.Photo)
		if err != nil {
			return
		}
		product.ProductPhotos = append(product.ProductPhotos, photo)
	}
	return
}
