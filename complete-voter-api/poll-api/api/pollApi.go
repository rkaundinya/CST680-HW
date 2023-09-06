package api

import (
	"fmt"
	"log"
	"net/http"
	"poll-api/poll"
	"strconv"

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

func (p *PollApi) AddPollOption(c *gin.Context) {
	pollID := c.Param("pollID")
	optionID := c.Param("optionID")
	optDescription := c.Param("description")

	pollIDuint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	poll, err := p.db.GetPoll(uint(pollIDuint))
	if err != nil {
		log.Println("Error finding poll with id ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	optionIDUint, err := strconv.ParseUint(optionID, 10, 32)
	if err != nil {
		log.Println("Error converting poll option id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = p.db.AddPollOption(uint(pollIDuint), uint(optionIDUint), optDescription)
	if err != nil {
		log.Println("Error adding poll option ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, poll)
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

func (p *PollApi) DeletePollOption(c *gin.Context) {
	pollID := c.Param("pollID")
	optionID := c.Param("optionID")

	pollIDuint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	poll, err := p.db.GetPoll(uint(pollIDuint))
	if err != nil {
		log.Println("Error finding poll with id ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	optionIDUint, err := strconv.ParseUint(optionID, 10, 32)
	if err != nil {
		log.Println("Error converting poll option id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = p.db.DeletePollOption(uint(pollIDuint), uint(optionIDUint))
	if err != nil {
		log.Println("Error deleting poll option ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (p *PollApi) GetPoll(c *gin.Context) {
	pollID := c.Param("pollID")

	pollIDuint, err := strconv.ParseUint(pollID, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to uint ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	poll, err := p.db.GetPoll(uint(pollIDuint))

	if err != nil {
		log.Println("Error getting poll")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (p *PollApi) DeletePoll(c *gin.Context) {
	pID := c.Param("pollID")

	pIDInt, err := strconv.ParseInt(pID, 10, 32)
	if err != nil {
		fmt.Println("vote ID int conversion failed")
		log.Println("Error converting vote id to int ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := p.db.DeletePoll(int(pIDInt)); err != nil {
		log.Println("failed to delete voter with ID " + fmt.Sprint(pID))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}
