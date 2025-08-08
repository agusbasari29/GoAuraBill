package model

import "time"

// RevenueReport merepresentasikan data pendapatan dalam periode tertentu.
type RevenueReport struct {
	Date         time.Time `json:"date"`
	TotalRevenue float64   `json:"total_revenue"`
}

// SummaryReport memberikan gambaran umum tentang data kunci.
type SummaryReport struct {
	TotalCustomers      int64   `json:"total_customers"`
	TotalActiveServices int64   `json:"total_active_services"`
	TotalRevenueToday   float64 `json:"total_revenue_today"`
	TotalRevenueMonth   float64 `json:"total_revenue_month"`
}