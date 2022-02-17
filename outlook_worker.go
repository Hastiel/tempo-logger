package main

import (
	"fmt"
	"log"
	"tempo-loger/pkg/enviroment"
	"tempo-loger/pkg/jira"
	"tempo-loger/pkg/outlook"
	"tempo-loger/pkg/service"
)

func process(env *enviroment.Environment) error {
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

	createParams, err := prepare(outlookEvents, *env)
	if err != nil {
		return err
	}
	log.Printf(createParams[0].OriginTaskId)
	return nil
}

func prepare(outlookEvents outlook.EventsRs, env enviroment.Environment) ([]jira.CreateParams, error) {
	var createParams []jira.CreateParams
	if len(outlookEvents.Value) > 0 {
		for _, val := range outlookEvents.Value {
			jiraTicket := service.ExtractTicketFromBody(val.Body.Content, env.JiraUrl)
			if "" == jiraTicket {
				jiraTicket = env.OutlookDefaultTask
			}
			startDate, err := service.ExtractDateTime(val.Start.DateTime)
			if err != nil {
				return nil, err
			}
			endDate, err := service.ExtractDateTime(val.End.DateTime)
			if err != nil {
				return nil, err
			}
			ev := endDate.Unix() - startDate.Unix()
			createParams = append(createParams, jira.CreateParams{
				BillableSeconds:  int(ev),
				Comment:          fmt.Sprintf("Участие во встрече \"%s\"", val.Subject),
				EndDate:          endDate,
				Started:          startDate,
				OriginTaskId:     jiraTicket,
				TimeSpentSeconds: int(ev),
			})
		}
	}
	return createParams, nil
}
