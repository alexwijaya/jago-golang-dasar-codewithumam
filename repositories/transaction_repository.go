package repositories

import (
	"context"

	"cashier/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionRepository interface {
	CreateTransactionWithStock(transaction model.Transaction, stockUpdates []model.StockUpdate) (model.Transaction, error)
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