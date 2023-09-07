package voters

import (
	"errors"
	"fmt"
	"time"
	"voters-api/cache"
	"voters-api/schema"
	"voters-api/utils"

	"github.com/gin-gonic/gin"
)

type VotersService struct {
	cache       *cache.Cache
	Hostname    string
	PollsApiUrl string
}

func NewVotersService() (*VotersService, error) {
	cache, err := cache.NewCache("voters")

	if err != nil {
		return nil, err
	}

	return &VotersService{
		cache:       cache,
		Hostname:    "",
		PollsApiUrl: utils.GetEnvironmentVariable("POLLS_API_URL", "http://localhost:8080"),
	}, nil
}

func (service *VotersService) FormatExternalVoter(voter schema.Voter) gin.H {
	fullUrl := fmt.Sprintf("%s%s", service.Hostname, service.GetVoterPath(voter.ID))
	return gin.H{
		"id":        voter.ID,
		"firstName": voter.FirstName,
		"lastName":  voter.LastName,
		"votes":     voter.VoteHistory,
		"links": gin.H{
			"get": schema.Link{
				Method: "GET",
				Url:    fullUrl,
			},
			"update": schema.Link{
				Method: "PATCH",
				Url:    fullUrl,
			},
			"delete": schema.Link{
				Method: "DELETE",
				Url:    fullUrl,
			},
		},
	}
}

func (service *VotersService) FormatExternalVoteRecord(id uint, voteRecord schema.VoteRecord) gin.H {
	pollUrl := fmt.Sprintf("%s/polls/%d", service.PollsApiUrl, voteRecord.PollID)
	fullUrl := fmt.Sprintf("%s%s%s", service.Hostname, service.GetVoterPath(id), service.GetVoteRecordSubPath(voteRecord.PollID))
	return gin.H{
		"pollId":    voteRecord.PollID,
		"pollLinks": utils.GenerateCRUDLinks(pollUrl),
		"voteDate":  voteRecord.VoteDate,
		"links":     utils.GenerateCRUDLinks(fullUrl),
	}
}

func (service *VotersService) GetVoterPath(id uint) string {
	return fmt.Sprintf("/voters/%d", id)
}

func (service *VotersService) GetVoteRecordSubPath(id uint) string {
	return fmt.Sprintf("/votes/%d", id)
}

/* Manage Voters */

func (service *VotersService) GetVoters() ([]schema.Voter, error) {
	voters := make([]schema.Voter, 0)

	pattern := fmt.Sprintf("%s*", service.cache.KeyPrefix)

	keys, err := service.cache.Client.Keys(service.cache.Client.Context(), pattern).Result()

	if err != nil {
		return voters, err
	}

	for _, key := range keys {
		var voter schema.Voter

		err := service.cache.GetItemFromRedis(key, &voter)

		if err != nil {
			continue
		}

		voters = append(voters, voter)
	}

	return voters, nil
}

func (service *VotersService) GetVoter(id uint) (schema.Voter, error) {
	var voter schema.Voter
	redisId := service.cache.KeyFromId(id)

	err := service.cache.GetItemFromRedis(redisId, &voter)

	if err != nil {
		return schema.Voter{}, errors.New("voter not found")
	}

	return voter, nil
}

func (service *VotersService) AddVoter(voter schema.Voter) (schema.Voter, error) {
	_, existingVoterError := service.GetVoter(voter.ID)

	if existingVoterError == nil {
		return schema.Voter{}, errors.New("voter already exists")
	}

	redisId := service.cache.KeyFromId(voter.ID)
	_, setErr := service.cache.JSONHelper.JSONSet(redisId, ".", voter)

	if setErr != nil {
		return schema.Voter{}, setErr
	}

	return service.GetVoter(voter.ID)
}

func (service *VotersService) replaceVoter(id uint, voter schema.Voter) (schema.Voter, error) {
	redisId := service.cache.KeyFromId(id)
	_, setErr := service.cache.JSONHelper.JSONSet(redisId, ".", voter)

	if setErr != nil {
		return schema.Voter{}, setErr
	}

	return service.GetVoter(id)
}

func (service *VotersService) UpdateVoter(id uint, voter schema.Voter) (schema.Voter, error) {
	existingVoter, err := service.GetVoter(id)

	if err != nil {
		return schema.Voter{}, err
	}

	existingVoter.FirstName = voter.FirstName
	existingVoter.LastName = voter.LastName

	return service.replaceVoter(existingVoter.ID, existingVoter)
}

func (service *VotersService) DeleteVoter(id uint) error {
	existingVoter, existingVoterError := service.GetVoter(id)

	if existingVoterError != nil {
		return errors.New("voter does not exist")
	}

	redisId := service.cache.KeyFromId(existingVoter.ID)
	_, delErr := service.cache.Client.Del(service.cache.Context, redisId).Result()

	if delErr != nil {
		return delErr
	}

	return nil
}

/* Manage Vote History */

func (service *VotersService) GetVoterHistory(id uint) (schema.VoteHistory, error) {
	voter, err := service.GetVoter(id)

	if err != nil {
		return nil, err
	}

	return voter.VoteHistory, nil
}

func (service *VotersService) GetVoteRecord(id uint, pollId uint) (*schema.VoteRecord, int, error) {
	voter, err := service.GetVoter(id)

	if err != nil {
		return &schema.VoteRecord{}, -1, err
	}

	return voter.GetVoteRecord(pollId)
}

func (service *VotersService) CreateVoteRecord(id uint, pollId uint) (*schema.VoteRecord, error) {
	voter, voterErr := service.GetVoter(id)

	if voterErr != nil {
		return &schema.VoteRecord{}, voterErr
	}

	_, _, voteErr := voter.GetVoteRecord(pollId)

	if voteErr == nil {
		return &schema.VoteRecord{}, errors.New("vote already exists for that poll")
	}

	newVote, _ := voter.AddVoteRecord(pollId)

	_, replaceErr := service.replaceVoter(id, voter)

	if replaceErr != nil {
		return &schema.VoteRecord{}, replaceErr
	}

	return newVote, nil
}

func (service *VotersService) UpdateVoteRecord(id uint, pollId uint) (*schema.VoteRecord, error) {
	voter, voterErr := service.GetVoter(id)

	if voterErr != nil {
		return &schema.VoteRecord{}, voterErr
	}

	vote, voteIndex, voteErr := voter.GetVoteRecord(pollId)

	if voteErr != nil {
		return &schema.VoteRecord{}, voteErr
	}

	vote.VoteDate = time.Now()

	voter.VoteHistory[voteIndex] = *vote

	voter.VoteHistory = utils.ShiftEnd(voter.VoteHistory, voteIndex)

	_, replaceErr := service.replaceVoter(id, voter)

	if replaceErr != nil {
		return &schema.VoteRecord{}, replaceErr
	}

	return vote, nil
}

func (service *VotersService) DeleteVoteRecord(id uint, pollId uint) error {
	voter, voterErr := service.GetVoter(id)

	if voterErr != nil {
		return voterErr
	}

	_, voteIndex, voteErr := voter.GetVoteRecord(pollId)

	if voteErr != nil {
		return voteErr
	}

	voter.VoteHistory = append(voter.VoteHistory[:voteIndex], voter.VoteHistory[voteIndex+1:]...)

	_, replaceErr := service.replaceVoter(id, voter)

	if replaceErr != nil {
		return replaceErr
	}

	return nil
}
