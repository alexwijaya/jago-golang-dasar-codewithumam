package services

import (
	"cashier/model"
	"cashier/repositories"
)

type ReportService interface {
	GetTodaysSalesSummary() (model.SalesSummary, error)
}

type reportService struct {
	transactionRepo repositories.TransactionRepository
}

func NewReportService(transactionRepo repositories.TransactionRepository) ReportService {
	return &reportService{
		transactionRepo: transactionRepo,
	}
}

func (s *reportService) GetTodaysSalesSummary() (model.SalesSummary, error) {
	return s.transactionRepo.GetTodaysSalesReport()
}