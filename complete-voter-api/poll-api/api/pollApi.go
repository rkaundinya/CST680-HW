package api

import (
	"log"
	"net/http"
	"poll-api/poll"

	"github.com/gin-gonic/gin"
)

type PollApi struct {
	db *poll.PollDB
}

func NewPollApi(location string) (*PollApi, error) {
	dbHandler, err := poll.NewWithCacheInstance(location)
	if err != nil {
		return nil, err
	}

	return &PollApi{
		db: dbHandler,
	}, nil
}

func (p *PollApi) AddPoll(c *gin.Context) {
	var newPoll poll.Poll

	if err := c.ShouldBindJSON(&newPoll); err != nil {
		log.Println("error binding poll json", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := p.db.AddPoll(newPoll); err != nil {
		log.Println("error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.JSON(http.StatusOK, newPoll)
}

func (p *PollApi) GetPolls(c *gin.Context) {
	polls, err := p.db.GetPolls()
	if err != nil {
		log.Println("Error retrieving voters")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, polls)
}
