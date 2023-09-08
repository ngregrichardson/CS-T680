package services

import (
	"errors"
	"fmt"
	"polls-api/cache"
	"polls-api/schema"
	"polls-api/utils"

	"github.com/gin-gonic/gin"
)

type PollsService struct {
	cache    *cache.Cache
	Hostname string
}

func NewPollsService() (*PollsService, error) {
	cache, err := cache.NewCache("polls")

	if err != nil {
		return nil, err
	}

	return &PollsService{
		cache:    cache,
		Hostname: "",
	}, nil
}

func (service *PollsService) GetPollPath(id uint) string {
	return fmt.Sprintf("/polls/%d", id)
}

func (service *PollsService) GetPollOptionSubPath(id uint) string {
	return fmt.Sprintf("/options/%d", id)
}

func (service *PollsService) FormatExternal(poll schema.Poll) gin.H {
	fullUrl := fmt.Sprintf("%s%s", service.Hostname, service.GetPollPath(poll.ID))

	externalOptions := make([]gin.H, 0)

	for _, option := range poll.Options {
		externalOptions = append(externalOptions, service.FormatExternalPollOption(poll.ID, option))
	}

	return gin.H{
		"id":       poll.ID,
		"title":    poll.Title,
		"question": poll.Question,
		"options":  externalOptions,
		"links":    utils.GenerateCRUDLinks(fullUrl),
	}
}

func (service *PollsService) FormatExternalPollOption(id uint, pollOption schema.PollOption) gin.H {
	fullUrl := fmt.Sprintf("%s%s%s", service.Hostname, service.GetPollPath(id), service.GetPollOptionSubPath(pollOption.ID))
	return gin.H{
		"id":    pollOption.ID,
		"title": pollOption.Title,
		"links": utils.GenerateCRUDLinks(fullUrl),
	}
}

/* Manage Polls */

func (service *PollsService) GetPolls() ([]schema.Poll, error) {
	polls := make([]schema.Poll, 0)

	pattern := fmt.Sprintf("%s*", service.cache.KeyPrefix)

	keys, err := service.cache.Client.Keys(service.cache.Client.Context(), pattern).Result()

	if err != nil {
		return polls, err
	}

	for _, key := range keys {
		var poll schema.Poll

		err := service.cache.GetItemFromRedis(key, &poll)

		if err != nil {
			continue
		}

		polls = append(polls, poll)
	}

	return polls, nil
}

func (service *PollsService) GetPoll(id uint) (schema.Poll, error) {
	var poll schema.Poll
	redisId := service.cache.KeyFromId(id)

	err := service.cache.GetItemFromRedis(redisId, &poll)

	if err != nil {
		return schema.Poll{}, errors.New("poll not found")
	}

	return poll, nil
}

func (service *PollsService) AddPoll(poll schema.Poll) (schema.Poll, error) {
	_, existingPollErr := service.GetPoll(poll.ID)

	if existingPollErr == nil {
		return schema.Poll{}, errors.New("poll already exists")
	}

	redisId := service.cache.KeyFromId(poll.ID)
	_, setErr := service.cache.JSONHelper.JSONSet(redisId, ".", poll)

	if setErr != nil {
		return schema.Poll{}, setErr
	}

	return poll, nil
}

func (service *PollsService) replacePoll(id uint, poll schema.Poll) (schema.Poll, error) {
	redisId := service.cache.KeyFromId(id)
	_, setErr := service.cache.JSONHelper.JSONSet(redisId, ".", poll)

	if setErr != nil {
		return schema.Poll{}, setErr
	}

	return poll, nil
}

func (service *PollsService) UpdatePoll(id uint, poll schema.Poll) (schema.Poll, error) {
	existingPoll, err := service.GetPoll(id)

	if err != nil {
		return schema.Poll{}, err
	}

	existingPoll.Title = poll.Title
	existingPoll.Question = poll.Question

	return service.replacePoll(existingPoll.ID, existingPoll)
}

func (service *PollsService) DeletePoll(id uint) error {
	existingPoll, existingPollErr := service.GetPoll(id)

	if existingPollErr != nil {
		return errors.New("poll does not exist")
	}

	redisId := service.cache.KeyFromId(existingPoll.ID)
	_, delErr := service.cache.Client.Del(service.cache.Context, redisId).Result()

	if delErr != nil {
		return delErr
	}

	return nil
}

/* Manage Poll Options */

func (service *PollsService) GetPollOptions(id uint) (schema.PollOptions, error) {
	poll, err := service.GetPoll(id)

	if err != nil {
		return nil, err
	}

	return poll.Options, nil
}

func (service *PollsService) GetPollOption(id uint, optionId uint) (*schema.PollOption, int, error) {
	poll, err := service.GetPoll(id)

	if err != nil {
		return &schema.PollOption{}, -1, err
	}

	return poll.GetPollOption(optionId)
}

func (service *PollsService) CreatePollOption(id uint, pollOption schema.PollOption) (*schema.PollOption, error) {
	poll, pollErr := service.GetPoll(id)

	if pollErr != nil {
		return &schema.PollOption{}, pollErr
	}

	_, _, existingPollOptionErr := poll.GetPollOption(pollOption.ID)

	if existingPollOptionErr == nil {
		return &schema.PollOption{}, errors.New("option already exists for that poll")
	}

	newPollOption, _ := poll.AddPollOption(pollOption)

	_, replaceErr := service.replacePoll(id, poll)

	if replaceErr != nil {
		return &schema.PollOption{}, replaceErr
	}

	return newPollOption, nil
}

func (service *PollsService) UpdatePollOption(id uint, optionId uint, pollOptionBody schema.PollOption) (*schema.PollOption, error) {
	poll, pollNotFoundError := service.GetPoll(id)

	if pollNotFoundError != nil {
		return &schema.PollOption{}, pollNotFoundError
	}

	_, pollOptionIndex, pollOptionError := poll.GetPollOption(optionId)

	if pollOptionError != nil {
		return &schema.PollOption{}, pollOptionError
	}

	pollOptionBody.ID = optionId

	poll.Options[pollOptionIndex] = pollOptionBody

	_, replaceErr := service.replacePoll(id, poll)

	if replaceErr != nil {
		return &schema.PollOption{}, replaceErr
	}

	return &poll.Options[pollOptionIndex], nil
}

func (service *PollsService) DeletePollOption(id uint, optionId uint) error {
	poll, pollNotFoundError := service.GetPoll(id)

	if pollNotFoundError != nil {
		return pollNotFoundError
	}

	_, pollOptionIndex, pollOptionNotFoundError := poll.GetPollOption(optionId)

	if pollOptionNotFoundError != nil {
		return pollOptionNotFoundError
	}

	poll.Options = append(poll.Options[:pollOptionIndex], poll.Options[pollOptionIndex+1:]...)

	_, replaceErr := service.replacePoll(id, poll)

	if replaceErr != nil {
		return replaceErr
	}

	return nil
}
