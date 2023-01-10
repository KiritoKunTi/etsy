package db_store

import (
	"database/sql"
	"github.com/TutorialEdge/realtime-chat-go-react/db"
	"os"
	"time"
)

type Product struct {
	ID                int                `json:"id,omitempty"`
	UserID            int                `json:"user_id"`
	User              db.User            `json:"user,omitempty"`
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
	ProductPhotos     []ProductPhoto     `json:"photos,omitempty"`
	IsActive          bool               `json:"-"`
	Comments          []ProductComment   `json:"comments,omitempty"`
}

func ProductByID(productID int) (product Product, err error) {
	err = db.DB.QueryRow("SELECT * FROM PRODUCTS WHERE ID=$1 AND IS_ACTIVE=TRUE", productID).Scan(
		&product.ID, &product.UserID, &product.CategoryID, &product.Name, &product.MainPhoto,
		&product.Price, &product.Amount, &product.Description, &product.AmountLikes, &product.AmountComments,
		&product.AmountRatings,
		&product.Rating, &product.IsActive, &product.CreatedAt,
	)
	return
}

func ProductByIDDetail(productID int) (product Product, err error) {
	product, err = ProductByID(productID)
	if err != nil {
		return
	}
	if product.User, err = db.UserByIDForPublic(product.UserID); err != nil {
		return
	}
	if product.Category, err = CategoryByID(product.CategoryID); err != nil {
		return
	}
	if product.ProductParameters, err = ProductParametersByProductID(productID); err != nil {
		return
	}
	if err = product.GetPhotos(); err != nil {
		return
	}
	if err = product.GetPhotos(); err != nil {
		return
	}
	if product.Comments, err = CommentsByProductID(product.ID); err != nil {
		return
	}
	return
}

func (product *Product) Update() (err error) {
	stmt, err := db.DB.Prepare("UPDATE PRODUCTS SET NAME=$1, PRICE=$2, AMOUNT=$3, DESCRIPTION=$4 WHERE ID=$5")
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
	st, err := db.DB.Prepare("INSERT INTO PRODUCTS(USER_ID, CATEGORY_ID, NAME, PRICE, AMOUNT, DESCRIPTION, CREATED_AT) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ID, CREATED_AT")
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
	stmt, err := db.DB.Prepare("DELETE FROM PRODUCTS WHERE ID=$1")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.ID)
	return
}

func (product *Product) UpdatePhoto(mainPhoto string) (err error) {
	os.Remove(product.MainPhoto)
	product.MainPhoto = mainPhoto
	_, err = db.DB.Exec("UPDATE PRODUCTS SET PHOTO=$1 WHERE ID=$2", product.MainPhoto, product.ID)
	return
}

func (product *Product) SetUser(userID int) {
	product.UserID = userID
	product.User, _ = db.UserByID(userID)
	product.User.HideInfo()
}

func ProductsByUserID(userID int) (products []Product, err error) {
	rows, err := db.DB.Query("SELECT * FROM PRODUCTS	WHERE USER_ID=$1 AND IS_ACTIVE", userID)
	if err != nil {
		return
	}
	defer rows.Close()
	products, err = QueryToSliceProducts(rows)
	return
}

func RecentProducts(amount int) (products []Product, err error) {
	rows, err := db.DB.Query("SELECT * FROM PRODUCTS WHERE IS_ACTIVE=TRUE ORDER BY CREATED_AT DESC LIMIT $1", amount)
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
			&product.Rating, &product.IsActive, &product.CreatedAt,
		)
		if err != nil {
			return
		}
		products = append(products, product)
	}
	return
}

func (product *Product) Deactivate() {
	product.IsActive = false
	db.DB.Exec("UPDATE PRODUCTS SET IS_ACTIVE=FALSE WHERE ID=$1", product.ID)
}
