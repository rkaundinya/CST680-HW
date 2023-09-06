package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"vote-api/schema"
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

	var voters = []schema.Voter{}
	votersPath := v.voterAPIURL + "/voters"

	_, err := v.apiClient.R().SetResult(&voters).Get(votersPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache with id: " + string(vID)})
		fmt.Println("Error getting voter v2")
		return
	}

	// Check if voter with ID exists
	var foundVoterID bool = false
	for _, voter := range voters {
		if voter.VoterID == vID {
			foundVoterID = true
			break
		}
	}

	// Early exit if no matching voter
	if !foundVoterID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache with id: " + string(vID)})
		fmt.Println("Error getting voter v2")
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

	var voters = []schema.Voter{}
	votersPath := v.voterAPIURL + "/voters"

	_, err := v.apiClient.R().SetResult(&voters).Get(votersPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache with id: " + string(vID)})
		fmt.Println("Error getting voter v2")
		return
	}

	// Check if voter with ID exists
	var foundVoterID bool = false
	for _, voter := range voters {
		if voter.VoterID == vID {
			foundVoterID = true
			break
		}
	}

	// Early exit if no matching voter
	if !foundVoterID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find voter in cache with id: " + string(vID)})
		fmt.Println("Error getting voter v2")
		return
	}

	pID := newVote.PollID
	optID := newVote.VoteValue

	var polls = []schema.Poll{}
	pollsPath := v.pollAPIURL + "/polls"

	_, err = v.apiClient.R().SetResult(&polls).Get(pollsPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll in cache with id: " + string(pID)})
		fmt.Println("Error getting poll")
		return
	}

	//Check if poll with ID and poll option with ID exist
	var foundPollID bool = false
	var foundPollOptID bool = false
	for _, poll := range polls {
		if poll.PollID == pID {
			foundPollID = true
			for _, option := range poll.PollOptions {
				if option.PollOptionID == optID {
					foundPollOptID = true
				}
			}
			break
		}
	}

	//Early exit if no matching poll
	if !foundPollID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll in cache with id: " + string(pID)})
		fmt.Println("Error getting poll: " + strconv.FormatUint(uint64(pID), 32))
		return
	}

	if !foundPollOptID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find poll option in cache with id: " + string(optID)})
		fmt.Println("Error getting poll option: " + strconv.FormatUint(uint64(optID), 32))
		return
	}

	//Otherwise safe to attempt adding vote
	if err := v.db.AddVote(newVote); err != nil {
		fmt.Println("Error adding vote")
		log.Println("error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

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

func (p *VoteApi) GetVote(c *gin.Context) {
	voteID := c.Param("voteID")

	voteIDuint, err := strconv.ParseUint(voteID, 10, 32)
	if err != nil {
		log.Println("Error converting vote id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	vote, err := p.db.GetVote(uint(voteIDuint))
	if err != nil {
		log.Println("Failed to get vote with ID " + string(voteID))
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, vote)
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
