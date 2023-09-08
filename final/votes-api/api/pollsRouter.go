package api

import (
	"net/http"
	"votes-api/services"
	"votes-api/utils"

	"github.com/gin-gonic/gin"
)

var (
	pollsService *services.PollsService
)

func init() {
	var err error
	pollsService, err = services.NewPollsService()

	if err != nil {
		panic(err)
	}
}

func GetPollsRouter(group *gin.RouterGroup) {
	group.GET("/health", getPollsAPIHealth)
	group.GET("", getPolls)
}

func getPollsAPIHealth(c *gin.Context) {
	health, err := pollsService.GetPollsAPIHealth()

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, health)
}

func getPolls(c *gin.Context) {
	polls, err := pollsService.GetPolls()

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, polls)
}
