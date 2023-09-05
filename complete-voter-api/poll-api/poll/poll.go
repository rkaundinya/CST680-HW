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
