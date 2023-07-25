package main

import (
	"flag"
	"fmt"
	"voter-api/server"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	server.Init(fmt.Sprintf("%s:%d", hostFlag, portFlag))
}
