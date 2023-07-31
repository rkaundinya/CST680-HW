package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"voter-api-starter/voter"

	"github.com/gin-gonic/gin"
)

type VoterApi struct {
	voterList voter.VoterList
}

// TODO make more robust error handling
func NewVoterApi() (*VoterApi, error) {
	return &VoterApi{
		voterList: voter.VoterList{
			Voters: make(map[uint]voter.Voter),
		},
	}, nil
}

func (v *VoterApi) AddVoter(c *gin.Context) {
	var newVoter voter.Voter

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding voter json", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, ok := v.voterList.Voters[newVoter.VoterID]
	if ok {
		log.Println("Warning - trying to add an already existing voter")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	v.voterList.Voters[newVoter.VoterID] = newVoter
}

func (v *VoterApi) AddVoterJson(c *gin.Context) {
	var newVoter voter.Voter

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, ok := v.voterList.Voters[newVoter.VoterID]
	if ok {
		log.Println("Warning - trying to add an already existing voter")
		return
	}

	v.voterList.Voters[newVoter.VoterID] = newVoter
}

func (v *VoterApi) DeleteVoter(c *gin.Context) {
	voterID := c.Param("voterID")

	voterIDuint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	delete(v.voterList.Voters, uint(voterIDuint))
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

	voter, ok := v.voterList.Voters[uint(voterIDuint)]

	if !ok {
		log.Println("Voter not found in voter list")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	for idx, poll := range voter.VoteHistory {
		if poll.PollID == uint(pollIDuint) {
			voter.VoteHistory[idx] = voter.VoteHistory[len(voter.VoteHistory)-1]
			voter.VoteHistory = voter.VoteHistory[:len(voter.VoteHistory)-1]
			v.voterList.Voters[uint(voterIDuint)] = voter
			return
		}
	}

	log.Println("No poll with ID found to delete")
	c.AbortWithStatus(http.StatusNotFound)
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

	voter, ok := v.voterList.Voters[uint(voterIDuint)]

	if !ok {
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

	voter.AddPoll(uint(pollIDuint))

	//Update voter list with up to date voter object
	v.voterList.Voters[uint(voterIDuint)] = voter
}

func (v *VoterApi) UpdateVoter(c *gin.Context) {
	var newVoter voter.Voter

	if err := c.ShouldBindJSON(&newVoter); err != nil {
		log.Println("Error binding voter json", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, ok := v.voterList.Voters[newVoter.VoterID]
	if !ok {
		log.Println("INVALID - no voter with id exists")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	v.voterList.Voters[newVoter.VoterID] = newVoter
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

	voter, ok := v.voterList.Voters[uint(voterIDuint)]

	if !ok {
		log.Println("Voter not found in voter list")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	for idx, poll := range voter.VoteHistory {
		if poll.PollID == uint(pollIDuint) {
			voter.VoteHistory[idx].VoteDate = time.Now()
			log.Println("Updated vote time")
			return
		}
	}

	log.Println("INVALID - no poll found with ID")
	c.AbortWithStatus(http.StatusNotFound)
}

func (v *VoterApi) GetVoterJson(c *gin.Context) {
	voterID := c.Param("voterID")

	voterIDuint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, ok := v.voterList.Voters[uint(voterIDuint)]

	if !ok {
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

	voter, ok := v.voterList.Voters[uint(voterIDuint)]

	if !ok {
		log.Println("Voter not found in voter list")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	pollIdx := -1

	for idx, poll := range voter.VoteHistory {
		if poll.PollID == uint(pollIDuint) {
			pollIdx = idx
			break
		}
	}

	if pollIdx == -1 {
		log.Println("Poll ID not found in voter polls")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter.VoteHistory[pollIdx])
}

func (v *VoterApi) GetVoterPollsJson(c *gin.Context) {
	voterID := c.Param("voterID")

	voterIDuint, err := strconv.ParseUint(voterID, 10, 32)
	if err != nil {
		log.Println("Error converting voter id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voter, ok := v.voterList.Voters[uint(voterIDuint)]

	if !ok {
		log.Println("Voter not found in voter list")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter.VoteHistory)
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
	c.JSON(http.StatusOK, v.voterList.Voters)
}
