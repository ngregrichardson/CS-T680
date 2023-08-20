package voters

import (
	"errors"
	"time"
)

type pollVote struct {
	PollID   uint      `json:"pollId"`
	VoteDate time.Time `json:"voteDate"`
}

type VoteHistory []pollVote

type Voter struct {
	VoterID     uint        `json:"voterId"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	VoteHistory VoteHistory `json:"votes"`
}

func NewVoter(id uint, firstName string, lastName string) (Voter, error) {
	voter := Voter{
		VoterID:     id,
		FirstName:   firstName,
		LastName:    lastName,
		VoteHistory: make([]pollVote, 0),
	}

	return voter, nil
}

func (v *Voter) GetPollVote(pollId uint) (*pollVote, int, error) {
	for i, vote := range v.VoteHistory {
		if vote.PollID == pollId {
			return &vote, i, nil
		}
	}

	return &pollVote{}, 0, errors.New("Vote not found")
}

func (v *Voter) AddPollVote(pollId uint) (*pollVote, error) {
	newPollVote := pollVote{
		PollID:   pollId,
		VoteDate: time.Now(),
	}

	v.VoteHistory = append(v.VoteHistory, newPollVote)

	return &newPollVote, nil
}
