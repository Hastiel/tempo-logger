package outlook

import "time"

type EventsRs struct {
	Value []struct {
		Subject        string `json:"Subject"`
		ResponseStatus struct {
			Response string    `json:"Response"`
			Time     time.Time `json:"Time"`
		} `json:"ResponseStatus"`
		Body struct {
			ContentType string `json:"ContentType"`
			Content     string `json:"Content"`
		} `json:"Body"`
		Start struct {
			DateTime string `json:"DateTime"`
			TimeZone string `json:"TimeZone"`
		} `json:"Start"`
		End struct {
			DateTime string `json:"DateTime"`
			TimeZone string `json:"TimeZone"`
		} `json:"End"`
	} `json:"value"`
}
