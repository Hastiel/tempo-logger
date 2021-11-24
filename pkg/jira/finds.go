package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type FindsRq struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Worker []string `json:"worker"`
}

type FindsRs []struct {
	BillableSeconds  int `json:"billableSeconds"`
	TimeSpentSeconds int `json:"timeSpentSeconds"`
}

var year, month, day = time.Now().Date()
var currentDate = fmt.Sprintf("%d-%d-%d", year, month, day)

func Find(login string, password string, jiraUserKey string, jiraUrl string, jiraTempoFindsUri string) (FindsRs, error) {
	findsRq := FindsRq{
		From:   currentDate,
		To:     currentDate,
		Worker: []string{jiraUserKey},
	}

	body, err := json.Marshal(findsRq)
	if err != nil {
		log.Println("Cannot marshal body for FindsRq")
	}

	url := fmt.Sprintf("%s/%s", jiraUrl, jiraTempoFindsUri)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Error!")
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(login, password)
	hc := &http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		fmt.Println("Error!")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	log.Println(string(data))

	var val FindsRs

	if err = json.Unmarshal(data, &val); err != nil {
		return nil, err
	}

	return val, nil
}
