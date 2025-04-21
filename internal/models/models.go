package models

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	gorm.Model
	MoyskladID      string         `gorm:"type:varchar(128);unique;index" json:"moyskladID"`
	Name            string         `gorm:"type:varchar(128);not null" json:"name"`
	Code            string         `gorm:"type:varchar(10)" json:"code"`
	Category        string         `gorm:"type:varchar(64);" json:"category"`
	Price           float64        `gorm:"default:0.0" json:"price"`
	Discount        float64        `gorm:"default:0.0;check:Discount >= 0 AND Discount <= 1" json:"discount"`
	Stock           int            `gorm:"default:0" json:"stock"`
	City            string         `gorm:"type:varchar(32)" json:"city" `
	PopularityCount uint64         `gorm:"default:0" json:"popularityCount"`
	Modification    []Modification `gorm:"foreignKey:MoyskladID;references:MoyskladID;constraint:OnDelete:CASCADE;" json:"modification"`
	Image           []Image        `gorm:"foreignKey:MoyskladID;references:MoyskladID;constraint:OnDelete:CASCADE;" json:"image"`
	Display         bool           `gorm:"default:true" json:"display"`
}

type Modification struct {
	gorm.Model
	Name              string              `gorm:"type:varchar(128);not null" json:"name"`
	ModID             string              `gorm:"type:varchar(191);unique" json:"modificationID"`
	MoyskladID        string              `gorm:"type:varchar(128);not null" json:"MoyskladID"`
	Code              string              `gorm:"type:varchar(9)" json:"code"`
	Stock             int                 `gorm:"default:0" json:"stock"`
	Characteristic    []Characteristic    `gorm:"foreignKey:ModID;references:ModID;constraint:OnDelete:CASCADE;" json:"characteristic"`
	ModificationImage []ModificationImage `gorm:"foreignKey:ModID;references:ModID;constraint:OnDelete:CASCADE;" json:"modificationImage"`
	Price             float64             `gorm:"default:0.0" json:"price"`
}

type Image struct {
	gorm.Model
	ImgID      string `gorm:"unique" json:"imageID"`
	MoyskladID string `gorm:"type:varchar(128);not null" json:"moyskladID"`
	ImageURL   string `gorm:"type:varchar(255)" json:"imageURL"`
	ImagePath  string `gorm:"type:varchar(255)" json:"imagePath"`
}

type Characteristic struct {
	gorm.Model
	ModID string `gorm:"type:varchar(128)" json:"modificationID"`
	Name  string `gorm:"type:varchar(128);not null" json:"name"`
	Value string `gorm:"type:varchar(128);not null" json:"value"`
}

type ModificationImage struct {
	gorm.Model
	ImgID     string `gorm:"unique" json:"imageID"`
	ModID     string `gorm:"type:varchar(128)" json:"modificationID"`
	ImageURL  string `gorm:"type:varchar(255);not null" json:"imageURL"`
	ImagePath string `gorm:"type:varchar(255)" json:"imagePath"`
}

type Feedback struct {
	gorm.Model
	Name      string `gorm:"type:varchar(64);not null" json:"name"`
	Telephone string `gorm:"type:varchar(20);not null" json:"telephone"`
	City      string `gorm:"type:varchar(32);not null" json:"city"`
}

type Order struct {
	gorm.Model
	CustomerInfoID uint         `gorm:"not null" json:"customerInfoId"`
	CustomerInfo   CustomerInfo `gorm:"foreignKey:CustomerInfoID" json:"customerInfo"`
	OrderItems     []OrderItem  `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"orderItems"`
	OrderDate      time.Time    `gorm:"autoCreateTime" json:"orderDate"`
	TotalAmount    float64      `gorm:"type:decimal(10,2);not null" json:"totalAmount"`
}

type CustomerInfo struct {
	gorm.Model
	FullName               string `gorm:"type:varchar(64);not null" json:"fullName"`
	TelephoneNumber        string `gorm:"type:varchar(20);not null" json:"telephoneNumber"`
	Email                  string `gorm:"type:varchar(64);not null" json:"email"`
	Comment                string `gorm:"type:varchar(64);default:''" json:"comment,omitempty"`
	City                   string `gorm:"type:varchar(32);not null" json:"city"`
	Street                 string `gorm:"type:varchar(32);not null" json:"street"`
	House                  string `gorm:"type:varchar(32);not null" json:"house"`
	Entrance               string `gorm:"type:varchar(32);not null" json:"entrance"`
	Floor                  string `gorm:"type:varchar(8);default:null" json:"floor"`
	Apartment              string `gorm:"type:varchar(8);default:null" json:"apartment"`
	AnotherFullName        string `gorm:"type:varchar(64);default:''" json:"anotherFullName,omitempty"`
	AnotherTelephoneNumber string `gorm:"type:varchar(20);default:''" json:"anotherTelephoneNumber,omitempty"`
}

// this is modification, not item
type OrderItem struct {
	gorm.Model
	OrderID         uint                  `gorm:"not null" json:"orderId"`
	Name            string                `gorm:"type:varchar(255);not null" json:"name"`
	MoyskladID      string                `gorm:"type:varchar(255);not null" json:"moyskladId"`
	Characteristics []CharacteristicOrder `gorm:"foreignKey:OrderItemID;constraint:OnDelete:CASCADE;" json:"characteristics"`
	Quantity        uint                  `gorm:"not null" json:"quantity"`
	Price           float64               `gorm:"type:decimal(10,2);not null" json:"price"`
}

type CharacteristicOrder struct {
	gorm.Model
	OrderItemID uint   `gorm:"not null" json:"orderItemId"`
	ModID       string `gorm:"type:varchar(191);not null" json:"modId"`
	Name        string `gorm:"type:varchar(128);not null" json:"name"`
	Value       string `gorm:"type:varchar(128);not null" json:"value"`
}
