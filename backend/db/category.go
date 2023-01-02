package db

type Category struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	AmountProducts int    `json:"amount_products"`
}

func CategoryByID(categoryID int) (category Category, err error) {
	err = DB.QueryRow("SELECT * FROM CATEGORIES WHERE ID=$1", categoryID).Scan(
		&category.ID, &category.Name, &category.AmountProducts,
	)
	return
}
