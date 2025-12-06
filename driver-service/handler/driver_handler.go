package handler

//HTTP request/response işleri

import (
	"net/http"
	"strconv"

	"bitaksi_taxihub/dto"
	"bitaksi_taxihub/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DriverHandler struct {
	svc      service.DriverService
	validate *validator.Validate
}

func NewDriverHandler(svc service.DriverService) *DriverHandler {
	return &DriverHandler{
		svc:      svc,
		validate: validator.New(),
	}
}

// CreateDriver godoc
// @Summary Yeni driver oluştur
// @Description Yeni driver kaydı ekler ve id döner
// @Tags drivers
// @Accept json
// @Produce json
// @Param driver body dto.CreateDriverRequest true "Driver info"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /drivers [post]
func (h *DriverHandler) CreateDriver(c *gin.Context) {
	var req dto.CreateDriverRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.svc.CreateDriver(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id.Hex()})
}

// UpdateDriver godoc
// @Summary Driver güncelle
// @Description ID'ye göre driver alanlarını günceller
// @Tags drivers
// @Accept json
// @Produce json
// @Param id path string true "Driver ID"
// @Param driver body dto.UpdateDriverRequest true "Update fields"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /drivers/{id} [put]
func (h *DriverHandler) UpdateDriver(c *gin.Context) {
	idStr := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateDriverRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateDriver(c.Request.Context(), objID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})

}

// ListDrivers godoc
// @Summary Driver listele
// @Description Sayfalı driver listesi döner
// @Tags drivers
// @Produce json
// @Param page query int false "Page" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} dto.ListDrivers
// @Failure 500 {object} map[string]string
// @Router /drivers [get]
func (h *DriverHandler) ListDrivers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		page = 20
	}
	resp, err := h.svc.ListDrivers(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// NearbyDrivers godoc
// @Summary Yakındaki taksileri getir
// @Description lat/lon ve taxiType’a göre 6km içindeki taksileri döner
// @Tags drivers
// @Produce json
// @Param lat query number true "Latitude"
// @Param lon query number true "Longitude"
// @Param taxiType query string true "Taxi Type" Enums(sari,buyuk,luks)
// @Success 200 {array} dto.NearbyDriverResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /drivers/nearby [get]
func (h *DriverHandler) NearbyDrivers(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	taxiType := c.Query("taxiType")

	if latStr == "" || lonStr == "" || taxiType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "lat, lon and taxiType are required",
		})
		return
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lat"})
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lon"})
		return
	}
	resp, err := h.svc.NearbyDrivers(c.Request.Context(), lat, lon, taxiType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)

}
