package utils

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	// Cfg ...
	Cfg *Config
)

// Config - ...
type Config struct {
	// Username - Jira username for authentication
	Username string `yaml:"username"`
	// Password - Jira password for authentication
	Password string `yaml:"password"`
	// EligibleUsersHistories - updates from these users are considered valid when finding tickets from the last business day
	EligibleUsersHistories []string `yaml:"eligible_users_histories" default:"gitlab"`
	// ExcludeJiraProjects - Jira projects that tickets will be ignore while creating report
	ExcludeJiraProjects []string `yaml:"exclude_jira_projects" default:"retest"`

	// CmdStartTime - command start time in UTC
	CmdStartTime time.Time
	// NumberOfDaysForGetTickets - number of days for which we take tickets for report
	NumberOfDaysForGetTickets int
}

// NewConfig - ...
func NewConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		Log.Debug().Err(err).Msgf("Can't open file %s", path)
		return nil, err
	}
	nowTime := time.Now().UTC()
	cfg := Config{
		CmdStartTime:              nowTime,
		NumberOfDaysForGetTickets: getNumberOfDaysForGetTickets(nowTime),
	}
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		Log.Debug().Err(err).Msgf("Failed to unmarshall file '%s' to yaml. Check that you fill username and password fields", path)
		return nil, err
	}
	// fill empty fields of config if found any
	cfg.fillEmptyFields()
	Log.Debug().Msgf("Successfully create config %v", cfg)
	return &cfg, nil
}

func (c *Config) fillEmptyFields() {
	typ := reflect.TypeOf(*c)
	if len(c.EligibleUsersHistories) == 0 {
		f, _ := typ.FieldByName("EligibleUsersHistories")
		c.EligibleUsersHistories = strings.Split(f.Tag.Get("default"), ",")
		c.EligibleUsersHistories = append(c.EligibleUsersHistories, c.Username)
		Log.Debug().Msgf("Found empty config field 'EligibleUsersHistories'. Fill it with %v", c.EligibleUsersHistories)
	}
	if len(c.ExcludeJiraProjects) == 0 {
		f, _ := typ.FieldByName("ExcludeJiraProjects")
		c.ExcludeJiraProjects = strings.Split(f.Tag.Get("default"), ",")
		Log.Debug().Msgf("Found empty config field 'ExcludeJiraProjects'. Fill it with %v", c.ExcludeJiraProjects)
	}
}

func getNumberOfDaysForGetTickets(cmdStartTime time.Time) int {
	// number of days for which we take tickets
	numberOfDaysForGetTickets := 1
	// if it's Monday we need take ticket for 3 days, because of holidays
	if cmdStartTime.Weekday() == time.Monday {
		numberOfDaysForGetTickets += 2
	} else if cmdStartTime.Weekday() == time.Sunday {
		numberOfDaysForGetTickets++
	}

	return numberOfDaysForGetTickets
}
