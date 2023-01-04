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
