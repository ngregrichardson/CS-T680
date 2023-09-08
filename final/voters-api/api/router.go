package api

import (
	"net/http"
	"time"
	"voters-api/middleware"
	"voters-api/schema"
	"voters-api/utils"
	"voters-api/voters"

	"github.com/gin-gonic/gin"
)

var (
	votersService *voters.VotersService
)

func init() {
	var err error
	votersService, err = voters.NewVotersService()

	if err != nil {
		panic(err)
	}
}

func GetRouter(group *gin.RouterGroup, hostname string) {
	votersService.Hostname = hostname

	group.GET("/health", getHealth)
	group.GET("", getVoters)
	group.POST("", createVoter)
	group.GET("/:id", getVoter)
	group.PUT("/:id", updateVoter)
	group.DELETE("/:id", deleteVoter)

	group.GET("/:id/votes", getVoterHistory)
	group.POST("/:id/votes", createVoteRecord)
	group.GET("/:id/votes/:pollId", getVoteRecord)
	group.PATCH("/:id/votes/:pollId", updateVoteRecord)
	group.DELETE("/:id/votes/:pollId", deleteVoteRecord)
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

func getVoters(c *gin.Context) {
	voters, err := votersService.GetVoters()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
	}

	externalVoters := make([]gin.H, 0)

	for _, voter := range voters {
		externalVoters = append(externalVoters, votersService.FormatExternalVoter(voter))
	}

	c.JSON(http.StatusOK, externalVoters)
}

func createVoter(c *gin.Context) {
	var voterBody schema.Voter = schema.BlankVoter()
	validationError := utils.ValidateBody(c, &voterBody)

	if validationError != nil {
		return
	}

	createdVoter, createVoterError := votersService.AddVoter(voterBody)

	if createVoterError != nil {
		c.AbortWithStatusJSON(http.StatusConflict, utils.ResponseError(createVoterError.Error()))
		return
	}

	c.JSON(http.StatusCreated, votersService.FormatExternalVoter(createdVoter))
}

func getVoter(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	voter, voterNotFoundError := votersService.GetVoter(id)

	if voterNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(voterNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, votersService.FormatExternalVoter(voter))
}

func updateVoter(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	var voterBody schema.Voter = schema.BlankVoter()
	validationError := utils.ValidateBody(c, &voterBody)

	if validationError != nil {
		return
	}

	updatedVoter, voterNotFoundError := votersService.UpdateVoter(id, voterBody)

	if voterNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(voterNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, votersService.FormatExternalVoter(updatedVoter))
}

func deleteVoter(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	voterNotFoundError := votersService.DeleteVoter(id)

	if voterNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(voterNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseMessage("voter deleted"))
}

func getVoterHistory(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	voterHistory, err := votersService.GetVoterHistory(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(err.Error()))
		return
	}

	externalVoteRecords := make([]gin.H, 0)

	for _, voteRecord := range voterHistory {
		externalVoteRecords = append(externalVoteRecords, votersService.FormatExternalVoteRecord(id, voteRecord))
	}

	c.JSON(http.StatusOK, externalVoteRecords)
}

func getVoteRecord(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	pollId, pollIdValidationError := utils.ValidateIntParam(c, "pollId")

	if pollIdValidationError != nil {
		return
	}

	voteRecord, _, voteRecordNotFoundError := votersService.GetVoteRecord(id, pollId)

	if voteRecordNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(voteRecordNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, votersService.FormatExternalVoteRecord(id, *voteRecord))
}

func createVoteRecord(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	var voteRecordBody schema.VoteRecord = schema.BlankVoteRecord()
	validationError := utils.ValidateBody(c, &voteRecordBody)

	if validationError != nil {
		return
	}

	createdVoteRecord, createVoteRecordError := votersService.CreateVoteRecord(id, voteRecordBody.PollID)

	if createVoteRecordError != nil {
		c.AbortWithStatusJSON(http.StatusConflict, utils.ResponseError(createVoteRecordError.Error()))
		return
	}

	c.JSON(http.StatusCreated, votersService.FormatExternalVoteRecord(id, *createdVoteRecord))
}

func updateVoteRecord(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	pollId, pollIdValidationError := utils.ValidateIntParam(c, "pollId")

	if pollIdValidationError != nil {
		return
	}

	voteRecord, voteRecordNotFoundError := votersService.UpdateVoteRecord(id, pollId)

	if voteRecordNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(voteRecordNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, votersService.FormatExternalVoteRecord(id, *voteRecord))
}

func deleteVoteRecord(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	pollId, pollIdValidationError := utils.ValidateIntParam(c, "pollId")

	if pollIdValidationError != nil {
		return
	}

	voteRecordNotFoundError := votersService.DeleteVoteRecord(id, pollId)

	if voteRecordNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(voteRecordNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseMessage("vote record deleted"))
}
