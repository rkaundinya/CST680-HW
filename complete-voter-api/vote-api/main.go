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
	hostFlag    string
	portFlag    uint
	cacheURL    string
	voterAPIURL string
	pollAPIURL  string
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.StringVar(&cacheURL, "c", "0.0.0.0:6379", "Default cache location")
	flag.StringVar(&voterAPIURL, "v", "http://localhost:1080", "Default voter API location")
	flag.StringVar(&pollAPIURL, "papi", "http://localhost:2080", "Default poll API location")
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
	voterAPIURL = envVarOrDefault("VOTER_API_URL", voterAPIURL)
	pollAPIURL = envVarOrDefault("POLL_API_URL", pollAPIURL)
	hostFlag = envVarOrDefault("VOTEAPI_HOST", hostFlag)
	pfNew, err := strconv.Atoi(envVarOrDefault("VOTEAPI_PORT", fmt.Sprintf("%d", portFlag)))
	// only update port if env var converts to int successfully - else use default
	if err == nil {
		portFlag = uint(pfNew)
	}
}

func main() {
	setupParams()

	apiHandler, err := api.NewVoteApi(cacheURL, pollAPIURL, voterAPIURL)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/votes", apiHandler.GetVotes)
	r.GET("/votes/voter/:VoterID", apiHandler.GetVoter)

	r.POST("/votes", apiHandler.AddVoteJson)
	r.POST("/votes/voteID/:voteID/voterID/:voterID/pollID/:pollID/voteVal/:voteVal", apiHandler.AddVote)
	r.DELETE("votes/vote/:voteID", apiHandler.DeleteVote)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
