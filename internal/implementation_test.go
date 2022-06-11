package internal

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/andygrunwald/go-jira"
	"github.com/golang/mock/gomock"
	"github.com/hihoak/auto-standup/pkg/utils"
	"github.com/hihoak/auto-standup/test"
)

func TestFromStrKeysToIssues(t *testing.T) {
	t.Parallel()
	mc := gomock.NewController(t)

	testIssueKey := "TEST-1000"
	testIssueKey1 := "TEST-2000"

	testIssuesKeys := []string{
		testIssueKey,
		testIssueKey1,
	}

	testIssues := []*jira.Issue{
		{Key: testIssueKey},
		{Key: testIssueKey1},
	}

	cases := []test.Case{
		{
			Name: "Error. Error while getting issue from Jira",
			FuncArguments: []interface{}{
				testIssuesKeys,
			},
			ExpectedError: fmt.Errorf("can't get issue %s. error: %w", testIssueKey, test.ErrorTest),
			Setup: func() (*test.MockClients, *utils.Config) {
				mockClients := test.InitDefaultMockClients(mc)
				mockClients.JiraMockClient.EXPECT().GetIssue(gomock.Any(), gomock.Any()).
					Return(nil, nil, test.ErrorTest)
				return mockClients, nil
			},
		},
		{
			Name: "Error. Got only only 1 of 2 tickets",
			FuncArguments: []interface{}{
				testIssuesKeys,
			},
			ExpectedError: fmt.Errorf("can't get issue %s. error: %w", testIssueKey1, test.ErrorTest),
			Setup: func() (*test.MockClients, *utils.Config) {
				mockClients := test.InitDefaultMockClients(mc)
				mockClients.JiraMockClient.EXPECT().GetIssue(gomock.Any(), gomock.Any()).
					Return(testIssues[0], nil, nil)
				mockClients.JiraMockClient.EXPECT().GetIssue(gomock.Any(), gomock.Any()).
					Return(nil, nil, test.ErrorTest)
				return mockClients, nil
			},
		},
		{
			Name: "Success. Get all tickets",
			FuncArguments: []interface{}{
				testIssuesKeys,
			},
			ExpectedResult: testIssues,
			Setup: func() (*test.MockClients, *utils.Config) {
				mockClients := test.InitDefaultMockClients(mc)
				for _, issue := range testIssues {
					mockClients.JiraMockClient.EXPECT().GetIssue(gomock.Any(), gomock.Any()).
						Return(issue, nil, nil)
				}
				return mockClients, nil
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockClients, _ := tc.Setup()
			impl := InitTestImplementator(mockClients)
			res, err := impl.FromStrKeysToIssues(context.Background(), tc.FuncArguments[0].([]string))
			tc.CheckCase(t, res, err)
		})
	}
}

func TestGetIssuesFromLastWorkDay(t *testing.T) {
	t.Parallel()
	mc := gomock.NewController(t)

	cfg := &utils.Config{
		NumberOfDaysForGetTickets: 1,
		Username:                  "test user",
	}

	validIssues := []*jira.Issue{
		{Key: "VALID-1000"},
		{Key: "VALIDTOO-1500"},
	}

	testIssues := []jira.Issue{
		*validIssues[0],
		{Key: "NOTVALID_PROJECT-2000"},
		*validIssues[1],
		{Key: "NOTVALID_ACTIVITY-2000"},
		{Key: "NOTVALID_ALL-2000"},
	}

	cases := []test.Case{
		{
			Name: "Error. Failed to search issues",
			Setup: func() (*test.MockClients, *utils.Config) {
				mockClients := test.InitDefaultMockClients(mc)
				mockClients.JiraMockClient.EXPECT().SearchIssue(gomock.Any(), gomock.Any()).
					Return(nil, nil, test.ErrorTest)
				return mockClients, cfg
			},
			ExpectedError: test.ErrorTest,
		},
		{
			Name: "Success. Got some issues from search and filter them",
			Setup: func() (*test.MockClients, *utils.Config) {
				mockClients := test.InitDefaultMockClients(mc)
				mockClients.JiraMockClient.EXPECT().SearchIssue(gomock.Any(), gomock.Any()).
					Return(testIssues, nil, nil)
				for _, issue := range testIssues {
					if strings.HasPrefix(issue.Key, "NOTVALID_PROJECT") || strings.HasPrefix(issue.Key, "NOTVALID_ALL") {
						mockClients.FiltersMock.EXPECT().FilterIssuesByProject(gomock.Any(), gomock.Any()).Return(false)
						continue
					} else {
						mockClients.FiltersMock.EXPECT().FilterIssuesByProject(gomock.Any(), gomock.Any()).Return(true)
					}
					if strings.HasPrefix(issue.Key, "NOTVALID_ACTIVITY") || strings.HasPrefix(issue.Key, "NOTVALID_ALL") {
						mockClients.FiltersMock.EXPECT().FilterIssueByActivity(gomock.Any(), gomock.Any()).Return(false)
					} else {
						mockClients.FiltersMock.EXPECT().FilterIssueByActivity(gomock.Any(), gomock.Any()).Return(true)
					}
				}
				return mockClients, cfg
			},
			ExpectedResult: validIssues,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockClients, cfg := tc.Setup()
			impl := InitTestImplementator(mockClients)
			res, err := impl.GetIssuesFromLastWorkDay(cfg)
			tc.CheckCase(t, res, err)
		})
	}
}

func TestTodoIssuesToReport(t *testing.T) {
	t.Parallel()
	mc := gomock.NewController(t)

	jiraIssues := []*jira.Issue{
		{
			Key: "TEST-1000",
			Fields: &jira.IssueFields{
				Summary: "for test",
			},
		},
		{
			Key: "TEST-2000",
			Fields: &jira.IssueFields{
				Summary: "another test",
			},
		},
	}

	jiraIssuesForEstimateCase := []*jira.Issue{
		{
			Key: "TEST-1000",
			Fields: &jira.IssueFields{
				Summary: "for test",
				// 2h 30m
				TimeEstimate: 9000,
			},
		},
		{
			Key: "TEST-2000",
			Fields: &jira.IssueFields{
				Summary: "another test",
				// 1w 1d 1h 1m 5s
				TimeEstimate: 694865,
			},
		},
		{
			Key: "TEST-3000",
			Fields: &jira.IssueFields{
				Summary: "another test",
				// no estimate
				TimeEstimate: 0,
			},
		},
	}

	cases := []test.Case{
		{
			Name: "Just convert",
			FuncArguments: []interface{}{
				jiraIssues,
				false,
			},
			Setup: func() (*test.MockClients, *utils.Config) {
				return test.InitDefaultMockClients(mc), nil
			},
			ExpectedResult: fmt.Sprintf("* [%s](%s) - %s\n", jiraIssues[0].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssues[0].Key), jiraIssues[0].Fields.Summary) +
				fmt.Sprintf("* [%s](%s) - %s\n", jiraIssues[1].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssues[1].Key), jiraIssues[1].Fields.Summary),
		},
		{
			Name: "Convert with estimated time",
			FuncArguments: []interface{}{
				jiraIssuesForEstimateCase,
				true,
			},
			Setup: func() (*test.MockClients, *utils.Config) {
				return test.InitDefaultMockClients(mc), nil
			},
			ExpectedResult: fmt.Sprintf("* [%s](%s) - %s", jiraIssuesForEstimateCase[0].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssuesForEstimateCase[0].Key), jiraIssuesForEstimateCase[0].Fields.Summary) + " [2h 30m]\n" +
				fmt.Sprintf("* [%s](%s) - %s", jiraIssuesForEstimateCase[1].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssuesForEstimateCase[1].Key), jiraIssuesForEstimateCase[1].Fields.Summary) + " [1w 1d 1h 1m 5s]\n" +
				fmt.Sprintf("* [%s](%s) - %s", jiraIssuesForEstimateCase[2].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssuesForEstimateCase[2].Key), jiraIssuesForEstimateCase[2].Fields.Summary) + " [no estimate]\n" +
				"*Суммарно запланировано времени: 1w 1d 3h 31m 5s*",
		},
		{
			Name: "Convert with estimated time, but no estimate time supplied",
			FuncArguments: []interface{}{
				jiraIssues,
				true,
			},
			Setup: func() (*test.MockClients, *utils.Config) {
				return test.InitDefaultMockClients(mc), nil
			},
			ExpectedResult: fmt.Sprintf("* [%s](%s) - %s", jiraIssues[0].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssues[0].Key), jiraIssues[0].Fields.Summary) + " [no estimate]\n" +
				fmt.Sprintf("* [%s](%s) - %s", jiraIssues[1].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssues[1].Key), jiraIssues[1].Fields.Summary) + " [no estimate]\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockClients, _ := tc.Setup()
			impl := InitTestImplementator(mockClients)
			res := impl.TodoIssuesToReport(tc.FuncArguments[0].([]*jira.Issue), tc.FuncArguments[1].(bool))
			tc.CheckCase(t, res, nil)
		})
	}
}

func TestDoneIssuesToReport(t *testing.T) {
	t.Parallel()
	mc := gomock.NewController(t)

	jiraIssues := []*jira.Issue{
		{
			Key: "TEST-1000",
			Fields: &jira.IssueFields{
				Summary: "for test",
			},
		},
		{
			Key: "TEST-2000",
			Fields: &jira.IssueFields{
				Summary: "another test",
			},
		},
	}

	jiraIssuesForEstimateCase := []*jira.Issue{
		{
			Key: "TEST-1000",
			Fields: &jira.IssueFields{
				Summary: "for test",
				Worklog: &jira.Worklog{
					Worklogs: []jira.WorklogRecord{
						{
							// 2h 30m
							TimeSpentSeconds: 9000,
						},
					},
				},
			},
		},
		{
			Key: "TEST-2000",
			Fields: &jira.IssueFields{
				Summary: "another test",
				Worklog: &jira.Worklog{
					Worklogs: []jira.WorklogRecord{
						{
							// 1w 1d 1h 1m 5s
							TimeSpentSeconds: 694865,
						},
					},
				},
			},
		},
		{
			Key: "TEST-3000",
			Fields: &jira.IssueFields{
				Summary: "another test",
				Worklog: &jira.Worklog{
					Worklogs: []jira.WorklogRecord{
						{
							// no estimate
							TimeSpentSeconds: 0,
						},
					},
				},
			},
		},
	}

	cases := []test.Case{
		{
			Name: "Just convert",
			FuncArguments: []interface{}{
				jiraIssues,
				false,
			},
			Setup: func() (*test.MockClients, *utils.Config) {
				return test.InitDefaultMockClients(mc), nil
			},
			ExpectedResult: fmt.Sprintf("* [%s](%s) - %s\n", jiraIssues[0].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssues[0].Key), jiraIssues[0].Fields.Summary) +
				fmt.Sprintf("* [%s](%s) - %s\n", jiraIssues[1].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssues[1].Key), jiraIssues[1].Fields.Summary),
		},
		{
			Name: "Convert with logged time",
			FuncArguments: []interface{}{
				jiraIssuesForEstimateCase,
				true,
			},
			Setup: func() (*test.MockClients, *utils.Config) {
				mocks := test.InitDefaultMockClients(mc)
				for _, issue := range jiraIssuesForEstimateCase {
					mocks.JiraMockClient.EXPECT().GetIssueLogTimeForTheLastWorkDay(gomock.Any(), gomock.Any()).
						Return(issue.Fields.Worklog.Worklogs[0].TimeSpentSeconds)
				}
				return mocks, nil
			},
			ExpectedResult: fmt.Sprintf("* [%s](%s) - %s", jiraIssuesForEstimateCase[0].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssuesForEstimateCase[0].Key), jiraIssuesForEstimateCase[0].Fields.Summary) + " [log: 2h 30m]\n" +
				fmt.Sprintf("* [%s](%s) - %s", jiraIssuesForEstimateCase[1].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssuesForEstimateCase[1].Key), jiraIssuesForEstimateCase[1].Fields.Summary) + " [log: 1w 1d 1h 1m 5s]\n" +
				fmt.Sprintf("* [%s](%s) - %s", jiraIssuesForEstimateCase[2].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssuesForEstimateCase[2].Key), jiraIssuesForEstimateCase[2].Fields.Summary) + " [log: no time]\n" +
				"*Суммарно залогировано времени: 1w 1d 3h 31m 5s*",
		},
		{
			Name: "Convert with logged time, but no log time supplied",
			FuncArguments: []interface{}{
				jiraIssues,
				true,
			},
			Setup: func() (*test.MockClients, *utils.Config) {
				mocks := test.InitDefaultMockClients(mc)
				for range jiraIssues {
					mocks.JiraMockClient.EXPECT().GetIssueLogTimeForTheLastWorkDay(gomock.Any(), gomock.Any()).
						Return(0)
				}
				return mocks, nil
			},
			ExpectedResult: fmt.Sprintf("* [%s](%s) - %s", jiraIssues[0].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssues[0].Key), jiraIssues[0].Fields.Summary) + " [log: no time]\n" +
				fmt.Sprintf("* [%s](%s) - %s", jiraIssues[1].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssues[1].Key), jiraIssues[1].Fields.Summary) + " [log: no time]\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockClients, cfg := tc.Setup()
			impl := InitTestImplementator(mockClients)
			res := impl.DoneIssuesToReport(cfg, tc.FuncArguments[0].([]*jira.Issue), tc.FuncArguments[1].(bool))
			tc.CheckCase(t, res, nil)
		})
	}
}

func TestConvertSecToJiraFormat(t *testing.T) {
	t.Parallel()
	mc := gomock.NewController(t)

	setupFunc := func() (*test.MockClients, *utils.Config) {
		return test.InitDefaultMockClients(mc), nil
	}

	cases := []test.Case{
		{
			Name: "Full time format",
			FuncArguments:[]interface{}{
				// 6w + 5d + 3h + 10m + 50s
				3628800 + 432000 + 10800 + 600 + 50,
			},
			Setup: setupFunc,
			ExpectedResult: "6w 5d 3h 10m 50s",
		},
		{
			Name: "Without hours",
			FuncArguments:[]interface{}{
				// 6w + 5d + 10m + 50s
				3628800 + 432000 + 600 + 50,
			},
			Setup: setupFunc,
			ExpectedResult: "6w 5d 10m 50s",
		},
		{
			Name: "Without weeks, hours",
			FuncArguments:[]interface{}{
				// 5d + 10m + 50s
				432000 + 600 + 50,
			},
			Setup: setupFunc,
			ExpectedResult: "5d 10m 50s",
		},
		{
			Name: "Only minutes and seconds",
			FuncArguments:[]interface{}{
				// 10m + 50s
				600 + 50,
			},
			Setup: setupFunc,
			ExpectedResult: "10m 50s",
		},
		{
			Name: "Empty case",
			FuncArguments:[]interface{}{
				// empty string
				0,
			},
			Setup: setupFunc,
			ExpectedResult: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockClients, _ := tc.Setup()
			impl := InitTestImplementator(mockClients)
			res := impl.ConvertSecToJiraFormat(tc.FuncArguments[0].(int))
			tc.CheckCase(t, res, nil)
		})
	}
}
