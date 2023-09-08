package services

import (
	"fmt"
	"votes-api/schema"
	"votes-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type VotersService struct {
	client   *resty.Client
	Hostname string
}

func NewVotersService() (*VotersService, error) {
	visibleHostname := utils.GetEnvironmentVariable("VOTERS_VISIBLE_API_URL", "http://localhost:1081/voters")
	hostname := utils.GetEnvironmentVariable("VOTERS_API_URL", "http://voters-api:1081/voters")
	return &VotersService{
		client:   resty.New().SetBaseURL(hostname),
		Hostname: visibleHostname,
	}, nil
}

func (service *VotersService) GetVotersAPIHealth() (gin.H, error) {
	var result gin.H
	err := utils.GetRequestError(service.client.R().SetResult(&result).Get("/health"))

	if err != nil {
		return gin.H{}, err
	}

	return result, nil
}

func (service *VotersService) GetVoters() ([]schema.Voter, error) {
	var result []schema.Voter
	err := utils.GetRequestError(service.client.R().SetResult(&result).Get("/"))

	if err != nil {
		return []schema.Voter{}, err
	}

	return result, nil
}

func (service *VotersService) GetVoter(id uint) (schema.Voter, error) {
	var result schema.Voter
	err := utils.GetRequestError(service.client.R().SetResult(&result).Get(fmt.Sprintf("/%d", id)))

	if err != nil {
		return schema.Voter{}, err
	}

	return result, nil
}

func (service *VotersService) CreateVoteRecord(id uint, pollId uint) (schema.VoteRecord, error) {
	var result schema.VoteRecord
	err := utils.GetRequestError(service.client.R().SetResult(&result).SetBody(gin.H{
		"pollId": pollId,
	}).Post(fmt.Sprintf("/%d/votes", id)))

	if err != nil {
		return schema.VoteRecord{}, err
	}

	return result, nil
}

func (service *VotersService) UpdateVoteRecord(id uint, pollId uint) (schema.VoteRecord, error) {
	var result schema.VoteRecord
	err := utils.GetRequestError(service.client.R().SetResult(&result).Patch(fmt.Sprintf("/%d/votes/%d", id, pollId)))

	if err != nil {
		return schema.VoteRecord{}, err
	}

	return result, nil
}

func (service *VotersService) DeleteVoteRecord(id uint, pollId uint) (schema.VoteRecord, error) {
	var result schema.VoteRecord
	err := utils.GetRequestError(service.client.R().SetResult(&result).Delete(fmt.Sprintf("/%d/votes/%d", id, pollId)))

	if err != nil {
		return schema.VoteRecord{}, err
	}

	return result, nil
}
