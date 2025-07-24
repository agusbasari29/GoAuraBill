package handler

import (
	"net/http"
	"strconv"

	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/service"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service service.CustomerService
}

func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

type CreateCustomerRequest struct {
	FullName  string `json:"full_name" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	ProfileID *uint  `json:"profile_id"` // ID paket layanan
}

type UpdateCustomerRequest struct {
	FullName  string `json:"full_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password,omitempty"` // Opsional saat update
	ProfileID *uint  `json:"profile_id"`
	IsActive  bool   `json:"is_active"`
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := model.User{
		FullName:  req.FullName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		ProfileID: req.ProfileID,
	}

	if err := h.service.CreateCustomer(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h *CustomerHandler) GetAll(c *gin.Context) {
	customers, err := h.service.GetAllCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	customer, err := h.service.GetCustomerByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil data customer yang ada
	customer, err := h.service.GetCustomerByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// Update data
	customer.FullName = req.FullName
	customer.Email = req.Email
	customer.ProfileID = req.ProfileID
	customer.IsActive = req.IsActive
	if req.Password != "" {
		customer.Password = req.Password // Password akan di-hash oleh hook
	}

	if err := h.service.UpdateCustomer(customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteCustomer(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}