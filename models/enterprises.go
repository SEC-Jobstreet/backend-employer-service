package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Enterprises struct {
	ID uuid.UUID `gorm:"primarykey" json:"id"`

	Name      string `json:"name"`
	Country   string `json:"country"`
	Address   string `json:"address"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Field     string `json:"field"`
	Size      string `json:"size"`
	Url       string `json:"url"`
	License   string `json:"license"`

	EmployerID   string `json:"employer_id"`
	EmployerRole string `json:"employer_role"`

	UpdatedAt int64 `gorm:"autoUpdateTime"`
	CreatedAt int64 `gorm:"autoCreateTime"`
}

func MigrateEnterprises(db *gorm.DB) error {
	err := db.AutoMigrate(&Enterprises{})
	return err
}
