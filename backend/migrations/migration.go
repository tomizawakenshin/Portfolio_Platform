package main

import (
	"backend/infra"
	"backend/models"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	if err := db.AutoMigrate(&models.User{}, &models.JobType{}, &models.Skill{}, &models.Post{}, &models.Image{}); err != nil {
		panic("Failed to migrate db")
	}
}
