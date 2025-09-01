package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"monarchmoney-sync-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestReceiveOrders_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/walmart/orders", ReceiveOrders)

	order := models.Order{
		OrderNumber: "123456789",
		OrderDate:   "2024-01-15",
		OrderTotal:  150.00,
		Items: []models.OrderItem{
			{
				Name:     "Great Value Milk",
				Price:    3.99,
				Quantity: 1,
			},
			{
				Name:     "Bounty Paper Towels",
				Price:    12.99,
				Quantity: 2,
			},
		},
	}

	jsonData, _ := json.Marshal(order)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/walmart/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Extension-Key", "test-secret")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.OrderResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "123456789", response.OrderID)
}

func TestReceiveOrders_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/walmart/orders", ReceiveOrders)

	invalidJSON := []byte(`{"invalid": json}`)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/walmart/orders", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Extension-Key", "test-secret")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Contains(t, response["message"], "Invalid JSON")
}

func TestReceiveOrders_MissingAuth(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware())
	router.POST("/api/walmart/orders", ReceiveOrders)

	order := models.Order{
		OrderNumber: "123456789",
		OrderDate:   "2024-01-15",
		OrderTotal:  150.00,
		Items: []models.OrderItem{
			{
				Name:     "Great Value Milk",
				Price:    3.99,
				Quantity: 1,
			},
		},
	}

	jsonData, _ := json.Marshal(order)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/walmart/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	// Note: NOT setting X-Extension-Key header
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Contains(t, response["message"], "Unauthorized")
}

func TestReceiveOrders_EmptyOrder(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/walmart/orders", ReceiveOrders)

	order := models.Order{
		// Missing required fields
	}

	jsonData, _ := json.Marshal(order)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/walmart/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Extension-Key", "test-secret")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Contains(t, response["message"], "validation")
}

func TestReceiveOrders_InvalidOrderTotal(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/walmart/orders", ReceiveOrders)

	order := models.Order{
		OrderNumber: "123456789",
		OrderDate:   "2024-01-15",
		OrderTotal:  -50.00, // Negative total
		Items: []models.OrderItem{
			{
				Name:     "Great Value Milk",
				Price:    3.99,
				Quantity: 1,
			},
		},
	}

	jsonData, _ := json.Marshal(order)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/walmart/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Extension-Key", "test-secret")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Contains(t, response["message"], "Invalid order total")
}
