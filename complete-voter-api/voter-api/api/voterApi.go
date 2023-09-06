package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"voter-api/voter"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type VoterApi struct {
	db         *voter.VoterDB
	voteAPIURL string
	apiClient  *resty.Client
}

// TODO make more robust error handling
func NewVoterApi(location string, inVoteApiURL string) (*VoterApi, error) {
	dbHandler, err := voter.NewWithCacheInstance(location)
	if err != nil {
		return nil, err
	}

	apiClient := resty.New()

	return &VoterApi{
		db:         dbHandler,
		voteAPIURL: inVoteApiURL,
		apiClient:  apiClient,
	}, nil
}

func (v *VoterApi) AddVoter(c *gin.Context) {
	var newVoter voter.Voter

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding voter json", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.db.AddVoter(newVoter); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.JSON(http.StatusOK, newVoter)
}

func (v *VoterApi) DeleteVoter(c *gin.Context) {
	voterID := c.Param("voterID")

	voterIDuint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.db.DeleteVoter(uint(voterIDuint)); err != nil {
		log.Println("failed to delete voter with ID " + fmt.Sprint(voterID))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

func (v *VoterApi) DeletePoll(c *gin.Context) {
	voterID := c.Param("voterID")
	pollID := c.Param("pollID")

	voterIDuint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollIDuint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = v.db.DeletePoll(uint(voterIDuint), uint(pollIDuint)); err != nil {
		log.Println("Failed to delete poll with ID " + fmt.Sprint(pollID) + " for voter ID " + fmt.Sprint(voterID))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

func (v *VoterApi) AddPoll(c *gin.Context) {
	voterID := c.Param("voterID")
	pollID := c.Param("pollID")

	voterIDuint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollIDuint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, err := v.db.GetVoter(uint(voterIDuint))

	if err != nil {
		log.Println("Voter not found in voter list")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	for _, poll := range voter.VoteHistory {
		if poll.PollID == uint(pollIDuint) {
			log.Println("INVALID - Trying to add duplicate poll ID")
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}

	err = v.db.AddPoll(uint(voterIDuint), uint(pollIDuint))
	if err != nil {
		log.Println("Failed adding poll with ID " + string(pollIDuint))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

func (v *VoterApi) UpdateVoter(c *gin.Context) {
	var newVoter voter.Voter

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding voter json", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := v.db.UpdateVoter(newVoter)

	if err != nil {
		log.Println("Failed to update voter with ID " + fmt.Sprint(newVoter.VoterID))
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func (v *VoterApi) UpdatePoll(c *gin.Context) {
	voterID := c.Param("voterID")
	pollID := c.Param("pollID")

	voterIDuint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollIDuint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = v.db.UpdatePoll(uint(voterIDuint), uint(pollIDuint))

	if err != nil {
		log.Println("Failed to update poll with ID " + string(pollID) + " for voter ID " + string(voterID))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func (v *VoterApi) GetVoterJson(c *gin.Context) {
	voterID := c.Param("voterID")

	voterIDuint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, err := v.db.GetVoter(uint(voterIDuint))

	if err != nil {
		log.Println("Voter not found in voter list")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter)
}

func (v *VoterApi) GetPollJson(c *gin.Context) {
	voterID := c.Param("voterID")
	pollID := c.Param("pollID")

	voterIDuint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollIDuint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	poll, err := v.db.GetPoll(uint(voterIDuint), uint(pollIDuint))

	if err != nil {
		log.Println("Error getting poll")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (v *VoterApi) GetVoterPollsJson(c *gin.Context) {
	voterID := c.Param("voterID")

	voterIDuint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	polls, err := v.db.GetVoterPolls(uint(voterIDuint))

	if err != nil {
		log.Println("Voter not found in voter list")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, polls)
}

func (v *VoterApi) HealthCheck(c *gin.Context) {
	log.Println("Health check received")
	c.JSON(http.StatusOK,
		gin.H{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             100,
			"users_processed":    1000,
			"errors_encountered": 10,
		})
}

// Need to make this gin compatible --- see todo-api code!
func (v *VoterApi) GetVoterListJson(c *gin.Context) {
	voters, err := v.db.GetVoters()
	if err != nil {
		log.Println("Error retrieving voters")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, voters)
}
