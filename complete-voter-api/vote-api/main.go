package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"vote-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Using flag driven CLI for now
var (
	hostFlag string
	portFlag uint
	cacheURL string
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.StringVar(&cacheURL, "c", "0.0.0.0:6379", "Default cache location")
	flag.UintVar(&portFlag, "p", 3080, "Default Port")

	flag.Parse()
}

func envVarOrDefault(envVar string, defaultVal string) string {
	envVal := os.Getenv(envVar)
	if envVal != "" {
		return envVal
	}

	return defaultVal
}

func setupParams() {
	//process command line flags
	processCmdLineFlags()

	//process env variables
	cacheURL = envVarOrDefault("REDIS_URL", cacheURL)
	hostFlag = envVarOrDefault("VOTEAPI_HOST", hostFlag)
	pfNew, err := strconv.Atoi(envVarOrDefault("VOTEAPI_PORT", fmt.Sprintf("%d", portFlag)))
	// only update port if env var converts to int successfully - else use default
	if err == nil {
		portFlag = uint(pfNew)
	}
}

func main() {
	setupParams()

	apiHandler, err := api.NewVoteApi(cacheURL)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/votes", apiHandler.GetVotes)

	r.POST("/votes", apiHandler.AddVote)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
