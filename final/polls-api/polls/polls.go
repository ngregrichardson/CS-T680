package polls

import (
	"errors"
	"fmt"
	"polls-api/cache"
	"polls-api/schema"

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

func (service *PollsService) FormatExternal(poll schema.Poll) gin.H {
	fullUrl := fmt.Sprintf("%s%s", service.Hostname, service.GetIdPath(poll.ID))
	return gin.H{
		"id":       poll.ID,
		"title":    poll.Title,
		"question": poll.Question,
		"options":  poll.Options,
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

func (service *PollsService) GetIdPath(id uint) string {
	return fmt.Sprintf("/polls/%d", id)
}

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

	return service.GetPoll(poll.ID)
}

func (service *PollsService) replacePoll(id uint, poll schema.Poll) (schema.Poll, error) {
	redisId := service.cache.KeyFromId(id)
	_, setErr := service.cache.JSONHelper.JSONSet(redisId, ".", poll)

	if setErr != nil {
		return schema.Poll{}, setErr
	}

	return service.GetPoll(id)
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
