package main

import (
	"flag"
	"fmt"
	"os"
	"voter-api-starter/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Using flag driven CLI for now
var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	apiHandler, err := api.NewVoterApi()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/voter-api", apiHandler.GetVoterListJson)
	r.GET("voter-api/voters/:voterID", apiHandler.GetVoterJson)
	r.GET("voter-api/voters/:voterID/polls", apiHandler.GetVoterPollsJson)
	r.GET("voter-api/voters/:voterID/polls/:pollID", apiHandler.GetPollJson)
	r.GET("voter-api/voters/health", apiHandler.HealthCheck)

	r.POST("/voter-api", apiHandler.AddVoterJson)
	r.POST("/voter-api/voters/:voterID/firstName/:firstName/lastName/:lastName", apiHandler.AddVoter)
	r.POST("/voter-api/voters/:voterID/polls/:pollID", apiHandler.AddPoll)

	r.PUT("voter-api/voters/:voterID", apiHandler.UpdateVoter)
	r.PUT("voter-api/voters/:voterID/polls/:pollID", apiHandler.UpdatePoll)
	r.DELETE("voter-api/voters/:voterID", apiHandler.DeleteVoter)
	r.DELETE("voter-api/voters/:voterID/polls/:pollID", apiHandler.DeletePoll)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
