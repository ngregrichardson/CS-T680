package schema

import "errors"

type Link struct {
	Method string `json:"method"`
	Url    string `json:"url"`
}

type PollOption struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type PollOptions []PollOption

type Poll struct {
	ID       uint        `json:"id"`
	Title    string      `json:"title"`
	Question string      `json:"question"`
	Options  PollOptions `json:"options"`
}

func BlankPoll() Poll {
	return Poll{
		Options: make(PollOptions, 0),
	}
}

func BlankPollOption() PollOption {
	return PollOption{}
}

func (p *Poll) GetPollOption(optionId uint) (*PollOption, int, error) {
	for i, vote := range p.Options {
		if vote.ID == optionId {
			return &vote, i, nil
		}
	}

	return &PollOption{}, 0, errors.New("vote record not found")
}

func (p *Poll) AddPollOption(pollOption PollOption) (*PollOption, error) {
	p.Options = append(p.Options, pollOption)

	return &pollOption, nil
}
