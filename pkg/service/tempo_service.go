package service

import (
	"log"
	"strings"
	"tempo-loger/pkg/jira"
	"time"
)

func ProcessParams(totalSpentSeconds, targetSpentSeconds int, client jira.JiraClient, params []jira.CreateParams) error {
	for y, param := range params {
		if totalSpentSeconds < targetSpentSeconds {
			if "" == strings.TrimSpace(param.OriginTaskId) {
				log.Printf("OriginalTaskId for %s is empty. Skipping event.", param.Comment)
				continue
			}
			neededSpentSeconds := targetSpentSeconds - totalSpentSeconds
			timeSpentSeconds := choiceAvailableSecondsToSpent(y, len(params), neededSpentSeconds, param.TimeSpentSeconds)
			param.BillableSeconds = timeSpentSeconds
			param.TimeSpentSeconds = timeSpentSeconds
			if err := client.Create(param); err != nil {
				log.Fatal("Error while execute Create: ", err)
			}
			log.Printf("Successful logged: task id = %s, comment = %s, spent time = %dsec (%d hours)",
				param.OriginTaskId, param.Comment, timeSpentSeconds, ConvertSecondsToHours(timeSpentSeconds))
			totalSpentSeconds += timeSpentSeconds
		} else {
			log.Println("Excellent day! All time has already been spent")
			return nil
		}
	}
	return nil
}

func CalculateTimeInTempo(daysRs jira.DaysSearchRs, jiraClient jira.JiraClient) (int, int, error) {
	targetSpentSeconds := daysRs[0].Days[0].RequiredSeconds
	if targetSpentSeconds <= 0 {
		log.Fatalf("requiredSeconds = %d. Today is not-working day!", targetSpentSeconds)
	}
	log.Printf("Target time to spend = %dsec (%d hours)", targetSpentSeconds, ConvertSecondsToHours(targetSpentSeconds))

	findsRs, err := jiraClient.Find(time.Now())
	if err != nil {
		return 0, 0, err
	}

	var totalSpentSeconds int
	for _, f := range findsRs {
		totalSpentSeconds += f.TimeSpentSeconds
	}
	log.Printf("Already spent time = %dsec (%d hours)", totalSpentSeconds, ConvertSecondsToHours(totalSpentSeconds))
	return targetSpentSeconds, totalSpentSeconds, nil
}

func choiceAvailableSecondsToSpent(i, length, neededSpentSeconds, secondsToSpent int) int {
	if i == length-1 {
		return neededSpentSeconds
	} else if neededSpentSeconds < secondsToSpent {
		return neededSpentSeconds
	} else {
		return secondsToSpent
	}
}
