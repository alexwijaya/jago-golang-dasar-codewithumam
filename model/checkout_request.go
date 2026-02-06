package model

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}