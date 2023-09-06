package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"voter-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Using flag driven CLI for now
var (
	hostFlag   string
	portFlag   uint
	cacheURL   string
	voteAPIURL string
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.StringVar(&cacheURL, "c", "0.0.0.0:6379", "Default cache location")
	flag.StringVar(&voteAPIURL, "v", "http://localhost:3080", "Default vote api location")
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
	voteAPIURL = envVarOrDefault("VOTE_API_URL", voteAPIURL)
	hostFlag = envVarOrDefault("VOTERAPI_HOST", hostFlag)
	pfNew, err := strconv.Atoi(envVarOrDefault("VOTERAPI_PORT", fmt.Sprintf("%d", portFlag)))
	// only update port if env var converts to int successfully - else use default
	if err == nil {
		portFlag = uint(pfNew)
	}
}

func main() {
	setupParams()

	apiHandler, err := api.NewVoterApi(cacheURL, voteAPIURL)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/voters", apiHandler.GetVoterListJson)
	r.GET("/voters/:voterID", apiHandler.GetVoterJson)
	r.GET("/voters/health", apiHandler.HealthCheck)

	r.POST("/voters", apiHandler.AddVoter)
	r.POST("/voters/:voterID/firstName/:firstName/lastName/:lastName", apiHandler.AddVoter)

	r.PUT("/voters/:voterID", apiHandler.UpdateVoter)
	r.DELETE("/voters/:voterID", apiHandler.DeleteVoter)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
