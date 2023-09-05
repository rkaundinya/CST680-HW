package main

import (
	"complete-voter-api/voterApi"
	"flag"
	"fmt"
	"os"
	"strconv"

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
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

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
	hostFlag = envVarOrDefault("VOTERAPI_HOST", hostFlag)
	pfNew, err := strconv.Atoi(envVarOrDefault("VOTERAPI_PORT", fmt.Sprintf("%d", portFlag)))
	// only update port if env var converts to int successfully - else use default
	if err == nil {
		portFlag = uint(pfNew)
	}
}

func main() {
	setupParams()

	apiHandler, err := voterApi.NewVoterApi()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/voters", apiHandler.GetVoterListJson)
	r.GET("/voters/:voterID", apiHandler.GetVoterJson)
	r.GET("/voters/:voterID/polls", apiHandler.GetVoterPollsJson)
	r.GET("/voters/:voterID/polls/:pollID", apiHandler.GetPollJson)
	r.GET("/voters/health", apiHandler.HealthCheck)

	r.POST("/voters", apiHandler.AddVoter)
	r.POST("/voters/:voterID/firstName/:firstName/lastName/:lastName", apiHandler.AddVoter)
	r.POST("/voters/:voterID/polls/:pollID", apiHandler.AddPoll)

	r.PUT("/voters/:voterID", apiHandler.UpdateVoter)
	r.PUT("/voters/:voterID/polls/:pollID", apiHandler.UpdatePoll)
	r.DELETE("/voters/:voterID", apiHandler.DeleteVoter)
	r.DELETE("/voters/:voterID/polls/:pollID", apiHandler.DeletePoll)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
