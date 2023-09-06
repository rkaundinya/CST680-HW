package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"vote-api/vote"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type VoteApi struct {
	db          *vote.VoteDB
	pollAPIURL  string
	voterAPIURL string
	apiClient   *resty.Client
}

func NewVoteApi(location string, inPollAPIURl string, inVoterAPIURL string) (*VoteApi, error) {
	dbHandler, err := vote.NewWithCacheInstance(location)
	if err != nil {
		return nil, err
	}

	apiClient := resty.New()

	return &VoteApi{
		db:          dbHandler,
		pollAPIURL:  inPollAPIURl,
		voterAPIURL: inVoterAPIURL,
		apiClient:   apiClient,
	}, nil
}

func AddVoteApi(location string) (*VoteApi, error) {
	dbHandler, err := vote.NewWithCacheInstance(location)
	if err != nil {
		return nil, err
	}

	return &VoteApi{
		db: dbHandler,
	}, nil
}

func (v *VoteApi) AddVoteJson(c *gin.Context) {
	var newVote vote.Vote

	if err := c.ShouldBindJSON(&newVote); err != nil {
		log.Println("error binding poll json", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	vID := newVote.VoterID

	_, err := v.db.GetVoter(string(vID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache with id: " + string(vID)})
		fmt.Println("Error getting voter")
		return
	}

	if err := v.db.AddVote(newVote); err != nil {
		fmt.Println("Error adding vote")
		log.Println("error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	fmt.Println("Successfully added vote for voter ID " + string(vID))
	c.JSON(http.StatusOK, newVote)
}

func (v *VoteApi) AddVote(c *gin.Context) {
	var newVote vote.Vote

	if err := c.ShouldBindJSON(&newVote); err != nil {
		log.Println("error binding poll json", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	vID := newVote.VoterID
	vIDString := strconv.FormatUint(uint64(vID), 10)
	fmt.Println("Voter ID: " + vIDString)

	_, err := v.db.GetVoter(vIDString)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache with id: " + string(vID)})
		fmt.Println("Error getting voter")
		return
	}

	pID := newVote.PollID
	pIDString := strconv.FormatUint(uint64(pID), 10)

	_, err = v.db.GetPoll(pIDString)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll in cache with id: " + string(pID)})
		fmt.Println("Error getting poll")
		return
	}

	if err := v.db.AddVote(newVote); err != nil {
		fmt.Println("Error adding vote")
		log.Println("error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	fmt.Println("Successfully added vote for voter ID " + string(vID))
	c.JSON(http.StatusOK, newVote)
}

func (p *VoteApi) GetVotes(c *gin.Context) {
	votes, err := p.db.GetVotes()
	if err != nil {
		log.Println("Error retrieving voters")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, votes)
}

func (v *VoteApi) DeleteVote(c *gin.Context) {
	vID := c.Param("voteID")

	vIDInt, err := strconv.ParseInt(vID, 10, 32)
	if err != nil {
		fmt.Println("vote ID int conversion failed")
		log.Println("Error converting vote id to int ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fmt.Print("Vote ID converted to int: ")
	fmt.Println(vIDInt)
	if err := v.db.DeleteVote(int(vIDInt)); err != nil {
		log.Println("failed to delete voter with ID " + fmt.Sprint(vID))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}

func (v *VoteApi) GetVoter(c *gin.Context) {
	vID := c.Param("VoterID")
	if vID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No voter id provided"})
		return
	}

	voter, err := v.db.GetVoter(vID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache with id=" + vID})
		return
	}

	c.JSON(http.StatusOK, voter)
}
