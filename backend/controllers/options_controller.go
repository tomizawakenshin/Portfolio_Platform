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
}

type OptionsController struct {
	jobTypeService services.IJobTypeService
	skillService   services.ISkillService
}

func NewOptionsController(
	jobTypeService services.IJobTypeService,
	skillService services.ISkillService,
) IOptionsController {
	return &OptionsController{
		jobTypeService: jobTypeService,
		skillService:   skillService,
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
