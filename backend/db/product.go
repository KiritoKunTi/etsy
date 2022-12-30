package db

type Product struct {
	ID                int                `json:"id"`
	UserID            int                `json:"user_id"`
	Name              string             `json:"name"`
	Price             int                `json:"price"`
	Amount            int                `json:"amount"`
	Description       string             `json:"description"`
	ProductParameters []ProductParameter `json:"product_parameters"`
}

func (product *Product) Create() (err error) {
	st, err := DB.Prepare("INSERT INTO PRODUCTS(USER_ID, NAME, PRICE, AMOUNT, DESCRIPTION) VALUES ($1, $2, $3, $4, $5) RETURNING ID")
	if err != nil {
		return
	}
	defer st.Close()
	err = st.QueryRow(
		product.UserID, product.Name, product.Price, product.Amount, product.Description,
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
		err = parameter.Create(product.ID)
		return
	}
	return
}

func (product *Product) Delete() (err error) {
	stmt, err := DB.Prepare("DELETE FROM PRODUCT_PARAMETERS WHERE PRODUCT_ID=$1")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.ID)
	stmt, err = DB.Prepare("DELETE FROM PRODUCTS WHERE ID=$1")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.ID)
	return
}
