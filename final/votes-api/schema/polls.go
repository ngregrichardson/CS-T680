package schema

import "errors"

type PollOption struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Links Links  `json:"links"`
}

type PollOptions []PollOption

type Poll struct {
	ID       uint        `json:"id"`
	Title    string      `json:"title"`
	Question string      `json:"question"`
	Options  PollOptions `json:"options"`
	Links    Links       `json:"links"`
}

func (p *Poll) GetPollOption(optionId uint) (*PollOption, int, error) {
	for i, vote := range p.Options {
		if vote.ID == optionId {
			return &vote, i, nil
		}
	}

	return &PollOption{}, 0, errors.New("poll option not found")
}
