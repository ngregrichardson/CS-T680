package schema

import "time"

type VoteRecord struct {
	PollID    uint      `json:"pollId"`
	VoteDate  time.Time `json:"voteDate"`
	Links     Links     `json:"links"`
	PollLinks Links     `json:"pollLinks"`
}

type VoteHistory []VoteRecord

type Voter struct {
	ID          uint        `json:"id"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	VoteHistory VoteHistory `json:"votes"`
	Links       Links       `json:"links"`
}
