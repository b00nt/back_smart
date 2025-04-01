package models

import (
	"gorm.io/gorm"
	"time"
)

type Feedback struct {
	gorm.Model
	Name      string `gorm:"size:64;not null"`
	Telephone string `gorm:"size:20;not null"`
	City      string `gorm:"size:32;not null"`
}

type CustomerInfo struct {
	gorm.Model
	FullName        string `gorm:"type:varchar(64);not null" json:"full_name"`
	TelephoneNumber string `gorm:"type:varchar(20);not null" json:"telephone_number"`
	Email           string `gorm:"type:varchar(64);not null" json:"email"`
	Comment         string `gorm:"type:varchar(64);default:''" json:"comment,omitempty"`

	City      string `gorm:"type:varchar(32);not null" json:"city"`
	Street    string `gorm:"type:varchar(32);not null" json:"street"`
	House     string `gorm:"type:varchar(32);not null" json:"house"`
	Entrance  string `gorm:"type:varchar(32);not null" json:"entrance"`
	Floor     string `gorm:"type:varchar(8);default:null" json:"floor"`
	Apartment string `gorm:"type:varchar(8);default:null" json:"apartment"`

	AnotherFullName        string `gorm:"type:varchar(64);default:''" json:"another_full_name,omitempty"`
	AnotherTelephoneNumber string `gorm:"type:varchar(20);default:''" json:"another_telephone_number,omitempty"`
}

// Order represents the order model
type Order struct {
	gorm.Model
	CustomerInfoID uint         `gorm:"not null" json:"customer_info_id"`
	CustomerInfo   CustomerInfo `gorm:"foreignKey:CustomerInfoID" json:"customer_info"`
	OrderDate      time.Time    `gorm:"autoCreateTime" json:"order_date"`
	TotalAmount    float64      `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	// Status         string    `gorm:"type:varchar(50);default:'Pending'" json:"status"` // Uncomment if needed
}

type OrderItem struct {
	gorm.Model
	OrderID                     uint                               `gorm:"not null" json:"order_id"`
	Order                       Order                              `gorm:"foreignKey:OrderID" json:"order"`
	Name                        string                             `gorm:"type:varchar(255);not null" json:"name"`
	MoyskladID                  string                             `gorm:"type:varchar(255);not null" json:"moysklad_id"`
	ModificationCharacteristics []ModificationCharacteristicsOrder `gorm:"foreignKey:OrderItemID;constraint:OnDelete:CASCADE;" json:"modification_characteristics"`
	Quantity                    uint                               `gorm:"not null" json:"quantity"`
	Price                       float64                            `gorm:"type:decimal(10,2);not null" json:"price"`
}

// ModificationCharacteristicsOrder represents the modification characteristics for an order item
type ModificationCharacteristicsOrder struct {
	gorm.Model
	OrderItemID uint   `gorm:"not null" json:"order_item_id"`
	ModID       string `gorm:"type:varchar(191);not null"`
	Name        string `gorm:"type:varchar(128);not null" json:"name"`
	Value       string `gorm:"type:varchar(128);not null" json:"value"`
}

// OrderItem represents the order item model
// type OrderItem struct {
// 	gorm.Model
// 	OrderID                     uint                               `gorm:"not null" json:"order_id"`
// 	Order                       Order                              `gorm:"foreignKey:OrderID" json:"order"`
// 	Name                        string                             `gorm:"type:varchar(255);not null" json:"name"`
// 	MoyskladID                  string                             `gorm:"type:varchar(255);not null" json:"moysklad_id"`
// 	ModificationCharacteristics []ModificationCharacteristicsOrder `gorm:"foreignKey:MoyskladID;references:MoyskladID;constraint:OnDelete:CASCADE;"`
// 	Quantity                    uint                               `gorm:"not null" json:"quantity"`
// 	Price                       float64                            `gorm:"type:decimal(10,2);not null" json:"price"`
// }
//
// type ModificationCharacteristicsOrder struct {
// 	gorm.Model
// 	ModID string `gorm:"type:varchar(191);not null"`
// 	Name  string `gorm:"type:varchar(128);not null"`
// 	Value string `gorm:"type:varchar(128);not null"`
// }
