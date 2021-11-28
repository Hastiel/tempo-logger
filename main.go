package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"tempo-loger/pkg/enviroment"
	"tempo-loger/pkg/jira"
	"time"
)

func main() {
	env, err := enviroment.NewEnviroment()
	if err != nil {
		log.Println("Error while loading .env file")
		log.Println(err)
		return
	}

	jiraClient := jira.New(env.Login, env.Password, env.JiraUrl, env.JiraTempoFindsUri, env.JiraTempoUserkeyUri, env.JiraTempoCreatesUri, env.JiraTempoDaysSearch)

	daysRs, err := jiraClient.GetDayInfo(time.Now())
	if err != nil {
		log.Println("Error while sending GetDayInfo-request")
		fmt.Println(err)
		return
	}

	if daysRs[0].Days[0].Type != "WORKING_DAY" {
		log.Println("Today is not-working day!")
		return
	}
	targetSpentSeconds := daysRs[0].Days[0].RequiredSeconds

	findsRs, err := jiraClient.Find(time.Now())
	if err != nil {
		log.Println("Error while sending Finds-request")
		fmt.Println(err)
		return
	}

	var totalSpentSeconds int
	for _, f := range findsRs {
		totalSpentSeconds += f.TimeSpentSeconds
	}

	fmt.Println(totalSpentSeconds)

	worklogs := strings.Split(env.Worklog, ";")
	for i, worklog := range worklogs {

		if totalSpentSeconds < targetSpentSeconds {
			randomSeconds := 0
			neededSpentSeconds := targetSpentSeconds - totalSpentSeconds
			if i < len(worklogs)-1 {
				randomSeconds = generateRandomInt(neededSpentSeconds)
			}
			currentWorklog := strings.Split(worklog, ",")

			origanalTaskId := strings.TrimSpace(currentWorklog[0])
			comment := strings.TrimSpace(currentWorklog[2])
			envTimeSpentHours, err := strconv.Atoi(strings.TrimSpace(currentWorklog[1]))
			if err != nil {
				log.Printf("Error while convert hours to int from .env file")
				return
			}
			envTimeSpentSeconds := envTimeSpentHours * 60 * 60
			timeSpentSeconds := chooceAvaliableSecondsToSpent(i, len(worklogs), neededSpentSeconds, randomSeconds, envTimeSpentSeconds)

			if err := jiraClient.Create(comment, origanalTaskId, time.Now(), timeSpentSeconds); err != nil {
				log.Println(err)
				return
			}

			totalSpentSeconds += timeSpentSeconds
		}
	}
}

func generateRandomInt(neededSpentSeconds int) int {
	min, max := 1, neededSpentSeconds/60/60
	rand.Seed(time.Now().UnixNano())
	if max-min <= 0 {
		return 1 * 60 * 60
	} else {
		return (rand.Intn(max-min) + min) * 60 * 60
	}
}

func chooceAvaliableSecondsToSpent(i, length, neededSpentSeconds, randomSeconds, envTimeSpentSeconds int) int {
	if i == length-1 {
		return neededSpentSeconds
	} else if randomSeconds > envTimeSpentSeconds {
		return envTimeSpentSeconds
	} else {
		return randomSeconds
	}
}
