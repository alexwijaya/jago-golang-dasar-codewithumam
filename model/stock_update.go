package model

type StockUpdate struct {
	ProductID int `json:"product_id"`
	NewStock  int `json:"new_stock"`
}