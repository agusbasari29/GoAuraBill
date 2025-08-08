package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/agusbasari29/GoAuraBill/config"
	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/repository"
	"github.com/agusbasari29/GoAuraBill/internal/util"
)

type PaymentService interface {
	CreateTripayCharge(transactionID uint, method string) (map[string]interface{}, error)
	HandleTripayCallback(payload map[string]interface{}, signature string) error
}

type paymentService struct {
	cfg          config.Config
	transRepo    repository.TransactionRepository
	customerRepo repository.CustomerRepository
	httpClient   *http.Client
}

func NewPaymentService(cfg config.Config, transRepo repository.TransactionRepository, customerRepo repository.CustomerRepository) PaymentService {
	return &paymentService{
		cfg:          cfg,
		transRepo:    transRepo,
		customerRepo: customerRepo,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *paymentService) CreateTripayCharge(transactionID uint, method string) (map[string]interface{}, error) {
	transaction, err := s.transRepo.GetByID(transactionID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil transaksi: %w", err)
	}
	customer, err := s.customerRepo.GetCustomerByID(transaction.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data pelanggan: %w", err)
	}

	merchantRef := fmt.Sprintf("INV-%d-%d", transaction.ID, time.Now().Unix())
	transaction.ReferenceID = merchantRef // Update reference ID
	s.transRepo.Update(transaction)

	signature := util.GenerateSignature(s.cfg.TripayMerchantCode, merchantRef, int(transaction.Amount), s.cfg.TripayPrivateKey)

	payload := map[string]interface{}{
		"method":         method,
		"merchant_ref":   merchantRef,
		"amount":         transaction.Amount,
		"customer_name":  customer.User.FullName,
		"customer_email": customer.User.Email,
		"customer_phone": customer.Phone,
		"order_items": []map[string]interface{}{
			{
				"sku":      fmt.Sprintf("PROFILE-%d", customer.ProfileID),
				"name":     transaction.Description,
				"price":    transaction.Amount,
				"quantity": 1,
			},
		},
		"callback_url": "https://domain-anda.com/api/payments/tripay-callback", // Ganti dengan URL callback Anda
		"return_url":   "https://domain-anda.com/payment-success",
		"expired_time": time.Now().Add(24 * time.Hour).Unix(),
		"signature":    signature,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", s.cfg.TripayApiUrl+"/transaction/create", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+s.cfg.TripayApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal mengirim permintaan ke Tripay: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tripayResponse map[string]interface{}
	json.Unmarshal(body, &tripayResponse)

	return tripayResponse, nil
}

func (s *paymentService) HandleTripayCallback(payload map[string]interface{}, signature string) error {
	// 1. Validasi signature callback
	jsonBody, _ := json.Marshal(payload)
	expectedSignature := hmac.New(sha256.New, []byte(s.cfg.TripayPrivateKey))
	expectedSignature.Write(jsonBody)
	if signature != hex.EncodeToString(expectedSignature.Sum(nil)) {
		return errors.New("invalid callback signature")
	}

	// 2. Proses payload
	merchantRef, _ := payload["merchant_ref"].(string)
	status, _ := payload["status"].(string)

	transaction, err := s.transRepo.GetByReference(merchantRef)
	if err != nil {
		return fmt.Errorf("transaksi tidak ditemukan: %w", err)
	}

	switch status {
	case "PAID":
		transaction.Status = model.TransactionStatusCompleted
		now := time.Now()
		transaction.ProcessedAt = &now
	case "EXPIRED", "FAILED":
		transaction.Status = model.TransactionStatusFailed
	}

	return s.transRepo.Update(transaction)
}
