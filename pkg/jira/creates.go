package jira

type CreateRq struct {
	BillableSeconds  int    `json:"billableSeconds"`
	Comment          string `json:"comment"`
	EndDate          string `json:"endDate"`
	Started          string `json:"started"`
	OriginTaskId     string `json:"originTaskId"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
	Worker           string `json:"worker"`
}
