package models

import (
	"time"

	"gorm.io/gorm"
)

type ModelWithDeletes struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&VerifyAccount{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&ForgetPassword{}); err != nil {
		return err
	}

	return nil
}
