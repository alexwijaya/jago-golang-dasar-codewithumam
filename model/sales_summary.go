package model

type SalesSummary struct {
	RevenueTotal       int                 `json:"total_revenue"`
	TransactionTotal   int                 `json:"total_transaksi"`
	BestSellingProduct *BestSellingProduct `json:"produk_terlaris"`
}

type BestSellingProduct struct {
	Name string `json:"nama"`
	Qty  int    `json:"qty_terjual"`
}