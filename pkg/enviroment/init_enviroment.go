package enviroment

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type LocalEnv interface {
	//Init() (enviroment, error)
	//validate() error
}

type enviroment struct {
	JiraUrl             string
	JiraTempoCreatesUri string
	JiraTempoFindsUri   string
	JiraTempoUserkeyUri string
	JiraTempoDaysSearch string
	Login               string
	Password            string
	Worklog             string
}

func NewEnviroment() (*enviroment, error) {
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

func readEnv() (*enviroment, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	env := &enviroment{
		strings.TrimSpace(os.Getenv("JIRA_URL")),
		strings.TrimSpace(os.Getenv("JIRA_TEMPO_CREATES_URI")),
		strings.TrimSpace(os.Getenv("JIRA_TEMPO_FINDS_URI")),
		strings.TrimSpace(os.Getenv("JIRA_TEMPO_USERKEY_URI")),
		strings.TrimSpace(os.Getenv("JIRA_TEMPO_DAYS_SEARCH")),
		strings.TrimSpace(os.Getenv("LOGIN")),
		strings.TrimSpace(os.Getenv("PASSWORD")),
		strings.TrimSpace(os.Getenv("WORKLOG")),
	}

	return env, nil
}

func validate(env *enviroment) error {
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
				return errors.New("Items in \"WORKLOG\" shuld not be empty!")
			}
		}
	}
	return nil
}