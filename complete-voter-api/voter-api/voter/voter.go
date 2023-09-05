package voter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type voterPoll struct {
	PollID   uint
	VoteDate time.Time
}

type Voter struct {
	VoterID     uint
	FirstName   string
	LastName    string
	VoteHistory []voterPoll
}

type VoterList struct {
	Voters map[uint]Voter //A map of VoterIDs as keys and Voter structs as values
}

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

type VoterDB struct {
	voterList VoterList

	//Redis cache connections
	cache
}

func New() (*VoterDB, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*VoterDB, error) {
	//connect to redis
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	//context used to coordinate between our go code and redis operations
	ctx := context.Background()

	//Recommended way to ensure our redis connection is working
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error() + "cache might not be availble, continuing...")
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	//return pointer to new VoterDB struct
	return &VoterDB{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

func isRedisNilError(err error) bool {
	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
}

// Redis stores keys with prefix defined in RedisKeyPrefix
func redisKeyFromId(id int) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

// Helper to get voter from voterlist given key
func (vDB *VoterDB) getItemFromRedis(key string, voterItem *Voter) error {
	itemObject, err := vDB.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(itemObject.([]byte), voterItem)
	if err != nil {
		return err
	}

	return nil
}

// constructor for VoterList struct
func NewVoter(id uint, fn, ln string) *Voter {
	return &Voter{
		FirstName:   fn,
		LastName:    ln,
		VoteHistory: []voterPoll{},
	}
}

func (v *VoterDB) AddVoter(newVoter Voter) error {
	// Check if voter with id already exists
	redisKey := redisKeyFromId(int(newVoter.VoterID))
	var existingVoter Voter
	if err := v.getItemFromRedis(redisKey, &existingVoter); err == nil {
		return errors.New("item already exists")
	}

	// Add item to database with JSON set
	if _, err := v.jsonHelper.JSONSet(redisKey, ".", newVoter); err != nil {
		return err
	}

	// Return nil if everything is working fine
	return nil
}

func (v *VoterDB) DeleteVoter(vID uint) error {
	redisKey := redisKeyFromId(int(vID))
	var existingVoter Voter
	if err := v.getItemFromRedis(redisKey, &existingVoter); err != nil {
		return errors.New("no voter with ID " + fmt.Sprint(vID) + "exists to delete")
	}

	if _, err := v.jsonHelper.JSONDel(redisKey, "."); err != nil {
		return err
	}

	return nil
}

func (v *VoterDB) GetVoter(vID uint) (*Voter, error) {
	redisKey := redisKeyFromId(int(vID))
	var existingVoter Voter
	if err := v.getItemFromRedis(redisKey, &existingVoter); err != nil {
		return nil, errors.New("No voter with ID exists")
	}

	return &existingVoter, nil
}

func (v *VoterDB) GetVoters() ([]Voter, error) {
	var voter Voter
	var voterList []Voter

	//Query redis for all items
	pattern := RedisKeyPrefix + "*"
	ks, _ := v.cacheClient.Keys(v.context, pattern).Result()
	for _, key := range ks {
		err := v.getItemFromRedis(key, &voter)
		if err != nil {
			return nil, err
		}
		voterList = append(voterList, voter)
	}

	return voterList, nil
}

// Currently only updates polls if user inputs new voter with at least 1 poll
// Else keeps polls of original voter
func (v *VoterDB) UpdateVoter(voter Voter) error {
	redisKey := redisKeyFromId(int(voter.VoterID))
	var existingVoter Voter

	if err := v.getItemFromRedis(redisKey, &existingVoter); err != nil {
		return errors.New("no existing voter with ID")
	}

	if _, err := v.jsonHelper.JSONSet(redisKey, ".FirstName", voter.FirstName); err != nil {
		return err
	}

	if _, err := v.jsonHelper.JSONSet(redisKey, ".LastName", voter.LastName); err != nil {
		return err
	}

	if len(voter.VoteHistory) != 0 {
		if _, err := v.jsonHelper.JSONSet(redisKey, ".VoteHistory", voter.VoteHistory); err != nil {
			return err
		}
	}

	return nil
}

func (v *VoterDB) GetPoll(voterID uint, pollID uint) (voterPoll, error) {
	existingVoter, err := v.GetVoter(voterID)
	if err != nil {
		return voterPoll{}, err
	}

	pollIdx := -1

	for idx, poll := range existingVoter.VoteHistory {
		if poll.PollID == pollID {
			pollIdx = idx
			break
		}
	}

	if pollIdx == -1 {
		return voterPoll{}, errors.New("Failed to find poll")
	}

	return existingVoter.VoteHistory[pollIdx], nil
}

func (v *VoterDB) GetVoterPolls(voterID uint) ([]voterPoll, error) {
	voter, err := v.GetVoter(voterID)
	if err != nil {
		return nil, err
	}

	return voter.VoteHistory, nil
}

func (v *VoterDB) AddPoll(vID uint, pollID uint) error {
	existingVoter, err := v.GetVoter(vID)
	if err != nil {
		return errors.New("Could not add poll to voter " + string(vID))
	}

	for _, poll := range existingVoter.VoteHistory {
		if poll.PollID == uint(pollID) {
			return errors.New("INVALID - Trying to add duplicate poll ID")
		}
	}

	//TODO - Update this to use JSONArrAppend with relative path that works
	redisKey := redisKeyFromId(int(vID))
	existingVoter.VoteHistory = append(existingVoter.VoteHistory, voterPoll{PollID: pollID, VoteDate: time.Now()})

	if _, err := v.jsonHelper.JSONSet(redisKey, ".", existingVoter); err != nil {
		return err
	}

	return nil
}

func (v *VoterDB) UpdatePoll(vID uint, pollID uint) error {
	existingVoter, err := v.GetVoter(vID)
	if err != nil {
		return errors.New("Could not add poll to voter " + fmt.Sprint(vID))
	}

	pollIdx := -1

	for idx, poll := range existingVoter.VoteHistory {
		if poll.PollID == uint(pollID) {
			pollIdx = idx
			break
		}
	}

	if pollIdx == -1 {
		return errors.New("no existing poll with ID " + fmt.Sprint(pollID))
	}

	redisKey := redisKeyFromId(int(vID))

	// TODO - haven't figured out how to error hanlde this properly
	v.jsonHelper.JSONArrPop(redisKey, ".VoteHistory", pollIdx)

	if _, err := v.jsonHelper.JSONArrAppend(redisKey, ".VoteHistory", voterPoll{PollID: pollID, VoteDate: time.Now()}); err != nil {
		return err
	}

	return nil
}

func (v *VoterDB) DeletePoll(vID uint, pollID uint) error {
	existingVoter, err := v.GetVoter(vID)
	if err != nil {
		return errors.New("Could not add poll to voter " + fmt.Sprint(vID))
	}

	pollIdx := -1

	for idx, poll := range existingVoter.VoteHistory {
		if poll.PollID == uint(pollID) {
			pollIdx = idx
			break
		}
	}

	if pollIdx == -1 {
		return errors.New("no existing poll with ID " + fmt.Sprint(pollID))
	}

	redisKey := redisKeyFromId(int(vID))
	v.jsonHelper.JSONArrPop(redisKey, ".VoteHistory", pollIdx)

	return nil
}

func (v *Voter) ToJson() string {
	b, _ := json.Marshal(v)
	return string(b)
}
