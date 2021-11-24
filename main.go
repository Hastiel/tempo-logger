package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tempo-loger/pkg/jira"

	//"tempo-loger/pkg/user"
	"time"

	"github.com/joho/godotenv"
)

type CreateRq struct {
	BillableSeconds  int    `json:"billableSeconds"`
	Comment          string `json:"comment"`
	EndDate          string `json:"endDate"`
	Started          string `json:"started"`
	OriginTaskId     string `json:"originTaskId"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
	Worker           string `json:"worker"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error while loading .env file")
	}

	targetSpentSeconds := 8 * 60 * 60

	jiraUrl := os.Getenv("JIRA_URL")
	jiraTempoCreatesUri := os.Getenv("JIRA_TEMPO_CREATES_URI")
	jiraTempoFindsUri := os.Getenv("JIRA_TEMPO_FINDS_URI")
	jiraTempoUserkeyUri := os.Getenv("JIRA_TEMPO_USERKEY_URI")
	//jiraUserKey := os.Getenv("JIRA_USER_KEY")
	login := os.Getenv("LOGIN")
	password := os.Getenv("PASSWORD")
	worklog := os.Getenv("WORKLOG")

	//user := user.NewUser(jiraUrl, login, password)

	// user := user.NewUser(jiraUrl, login, password)

	jiraUserKey, err := jira.GetJiraUserKey(login, password, jiraUrl, jiraTempoUserkeyUri)
	if err != nil {
		log.Println("Error while sending GetJiraUserKey-request")
	}

	worklogs := strings.Split(worklog, ";")

	year, month, day := time.Now().Date()
	currentDate := fmt.Sprintf("%d-%d-%d", year, month, day)

	findsRs, err := jira.Find(login, password, jiraUserKey, jiraUrl, jiraTempoFindsUri)
	if err != nil {
		log.Println("Error while sending Finds-request")
	}

	totalSpentSeconds := 0
	for i := 0; i < len(findsRs); i++ {
		//var item tempo_worklog.FindsItem = findsRs
		totalSpentSeconds += findsRs[i].TimeSpentSeconds
	}

	fmt.Println(totalSpentSeconds)

	for i := 0; i < len(worklogs); i++ {
		if totalSpentSeconds < targetSpentSeconds {
			randomSeconds := 0
			neededSpentSeconds := targetSpentSeconds - totalSpentSeconds
			if i < len(worklogs)-1 {
				min, max := 1, neededSpentSeconds/60/60
				rand.Seed(time.Now().UnixNano())
				if max-min <= 0 {
					randomSeconds = 1 * 60 * 60
				} else {
					randomSeconds = (rand.Intn(max-min) + min) * 60 * 60
				}
			}
			currentWorklog := strings.Split(worklogs[i], ",")

			origanalTaskId := strings.TrimSpace(currentWorklog[0])
			envTimeSpendHours, err := strconv.Atoi(strings.TrimSpace(currentWorklog[1]))
			if err != nil {
				log.Printf("Error while convert hours to int from .env file")
			}
			envTimeSpendSeconds := envTimeSpendHours * 60 * 60
			var timeSpendSeconds int
			if i == len(worklogs)-1 {
				timeSpendSeconds = neededSpentSeconds
			} else if randomSeconds > envTimeSpendSeconds {
				timeSpendSeconds = envTimeSpendSeconds
			} else {
				timeSpendSeconds = randomSeconds
			}
			comment := strings.TrimSpace(currentWorklog[2])

			createRq := CreateRq{
				BillableSeconds:  timeSpendSeconds,
				Comment:          comment,
				EndDate:          currentDate,
				Started:          currentDate,
				OriginTaskId:     origanalTaskId,
				TimeSpentSeconds: timeSpendSeconds,
				Worker:           jiraUserKey,
			}

			body, err := json.Marshal(createRq)
			if err != nil {
				log.Printf("Error while deserialize CreateWorklog")
			}

			url := fmt.Sprintf("%s/%s", jiraUrl, jiraTempoCreatesUri)
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
				fmt.Println("Error!")
			}
			fmt.Println(string(data))

			totalSpentSeconds += timeSpendSeconds

		}
	}

}
