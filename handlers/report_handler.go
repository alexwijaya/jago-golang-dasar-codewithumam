package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"cashier/services"
)

type ReportHandler struct {
	service services.ReportService
}

func NewReportHandler(service services.ReportService) *ReportHandler {
	return &ReportHandler{service}
}

func (h *ReportHandler) GetTodaysSalesReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	summary, err := h.service.GetTodaysSalesReport()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch today's sales report"})
		return
	}
	
	json.NewEncoder(w).Encode(summary)
}

func (h *ReportHandler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	
	if startDateStr == "" || endDateStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "start_date and end_date parameters are required"})
		return
	}
	
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}
	
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		return
	}
	
	endDate = endDate.Add(24 * time.Hour)
	
	if startDate.After(endDate) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "start_date must be before or equal to end_date"})
		return
	}
	
	summary, err := h.service.GetSalesReportByDateRange(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch sales report"})
		return
	}
	
	json.NewEncoder(w).Encode(summary)
}