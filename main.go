package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"tempo-loger/pkg/enviroment"
	"tempo-loger/pkg/jira"
	"time"
)

func main() {
	logFile, err := os.OpenFile(fmt.Sprintf("%s.log", time.Now().Format("2006-01")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error while init .log file: ", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	env, err := enviroment.NewEnvironment()
	if err != nil {
		log.Fatal("Error while loading .env file: ", err)
	}

	jiraClient := jira.New(env.Login, env.Password, env.JiraUrl, env.JiraTempoFindsUri, env.JiraTempoUserkeyUri, env.JiraTempoCreatesUri, env.JiraTempoDaysSearch)

	daysRs, err := jiraClient.GetDayInfo(time.Now())
	if err != nil {
		log.Fatal("Error while sending GetDayInfo-request: ", err)
	}

	if daysRs[0].Days[0].Type != "WORKING_DAY" {
		log.Fatal("Today is not-working day!")
	}
	targetSpentSeconds := daysRs[0].Days[0].RequiredSeconds
	log.Printf("Target time to spend = %dsec (%d hours)", targetSpentSeconds, targetSpentSeconds/60/60)

	findsRs, err := jiraClient.Find(time.Now())
	if err != nil {
		log.Fatal("Error while sending Finds-request: ", err)
	}

	var totalSpentSeconds int
	for _, f := range findsRs {
		totalSpentSeconds += f.TimeSpentSeconds
	}
	log.Printf("Already spent time = %dsec (%d hours)", totalSpentSeconds, totalSpentSeconds/60/60)

	worklogs := strings.Split(env.Worklog, ";")
	for i, worklog := range worklogs {

		if totalSpentSeconds < targetSpentSeconds {
			randomSeconds := 0
			neededSpentSeconds := targetSpentSeconds - totalSpentSeconds
			if i < len(worklogs)-1 {
				randomSeconds = generateRandomInt(neededSpentSeconds)
			}
			currentWorklog := strings.Split(worklog, ",")

			originalTaskId := strings.TrimSpace(currentWorklog[0])
			comment := strings.TrimSpace(currentWorklog[2])
			envTimeSpentHours, err := strconv.Atoi(strings.TrimSpace(currentWorklog[1]))
			if err != nil {
				log.Fatal("Error while convert hours to int from .env file: ", err)
			}
			envTimeSpentSeconds := envTimeSpentHours * 60 * 60
			timeSpentSeconds := choiceAvailableSecondsToSpent(i, len(worklogs), neededSpentSeconds, randomSeconds, envTimeSpentSeconds)

			if err := jiraClient.Create(comment, originalTaskId, time.Now(), timeSpentSeconds); err != nil {
				log.Fatal("Error while execute Create: ", err)
			}
			log.Printf("Successful logged: task id = %s, comment = %s, spent time = %dsec (%d hours)", originalTaskId, comment, timeSpentSeconds, timeSpentSeconds/60/60)

			totalSpentSeconds += timeSpentSeconds
		} else {
			log.Println("Excellent day! All time has already been spent")
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

func choiceAvailableSecondsToSpent(i, length, neededSpentSeconds, randomSeconds, envTimeSpentSeconds int) int {
	if i == length-1 {
		return neededSpentSeconds
	} else if randomSeconds > envTimeSpentSeconds {
		return envTimeSpentSeconds
	} else {
		return randomSeconds
	}
}
