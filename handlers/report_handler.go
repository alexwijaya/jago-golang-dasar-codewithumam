package handlers

import (
	"encoding/json"
	"net/http"

	"cashier/services"
)

type ReportHandler struct {
	service services.ReportService
}

func NewReportHandler(service services.ReportService) *ReportHandler {
	return &ReportHandler{service}
}

func (h *ReportHandler) GetTodaysSales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	summary, err := h.service.GetTodaysSalesSummary()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch today's sales report"})
		return
	}
	
	json.NewEncoder(w).Encode(summary)
}