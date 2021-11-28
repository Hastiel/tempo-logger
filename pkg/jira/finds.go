package jira

type FindsRq struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Worker []string `json:"worker"`
}

type FindsRs []struct {
	BillableSeconds  int `json:"billableSeconds"`
	TimeSpentSeconds int `json:"timeSpentSeconds"`
}
