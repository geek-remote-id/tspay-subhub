package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetDailySalesReport() (*models.SalesReport, error) {
	return s.repo.GetDailySalesReport()
}

func (s *ReportService) GetSalesReportByDateRange(startDate, endDate string) (*models.SalesReport, error) {
	return s.repo.GetSalesReportByDateRange(startDate, endDate)
}
