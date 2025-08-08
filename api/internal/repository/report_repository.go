package repository

import (
	"time"

	"github.com/agusbasari29/GoAuraBill/internal/model"
	"gorm.io/gorm"
)

type ReportRepository interface {
	// GetRevenueReport mengambil data pendapatan berdasarkan rentang tanggal dan periode agregasi.
	GetRevenueReport(startDate, endDate time.Time, period string) ([]model.RevenueReport, error)
	// GetSummaryReport mengambil ringkasan metrik kunci bisnis.
	GetSummaryReport() (*model.SummaryReport, error)
}
type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

// GetRevenueReport: Mengambil total pendapatan yang diselesaikan (completed)
// dari transaksi berjenis 'topup' atau 'payment' dalam rentang tanggal tertentu,
// diagregasi berdasarkan hari atau bulan.
func (r *reportRepository) GetRevenueReport(startDate, endDate time.Time, period string) ([]model.RevenueReport, error) {
	var results []model.RevenueReport

	// Menentukan fungsi SQL untuk memotong tanggal berdasarkan periode
	// DATE_TRUNC('day', processed_at) untuk harian, DATE_TRUNC('month', processed_at) untuk bulanan
	dateTrunc := "day"
	if period == "monthly" {
		dateTrunc = "month"
	}
	err := r.db.Model(&model.Transaction{}).
		Select("DATE_TRUNC(?, processed_at) as date, SUM(amount) as total_revenue", dateTrunc).
		Where("status = ? AND type IN ? AND processed_at BETWEEN ? AND ?",
			model.TransactionStatusCompleted, // Hanya transaksi yang sudah selesai
			[]model.TransactionType{model.TransactionTypePayment, model.TransactionTypeTopUp}, // Hanya jenis transaksi yang menghasilkan pendapatan
					startDate,
					endDate).
		Group("date").       // Agregasi berdasarkan tanggal/bulan yang dipotong
		Order("date ASC").   // Urutkan berdasarkan tanggal/bulan
		Scan(&results).Error // Masukkan hasil ke struct RevenueReport
	return results, err
}

// GetSummaryReport: Mengambil metrik ringkasan seperti total pelanggan, layanan aktif,
// dan pendapatan hari ini/bulan ini.
func (r *reportRepository) GetSummaryReport() (*model.SummaryReport, error) {
	var summary model.SummaryReport
	// Kueri 1: Hitung total pelanggan terdaftar (dari tabel users dengan role 'customer')
	r.db.Model(&model.User{}).Where("role = ?", "customer").Count(&summary.TotalCustomers)
	// Kueri 2: Hitung total layanan aktif (dari tabel customers dengan status 'active')
	r.db.Model(&model.Customer{}).Where("status = ?", "active").Count(&summary.TotalActiveServices)
	// Kueri 3: Hitung pendapatan hari ini (transaksi completed hari ini)
	todayStart := time.Now().Truncate(24 * time.Hour)          // Awal hari ini (00:00:00)
	todayEnd := todayStart.Add(24*time.Hour - time.Nanosecond) // Akhir hari ini (23:59:59.999...)
	r.db.Model(&model.Transaction{}).
		Where("status = ? AND processed_at BETWEEN ? AND ?", model.TransactionStatusCompleted, todayStart, todayEnd).
		Select("COALESCE(SUM(amount), 0)"). // Gunakan COALESCE untuk memastikan 0 jika tidak ada transaksi
		Row().Scan(&summary.TotalRevenueToday)
	// Kueri 4: Hitung pendapatan bulan ini (transaksi completed sejak awal bulan ini)
	monthStart := time.Now().AddDate(0, 0, -time.Now().Day()+1).Truncate(24 * time.Hour) // Awal bulan ini
	r.db.Model(&model.Transaction{}).
		Where("status = ? AND processed_at >= ?", model.TransactionStatusCompleted, monthStart).
		Select("COALESCE(SUM(amount), 0)").
		Row().Scan(&summary.TotalRevenueMonth)
	return &summary, nil
}
