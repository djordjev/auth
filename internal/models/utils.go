package models

import "gorm.io/gorm"

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
