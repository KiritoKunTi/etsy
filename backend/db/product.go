package db

import (
	"database/sql"
	"time"
)

type Product struct {
	ID                int                `json:"id"`
	UserID            int                `json:"user_id"`
	User              User               `json:"user"`
	CategoryID        int                `json:"category_id"`
	Category          Category           `json:"category"`
	Name              string             `json:"name"`
	Price             int                `json:"price"`
	Amount            int                `json:"amount"`
	Description       string             `json:"description"`
	AmountLikes       int                `json:"amount_likes"`
	AmountComments    int                `json:"amount_comments"`
	AmountRatings     int                `json:"amount_ratings"`
	Rating            float64            `json:"rating"`
	CreatedAt         string             `json:"created_at"`
	ProductParameters []ProductParameter `json:"product_parameters"`
}

func ProductsByCategoryID(categoryID int) (products []Product, err error) {
	rows, err := DB.Query("SELECT * FROM PRODUCTS	WHERE CATEGORY_ID=$1", categoryID)
	if err != nil {
		return
	}
	defer rows.Close()
	products, err = IterateFromQuery(rows)
	return
}

func ProductsByUserID(userID int) (products []Product, err error) {
	rows, err := DB.Query("SELECT * FROM PRODUCTS	WHERE USER_ID=$1", userID)
	if err != nil {
		return
	}
	defer rows.Close()
	products, err = IterateFromQuery(rows)
	return
}

func RecentProducts(amount int) (products []Product, err error) {
	rows, err := DB.Query("SELECT * FROM PRODUCTS ORDER BY CREATED_AT DESC LIMIT $1", amount)
	if err != nil {
		return
	}
	defer rows.Close()
	products, err = IterateFromQuery(rows)
	return
}

func IterateFromQuery(rows *sql.Rows) (products []Product, err error) {
	for rows.Next() {
		var product Product
		err = rows.Scan(
			&product.ID, &product.UserID, &product.CategoryID, &product.Name, &product.Price, &product.Amount,
			&product.Description, &product.AmountLikes, &product.AmountComments, &product.AmountRatings,
			&product.Rating, &product.CreatedAt,
		)
		if err != nil {
			return
		}
		product.User, err = UserByIDForPublic(product.UserID)
		if err != nil {
			return
		}
		product.Category, err = CategoryByID(product.CategoryID)
		if err != nil {
			return
		}
		products = append(products, product)
	}
	return
}

func ProductByID(productID int) (product Product, err error) {
	err = DB.QueryRow("SELECT * FROM PRODUCTS WHERE ID=$1", productID).Scan(
		&product.ID, &product.UserID, &product.CategoryID, &product.Name, &product.Price, &product.Amount,
		&product.Description, &product.AmountLikes, &product.AmountComments, &product.AmountRatings,
		&product.Rating, &product.CreatedAt,
	)
	if err != nil {
		return
	}
	product.User, err = UserByIDForPublic(product.UserID)
	if err != nil {
		return
	}
	product.Category, err = CategoryByID(product.CategoryID)
	if err != nil {
		return
	}
	product.ProductParameters, err = ProductParametersByProductID(productID)
	return
}

func (product *Product) Create() (err error) {
	st, err := DB.Prepare("INSERT INTO PRODUCTS(USER_ID, CATEGORY_ID, NAME, PRICE, AMOUNT, DESCRIPTION, CREATED_AT) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ID")
	if err != nil {
		return
	}
	defer st.Close()
	err = st.QueryRow(
		product.UserID, product.CategoryID, product.Name, product.Price, product.Amount, product.Description,
		time.Now(),
	).Scan(&product.ID)
	if err != nil {
		return
	}
	err = product.CreateParameters()
	if err != nil {
		product.Delete()
		return
	}
	return
}
func (product *Product) CreateParameters() (err error) {
	for _, parameter := range product.ProductParameters {
		parameter.ProductID = product.ID
		err = parameter.Create()
		return
	}
	return
}

func (product *Product) Delete() (err error) {
	err = product.DeleteParameters()
	if err != nil {
		return
	}
	stmt, err := DB.Prepare("DELETE FROM PRODUCTS WHERE ID=$1")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.ID)
	return
}

func (product *Product) DeleteParameters() (err error) {
	stmt, err := DB.Prepare("DELETE FROM PRODUCT_PARAMETERS WHERE PRODUCT_ID=$1")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.ID)
	return
}
