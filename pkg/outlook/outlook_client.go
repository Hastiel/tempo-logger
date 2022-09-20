package outlook

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Azure/go-ntlmssp"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type OutlookClient interface {
	GetEvents(startDate, endDate time.Time) (EventsRs, error)
}

type outlookClient struct {
	login     string
	pwd       string
	host      string
	eventPath string
}

func New(login, pwd, host, eventPath string) OutlookClient {
	return &outlookClient{
		login:     login,
		pwd:       pwd,
		host:      host,
		eventPath: eventPath,
	}
}

func (o *outlookClient) GetEvents(startDate, endDate time.Time) (EventsRs, error) {
	layout := "2006-01-02T15:04:05"
	var eventsRs EventsRs
	url := fmt.Sprintf("%s/%s?startDateTime=%s&endDateTime=%s&$select=Subject,ResponseStatus,Body,Start,End",
		o.host, o.eventPath, startDate.Format(layout), endDate.Format(layout))
	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{
				TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
			},
		},
	}

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.SetBasicAuth(o.login, o.pwd)
	res, err := client.Do(req)
	if err != nil {
		return eventsRs, err
	}
	rsBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if err = json.Unmarshal(rsBody, &eventsRs); err != nil {
		return eventsRs, err
	}
	bodyString := string(rsBody)
	log.Println(bodyString)
	return eventsRs, nil
}
