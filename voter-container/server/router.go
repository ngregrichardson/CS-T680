package server

import (
	"net/http"
	"voter-api/api"
	"voter-api/api/responses"
	"voter-api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.Stats())
	router.Use(cors.Default())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, responses.NewHttpNotFoundResponseWithMessage("Not Found"))
	})

	v1Group := router.Group("/v1")

	api.GetVoterRouter(v1Group)

	return router
}
