package jira

import "time"

type CreateRq struct {
	BillableSeconds  int    `json:"billableSeconds"`
	Comment          string `json:"comment"`
	EndDate          string `json:"endDate"`
	Started          string `json:"started"`
	OriginTaskId     string `json:"originTaskId"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
	Worker           string `json:"worker"`
}

type CreateParams struct {
	BillableSeconds  int
	Comment          string
	EndDate          time.Time
	Started          time.Time
	OriginTaskId     string
	TimeSpentSeconds int
}
