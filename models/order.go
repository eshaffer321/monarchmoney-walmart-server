// Package models contains data models for the Walmart-Monarch sync backend.
package models

import "time"

// Order represents a Walmart order received from the Chrome extension.
type Order struct {
	OrderNumber string      `json:"orderNumber" binding:"required"`
	OrderDate   string      `json:"orderDate" binding:"required"`
	OrderTotal  float64     `json:"orderTotal" binding:"required"`
	Items       []OrderItem `json:"items" binding:"required"`
}

// OrderItem represents an individual item within a Walmart order.
type OrderItem struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Quantity int     `json:"quantity" binding:"required"`
	Category string  `json:"category,omitempty"`
}

// OrderResponse represents the API response after processing an order.
type OrderResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
	OrderID   string    `json:"orderId,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
