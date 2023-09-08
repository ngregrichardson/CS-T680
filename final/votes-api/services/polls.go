package services

import (
	"fmt"
	"votes-api/schema"
	"votes-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type PollsService struct {
	client   *resty.Client
	Hostname string
}

func NewPollsService() (*PollsService, error) {
	visibleHostname := utils.GetEnvironmentVariable("POLLS_VISIBLE_API_URL", "http://localhost:1080/polls")
	hostname := utils.GetEnvironmentVariable("POLLS_API_URL", "http://polls-api:1080/polls")
	return &PollsService{
		client:   resty.New().SetBaseURL(hostname),
		Hostname: visibleHostname,
	}, nil
}

func (service *PollsService) GetPollsAPIHealth() (gin.H, error) {
	var result gin.H
	err := utils.GetRequestError(service.client.R().SetResult(&result).Get("/health"))

	if err != nil {
		return gin.H{}, err
	}

	return result, nil
}

func (service *PollsService) GetPolls() ([]schema.Poll, error) {
	var result []schema.Poll
	err := utils.GetRequestError(service.client.R().SetResult(&result).Get("/"))

	if err != nil {
		return []schema.Poll{}, err
	}

	return result, nil
}

func (service *PollsService) GetPoll(id uint) (schema.Poll, error) {
	var result schema.Poll
	err := utils.GetRequestError(service.client.R().SetResult(&result).Get(fmt.Sprintf("/%d", id)))

	if err != nil {
		return schema.Poll{}, err
	}

	return result, nil
}
