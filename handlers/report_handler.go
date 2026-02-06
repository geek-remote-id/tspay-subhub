package handlers

import (
	"net/http"

	"kasir-api/services"
	"kasir-api/utils"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetDailySalesReport godoc
// @Summary      Get today's sales report
// @Description  Get sales report for today including total revenue, transaction count, and top-selling product
// @Tags         report
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /report/hari-ini [get]
func (h *ReportHandler) GetDailySalesReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetDailySalesReport()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Response{
			Status:  "failed",
			Message: "Failed to fetch daily sales report: " + err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Status:  "success",
		Message: "Daily sales report retrieved successfully",
		Data:    report,
	})
}

// GetSalesReportByDateRange godoc
// @Summary      Get sales report by date range
// @Description  Get sales report for a specific date range including total revenue, transaction count, and top-selling product
// @Tags         report
// @Accept       json
// @Produce      json
// @Param        start_date  query     string  true  "Start date (YYYY-MM-DD)"
// @Param        end_date    query     string  true  "End date (YYYY-MM-DD)"
// @Success      200         {object}  utils.Response
// @Failure      400         {object}  utils.Response
// @Failure      500         {object}  utils.Response
// @Router       /report [get]
func (h *ReportHandler) GetSalesReportByDateRange(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{
			Status:  "failed",
			Message: "start_date and end_date query parameters are required",
		})
		return
	}

	// Add time component to dates
	startDateTime := startDate + " 00:00:00"
	endDateTime := endDate + " 23:59:59"

	report, err := h.service.GetSalesReportByDateRange(startDateTime, endDateTime)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Response{
			Status:  "failed",
			Message: "Failed to fetch sales report: " + err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Status:  "success",
		Message: "Sales report retrieved successfully",
		Data:    report,
	})
}
