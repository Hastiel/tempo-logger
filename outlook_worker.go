package main

import (
	"fmt"
	"log"
	"strings"
	"tempo-loger/pkg/enviroment"
	"tempo-loger/pkg/jira"
	"tempo-loger/pkg/outlook"
	"tempo-loger/pkg/service"
	"time"
)

func process(env *enviroment.Environment, createParams *[]jira.CreateParams) error {
	outlookClient := outlook.New(env.Login, env.Password, env.OutlookUrl, env.OutlookEventPath)
	startDate, err := concatCurrantDateWithTime("00:00:00")
	if err != nil {
		log.Fatal("Cannot parse startDateTime for outlook events: ", err)
	}
	endDate, err := concatCurrantDateWithTime("23:59:59")
	if err != nil {
		log.Fatal("Cannot parse endDateTime for outlook events: ", err)
	}
	outlookEvents, err := outlookClient.GetEvents(startDate, endDate)
	if err != nil {
		log.Fatal("Error while getting outlook events: ", err)
	}
	filterEvents(&outlookEvents)
	if err := prepare(outlookEvents, *env, createParams); err != nil {
		return err
	}
	return nil
}

func filterEvents(rs *outlook.EventsRs) {
	for i, val := range rs.Value {
		if "organizer" != strings.ToLower(val.ResponseStatus.Response) &&
			"accepted" != strings.ToLower(val.ResponseStatus.Response) &&
			"tentativelyaccepted" != strings.ToLower(val.ResponseStatus.Response) {
			rs.Value = append(rs.Value[:i], rs.Value[i+1:]...)
			i--
		}
	}
}

func prepare(outlookEvents outlook.EventsRs, env enviroment.Environment, createParams *[]jira.CreateParams) error {
	if len(outlookEvents.Value) > 0 {
		for _, val := range outlookEvents.Value {
			jiraTicket := service.ExtractTicketFromBody(val.Body.Content, env.JiraUrl)
			if "" == jiraTicket {
				jiraTicket = env.OutlookDefaultTask
			}
			startDate, err := service.ExtractDateTime(val.Start.DateTime)
			if err != nil {
				return err
			}
			endDate, err := service.ExtractDateTime(val.End.DateTime)
			if err != nil {
				return err
			}
			eventDuration := endDate.Unix() - startDate.Unix()
			*createParams = append(*createParams, jira.CreateParams{
				BillableSeconds:  int(eventDuration),
				Comment:          fmt.Sprintf("Участие во встрече \"%s\"", val.Subject),
				EndDate:          endDate,
				Started:          startDate,
				OriginTaskId:     jiraTicket,
				TimeSpentSeconds: int(eventDuration),
			})
		}
	}
	return nil
}

func concatCurrantDateWithTime(timeVal string) (time.Time, error) {
	year, month, day := time.Now().Date()
	dateTime, err := time.Parse("2 January 2006 15:04:05", fmt.Sprintf("%d %s %d %s", day, month, year, timeVal))
	if err != nil {
		return time.Time{}, err
	}
	return dateTime, nil
}
