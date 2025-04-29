package main

import (
	"backend/config"
	portfolioInfra "backend/infrastructure/portfolio"
	userInfra "backend/infrastructure/user"
	"backend/models"
)

func main() {
	config.Initialize()
	db := config.SetupDB()

	if err := db.AutoMigrate(&userInfra.UserModel{}, &models.JobType{}, &models.Skill{}, &models.Genre{}, &portfolioInfra.PostModel{}, &portfolioInfra.ImageModel{}); err != nil {
		panic("Failed to migrate db")
	}
}
