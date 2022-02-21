package main

import (
	"math/rand"
	"strconv"
	"strings"
	"tempo-loger/pkg/enviroment"
	"tempo-loger/pkg/jira"
	"tempo-loger/pkg/service"
	"time"
)

func processWorklogs(env enviroment.Environment, createParams *[]jira.CreateParams) error {
	worklogs := strings.Split(env.Worklog, ";")
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(worklogs)) {
		currentWorklog := strings.Split(worklogs[i], ",")
		envOriginalTaskId := strings.TrimSpace(currentWorklog[0])
		envComment := strings.TrimSpace(currentWorklog[2])
		envTimeSpentHours, err := strconv.Atoi(strings.TrimSpace(currentWorklog[1]))
		if err != nil {
			return err
		}
		envTimeSpentSeconds := service.ConvertHoursToSeconds(envTimeSpentHours)
		randomSeconds := generateRandomInt(envTimeSpentSeconds)
		secondsToSpent := roundSecondsToSpent(randomSeconds)

		*createParams = append(*createParams, jira.CreateParams{
			BillableSeconds:  secondsToSpent,
			Comment:          envComment,
			EndDate:          time.Now(),
			Started:          time.Now(),
			OriginTaskId:     envOriginalTaskId,
			TimeSpentSeconds: secondsToSpent,
		})
	}
	return nil
}

func roundSecondsToSpent(randomSeconds int) int {
	roundedHours := service.ConvertSecondsToHours(randomSeconds)
	if roundedHours == 0 {
		roundedHours = 1
	}
	return service.ConvertHoursToSeconds(roundedHours)
}

func generateRandomInt(neededSpentSeconds int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(neededSpentSeconds)
}
