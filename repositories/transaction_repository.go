package repositories

import (
	"context"
	"time"

	"cashier/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionRepository interface {
	CreateTransactionWithStock(transaction model.Transaction, stockUpdates []model.StockUpdate) (model.Transaction, error)
	GetTodaysSalesReport() (model.SalesSummary, error)
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) CreateTransactionWithStock(transaction model.Transaction, stockUpdates []model.StockUpdate) (model.Transaction, error) {
	ctx := context.Background()
	
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return transaction, err
	}
	defer tx.Rollback(ctx)
	
	err = tx.QueryRow(ctx, 
		"INSERT INTO transactions (total_amount, created_at) VALUES ($1, $2) RETURNING id",
		transaction.TotalAmount, transaction.CreatedAt).Scan(&transaction.ID)
	if err != nil {
		return transaction, err
	}
	
	for i, detail := range transaction.Details {
		err = tx.QueryRow(ctx,
			"INSERT INTO transaction_details (transaction_id, product_id, product_name, quantity, subtotal) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			transaction.ID, detail.ProductID, detail.ProductName, detail.Quantity, detail.Subtotal).Scan(&transaction.Details[i].ID)
		if err != nil {
			return transaction, err
		}
		transaction.Details[i].TransactionID = transaction.ID
	}
	
	for _, stockUpdate := range stockUpdates {
		_, err = tx.Exec(ctx,
			"UPDATE products SET stock = $1 WHERE id = $2",
			stockUpdate.NewStock, stockUpdate.ProductID)
		if err != nil {
			return transaction, err
		}
	}
	
	if err = tx.Commit(ctx); err != nil {
		return transaction, err
	}
	
	return transaction, nil
}

func (r *transactionRepository) GetTodaysSalesReport() (model.SalesSummary, error) {
	var summary model.SalesSummary
	
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	err := r.db.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(total_amount), 0) as total_revenue, COUNT(*) as total_transactions 
		 FROM transactions 
		 WHERE created_at >= $1 AND created_at < $2`,
		startOfDay, endOfDay).Scan(&summary.RevenueTotal, &summary.TransactionTotal)
	if err != nil {
		return summary, err
	}
	
	var bestProduct model.BestSellingProduct
	err = r.db.QueryRow(context.Background(),
		`SELECT td.product_name, SUM(td.quantity) as total_qty
		 FROM transaction_details td
		 JOIN transactions t ON td.transaction_id = t.id
		 WHERE t.created_at >= $1 AND t.created_at < $2
		 GROUP BY td.product_name
		 ORDER BY total_qty DESC
		 LIMIT 1`,
		startOfDay, endOfDay).Scan(&bestProduct.Name, &bestProduct.Qty)
	
	if err != nil {
		if err.Error() == "no rows in result set" {
			summary.BestSellingProduct = nil
		} else {
			return summary, err
		}
	} else {
		summary.BestSellingProduct = &bestProduct
	}
	
	return summary, nil
}