package schema

type Link struct {
	Method string `json:"method"`
	Url    string `json:"url"`
}

type pollOption struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type PollOptions []pollOption

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
