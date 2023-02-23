package initializers

import "github.com/maantos/jwtAuth/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
