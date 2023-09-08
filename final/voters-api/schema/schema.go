package schema

import (
	"errors"
	"time"
)

type VoteRecord struct {
	PollID   uint      `json:"pollId"`
	VoteDate time.Time `json:"voteDate"`
}

type VoteHistory []VoteRecord

type Voter struct {
	ID          uint        `json:"id"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	VoteHistory VoteHistory `json:"votes"`
}

func BlankVoter() Voter {
	return Voter{
		VoteHistory: make(VoteHistory, 0),
	}
}

func BlankVoteRecord() VoteRecord {
	return VoteRecord{
		PollID: 0,
	}
}

func (v *Voter) GetVoteRecord(pollId uint) (*VoteRecord, int, error) {
	for i, vote := range v.VoteHistory {
		if vote.PollID == pollId {
			return &vote, i, nil
		}
	}

	return &VoteRecord{}, 0, errors.New("vote record not found")
}

func (v *Voter) AddVoteRecord(pollId uint) (*VoteRecord, error) {
	newPollVote := VoteRecord{
		PollID:   pollId,
		VoteDate: time.Now(),
	}

	v.VoteHistory = append(v.VoteHistory, newPollVote)

	return &newPollVote, nil
}
