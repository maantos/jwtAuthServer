package db

import "github.com/maantos/jwtAuth/pkg/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
