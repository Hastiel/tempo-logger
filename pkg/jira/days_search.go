package jira

type DaysSearchRq struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	UserKeys []string `json:"userKeys"`
}

type DaysSearchRs []struct {
	User string `json:"user"`
	Days []struct {
		Date            string `json:"date"`
		DayOpen         bool   `json:"dayOpen"`
		RequiredSeconds int    `json:"requiredSeconds"`
		Type            string `json:"type"`
	} `json:"days"`
}
