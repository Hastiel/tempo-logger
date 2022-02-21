package main

import (
	"fmt"
	"log"
	"os"
	"tempo-loger/pkg/enviroment"
	"tempo-loger/pkg/jira"
	"tempo-loger/pkg/service"
	"time"
)

func main() {
	logFile, err := os.OpenFile(fmt.Sprintf("%s.log", time.Now().Format("2006-01")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error while init .log file: ", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {

		}
	}(logFile)
	log.SetOutput(logFile)

	env, err := enviroment.NewEnvironment()
	if err != nil {
		log.Fatal("Error while loading .env file: ", err)
	}

	var createParams []jira.CreateParams
	if "true" == env.OutlookLoggingEnabled {
		if err := process(env, &createParams); err != nil {
			log.Fatal("Error while process outlook_worker: ", err)
		}
	}

	if err := processWorklogs(*env, &createParams); err != nil {
		log.Fatal("Error while process worklog_worker: ", err)
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
	log.Printf("Target time to spend = %dsec (%d hours)", targetSpentSeconds, service.ConvertSecondsToHours(targetSpentSeconds))

	findsRs, err := jiraClient.Find(time.Now())
	if err != nil {
		log.Fatal("Error while sending Finds-request: ", err)
	}

	var totalSpentSeconds int
	for _, f := range findsRs {
		totalSpentSeconds += f.TimeSpentSeconds
	}
	log.Printf("Already spent time = %dsec (%d hours)", totalSpentSeconds, service.ConvertSecondsToHours(totalSpentSeconds))

	if err := service.ProcessParams(totalSpentSeconds, targetSpentSeconds, jiraClient, createParams); err != nil {
		log.Fatal("Error while process prepared worklogs slice: ", err)
	} else {
		log.Println("finished")
	}
}
