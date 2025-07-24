package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/agusbasari29/GoAuraBill/internal/model"
    "github.com/agusbasari29/GoAuraBill/internal/service"
)

type RouterHandler struct {
    service service.RouterService
}

func NewRouterHandler(service service.RouterService) *RouterHandler {
    return &RouterHandler{service: service}
}

func (h *RouterHandler) Create(c *gin.Context) {
    var router model.Router
    if err := c.ShouldBindJSON(&router); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.CreateRouter(&router); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, router)
}

func (h *RouterHandler) GetAll(c *gin.Context) {
    routers, err := h.service.GetAllRouters()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, routers)
}

func (h *RouterHandler) GetByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    router, err := h.service.GetRouterByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Router not found"})
        return
    }
    c.JSON(http.StatusOK, router)
}

func (h *RouterHandler) Update(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var router model.Router
    if err := c.ShouldBindJSON(&router); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    router.ID = uint(id)
    if err := h.service.UpdateRouter(&router); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, router)
}

func (h *RouterHandler) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.service.DeleteRouter(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Router deleted successfully"})
}