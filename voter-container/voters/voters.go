package voters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"voter-api/voters/utils"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "redis:6379"
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

type VoterCache struct {
	cache
}

func NewVoterCache() (*VoterCache, error) {
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}

	return newVoterCacheInstance(redisUrl)
}

func newVoterCacheInstance(url string) (*VoterCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: url,
	})

	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	return &VoterCache{
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

func redisKeyFromId(id uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (vc *VoterCache) getItemFromRedis(key string, item *Voter) error {
	itemObject, err := vc.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(itemObject.([]byte), item)
	if err != nil {
		return err
	}

	return nil
}

func (vc *VoterCache) AddVoter(id uint, voter Voter) (Voter, error) {
	_, err := vc.GetVoter(id)

	if err == nil {
		return Voter{}, errors.New("Voter already exists")
	}

	voter.VoterID = id

	redisId := redisKeyFromId(id)
	_, setErr := vc.jsonHelper.JSONSet(redisId, ".", voter)

	if setErr != nil {
		return Voter{}, setErr
	}

	return voter, nil
}

func (vc *VoterCache) GetVoter(id uint) (Voter, error) {
	var voter Voter
	redisId := redisKeyFromId(id)

	err := vc.getItemFromRedis(redisId, &voter)

	if err != nil {
		return Voter{}, errors.New("Voter not found")
	}

	return voter, nil
}

func (vc *VoterCache) replaceVoter(id uint, voter Voter) (Voter, error) {
	redisId := redisKeyFromId(id)
	_, setErr := vc.jsonHelper.JSONSet(redisId, ".", voter)

	if setErr != nil {
		return Voter{}, setErr
	}

	return voter, nil
}

func (vc *VoterCache) UpdateVoter(id uint, voter Voter) (Voter, error) {
	existingVoter, err := vc.GetVoter(id)

	if err != nil {
		return Voter{}, err
	}

	existingVoter.FirstName = voter.FirstName
	existingVoter.LastName = voter.LastName

	return vc.replaceVoter(id, existingVoter)
}

func (vc *VoterCache) DeleteVoter(id uint) error {
	_, err := vc.GetVoter(id)

	if err != nil {
		return err
	}

	_, deleteErr := vc.cacheClient.Del(vc.context, redisKeyFromId(id)).Result()

	if deleteErr != nil {
		return errors.New("Voter not found")
	}

	return nil
}

func (vc *VoterCache) GetVoterHistory(id uint) (VoteHistory, error) {
	voter, err := vc.GetVoter(id)

	if err != nil {
		return nil, err
	}

	return voter.VoteHistory, nil
}

func (vc *VoterCache) GetVoterVote(id uint, pollId uint) (*pollVote, int, error) {
	voter, err := vc.GetVoter(id)

	if err != nil {
		return &pollVote{}, 0, err
	}

	return voter.GetPollVote(pollId)
}

func (vc *VoterCache) CreateVoterVote(id uint, pollId uint) (*pollVote, error) {
	voter, voterErr := vc.GetVoter(id)

	if voterErr != nil {
		return &pollVote{}, voterErr
	}

	_, _, voteErr := voter.GetPollVote(pollId)

	if voteErr == nil {
		return &pollVote{}, errors.New("Vote for that poll already exists")
	}

	newVote, _ := voter.AddPollVote(pollId)

	_, replaceErr := vc.replaceVoter(id, voter)

	if replaceErr != nil {
		return &pollVote{}, replaceErr
	}

	return newVote, nil
}

func (vc *VoterCache) UpdateVoterVote(id uint, pollId uint) (*pollVote, error) {
	voter, voterErr := vc.GetVoter(id)

	if voterErr != nil {
		return &pollVote{}, voterErr
	}

	vote, voteIndex, voteErr := voter.GetPollVote(pollId)

	if voteErr != nil {
		return &pollVote{}, voteErr
	}

	vote.VoteDate = time.Now()

	voter.VoteHistory[voteIndex] = *vote

	voter.VoteHistory = utils.ShiftEnd(voter.VoteHistory, voteIndex)

	_, replaceErr := vc.replaceVoter(id, voter)

	if replaceErr != nil {
		return &pollVote{}, replaceErr
	}

	return vote, nil
}

func (vc *VoterCache) DeleteVoterVote(id uint, pollId uint) error {
	voter, voterErr := vc.GetVoter(id)

	if voterErr != nil {
		return voterErr
	}

	_, voteIndex, voteErr := voter.GetPollVote(pollId)

	if voteErr != nil {
		return voteErr
	}

	voter.VoteHistory = append(voter.VoteHistory[:voteIndex], voter.VoteHistory[voteIndex+1:]...)

	_, replaceErr := vc.replaceVoter(id, voter)

	if replaceErr != nil {
		return replaceErr
	}

	return nil
}

func (v *VoterCache) GetVoters() ([]Voter, error) {
	voters := make([]Voter, 0)

	pattern := fmt.Sprintf("%s*", RedisKeyPrefix)

	keys, err := v.cacheClient.Keys(v.context, pattern).Result()

	if err != nil {
		return voters, err
	}

	for _, key := range keys {
		var voter Voter
		err := v.getItemFromRedis(key, &voter)

		if err != nil {
			return voters, err
		}

		voters = append(voters, voter)
	}

	return voters, nil
}
