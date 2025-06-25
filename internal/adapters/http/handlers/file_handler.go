package handlers

import (
    "net/http"
    "backend/internal/app/file"
    "github.com/gin-gonic/gin"
)

type FileHandler struct {
    service *file.Service
}

func NewFileHandler(service *file.Service) *FileHandler {
    return &FileHandler{service: service}
}

func (h *FileHandler) GetFileURL(c *gin.Context) {
    key := c.Param("key")
    
    url, err := h.service.GetFileURL(c.Request.Context(), key)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get file URL",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "url": url,
    })
}