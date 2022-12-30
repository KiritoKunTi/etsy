package db

type ProductParameter struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func (parameter *ProductParameter) Create(product_id int) (err error) {
	stmt, err := DB.Prepare("INSERT INTO PRODUCT_PARAMETERS(PRODUCT_ID, KEY, VALUE) VALUES ($1, $2, $3) RETURNING ID")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(product_id, parameter.Key, parameter.Value).Scan(&parameter.ID)
	return
}
