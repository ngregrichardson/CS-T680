package main

import (
	"flag"
	"fmt"
	"net/http"
	"votes-api/api"
	"votes-api/middleware"
	"votes-api/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
	tlsFlag  bool
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default port")
	flag.BoolVar(&tlsFlag, "t", false, "Enable TLS")

	flag.Parse()
}

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.Stats())
	router.Use(cors.Default())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, utils.ResponseError("route not found"))
	})

	api.GetVotesRouter(router.Group("/votes"), utils.FormatHostname(portFlag, tlsFlag))

	api.GetPollsRouter(router.Group("/polls"))

	api.GetVotersRouter(router.Group("/voters"))

	return router
}

func main() {
	processCmdLineFlags()

	r := NewRouter()

	r.Run(fmt.Sprintf("%s:%d", hostFlag, portFlag))
}
