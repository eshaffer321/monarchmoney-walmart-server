package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"monarchmoney-sync-backend/models"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates requests using the X-Extension-Key header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		extensionKey := c.GetHeader("X-Extension-Key")
		expectedKey := os.Getenv("EXTENSION_SECRET_KEY")

		// For testing, allow "test-secret" when env var is not set
		if expectedKey == "" {
			expectedKey = "test-secret"
		}

		if extensionKey == "" || extensionKey != expectedKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Unauthorized: Missing or invalid extension key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ReceiveOrders handles incoming Walmart order data from the Chrome extension.
func ReceiveOrders(c *gin.Context) {
	// Get Sentry hub from context if available
	hub := sentrygin.GetHubFromContext(c)

	var order models.Order

	// Parse JSON request body
	if err := c.ShouldBindJSON(&order); err != nil {
		// Capture validation errors to Sentry
		if hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetLevel(sentry.LevelWarning)
				scope.SetContext("order", map[string]interface{}{
					"error": err.Error(),
					"body":  c.Request.Body,
				})
				hub.CaptureMessage("Invalid order JSON received")
			})
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Invalid JSON or validation error: %v", err),
		})
		return
	}

	// Additional validation
	if order.OrderTotal <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid order total: must be positive",
		})
		return
	}

	// Validate items
	if len(order.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Order must contain at least one item",
		})
		return
	}

	for _, item := range order.Items {
		if item.Price < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid item price: must be non-negative",
			})
			return
		}
		if item.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid item quantity: must be positive",
			})
			return
		}
	}

	// Log the received order
	log.Printf("Received Walmart order: %s, Total: $%.2f, Items: %d\n",
		order.OrderNumber, order.OrderTotal, len(order.Items))

	// Track successful order in Sentry
	if hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelInfo)
			scope.SetContext("order", map[string]interface{}{
				"order_number": order.OrderNumber,
				"total":        order.OrderTotal,
				"items_count":  len(order.Items),
			})
			scope.SetTag("order.source", "walmart")
			hub.CaptureMessage("Order received successfully")
		})
	}

	// TODO: Process order with Monarch Money SDK
	// For now, just acknowledge receipt

	response := models.OrderResponse{
		Status:    "success",
		Message:   "Order received successfully",
		OrderID:   order.OrderNumber,
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}
