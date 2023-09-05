package vote

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
	RedisKeyPrefix = "vote:"
)

type cache struct {
	client  *redis.Client
	helper  *rejson.Handler
	context context.Context
}

type Vote struct {
	VoteID    uint
	VoterID   uint
	PollID    uint
	VoteValue uint
}

type VoteDB struct {
	//Redis cache connections
	cache
	//API connections
	voterAPIURL string
	pollAPIURL  string
}

func NewWithCacheInstance(location string) (*VoteDB, error) {
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
	return &VoteDB{
		cache: cache{
			client:  client,
			helper:  jsonHelper,
			context: ctx,
		},
	}, nil
}

func NewVote(voteID uint, voterID uint, pollID uint, voteVal uint) *Vote {
	return &Vote{
		VoteID:    voteID,
		VoterID:   voterID,
		PollID:    pollID,
		VoteValue: voteVal,
	}
}

func (v *VoteDB) AddVote(newVote Vote) error {
	//Check if vote with id already exists
	redisKey := RedisKeyFromId(int(newVote.VoteID), RedisKeyPrefix)
	var existingVote Vote
	if err := v.getItemFromRedis(redisKey, &existingVote); err == nil {
		return errors.New("item already exists")
	}

	//Add item to database with JSON set
	if _, err := v.cache.helper.JSONSet(redisKey, ".", newVote); err != nil {
		return err
	}

	// Return nil if everything is working fine
	return nil
}

func (v *VoteDB) GetVotes() ([]Vote, error) {
	var vote Vote
	var voteList []Vote

	//Query redis for all items
	pattern := RedisKeyPrefix + "*"
	ks, _ := v.client.Keys(v.context, pattern).Result()
	for _, key := range ks {
		err := v.getItemFromRedis(key, &vote)
		if err != nil {
			return nil, err
		}
		voteList = append(voteList, vote)
	}

	return voteList, nil
}

func (vDB *VoteDB) getItemFromRedis(key string, voteItem *Vote) error {
	voteObj, err := vDB.cache.helper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(voteObj.([]byte), voteItem)
	if err != nil {
		return err
	}

	return nil
}

func RedisKeyFromId(id int, prefix string) string {
	return fmt.Sprintf("%s%d", prefix, id)
}
