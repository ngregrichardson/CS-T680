package api

import (
	"net/http"
	"polls-api/middleware"
	"polls-api/polls"
	"polls-api/schema"
	"polls-api/utils"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	pollsService *polls.PollsService
)

func init() {
	var err error
	pollsService, err = polls.NewPollsService()

	if err != nil {
		panic(err)
	}
}

func GetRouter(group *gin.RouterGroup, hostname string) {
	pollsService.Hostname = hostname

	group.GET("/health", getHealth)
	group.GET("", getPolls)
	group.POST("", createPoll)
	group.GET("/:id", getPoll)
	group.PUT("/:id", updatePoll)
	group.DELETE("/:id", deletePoll)

	group.GET("/:id/options", getPollOptions)
	group.POST("/:id/options", createPollOption)
	group.GET("/:id/options/:optionId", getPollOption)
	group.PATCH("/:id/options/:optionId", updatePollOption)
	group.DELETE("/:id/options/:optionId", deletePollOption)
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

func getPolls(c *gin.Context) {
	polls, err := pollsService.GetPolls()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
	}

	externalPolls := make([]gin.H, 0)

	for _, poll := range polls {
		externalPolls = append(externalPolls, pollsService.FormatExternal(poll))
	}

	c.JSON(http.StatusOK, externalPolls)
}

func createPoll(c *gin.Context) {
	var pollBody schema.Poll = schema.BlankPoll()
	validationError := utils.ValidateBody(c, &pollBody)

	if validationError != nil {
		return
	}

	createdPoll, createPollError := pollsService.AddPoll(pollBody)

	if createPollError != nil {
		c.AbortWithStatusJSON(http.StatusConflict, utils.ResponseError(createPollError.Error()))
		return
	}

	c.JSON(http.StatusCreated, pollsService.FormatExternal(createdPoll))
}

func getPoll(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	poll, pollNotFoundError := pollsService.GetPoll(id)

	if pollNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(pollNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, pollsService.FormatExternal(poll))
}

func updatePoll(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	var pollBody schema.Poll = schema.BlankPoll()
	validationError := utils.ValidateBody(c, &pollBody)

	if validationError != nil {
		return
	}

	updatedPoll, pollNotFoundError := pollsService.UpdatePoll(id, pollBody)

	if pollNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(pollNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, pollsService.FormatExternal(updatedPoll))
}

func deletePoll(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	pollNotFoundError := pollsService.DeletePoll(id)

	if pollNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(pollNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseMessage("poll deleted"))
}

func getPollOptions(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	pollOptions, err := pollsService.GetPollOptions(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(err.Error()))
		return
	}

	externalPollOptions := make([]gin.H, 0)

	for _, pollOption := range pollOptions {
		externalPollOptions = append(externalPollOptions, pollsService.FormatExternalPollOption(id, pollOption))
	}

	c.JSON(http.StatusOK, externalPollOptions)
}

func getPollOption(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	optionId, optionIdValidationError := utils.ValidateIntParam(c, "optionId")

	if optionIdValidationError != nil {
		return
	}

	pollOption, _, pollOptionNotFoundError := pollsService.GetPollOption(id, optionId)

	if pollOptionNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(pollOptionNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, pollsService.FormatExternalPollOption(id, *pollOption))
}

func createPollOption(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	var pollOptionBody schema.PollOption = schema.BlankPollOption()
	validationError := utils.ValidateBody(c, &pollOptionBody)

	if validationError != nil {
		return
	}

	createPollOption, createPollOptionError := pollsService.CreatePollOption(id, pollOptionBody)

	if createPollOptionError != nil {
		c.AbortWithStatusJSON(http.StatusConflict, utils.ResponseError(createPollOptionError.Error()))
		return
	}

	c.JSON(http.StatusCreated, pollsService.FormatExternalPollOption(id, *createPollOption))
}

func updatePollOption(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	optionId, optionIdValidationError := utils.ValidateIntParam(c, "optionId")

	if optionIdValidationError != nil {
		return
	}

	var pollOptionBody schema.PollOption = schema.BlankPollOption()
	validationError := utils.ValidateBody(c, &pollOptionBody)

	if validationError != nil {
		return
	}

	pollOption, pollOptionNotFoundError := pollsService.UpdatePollOption(id, optionId, pollOptionBody)

	if pollOptionNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(pollOptionNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, pollsService.FormatExternalPollOption(id, *pollOption))
}

func deletePollOption(c *gin.Context) {
	id, idValidationError := utils.ValidateIntParam(c, "id")

	if idValidationError != nil {
		return
	}

	optionId, optionIdValidationError := utils.ValidateIntParam(c, "optionId")

	if optionIdValidationError != nil {
		return
	}

	pollOptionNotFoundError := pollsService.DeletePollOption(id, optionId)

	if pollOptionNotFoundError != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.ResponseError(pollOptionNotFoundError.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.ResponseMessage("poll option deleted"))
}
