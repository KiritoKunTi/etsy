package db_store

import "github.com/TutorialEdge/realtime-chat-go-react/db"

type Category struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	AmountProducts int       `json:"amount_products"`
	Products       []Product `json:"products,omitempty"`
}

func CategoryByID(categoryID int) (category Category, err error) {
	err = db.DB.QueryRow("SELECT * FROM CATEGORIES WHERE ID=$1", categoryID).Scan(
		&category.ID, &category.Name, &category.AmountProducts,
	)
	return
}

func AllCategories() (categories []Category, err error) {
	rows, err := db.DB.Query("SELECT * FROM CATEGORIES")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.Name, &category.AmountProducts)
		if err != nil {
			return
		}
		categories = append(categories, category)
	}
	return
}

func AllCategoriesWithProducts() (categories []Category, err error) {
	categories, err = AllCategories()
	if err != nil {
		return
	}
	for ind, category := range categories {
		category.Products, err = ProductsByCategoryID(category.ID)
		if err != nil {
			return
		}
		categories[ind] = category
	}
	return
}

func ProductsByCategoryID(categoryID int) (products []Product, err error) {
	rows, err := db.DB.Query(
		"SELECT * FROM PRODUCTS	WHERE CATEGORY_ID=$1 AND IS_ACTIVE=TRUE", categoryID,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	products, err = QueryToSliceProducts(rows)
	return
}
