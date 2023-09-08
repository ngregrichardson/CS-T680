package schema

type Link struct {
	Method string `json:"method"`
	Url    string `json:"url"`
}

type Links struct {
	Get    Link `json:"get"`
	Update Link `json:"update"`
	Delete Link `json:"delete"`
}
