package models

import (
	"gorm.io/gorm"
	"time"
)

type Feedback struct {
	gorm.Model
	Name      string `gorm:"size:64;not null"`
	Telephone string `gorm:"size:12;not null"`
	City      string `gorm:"size:32;not null"`
}

type CustomerInfo struct {
	gorm.Model
	FullName        string `gorm:"type:varchar(50);not null" json:"full_name"`
	TelephoneNumber string `gorm:"type:varchar(15);not null" json:"telephone_number"`
	Email           string `gorm:"type:varchar(50);not null" json:"email"`
	Comment         string `gorm:"type:varchar(50);default:''" json:"comment,omitempty"`

	City      string `gorm:"type:varchar(20);not null" json:"city"`
	Street    string `gorm:"type:varchar(60);not null" json:"street"`
	House     string `gorm:"type:varchar(10);not null" json:"house"`
	Entrance  string `gorm:"type:varchar(5);not null" json:"entrance"`
	Floor     string `gorm:"type:int;default:null" json:"floor"`     // Use pointer for nullable
	Apartment string `gorm:"type:int;default:null" json:"apartment"` // Use pointer for nullable

	AnotherFullName        string `gorm:"type:varchar(50);default:''" json:"another_full_name,omitempty"`
	AnotherTelephoneNumber string `gorm:"type:varchar(15);default:''" json:"another_telephone_number,omitempty"`
}

// Order represents the order model
type Order struct {
	gorm.Model // This includes fields like ID, CreatedAt, UpdatedAt, DeletedAt

	CustomerInfoID uint         `gorm:"not null" json:"customer_info_id"`               // Foreign key
	CustomerInfo   CustomerInfo `gorm:"foreignKey:CustomerInfoID" json:"customer_info"` // Relationship
	OrderDate      time.Time    `gorm:"autoCreateTime" json:"order_date"`
	TotalAmount    float64      `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	// Status         string    `gorm:"type:varchar(50);default:'Pending'" json:"status"` // Uncomment if needed
}

// OrderItem represents the order item model
type OrderItem struct {
	gorm.Model         // This includes fields like ID, CreatedAt, UpdatedAt, DeletedAt
	OrderID    uint    `gorm:"not null" json:"order_id"`        // Foreign key
	Order      Order   `gorm:"foreignKey:OrderID" json:"order"` // Relationship
	MoyskladID string  `gorm:"type:varchar(255);not null" json:"moysklad_id"`
	Quantity   uint    `gorm:"not null" json:"quantity"`
	Price      float64 `gorm:"type:decimal(10,2);not null" json:"price"`
}
