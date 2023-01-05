package db

type ProductPhoto struct {
	ID        int
	ProductID int
	Photo     string
}

func (photo *ProductPhoto) Create() (err error) {
	stmt, err := DB.Prepare("INSERT INTO PRODUCT_PHOTOS(PRODUCT_ID, PHOTO) VALUES ($1, $2) RETURNING ID")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(photo.ProductID, photo.Photo).Scan(&photo.ID)
	return
}

func (photo *ProductPhoto) Delete() (err error) {
	_, err = DB.Exec("DELETE FROM PRODUCT_PHOTO WHERE ID=$1", photo.ID)
	return
}

func (product *Product) CreatePhotos() (err error) {
	for _, photo := range product.ProductPhotos {
		photo.ProductID = product.ID
		if err = photo.Create(); err != nil {
			return
		}
	}
	return
}
