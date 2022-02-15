package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"tempo-loger/pkg/enviroment"
	"tempo-loger/pkg/jira"
	"tempo-loger/pkg/outlook"
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

	var outlookEvents outlook.EventsRs
	if "true" == env.OutlookLoggingEnabled {
		outlookClient := outlook.New(env.Login, env.Password, env.OutlookUrl, env.OutlookEventPath)
		startDate, err := concatCurrantDateWithTime("00:00:00")
		if err != nil {
			log.Fatal("Cannot parse startDateTime for outlook events: ", err)
		}
		endDate, err := concatCurrantDateWithTime("23:59:59")
		if err != nil {
			log.Fatal("Cannot parse endDateTime for outlook events: ", err)
		}
		outlookEvents, err = outlookClient.GetEvents(startDate, endDate)
		if err != nil {
			log.Fatal("Error while getting outlook events: ", err)
		}
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

	r := rand.New(rand.NewSource(time.Now().Unix()))
	counter := 0
	for _, i := range r.Perm(len(worklogs)) {
		if totalSpentSeconds < targetSpentSeconds {
			if len(outlookEvents.Value) > 0 {
				for _, val := range outlookEvents.Value {
					startDate, err := time.Parse("2006-01-02T15:04:05", val.Start.DateTime)
					if err != nil {
						log.Fatal("Error while parse startDate in outlookEventRs: ", err)
					}
					endDate, err := time.Parse("2006-01-02T15:04:05", val.End.DateTime)
					if err != nil {
						log.Fatal("Error while parse endDate in outlookEventRs: ", err)
					}
					ev := endDate.Unix() - startDate.Unix()
					if err := jiraClient.Create(fmt.Sprintf("Участие во встрече \"%s\"", val.Subject), env.OutlookDefaultTask, time.Now(), int(ev)); err != nil {
						log.Fatal("Error while execute Create: ", err)
					}
					totalSpentSeconds += int(ev)
				}
				log.Fatal("Ha-Ha-Ha")
			}
			randomSeconds := 0
			neededSpentSeconds := targetSpentSeconds - totalSpentSeconds
			if counter < len(worklogs)-1 {
				randomSeconds = generateRandomInt(neededSpentSeconds)
			}
			currentWorklog := strings.Split(worklogs[i], ",")

			originalTaskId := strings.TrimSpace(currentWorklog[0])
			comment := strings.TrimSpace(currentWorklog[2])
			envTimeSpentHours, err := strconv.Atoi(strings.TrimSpace(currentWorklog[1]))
			if err != nil {
				log.Fatal("Error while convert hours to int from .env file: ", err)
			}
			envTimeSpentSeconds := convertHoursToSeconds(envTimeSpentHours)
			timeSpentSeconds := roundAdnChoiceAvailableSecondsToSpent(counter, len(worklogs), neededSpentSeconds, randomSeconds, envTimeSpentSeconds)

			if err := jiraClient.Create(comment, originalTaskId, time.Now(), timeSpentSeconds); err != nil {
				log.Fatal("Error while execute Create: ", err)
			}
			log.Printf("Successful logged: task id = %s, comment = %s, spent time = %dsec (%d hours)", originalTaskId, comment, timeSpentSeconds, timeSpentSeconds/60/60)

			totalSpentSeconds += timeSpentSeconds
			counter++
		} else {
			log.Println("Excellent day! All time has already been spent")
			break
		}
	}
}

func generateRandomInt(neededSpentSeconds int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(neededSpentSeconds)
}

func roundAdnChoiceAvailableSecondsToSpent(i, length, neededSpentSeconds, randomSeconds, envTimeSpentSeconds int) int {
	roundedHours := convertSecondsToHours(randomSeconds)
	if roundedHours == 0 {
		roundedHours = 1
	}
	seconds := convertHoursToSeconds(roundedHours)
	if i == length-1 {
		return neededSpentSeconds
	} else if seconds > envTimeSpentSeconds {
		return envTimeSpentSeconds
	} else {
		return seconds
	}
}

func convertSecondsToHours(val int) int {
	rounded := math.Round(float64(val) / 60 / 60)
	return int(rounded)
}

func convertHoursToSeconds(val int) int {
	return val * 60 * 60
}

func concatCurrantDateWithTime(timeVal string) (time.Time, error) {
	year, month, day := time.Now().Date()
	dateTime, err := time.Parse("2 January 2006 15:04:05", fmt.Sprintf("%d %s %d %s", day, month, year, timeVal))
	if err != nil {
		return time.Time{}, err
	}
	return dateTime, nil
}
