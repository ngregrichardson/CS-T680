package api

import (
	"net/http"
	"votes-api/services"
	"votes-api/utils"

	"github.com/gin-gonic/gin"
)

var (
	votersService *services.VotersService
)

func init() {
	var err error
	votersService, err = services.NewVotersService()

	if err != nil {
		panic(err)
	}
}

func GetVotersRouter(group *gin.RouterGroup) {
	group.GET("/health", getVotersAPIHealth)
	group.GET("", getVoters)
}

func getVotersAPIHealth(c *gin.Context) {
	health, err := votersService.GetVotersAPIHealth()

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, health)
}

func getVoters(c *gin.Context) {
	voters, err := votersService.GetVoters()

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, voters)
}
