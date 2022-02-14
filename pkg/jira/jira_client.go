package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type JiraClient interface {
	Find(date time.Time) (FindsRs, error)
	Create(comment, originalTaskId string, currentDate time.Time, timeSpentSeconds int) error
	GetDayInfo(date time.Time) (DaysSearchRs, error)
}

type jiraClient struct {
	login      string
	pwd        string
	key        string
	host       string
	findPath   string
	keyPath    string
	createPath string
	daysPath   string
}

func New(login, pwd, host, findPath, keyPath, createPath, daysPath string) JiraClient {
	return &jiraClient{
		login:      login,
		pwd:        pwd,
		host:       host,
		findPath:   findPath,
		keyPath:    keyPath,
		createPath: createPath,
		daysPath:   daysPath,
	}
}

func (j *jiraClient) Create(comment, originalTaskId string, currentDate time.Time, timeSpentSeconds int) error {
	key, err := j.getUserKey()
	if err != nil {
		return err
	}

	dateStr := currentDate.Format("2006-01-02")

	rqStruct := CreateRq{
		BillableSeconds:  timeSpentSeconds,
		Comment:          comment,
		EndDate:          dateStr,
		Started:          dateStr,
		OriginTaskId:     originalTaskId,
		TimeSpentSeconds: timeSpentSeconds,
		Worker:           key,
	}

	body, err := json.Marshal(rqStruct)
	if err != nil {
		return err
	}

	rqUrl := fmt.Sprintf("%s/%s", j.host, j.createPath)
	req, err := http.NewRequest(http.MethodPost, rqUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(j.login, j.pwd)
	hc := &http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("unknown error: %s", string(data))
	}

	return nil
}

func (j *jiraClient) Find(date time.Time) (FindsRs, error) {
	key, err := j.getUserKey()
	if err != nil {
		return nil, err
	}

	dateStr := date.Format("2006-01-02")

	findsRq := FindsRq{
		From:   dateStr,
		To:     dateStr,
		Worker: []string{key},
	}

	body, err := json.Marshal(findsRq)
	if err != nil {
		log.Println("Cannot marshal body for FindsRq")
		return nil, err
	}

	rqUrl := fmt.Sprintf("%s/%s", j.host, j.findPath)
	req, err := http.NewRequest(http.MethodPost, rqUrl, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(j.login, j.pwd)
	hc := &http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var val FindsRs
	if err = json.Unmarshal(data, &val); err != nil {
		return nil, err
	}

	return val, nil
}

func (j *jiraClient) GetDayInfo(date time.Time) (DaysSearchRs, error) {
	var rsStruct DaysSearchRs

	key, err := j.getUserKey()
	if err != nil {
		return rsStruct, err
	}

	dateStr := date.Format("2006-01-02")
	rqStruct := DaysSearchRq{
		From:     dateStr,
		To:       dateStr,
		UserKeys: []string{key},
	}

	body, err := json.Marshal(rqStruct)
	if err != nil {
		return rsStruct, err
	}

	rqUrl := fmt.Sprintf("%s/%s", j.host, j.daysPath)
	req, err := http.NewRequest(http.MethodPost, rqUrl, bytes.NewBuffer(body))
	if err != nil {
		return rsStruct, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(j.login, j.pwd)
	hc := &http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		return rsStruct, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return rsStruct, err
	}

	if err = json.Unmarshal(data, &rsStruct); err != nil {
		return rsStruct, err
	}

	return rsStruct, nil
}

func (j *jiraClient) getUserKey() (string, error) {
	if len(j.key) > 0 {
		return j.key, nil
	}

	params := url.Values{}
	params.Add("username", j.login)
	rqUrl := fmt.Sprintf("%s/%s?%s", j.host, j.keyPath, params.Encode())

	req, err := http.NewRequest(http.MethodGet, rqUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(j.login, j.pwd)
	hc := &http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var val GetJiraUserKeyRs

	if err = json.Unmarshal(data, &val); err != nil {
		return "", err
	}

	j.key = val.Key
	return val.Key, nil
}
