package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/agusbasari29/GoAuraBill/internal/model" // Ganti dengan path modul Anda
	"github.com/agusbasari29/GoAuraBill/internal/service" // Ganti dengan path modul Anda
)
type AuthHandler struct {
	authService service.AuthService
}
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}
// RegisterRequest struct untuk binding input registrasi
type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
// LoginRequest struct untuk binding input login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := &model.User{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // Password akan di-hash di BeforeSave hook
		Role:     "customer",   // Default role
	}
	if err := h.authService.RegisterUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.authService.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}