package handler

import (
	"net/http"
	"strconv"

	"github.com/agusbasari29/GoAuraBill/internal/service"
	"github.com/gin-gonic/gin"
)

type VoucherHandler struct {
	service service.VoucherService
}

func NewVoucherHandler(service service.VoucherService) *VoucherHandler {
	return &VoucherHandler{service: service}
}

type GenerateRequest struct {
	Quantity  int  `json:"quantity" binding:"required,min=1"`
	ProfileID uint `json:"profile_id" binding:"required"`
}

func (h *VoucherHandler) Generate(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vouchers, err := h.service.GenerateVouchers(req.Quantity, req.ProfileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate vouchers: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Vouchers generated successfully",
		"count":   len(vouchers),
	})
}
func (h *VoucherHandler) GetAll(c *gin.Context) {
	vouchers, err := h.service.GetAllVouchers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vouchers)
}
func (h *VoucherHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	voucher, err := h.service.GetVoucherByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voucher not found"})
		return
	}
	c.JSON(http.StatusOK, voucher)
}
func (h *VoucherHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.service.DeleteVoucher(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Voucher deleted successfully"})
}
