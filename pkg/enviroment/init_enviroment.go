package enviroment

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Environment struct {
	JiraUrl               string
	JiraTempoCreatesUri   string
	JiraTempoFindsUri     string
	JiraTempoUserkeyUri   string
	JiraTempoDaysSearch   string
	Login                 string
	Password              string
	Worklog               string
	OutlookLoggingEnabled string
	OutlookUrl            string
	OutlookEventPath      string
	OutlookDefaultTask    string
}

func NewEnvironment() (*Environment, error) {
	env, err := readEnv()

	if err != nil {
		log.Println("Error while loading .env file")
		return nil, err
	}

	errValidate := validate(env)
	if errValidate != nil {
		return nil, errValidate
	}

	return env, nil
}

func readEnv() (*Environment, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	env := &Environment{
		strings.TrimSpace(os.Getenv("JIRA_URL")),
		strings.TrimSpace(os.Getenv("JIRA_TEMPO_CREATES_URI")),
		strings.TrimSpace(os.Getenv("JIRA_TEMPO_FINDS_URI")),
		strings.TrimSpace(os.Getenv("JIRA_TEMPO_USERKEY_URI")),
		strings.TrimSpace(os.Getenv("JIRA_TEMPO_DAYS_SEARCH")),
		strings.TrimSpace(os.Getenv("LOGIN")),
		strings.TrimSpace(os.Getenv("PASSWORD")),
		strings.TrimSpace(os.Getenv("WORKLOG")),
		strings.TrimSpace(os.Getenv("OUTLOOK_LOGGING_ENABLED")),
		strings.TrimSpace(os.Getenv("OUTLOOK_URL")),
		strings.TrimSpace(os.Getenv("OUTLOOK_EVENT_URI")),
		strings.TrimSpace(os.Getenv("OUTLOOK_DEFAULT_TASK_FOR_LOGGING")),
	}

	return env, nil
}

func validate(env *Environment) error {
	if "" == env.JiraUrl {
		return errors.New(genEnvErrMessage("JIRA_URL"))
	}
	if "" == env.JiraTempoCreatesUri {
		return errors.New(genEnvErrMessage("JIRA_TEMPO_CREATES_URI"))
	}
	if "" == env.JiraTempoFindsUri {
		return errors.New(genEnvErrMessage("JIRA_TEMPO_FINDS_URI"))
	}
	if "" == env.JiraTempoUserkeyUri {
		return errors.New(genEnvErrMessage("JIRA_TEMPO_USERKEY_URI"))
	}
	if "" == env.JiraTempoDaysSearch {
		return errors.New(genEnvErrMessage("JIRA_TEMPO_DAYS_SEARCH"))
	}
	if "" == env.Login {
		return errors.New(genEnvErrMessage("LOGIN"))
	}
	if "" == env.Password {
		return errors.New(genEnvErrMessage("PASSWORD"))
	}
	if "" == env.Worklog {
		return errors.New(genEnvErrMessage("WORKLOG"))
	}
	if "" == env.OutlookLoggingEnabled {
		return errors.New(genEnvErrMessage("OUTLOOK_LOGGING_ENABLED"))
	}
	if "" == env.OutlookUrl {
		return errors.New(genEnvErrMessage("OUTLOOK_URL"))

	}
	if "" == env.OutlookEventPath {
		return errors.New(genEnvErrMessage("OUTLOOK_EVENT_URI"))

	} else {
		err := validateWorklogItems(env.Worklog)
		if err != nil {
			return err
		}
	}

	return nil
}

func genEnvErrMessage(envName string) string {
	return fmt.Sprintf("Enviroment value \"%s\" is empty", envName)
}

func validateWorklogItems(worklog string) error {
	worklogs := strings.Split(worklog, ";")
	for _, workBlock := range worklogs {
		work := strings.Split(workBlock, ",")
		if len(work) < 3 {
			return errors.New("Each block of \"WORKLOG\" env param should contain 3 elements!")
		}
		for _, workItem := range work {
			if "" == strings.TrimSpace(workItem) {
				return errors.New("Items in \"WORKLOG\" should not be empty!")
			}
		}
	}
	return nil
}
