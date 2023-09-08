package schema

type Vote struct {
	ID       uint `json:"id"`
	VoterID  uint `json:"voterId"`
	PollID   uint `json:"pollId"`
	OptionID uint `json:"optionId"`
}

type VoteRelations struct {
	Poll   Poll       `json:"poll"`
	Voter  Voter      `json:"voter"`
	Option PollOption `json:"option"`
}

func BlankVote() Vote {
	return Vote{
		ID:       0,
		VoterID:  0,
		PollID:   0,
		OptionID: 0,
	}
}
