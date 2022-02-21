package service

import (
	"log"
	"strings"
	"tempo-loger/pkg/jira"
)

func ProcessParams(totalSpentSeconds, targetSpentSeconds int, client jira.JiraClient, params []jira.CreateParams) error {
	for y, param := range params {
		if totalSpentSeconds < targetSpentSeconds {
			if "" == strings.TrimSpace(param.OriginTaskId) {
				log.Printf("OriginalTaskId for %s is empty. Skipping event.", param.Comment)
				continue
			}
			neededSpentSeconds := targetSpentSeconds - totalSpentSeconds
			timeSpentSeconds := roundAdnChoiceAvailableSecondsToSpent(y, len(params), neededSpentSeconds, param.TimeSpentSeconds)
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

func roundAdnChoiceAvailableSecondsToSpent(i, length, neededSpentSeconds, secondsToSpent int) int {
	if i == length-1 {
		return neededSpentSeconds
	} else {
		return secondsToSpent
	}
}
