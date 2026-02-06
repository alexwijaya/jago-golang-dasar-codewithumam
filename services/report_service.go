package services

import (
	"time"

	"cashier/model"
	"cashier/repositories"
)

type ReportService interface {
	GetTodaysSalesReport() (model.SalesSummary, error)
	GetSalesReportByDateRange(startDate, endDate time.Time) (model.SalesSummary, error)
}

type reportService struct {
	transactionRepo repositories.TransactionRepository
}

func NewReportService(transactionRepo repositories.TransactionRepository) ReportService {
	return &reportService{
		transactionRepo: transactionRepo,
	}
}

func (s *reportService) GetTodaysSalesReport() (model.SalesSummary, error) {
	return s.transactionRepo.GetTodaysSalesReport()
}

func (s *reportService) GetSalesReportByDateRange(startDate, endDate time.Time) (model.SalesSummary, error) {
	return s.transactionRepo.GetSalesReportByDateRange(startDate, endDate)
}