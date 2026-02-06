package services

import (
	"errors"
	"fmt"
	"time"

	"cashier/model"
	"cashier/repositories"
)

type TransactionService interface {
	ProcessCheckout(request model.CheckoutRequest) (model.Transaction, error)
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
	productRepo     repositories.ProductRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository, productRepo repositories.ProductRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
	}
}

func (s *transactionService) ProcessCheckout(request model.CheckoutRequest) (model.Transaction, error) {
	var transaction model.Transaction

	if len(request.Items) == 0 {
		return transaction, errors.New("Checkout request must contain at least one item")
	}
	
	totalAmount := 0
	var details []model.TransactionDetail
	var stockUpdates []model.StockUpdate

	for _, item := range request.Items {
		if item.Quantity <= 0 {
			return transaction, fmt.Errorf("Quantity must be greater than 0 for product ID %d", item.ProductID)
		}

		product, err := s.productRepo.FindRawByID(item.ProductID)
		if err != nil {
			return transaction, fmt.Errorf("Product with ID %d not found", item.ProductID)
		}
		
		if product.Stock < item.Quantity {
			return transaction, fmt.Errorf("Insufficient stock for product '%s'. Available: %d, Requested: %d", 
				product.Name, product.Stock, item.Quantity)
		}
		
		subtotal := product.Price * item.Quantity
		totalAmount += subtotal
		
		detail := model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		}
		details = append(details, detail)
		
		newStock := product.Stock - item.Quantity
		stockUpdate := model.StockUpdate{
			ProductID: item.ProductID,
			NewStock:  newStock,
		}
		stockUpdates = append(stockUpdates, stockUpdate)
	}
	
	transaction = model.Transaction{
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		Details:     details,
	}
	
	createdTransaction, err := s.transactionRepo.CreateTransactionWithStock(transaction, stockUpdates)
	if err != nil {
		return transaction, fmt.Errorf("Failed to process checkout: %v", err)
	}
	
	return createdTransaction, nil
}