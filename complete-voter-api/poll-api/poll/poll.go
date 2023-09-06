package poll

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

const (
	RedisKeyPrefix = "poll:"
)

type cache struct {
	client  *redis.Client
	helper  *rejson.Handler
	context context.Context
}

type pollOption struct {
	PollOptionID   uint
	PollOptionText string
}

type Poll struct {
	PollID       uint
	PollTitle    string
	PollQuestion string
	PollOptions  []pollOption
}

type PollDB struct {
	//Redis cache connections
	cache
	votesAPIURL string
}

func NewWithCacheInstance(location string) (*PollDB, error) {
	//connect to redis
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	//context used to coordinate between go code and redis operations
	ctx := context.Background()

	//Recommended way to ensure redis connection is working
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error() + "cache might noe be avaialble, continuing...")
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	//return pointer to new PollDB struct
	return &PollDB{
		cache: cache{
			client:  client,
			helper:  jsonHelper,
			context: ctx,
		},
	}, nil
}

func NewPoll(pollID uint, title string, question string, options []pollOption) *Poll {
	return &Poll{
		PollID:       pollID,
		PollTitle:    title,
		PollQuestion: question,
		PollOptions:  options,
	}
}

func (p *PollDB) AddPoll(newPoll Poll) error {
	//Check if poll with id already exists
	redisKey := RedisKeyFromId(int(newPoll.PollID), RedisKeyPrefix)
	var existingPoll Poll
	if err := p.getItemFromRedis(redisKey, &existingPoll); err == nil {
		return errors.New("item already exists")
	}

	//Add item to database with JSON set
	if _, err := p.cache.helper.JSONSet(redisKey, ".", newPoll); err != nil {
		return err
	}

	// Return nil if everything is working fine
	return nil
}

func (p *PollDB) AddPollOption(pollID uint, optionId uint, body string) error {
	//Check if poll with id already exists
	redisKey := RedisKeyFromId(int(pollID), RedisKeyPrefix)
	var existingPoll Poll

	if err := p.getItemFromRedis(redisKey, &existingPoll); err != nil {
		fmt.Println("failed to get poll from redis")
		return errors.New("poll does not exist")
	}

	for _, option := range existingPoll.PollOptions {
		if option.PollOptionID == optionId {
			fmt.Println("Trying to re-add existing poll option - not allowed")
			return errors.New("poll option with id already exists")
		}
	}

	if err := p.getItemFromRedis(redisKey, &existingPoll); err != nil {
		fmt.Println("Error finding poll")
		return err
	}

	option := pollOption{PollOptionID: optionId, PollOptionText: body}
	existingPoll.PollOptions = append(existingPoll.PollOptions, option)
	fmt.Print("printing poll options")
	fmt.Println(existingPoll.PollOptions[0])

	if _, err := p.cache.helper.JSONSet(redisKey, ".", existingPoll); err != nil {
		fmt.Println("Error adding poll option")
		return err
	}

	return nil
}

func (p *PollDB) DeletePollOption(pollID uint, optionID uint) error {
	redisKey := RedisKeyFromId(int(pollID), RedisKeyPrefix)
	var existingPoll Poll

	if err := p.getItemFromRedis(redisKey, &existingPoll); err != nil {
		fmt.Println("failed to get poll from redis")
		return errors.New("poll does not exist")
	}

	optionIdx := -1

	for idx, option := range existingPoll.PollOptions {
		if option.PollOptionID == optionID {
			optionIdx = idx
			break
		}
	}

	if optionIdx == -1 {
		return errors.New("no poll option with ID " + fmt.Sprint(optionID))
	}

	p.helper.JSONArrPop(redisKey, ".PollOptions", optionIdx)

	return nil
}

func (p *PollDB) GetPolls() ([]Poll, error) {
	var poll Poll
	var voterList []Poll

	//Query redis for all items
	pattern := RedisKeyPrefix + "*"
	ks, _ := p.client.Keys(p.context, pattern).Result()
	for _, key := range ks {
		err := p.getItemFromRedis(key, &poll)
		if err != nil {
			return nil, err
		}
		voterList = append(voterList, poll)
	}

	return voterList, nil
}

func (p *PollDB) GetPoll(pollID uint) (Poll, error) {
	polls, err := p.GetPolls()
	if err != nil {
		return Poll{}, err
	}

	pollIdx := -1

	for idx, poll := range polls {
		if poll.PollID == pollID {
			pollIdx = idx
			break
		}
	}

	if pollIdx == -1 {
		return Poll{}, errors.New("Failed to find poll")
	}

	return polls[pollIdx], nil
}

func (pDB *PollDB) DeletePoll(pID int) error {
	redisKey := RedisKeyFromId(pID, RedisKeyPrefix)
	var existingPoll Poll
	if err := pDB.getItemFromRedis(redisKey, &existingPoll); err != nil {
		return errors.New("no voter with ID " + fmt.Sprint(pID) + "exists to delete")
	}

	if _, err := pDB.cache.helper.JSONDel(redisKey, "."); err != nil {
		return err
	}

	return nil
}

func (pDB *PollDB) getItemFromRedis(key string, pollItem *Poll) error {
	pollObj, err := pDB.cache.helper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(pollObj.([]byte), pollItem)
	if err != nil {
		return err
	}

	return nil
}

func RedisKeyFromId(id int, prefix string) string {
	return fmt.Sprintf("%s%d", prefix, id)
}
