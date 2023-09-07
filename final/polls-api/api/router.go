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
	group.PATCH("/:id", updatePoll)
	group.DELETE("/:id", deletePoll)
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
