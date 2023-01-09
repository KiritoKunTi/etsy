package db_store

import "github.com/TutorialEdge/realtime-chat-go-react/db"

type ProductParameter struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func (parameter *ProductParameter) Create() (err error) {
	stmt, err := db.DB.Prepare("INSERT INTO PRODUCT_PARAMETERS(PRODUCT_ID, KEY, VALUE) VALUES ($1, $2, $3) RETURNING ID")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(parameter.ProductID, parameter.Key, parameter.Value).Scan(&parameter.ID)
	return
}

func ProductParametersByProductID(productID int) (parameters []ProductParameter, err error) {
	rows, err := db.DB.Query("SELECT * FROM PRODUCT_PARAMETERS WHERE PRODUCT_ID=$1", productID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var parameter ProductParameter
		err = rows.Scan(&parameter.ID, &parameter.ProductID, &parameter.Key, &parameter.Value)
		if err != nil {
			return
		}
		parameters = append(parameters, parameter)
	}
	return
}

func (product *Product) DeleteParameters() (err error) {
	stmt, err := db.DB.Prepare("DELETE FROM PRODUCT_PARAMETERS WHERE PRODUCT_ID=$1")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.ID)
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
