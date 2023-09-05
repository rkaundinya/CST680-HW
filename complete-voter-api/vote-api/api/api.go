package api

import (
	"log"
	"net/http"
	"vote-api/vote"

	"github.com/gin-gonic/gin"
)

type VoteApi struct {
	db *vote.VoteDB
}

func NewVoteApi(location string) (*VoteApi, error) {
	dbHandler, err := vote.NewWithCacheInstance(location)
	if err != nil {
		return nil, err
	}

	return &VoteApi{
		db: dbHandler,
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

func (p *VoteApi) AddVote(c *gin.Context) {
	var newVote vote.Vote

	if err := c.ShouldBindJSON(&newVote); err != nil {
		log.Println("error binding poll json", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := p.db.AddVote(newVote); err != nil {
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
