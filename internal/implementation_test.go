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
			t.Parallel()
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
			t.Parallel()
			mockClients, cfg := tc.Setup()
			impl := InitTestImplementator(mockClients)
			res, err := impl.GetIssuesFromLastWorkDay(cfg)
			tc.CheckCase(t, res, err)
		})
	}
}

func TestIssuesToStr(t *testing.T) {
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

	cases := []test.Case{
		{
			Name: "Just convert",
			FuncArguments: []interface{}{
				jiraIssues,
			},
			Setup: func() (*test.MockClients, *utils.Config) {
				return test.InitDefaultMockClients(mc), nil
			},
			ExpectedResult: fmt.Sprintf("* [%s](%s) - %s\n", jiraIssues[0].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssues[0].Key), jiraIssues[0].Fields.Summary) +
				fmt.Sprintf("* [%s](%s) - %s\n", jiraIssues[1].Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", jiraIssues[1].Key), jiraIssues[1].Fields.Summary),
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			mockClients, _ := tc.Setup()
			impl := InitTestImplementator(mockClients)
			res := impl.IssuesToStr(tc.FuncArguments[0].([]*jira.Issue))
			tc.CheckCase(t, res, nil)
		})
	}
}
