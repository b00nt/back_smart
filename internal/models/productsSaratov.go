package models

import (
	"gorm.io/gorm"
)

type CategorySaratov struct {
	gorm.Model
	Name string `gorm:"size:32"`
}

type ProductsSaratov struct {
	gorm.Model
	MoyskladID    string         `gorm:"unique"`
	Code          string         `gorm:"type:varchar(9)"`
	Name          string         `gorm:"type:varchar(64);not null"`
	CategoryID    uint           `gorm:"not null"`
	Category      Category       `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;"`
	Modification  []Modification `gorm:"foreignKey:MoyskladID;references:MoyskladID;constraint:OnDelete:CASCADE;"`
	Price         float64
	Discount      float64         `gorm:"default:0.0;check:Discount >= 0 AND Discount <= 1"`
	Popular       uint64          `gorm:"default:0"`
	Stock         int             `gorm:"default:0"`
	ProductImages []ProductImages `gorm:"foreignKey:MoyskladID;references:MoyskladID;constraint:OnDelete:CASCADE;"`
}

type ProductImagesSaratov struct {
	gorm.Model
	ImgID      string `gorm:"unique"`
	MoyskladID string `gorm:"not null"`
	ImageURL   string `gorm:"type:varchar(255)"`
}

type ModificationSaratov struct {
	gorm.Model
	ModID                       string                        `gorm:"unique"`
	MoyskladID                  string                        `gorm:"not null"`
	Code                        string                        `gorm:"type:varchar(9)"`
	Stock                       int                           `gorm:"default:0"`
	ModificationCharacteristics []ModificationCharacteristics `gorm:"foreignKey:ModID;references:ModID;constraint:OnDelete:CASCADE;"`
	ModificationImages          []ModificationImages          `gorm:"foreignKey:ModID;references:ModID;constraint:OnDelete:CASCADE;"`
	SalePrices                  float64
}

type ModificationImagesSaratov struct {
	gorm.Model
	ImgID    string `gorm:"unique"`
	ModID    string `gorm:"not null"`
	ImageURL string `gorm:"type:varchar(255);not null"`
}

type ModificationCharacteristicsSaratov struct {
	gorm.Model
	ModID string `gorm:"not null"`
	Name  string `gorm:"type:varchar(64);not null"`
	Value string `gorm:"type:varchar(64);not null"`
}
