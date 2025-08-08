package handler

import (
	"net/http"
	"strconv"

	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/service"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service}
}

type CreateTransactionRequest struct {
	CustomerID  uint       `json:"customer_id" binding:"required"`
	Amount      float64    `json:"amount" binding:"required,gt=0"`
	Type        string     `json:"type" binding:"required,oneof=topup payment refund adjustment voucher"`
	Description string     `json:"description"`
	ReferenceID string     `json:"reference_id"`
	Metadata    model.JSON `json:"metadata"`
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	txn := model.Transaction{
		CustomerID:  req.CustomerID,
		Amount:      req.Amount,
		Type:        model.TransactionType(req.Type),
		Description: req.Description,
		ReferenceID: req.ReferenceID,
		Metadata:    req.Metadata,
		Status:      model.TransactionStatusPending,
	}
	if err := h.service.CreateTransaction(&txn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, txn)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	txn, err := h.service.GetTransaction(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	c.JSON(http.StatusOK, txn)
}

func (h *TransactionHandler) GetCustomerTransactions(c *gin.Context) {
	customerID, err := strconv.Atoi(c.Param("customer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}
	txns, err := h.service.GetCustomerTransactions(uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, txns)
}

func (h *TransactionHandler) ProcessPayment(c *gin.Context) {
	txnID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}
	var req struct {
		ReferenceID string `json:"reference_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.ProcessPayment(uint(txnID), req.ReferenceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment processed successfully"})
}

func (h *TransactionHandler) CancelTransaction(c *gin.Context) {
	txnID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}
	if err := h.service.CancelTransaction(uint(txnID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction cancelled"})
}
