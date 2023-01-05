package db

type Category struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	AmountProducts int       `json:"amount_products"`
	Products       []Product `json:"products,omitempty"`
}

func CategoryByID(categoryID int) (category Category, err error) {
	err = DB.QueryRow("SELECT * FROM CATEGORIES WHERE ID=$1", categoryID).Scan(
		&category.ID, &category.Name, &category.AmountProducts,
	)
	return
}

func AllCategories() (categories []Category, err error) {
	rows, err := DB.Query("SELECT * FROM CATEGORIES")
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

func AllCategoriesWithGivenAmountProducts(amo int) (categories []Category, err error) {
	categories, err = AllCategories()
	if err != nil {
		return
	}
	for ind, category := range categories {
		category.Products, err = ProductsByCategoryID(category.ID, amo)
		if err != nil {
			return
		}
		categories[ind] = category
	}
	return
}
