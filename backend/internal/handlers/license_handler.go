package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// LicenseHandler handles license-related requests
type LicenseHandler struct {
	*BaseHandler
	licenseService *services.LicenseService
}

// NewLicenseHandler creates a new license handler
func NewLicenseHandler(db *gorm.DB, cfg *config.Config) *LicenseHandler {
	return &LicenseHandler{
		BaseHandler:    NewBaseHandler(db, cfg),
		licenseService: services.LicenseServiceInstance,
	}
}

// CreateLicense creates a new license
func (h *LicenseHandler) CreateLicense(c *gin.Context) {
	var req services.CreateLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	license, err := h.licenseService.CreateLicense(&req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, license)
}

// GetLicense gets a license by ID
func (h *LicenseHandler) GetLicense(c *gin.Context) {
	licenseID := c.Param("id")

	license, err := h.licenseService.GetLicense(licenseID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "License not found")
		return
	}

	h.sendSuccess(c, license)
}

// GetLicenseByKey gets a license by key
func (h *LicenseHandler) GetLicenseByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		h.sendError(c, http.StatusBadRequest, "License key is required")
		return
	}

	license, err := h.licenseService.GetLicenseByKey(key)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "License not found")
		return
	}

	h.sendSuccess(c, license)
}

// ValidateLicense validates a license
func (h *LicenseHandler) ValidateLicense(c *gin.Context) {
	var req services.ValidateLicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	license, err := h.licenseService.ValidateLicense(&req)
	if err != nil {
		h.sendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.sendSuccess(c, license)
}

// UpdateLicense updates a license
func (h *LicenseHandler) UpdateLicense(c *gin.Context) {
	licenseID := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	license, err := h.licenseService.UpdateLicense(licenseID, updates)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, license)
}

// RevokeLicense revokes a license
func (h *LicenseHandler) RevokeLicense(c *gin.Context) {
	licenseID := c.Param("id")

	if err := h.licenseService.RevokeLicense(licenseID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "License revoked successfully"})
}

// ListLicenses lists licenses with pagination
func (h *LicenseHandler) ListLicenses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	orgID := c.Query("org_id")

	licenses, total, err := h.licenseService.ListLicenses(page, limit, &orgID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"licenses": licenses,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// CheckLicenseUsage checks license usage
func (h *LicenseHandler) CheckLicenseUsage(c *gin.Context) {
	licenseID := c.Param("id")

	usage, err := h.licenseService.CheckLicenseUsage(licenseID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, usage)
}

// GenerateOfflineLicense generates an offline license
func (h *LicenseHandler) GenerateOfflineLicense(c *gin.Context) {
	licenseID := c.Param("id")

	var req struct {
		PrivateKeyPath string `json:"private_key_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	offlineLicense, err := h.licenseService.GenerateOfflineLicense(licenseID, req.PrivateKeyPath)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"offline_license": offlineLicense})
}

// ValidateOfflineLicense validates an offline license
func (h *LicenseHandler) ValidateOfflineLicense(c *gin.Context) {
	var req struct {
		OfflineLicense string `json:"offline_license" binding:"required"`
		PublicKeyPath  string `json:"public_key_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	license, err := h.licenseService.ValidateOfflineLicense(req.OfflineLicense, req.PublicKeyPath)
	if err != nil {
		h.sendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.sendSuccess(c, license)
}

// GetLicenseStats gets license statistics
func (h *LicenseHandler) GetLicenseStats(c *gin.Context) {
	stats, err := h.licenseService.GetLicenseStats()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}
