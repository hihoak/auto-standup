package jirer

import (
	"github.com/andygrunwald/go-jira"
	"github.com/hihoak/auto-standup/pkg/utils"
	"github.com/hihoak/auto-standup/test"
	"testing"
	"time"
)

func TestGetIssueLogTimeForTheLastWorkDay(t *testing.T) {
	issueWorklogNil := jira.Issue{
		Fields: &jira.IssueFields{
			Worklog: nil,
		},
	}

	issueWorklogsNil := jira.Issue{
		Fields: &jira.IssueFields{
			Worklog: &jira.Worklog{
				Worklogs: nil,
			},
		},
	}

	invalidTimes := []jira.Time{
		jira.Time(time.Date(2000, time.March, 4, 14, 30, 0, 0, time.UTC)),
		jira.Time(time.Date(2000, time.March, 5, 6, 30, 0, 0, time.UTC)),
	}
	validTimes := []jira.Time{
		jira.Time(time.Date(2000, time.March, 5, 14, 30, 0, 0, time.UTC)),
		jira.Time(time.Date(2000, time.March, 6, 6, 30, 0, 0, time.UTC)),
	}

	issue := jira.Issue{
		Fields: &jira.IssueFields{
			Worklog: &jira.Worklog{
				Worklogs: []jira.WorklogRecord{
					{
						Created: &invalidTimes[0],
						// 1m
						TimeSpentSeconds: 60,
					},
					{
						Created: &invalidTimes[1],
						// 1h
						TimeSpentSeconds: 3600,
					},
					{
						Created: &validTimes[0],
						// 1d
						TimeSpentSeconds: 3600 * 24,
					},
					{
						Created: &validTimes[0],
						// 1w
						TimeSpentSeconds: 3600 * 24 * 7,
					},
				},
			},
		},
	}

	cmdStartTime := time.Date(2000, time.March, 6, 12, 0, 0, 0, time.UTC)

	testCfg := &utils.Config{
		NumberOfDaysForGetTickets: 1,
		CmdStartTime: cmdStartTime,
	}

	cases := []test.Case{
		{
			Name: "Error. Worklog is nil",
			FuncArguments: []interface{}{
				issueWorklogNil,
			},
			ExpectedResult: 0,
		},
		{
			Name: "Error. Woklogs is nil",
			FuncArguments: []interface{}{
				issueWorklogsNil,
			},
			ExpectedResult: 0,
		},
		{
			Name: "Success. All ok",
			FuncArguments: []interface{}{
				issue,
			},
			ExpectedResult: 3600 * 24 * 8,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			cl := Client{}
			res := cl.GetIssueLogTimeForTheLastWorkDay(testCfg, tc.FuncArguments[0].(jira.Issue))
			tc.CheckCase(t, res, nil)
		})
	}
}
