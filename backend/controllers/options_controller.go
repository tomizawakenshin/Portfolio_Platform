// controllers/options_controller.go

package controllers

import (
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IOptionsController interface {
	GetJobTypes(ctx *gin.Context)
	GetSkills(ctx *gin.Context)
	GetGenre(ctx *gin.Context)
}

type OptionsController struct {
	jobTypeService services.IJobTypeService
	skillService   services.ISkillService
	genreService   services.IGenreService
}

func NewOptionsController(
	jobTypeService services.IJobTypeService,
	skillService services.ISkillService,
	genreService services.IGenreService,
) IOptionsController {
	return &OptionsController{
		jobTypeService: jobTypeService,
		skillService:   skillService,
		genreService:   genreService,
	}
}

func (c *OptionsController) GetJobTypes(ctx *gin.Context) {
	jobTypes, err := c.jobTypeService.GetAllJobTypes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get job types"})
		return
	}

	// `Name` フィールドのみを返すようにマッピング
	var jobTypeNames []string
	for _, jt := range jobTypes {
		jobTypeNames = append(jobTypeNames, jt.Name)
	}

	ctx.JSON(http.StatusOK, gin.H{"jobTypes": jobTypeNames})
}

func (c *OptionsController) GetSkills(ctx *gin.Context) {
	skills, err := c.skillService.GetAllSkills()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get skills"})
		return
	}

	// `Name` フィールドのみを返すようにマッピング
	var skillNames []string
	for _, skill := range skills {
		skillNames = append(skillNames, skill.Name)
	}

	ctx.JSON(http.StatusOK, gin.H{"skills": skillNames})
}

func (c *OptionsController) GetGenre(ctx *gin.Context) {
	genres, err := c.genreService.GetAllGenre()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get genres"})
		return
	}

	// `Name` フィールドのみを返すようにマッピング
	var genreNames []string
	for _, genre := range genres {
		genreNames = append(genreNames, genre.Name)
	}

	ctx.JSON(http.StatusOK, gin.H{"genres": genreNames})
}
