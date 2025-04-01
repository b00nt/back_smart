// internal/models/productsSaratov.go
package models

import (
	"gorm.io/gorm"
)

type ProductsSaratov struct {
	gorm.Model
	MoyskladID    string                 `gorm:"type:varchar(128);unique;index"`
	Code          string                 `gorm:"type:varchar(9)"`
	Name          string                 `gorm:"type:varchar(128);not null"`
	CategoryID    uint                   `gorm:"not null"`
	Category      Category               `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;"`
	Modification  []ModificationSaratov  `gorm:"foreignKey:MoyskladID;references:MoyskladID;constraint:OnDelete:CASCADE;"`
	Price         float64                `gorm:"default:0.0"`
	Discount      float64                `gorm:"default:0.0;check:Discount >= 0 AND Discount <= 1"`
	Popular       uint64                 `gorm:"default:0"`
	Stock         int                    `gorm:"default:0"`
	ProductImages []ProductImagesSaratov `gorm:"foreignKey:MoyskladID;references:MoyskladID;constraint:OnDelete:CASCADE;"`
	Display       bool                   `gorm:"default:true"`
}

type ModificationImagesSaratov struct {
	gorm.Model
	ImgID    string `gorm:"unique"`
	ModID    string `gorm:"type:varchar(191)"`
	ImageURL string `gorm:"type:varchar(255);not null"`
}

type ModificationCharacteristicsSaratov struct {
	gorm.Model
	ModID string `gorm:"type:varchar(191);not null;index:idx_modification_characteristics,unique"`
	Name  string `gorm:"type:varchar(128);not null;index:idx_modification_characteristics,unique"`
	Value string `gorm:"type:varchar(128);not null;index:idx_modification_characteristics,unique"`
}

type ModificationSaratov struct {
	gorm.Model
	Name                        string                               `gorm:"type:varchar(128);not null"`
	ModID                       string                               `gorm:"type:varchar(191);unique"`
	MoyskladID                  string                               `gorm:"type:varchar(128);not null"`
	Code                        string                               `gorm:"type:varchar(9)"`
	Stock                       int                                  `gorm:"default:0"`
	ModificationCharacteristics []ModificationCharacteristicsSaratov `gorm:"foreignKey:ModID;references:ModID;constraint:OnDelete:CASCADE;"`
	ModificationImages          []ModificationImagesSaratov          `gorm:"foreignKey:ModID;references:ModID;constraint:OnDelete:CASCADE;"`
	SalePrices                  float64                              `gorm:"default:0.0"`
	Display                     bool                                 `gorm:"default:true"`
}

type ProductImagesSaratov struct {
	gorm.Model
	ImgID      string `gorm:"unique"`
	MoyskladID string `gorm:"type:varchar(128);not null"`
	ImageURL   string `gorm:"type:varchar(255)"`
}
