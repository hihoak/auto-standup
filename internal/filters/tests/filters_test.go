package tests

import (
	"github.com/andygrunwald/go-jira"
	"github.com/golang/mock/gomock"
	"github.com/hihoak/auto-standup/internal"
	"github.com/hihoak/auto-standup/internal/filters"
	"github.com/hihoak/auto-standup/pkg/utils"
	"github.com/hihoak/auto-standup/test"
	"testing"
	"time"
)

func TestFilterIssuesByProject(t *testing.T) {
	t.Parallel()

	mc := gomock.NewController(t)

	testCfg := &utils.Config{
		ExcludeJiraProjects: []string{
			"eXCLude",
			"TEST",
		},
	}

	cases := []test.Case{
		{
			Name: "Error. 1) Issue included in Excluded projects",
			FuncArguments: []interface{}{
				&jira.Issue{Key: "EXCLUDE-1000"},
			},
			ExpectedResult: false,
			Setup: func() (*internal.Implementator, *utils.Config) {
				return test.InitTestImplementator(test.InitDefaultMockClients(mc)), testCfg
			},
		},
		{
			Name: "Error. 2) Issue included in Excluded projects",
			FuncArguments: []interface{}{
				&jira.Issue{Key: "exclude-1000"},
			},
			ExpectedResult: false,
			Setup: func() (*internal.Implementator, *utils.Config) {
				return test.InitTestImplementator(test.InitDefaultMockClients(mc)), testCfg
			},
		},
		{
			Name: "Error. 3) Issue included in Excluded projects",
			FuncArguments: []interface{}{
				&jira.Issue{Key: "TEST-1000"},
			},
			ExpectedResult: false,
			Setup: func() (*internal.Implementator, *utils.Config) {
				return test.InitTestImplementator(test.InitDefaultMockClients(mc)), testCfg
			},
		},
		{
			Name: "Error. 4) Issue included in Excluded projects",
			FuncArguments: []interface{}{
				&jira.Issue{Key: "eXclude-1000"},
			},
			ExpectedResult: false,
			Setup: func() (*internal.Implementator, *utils.Config) {
				return test.InitTestImplementator(test.InitDefaultMockClients(mc)), testCfg
			},
		},
		{
			Name: "Success. 1) Issue included not in Excluded projects",
			FuncArguments: []interface{}{
				&jira.Issue{Key: "VALID-1000"},
			},
			ExpectedResult: true,
			Setup: func() (*internal.Implementator, *utils.Config) {
				return test.InitTestImplementator(test.InitDefaultMockClients(mc)), testCfg
			},
		},
		{
			Name: "Success. 2) Issue included not in Excluded projects",
			FuncArguments: []interface{}{
				&jira.Issue{Key: "Valid-1000"},
			},
			ExpectedResult: true,
			Setup: func() (*internal.Implementator, *utils.Config) {
				return test.InitTestImplementator(test.InitDefaultMockClients(mc)), testCfg
			},
		},
		{
			Name: "Success. 3) Issue included not in Excluded projects",
			FuncArguments: []interface{}{
				&jira.Issue{Key: "ANOtHErVAlID-1000"},
			},
			ExpectedResult: true,
			Setup: func() (*internal.Implementator, *utils.Config) {
				return test.InitTestImplementator(test.InitDefaultMockClients(mc)), testCfg
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			_, cfg := tc.Setup()
			cl := filters.Filters{}
			res := cl.FilterIssuesByProject(cfg, tc.FuncArguments[0].(*jira.Issue))
			tc.CheckCase(t, res, nil)
		})
	}
}

func TestFilterIssueByActivity(t *testing.T) {
	t.Parallel()
	mc := gomock.NewController(t)

	notValidUser := "not valid user"
	validUser := "valid user"

	cmdStartTime := time.Date(2000, time.March, 6, 12, 0, 0, 0, time.UTC)
	validTimes := []string{
		"2000-03-05T14:30:00.000+0000",
		"2000-03-05T22:30:00.000+0000",
		"2000-03-06T06:30:00.000+0000",
		"2000-03-06T11:59:00.000+0000",
	}

	cfg := &utils.Config{
		EligibleUsersHistories: []string{
			validUser,
			"another good",
		},
		CmdStartTime: cmdStartTime,
		NumberOfDaysForGetTickets: 1,
	}

	setupFunc := func() (*internal.Implementator, *utils.Config) {
		return test.InitTestImplementator(test.InitDefaultMockClients(mc)), cfg
	}

	cases := []test.Case{
		{
			Name: "False. Issue doesn't contains valid history. Changelog is nil",
			FuncArguments: []interface{}{
				&jira.Issue{Changelog: nil},
			},
			Setup: setupFunc,
			ExpectedResult: false,
		},
		{
			Name: "False. Issue doesn't contains valid history. Changelog is empty",
			FuncArguments: []interface{}{
				&jira.Issue{Changelog: &jira.Changelog{}},
			},
			Setup: setupFunc,
			ExpectedResult: false,
		},
		{
			Name: "False. Issue doesn't contains valid history. Histories is nil",
			FuncArguments: []interface{}{
				&jira.Issue{
					Changelog: &jira.Changelog{
						Histories: nil,
					},
				},
			},
			Setup: setupFunc,
			ExpectedResult: false,
		},
		{
			Name: "False. Issue doesn't contains valid history. Histories is empty",
			FuncArguments: []interface{}{
				&jira.Issue{
					Changelog: &jira.Changelog{
						Histories: []jira.ChangelogHistory{},
					},
				},
			},
			Setup: setupFunc,
			ExpectedResult: false,
		},
		{
			Name: "False. Issue doesn't contains valid history. Have Histories, but no one from eligible user",
			FuncArguments: []interface{}{
				&jira.Issue{
					Changelog: &jira.Changelog{
						Histories: []jira.ChangelogHistory{
							{
								Author: jira.User{
									Name: notValidUser,
								},
							},
							{
								Author: jira.User{
									Name: "another not valid user",
								},
							},
						},
					},
				},
			},
			Setup: setupFunc,
			ExpectedResult: false,
		},
		{
			Name: "False. Issue doesn't contains valid history. Have Histories, found from valid user, but time is expired",
			FuncArguments: []interface{}{
				&jira.Issue{
					Changelog: &jira.Changelog{
						Histories: []jira.ChangelogHistory{
							{
								Author: jira.User{
									Name: validUser,
								},
								Created: cmdStartTime.AddDate(0, 0, 7).Format(time.RFC3339),
							},
							{
								Author: jira.User{
									Name: "another not valid user",
								},
							},
						},
					},
				},
			},
			Setup: setupFunc,
			ExpectedResult: false,
		},
		{
			Name: "True. 1) Issue contains valid history",
			FuncArguments: []interface{}{
				&jira.Issue{
					Changelog: &jira.Changelog{
						Histories: []jira.ChangelogHistory{
							{
								Author: jira.User{
									Name: validUser,
								},
								Created: validTimes[0],
							},
						},
					},
				},
			},
			Setup: setupFunc,
			ExpectedResult: true,
		},
		{
			Name: "True. 2) Issue contains valid history",
			FuncArguments: []interface{}{
				&jira.Issue{
					Changelog: &jira.Changelog{
						Histories: []jira.ChangelogHistory{
							{
								Author: jira.User{
									Name: validUser,
								},
								Created: validTimes[1],
							},
						},
					},
				},
			},
			Setup: setupFunc,
			ExpectedResult: true,
		},
		{
			Name: "True. 3) Issue contains valid history",
			FuncArguments: []interface{}{
				&jira.Issue{
					Changelog: &jira.Changelog{
						Histories: []jira.ChangelogHistory{
							{
								Author: jira.User{
									Name: validUser,
								},
								Created: validTimes[2],
							},
						},
					},
				},
			},
			Setup: setupFunc,
			ExpectedResult: true,
		},
		{
			Name: "True. 4) Issue contains valid history",
			FuncArguments: []interface{}{
				&jira.Issue{
					Changelog: &jira.Changelog{
						Histories: []jira.ChangelogHistory{
							{
								Author: jira.User{
									Name: validUser,
								},
								Created: validTimes[3],
							},
						},
					},
				},
			},
			Setup: setupFunc,
			ExpectedResult: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			_, cfg := tc.Setup()
			cl := filters.Filters{}
			res := cl.FilterIssueByActivity(cfg, tc.FuncArguments[0].(*jira.Issue))
			tc.CheckCase(t, res, nil)
		})
	}
}
