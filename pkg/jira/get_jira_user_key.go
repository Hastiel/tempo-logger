package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type GetJiraUserKeyRs struct {
	Key string `json:"key"`
}

func GetJiraUserKey(login string, password string, jiraUrl string, jiraTempoUserkeyUri string) (string, error) {
	params := url.Values{}
	params.Add("username", login)
	url := fmt.Sprintf("%s/%s?%s", jiraUrl, jiraTempoUserkeyUri, params.Encode())
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
		return "", err
	}

	var val GetJiraUserKeyRs

	if err = json.Unmarshal(data, &val); err != nil {
		return "", err
	}

	log.Println(string(data))
	return val.Key, nil
}
