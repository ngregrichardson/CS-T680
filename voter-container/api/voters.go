package api

import (
	"net/http"
	"time"
	"voter-api/api/responses"
	"voter-api/api/utils"
	"voter-api/middleware"

	"voter-api/voters"

	"github.com/gin-gonic/gin"
)

var (
	voterList *voters.VoterCache
)

func init() {
	voterList, _ = voters.NewVoterCache()
}

func GetVoterRouter(parentGroup *gin.RouterGroup) {
	group := parentGroup.Group("/voters")

	group.GET("", getVoters)
	group.GET("/:id", getVoter)
	group.POST("/:id", createVoter)
	group.PUT("/:id", updateVoter)
	group.DELETE("/:id", deleteVoter)

	group.GET("/:id/votes", getVoterHistory)
	group.GET("/:id/votes/:pollId", getVoterVote)
	group.POST("/:id/votes/:pollId", createVoterVote)
	group.PUT("/:id/votes/:pollId", updateVoterVote)
	group.DELETE("/:id/votes/:pollId", deleteVoterVote)

	group.GET("/health", getHealth)
}

func getVoters(c *gin.Context) {
	voters, err := voterList.GetVoters()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.NewHttpInternalServerErrorResponse(err.Error()))
	}

	c.JSON(http.StatusOK, responses.NewHttpOkResponse(voters))
}

func getVoter(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	voter, voterNotFoundError := voterList.GetVoter(id)

	if voterNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, responses.NewHttpNotFoundResponseWithMessage(voterNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.NewHttpOkResponse(voter))
}

func createVoter(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	var voterBody voters.Voter

	if err := c.ShouldBindJSON(&voterBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.NewHttpBadRequestResponseWithMessage(err.Error()))
		return
	}

	newVoter, voterCreationError := voters.NewVoter(voterBody.VoterID, voterBody.FirstName, voterBody.LastName)

	if voterCreationError != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.NewHttpInternalServerErrorResponseWithMessage(voterCreationError.Error()))
		return
	}

	createdVoter, voterAddError := voterList.AddVoter(id, newVoter)

	if voterAddError != nil {
		c.AbortWithStatusJSON(http.StatusConflict, responses.NewHttpConflictResponseWithMessage(voterAddError.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.NewHttpCreatedResponse(createdVoter))
}

func updateVoter(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	var voterBody voters.Voter

	if err := c.ShouldBindJSON(&voterBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.NewHttpBadRequestResponseWithMessage(err.Error()))
		return
	}

	updatedVoter, voterUpdateError := voterList.UpdateVoter(id, voterBody)

	if voterUpdateError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, responses.NewHttpNotFoundResponseWithMessage(voterUpdateError.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.NewHttpOkResponse(updatedVoter))
}

func deleteVoter(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	voterDeleteError := voterList.DeleteVoter(id)

	if voterDeleteError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, responses.NewHttpNotFoundResponseWithMessage(voterDeleteError.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.NewHttpOkResponseWithMessage("Voter deleted", nil))
}

func getVoterHistory(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	voteHistory, err := voterList.GetVoterHistory(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, responses.NewHttpNotFoundResponseWithMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.NewHttpOkResponse(voteHistory))
}

func getVoterVote(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")
	pollId, pollIdValidationError := utils.ValidateIntParam(c, "pollId")

	if idValidationError != nil || pollIdValidationError != nil {
		return
	}

	voterVote, _, err := voterList.GetVoterVote(id, pollId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, responses.NewHttpNotFoundResponseWithMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.NewHttpOkResponse(voterVote))
}

func createVoterVote(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")
	pollId, pollIdValidationError := utils.ValidateIntParam(c, "pollId")

	if idValidationError != nil || pollIdValidationError != nil {
		return
	}

	newVoterVote, err := voterList.CreateVoterVote(id, pollId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, responses.NewHttpConflictResponseWithMessage(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.NewHttpCreatedResponse(newVoterVote))
}

func updateVoterVote(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")
	pollId, pollIdValidationError := utils.ValidateIntParam(c, "pollId")

	if idValidationError != nil || pollIdValidationError != nil {
		return
	}

	updatedVoterVote, err := voterList.UpdateVoterVote(id, pollId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, responses.NewHttpNotFoundResponseWithMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.NewHttpOkResponse(updatedVoterVote))
}

func deleteVoterVote(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")
	pollId, pollIdValidationError := utils.ValidateIntParam(c, "pollId")

	if idValidationError != nil || pollIdValidationError != nil {
		return
	}

	err := voterList.DeleteVoterVote(id, pollId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, responses.NewHttpNotFoundResponseWithMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.NewHttpOkResponseWithMessage("Vote deleted", nil))
}

func getHealth(c *gin.Context) {
	hang, hangValidationError := utils.ValidateOptionalIntQuery(c, "hang")

	if hangValidationError != nil {
		return
	}

	if hang > 0 {
		time.Sleep(time.Duration(hang) * time.Second)
	}

	c.JSON(http.StatusOK, responses.NewHttpOkResponse(gin.H{
		"version": "1.0.0",
		"stats":   middleware.GetStats(),
	}))
}
