package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func ExtractTicketFromBody(body, jiraUrl string) string {
	urlWithPath := fmt.Sprintf("%s/browse/", jiraUrl)
	pattern := regexp.MustCompile(fmt.Sprintf("(%s)([aA-zZ]{1,10}-\\d{1,6})", urlWithPath))
	s := pattern.FindString(body)
	return strings.TrimSuffix(s, urlWithPath)
}

func ExtractDateTime(val string) (time.Time, error) {
	res, err := time.Parse("2006-01-02T15:04:05", val)
	if err != nil {
		return time.Time{}, err
	}
	return res, nil
}
