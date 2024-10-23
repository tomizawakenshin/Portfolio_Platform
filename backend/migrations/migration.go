package main

import (
	"backend/infra"
	"backend/models"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	if err := db.AutoMigrate(&models.User{}, &models.JobType{}, &models.Skill{}); err != nil {
		panic("Failed to migrate db")
	}
}
