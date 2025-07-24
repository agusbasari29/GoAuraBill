package handler

import (
	"net/http"
	"strconv"

	"github.com/agusbasari29/GoAuraBill/internal/model"
	"github.com/agusbasari29/GoAuraBill/internal/service"
	"github.com/gin-gonic/gin"
)

type ServiceProfileHandler struct {
	service service.ServiceProfileService
}

func NewServiceProfileHandler(service service.ServiceProfileService) *ServiceProfileHandler {
	return &ServiceProfileHandler{service: service}
}

type CreateProfileRequest struct {
	Name         string  `json:"name" binding:"required"`
	DownloadRate uint    `json:"download_rate" binding:"required"`
	UploadRate   uint    `json:"upload_rate" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
	ValidityDays int     `json:"validity_days" binding:"required"`
}

func (h *ServiceProfileHandler) Create(c *gin.Context) {
	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile := model.ServiceProfile{
		Name:         req.Name,
		DownloadRate: req.DownloadRate,
		UploadRate:   req.UploadRate,
		Price:        req.Price,
		ValidityDays: req.ValidityDays,
	}

	if err := h.service.CreateProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

func (h *ServiceProfileHandler) GetAll(c *gin.Context) {
	profiles, err := h.service.GetAllProfiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

func (h *ServiceProfileHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	profile, err := h.service.GetProfileByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *ServiceProfileHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var profile model.ServiceProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	profile.ID = uint(id)

	if err := h.service.UpdateProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *ServiceProfileHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteProfile(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}