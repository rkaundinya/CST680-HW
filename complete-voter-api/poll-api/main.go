package main

import (
	"flag"
	"fmt"
	"os"
	"poll-api/api"
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
	flag.UintVar(&portFlag, "p", 2080, "Default Port")

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

	apiHandler, err := api.NewPollApi(cacheURL)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/polls", apiHandler.GetPolls)

	r.POST("/polls", apiHandler.AddPoll)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
