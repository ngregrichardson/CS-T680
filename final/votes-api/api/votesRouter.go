package api

import (
	"net/http"
	"time"
	"votes-api/middleware"
	"votes-api/schema"
	"votes-api/services"
	"votes-api/utils"

	"github.com/gin-gonic/gin"
)

var (
	votesService *services.VotesService
)

func init() {
	var err error
	votesService, err = services.NewVotesService()

	if err != nil {
		panic(err)
	}
}

func GetVotesRouter(group *gin.RouterGroup, hostname string) {
	votesService.Hostname = hostname

	group.GET("/health", getHealth)
	group.GET("/:id", getVote)
	group.POST("", createVote)
	group.PATCH("/:id", updateVote)
	group.DELETE("/:id", deleteVote)
}

func getHealth(c *gin.Context) {
	hang, hangValidationError := utils.ValidateOptionalIntQuery(c, "hang")

	if hangValidationError != nil {
		return
	}

	if hang > 0 {
		time.Sleep(time.Duration(hang) * time.Second)
	}

	c.JSON(http.StatusOK, gin.H{
		"version": "1.0.0",
		"stats":   middleware.GetStats(),
	})
}

func getVote(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	vote, voteRelations, getVoteError := votesService.GetVote(id)

	if getVoteError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(getVoteError.Error()))
		return
	}

	c.JSON(http.StatusOK, votesService.FormatExternalVote(vote, voteRelations))
}

func createVote(c *gin.Context) {
	var voteBody schema.Vote = schema.BlankVote()
	validationError := utils.ValidateBody(c, &voteBody)

	if validationError != nil {
		return
	}

	createdVote, voteRelations, createVoteError := votesService.AddVote(voteBody)

	if createVoteError != nil {
		c.AbortWithStatusJSON(http.StatusConflict, utils.ResponseError(createVoteError.Error()))
		return
	}

	c.JSON(http.StatusCreated, votesService.FormatExternalVote(createdVote, voteRelations))
}

func updateVote(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	var voteBody schema.Vote = schema.BlankVote()
	validationError := utils.ValidateBody(c, &voteBody)

	if validationError != nil {
		return
	}

	updatedVote, voteRelations, updatedVoteError := votesService.UpdateVote(id, voteBody)

	if updatedVoteError != nil {
		c.AbortWithStatusJSON(http.StatusConflict, utils.ResponseError(updatedVoteError.Error()))
		return
	}

	c.JSON(http.StatusOK, votesService.FormatExternalVote(updatedVote, voteRelations))
}

func deleteVote(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	deleteVoteError := votesService.DeleteVote(id)

	if deleteVoteError != nil {
		c.AbortWithStatusJSON(http.StatusConflict, utils.ResponseError(deleteVoteError.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseMessage("vote deleted"))
}
