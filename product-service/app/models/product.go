package models

// Product is the model for Users on DB
type Product struct {
	ID          int64   `json:"product_id"`
	Name        string  `json:"product_name"`
	Description string  `json:"product_description"`
	Prise       float64 `json:"product_prise"`
}
