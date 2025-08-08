package service

import (
	"time"

	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/repository"
)

type ReportService interface {
	GetRevenueReport(period string) ([]model.RevenueReport, error)
	GetSummaryReport() (*model.SummaryReport, error)
}
type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

// GetRevenueReport: Menentukan rentang tanggal berdasarkan parameter 'period'.
// 'daily' akan mengambil data 30 hari terakhir.
// 'monthly' akan mengambil data 12 bulan terakhir.
func (s *reportService) GetRevenueReport(period string) ([]model.RevenueReport, error) {
	now := time.Now()
	var startDate, endDate time.Time
	if period == "monthly" {
		// Laporan 12 bulan terakhir
		startDate = now.AddDate(-1, 0, 0) // 1 tahun yang lalu
		endDate = now
	} else {
		// Laporan 30 hari terakhir (default jika period bukan 'monthly')
		period = "daily"                   // Pastikan period diset ke 'daily'
		startDate = now.AddDate(0, 0, -30) // 30 hari yang lalu
		endDate = now
	}
	return s.repo.GetRevenueReport(startDate, endDate, period)
}
func (s *reportService) GetSummaryReport() (*model.SummaryReport, error) {
	return s.repo.GetSummaryReport()
}
