package config

import (
	"authentication-service/models"

	"gorm.io/gorm"
)

func UpdateDB(db *gorm.DB) {
	db.AutoMigrate(models.User{})
}
