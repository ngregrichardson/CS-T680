package services

import (
	"errors"
	"fmt"
	"votes-api/cache"
	"votes-api/schema"
	"votes-api/utils"

	"github.com/gin-gonic/gin"
)

var (
	pollsService  *PollsService
	votersService *VotersService
)

func init() {
	var psErr, vsErr error
	pollsService, psErr = NewPollsService()

	if psErr != nil {
		panic(psErr)
	}

	votersService, vsErr = NewVotersService()

	if vsErr != nil {
		panic(vsErr)
	}
}

type VotesService struct {
	cache    *cache.Cache
	Hostname string
}

func NewVotesService() (*VotesService, error) {
	cache, err := cache.NewCache("votes")

	if err != nil {
		return nil, err
	}

	return &VotesService{
		cache:    cache,
		Hostname: "",
	}, nil
}

func (service *VotesService) GetVotePath(id uint) string {
	return fmt.Sprintf("/votes/%d", id)
}

func (service *VotesService) FormatExternalVote(vote schema.Vote, voteRelations schema.VoteRelations) gin.H {
	fullUrl := fmt.Sprintf("%s%s", service.Hostname, service.GetVotePath(vote.ID))

	return gin.H{
		"id":       vote.ID,
		"voterId":  vote.VoterID,
		"pollId":   vote.PollID,
		"optionId": vote.OptionID,
		"links":    utils.GenerateCRUDLinks(fullUrl),
		"poll":     voteRelations.Poll,
		"voter":    voteRelations.Voter,
		"option":   voteRelations.Option,
	}
}

func (service *VotesService) AddVote(vote schema.Vote) (schema.Vote, schema.VoteRelations, error) {
	voteRelations, relationErr := service.getVoteRelations(vote)

	if relationErr != nil {
		return schema.Vote{}, schema.VoteRelations{}, relationErr
	}

	redisId := service.cache.KeyFromId(vote.ID)
	_, setErr := service.cache.JSONHelper.JSONSet(redisId, ".", vote)

	if setErr != nil {
		return schema.Vote{}, schema.VoteRelations{}, setErr
	}

	_, voteRecordError := votersService.CreateVoteRecord(vote.VoterID, vote.PollID)

	if voteRecordError != nil {
		return schema.Vote{}, schema.VoteRelations{}, voteRecordError
	}

	return vote, voteRelations, nil
}

func (service *VotesService) getVoteRelations(vote schema.Vote) (schema.VoteRelations, error) {
	voter, voterError := votersService.GetVoter(vote.VoterID)

	if voterError != nil {
		return schema.VoteRelations{}, voterError
	}

	poll, pollError := pollsService.GetPoll(vote.PollID)

	if pollError != nil {
		return schema.VoteRelations{}, pollError
	}

	hasOption := false
	existingOption := schema.PollOption{}

	for _, option := range poll.Options {
		if option.ID == vote.OptionID {
			hasOption = true
			existingOption = option
			break
		}
	}

	if !hasOption {
		return schema.VoteRelations{}, fmt.Errorf("option %d does not exist for poll %d", vote.OptionID, vote.PollID)
	}

	return schema.VoteRelations{
		Poll:   poll,
		Voter:  voter,
		Option: existingOption,
	}, nil
}

func (service *VotesService) GetVote(id uint) (schema.Vote, schema.VoteRelations, error) {
	var vote schema.Vote
	redisId := service.cache.KeyFromId(id)

	err := service.cache.GetItemFromRedis(redisId, &vote)

	if err != nil {
		return schema.Vote{}, schema.VoteRelations{}, errors.New("poll not found")
	}

	voteRelations, relationErr := service.getVoteRelations(vote)

	if relationErr != nil {
		return schema.Vote{}, schema.VoteRelations{}, relationErr
	}

	return vote, voteRelations, nil
}

func (service *VotesService) replaceVote(id uint, vote schema.Vote) (schema.Vote, error) {
	redisId := service.cache.KeyFromId(id)
	_, setErr := service.cache.JSONHelper.JSONSet(redisId, ".", vote)

	if setErr != nil {
		return schema.Vote{}, setErr
	}

	return vote, nil
}

func (service *VotesService) UpdateVote(id uint, vote schema.Vote) (schema.Vote, schema.VoteRelations, error) {
	existingVote, voteRelations, voteNotFoundError := service.GetVote(id)

	if voteNotFoundError != nil {
		return schema.Vote{}, schema.VoteRelations{}, voteNotFoundError
	}

	existingOption, _, optionNotFoundError := voteRelations.Poll.GetPollOption(vote.OptionID)

	if optionNotFoundError != nil {
		return schema.Vote{}, schema.VoteRelations{}, optionNotFoundError
	}

	existingVote.OptionID = vote.OptionID
	voteRelations.Option = *existingOption

	res, replaceVoteError := service.replaceVote(existingVote.ID, existingVote)

	if replaceVoteError != nil {
		return schema.Vote{}, schema.VoteRelations{}, replaceVoteError
	}

	_, updateVoteRecordError := votersService.UpdateVoteRecord(existingVote.VoterID, existingVote.PollID)

	return res, voteRelations, updateVoteRecordError
}

func (service *VotesService) DeleteVote(id uint) error {
	existingVote, _, voteNotFoundError := service.GetVote(id)

	if voteNotFoundError != nil {
		return voteNotFoundError
	}

	redisId := service.cache.KeyFromId(existingVote.ID)
	_, delErr := service.cache.Client.Del(service.cache.Context, redisId).Result()

	if delErr != nil {
		return delErr
	}

	_, deleteVoteRecordError := votersService.DeleteVoteRecord(existingVote.VoterID, existingVote.PollID)

	return deleteVoteRecordError
}
