package db

import (
	"database/sql"
	"os"
	"time"
)

type Product struct {
	ID                int                `json:"id,omitempty"`
	UserID            int                `json:"user_id"`
	User              User               `json:"user,omitempty"`
	CategoryID        int                `json:"category_id"`
	Category          Category           `json:"category,omitempty"`
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
	MainPhoto         string             `json:"photo"`
	ProductPhotos     []ProductPhoto     `json:"photos"`
}

func (product *Product) Update() (err error) {
	stmt, err := DB.Prepare("UPDATE PRODUCTS SET NAME=$1, PRICE=$2, AMOUNT=$3, DESCRIPTION=$4 WHERE ID=$5")
	if prod, err := ProductByID(product.ID); os.Remove(prod.MainPhoto) != nil || err != nil {
		return err
	}
	if err != nil {
		return
	}
	_, err = stmt.Exec(product.Name, product.Price, product.Amount, product.Description, product.ID)
	if err != nil {
		return
	}
	err = product.DeleteParameters()
	if err != nil {
		return
	}
	err = product.CreateParameters()
	if err != nil {
		return
	}
	return
}

func (product *Product) Create() (err error) {
	st, err := DB.Prepare("INSERT INTO PRODUCTS(USER_ID, CATEGORY_ID, NAME, PRICE, AMOUNT, DESCRIPTION, CREATED_AT) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ID, CREATED_AT")
	if err != nil {
		return
	}
	defer st.Close()
	err = st.QueryRow(
		product.UserID, product.CategoryID, product.Name, product.Price, product.Amount, product.Description,
		time.Now(),
	).Scan(&product.ID, &product.CreatedAt)
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

func (product *Product) UpdatePhoto() (err error) {
	_, err = DB.Exec("UPDATE PRODUCTS SET PHOTO=$1 WHERE ID=$2", product.MainPhoto, product.ID)
	return
}

func (product *Product) SetUser(userID int) {
	product.UserID = userID
	product.User, _ = UserByID(userID)
	product.User.HideInfo()
}

func ProductByID(productID int) (product Product, err error) {
	err = DB.QueryRow("SELECT * FROM PRODUCTS WHERE ID=$1", productID).Scan(
		&product.ID, &product.UserID, &product.CategoryID, &product.Name, &product.MainPhoto,
		&product.Price, &product.Amount, &product.Description, &product.AmountLikes, &product.AmountComments,
		&product.AmountRatings,
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

func ProductsByCategoryID(categoryID int, amo int) (products []Product, err error) {
	rows, err := DB.Query("SELECT * FROM PRODUCTS	WHERE CATEGORY_ID=$1 LIMIT $2", categoryID, amo)
	if err != nil {
		return
	}
	defer rows.Close()
	products, err = QueryToSliceProducts(rows)
	return
}

func ProductsByUserID(userID int) (products []Product, err error) {
	rows, err := DB.Query("SELECT * FROM PRODUCTS	WHERE USER_ID=$1", userID)
	if err != nil {
		return
	}
	defer rows.Close()
	products, err = QueryToSliceProducts(rows)
	return
}

func RecentProducts(amount int) (products []Product, err error) {
	rows, err := DB.Query("SELECT * FROM PRODUCTS ORDER BY CREATED_AT DESC LIMIT $1", amount)
	if err != nil {
		return
	}
	defer rows.Close()
	products, err = QueryToSliceProducts(rows)
	return
}

func QueryToSliceProducts(rows *sql.Rows) (products []Product, err error) {
	for rows.Next() {
		var product Product
		err = rows.Scan(
			&product.ID, &product.UserID, &product.CategoryID, &product.Name, &product.MainPhoto, &product.Price,
			&product.Amount,
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
