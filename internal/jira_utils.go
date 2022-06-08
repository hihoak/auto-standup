package internal

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/hihoak/auto-standup/pkg/utils"
)

func (i *Implementator) IssuesToStr(issues []*jira.Issue) string {
	strIssues := ""
	for _, issue := range issues {
		strIssues += fmt.Sprintf("* [%s](%s) - %s\n", issue.Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", issue.Key), issue.Fields.Summary)
	}
	return strIssues
}

func (i *Implementator) GetIssuesFromLastWorkDay(cfg *utils.Config) ([]*jira.Issue, error) {
	jql := fmt.Sprintf("updatedDate >= \"-%dd\" AND assignee = %s", cfg.NumberOfDaysForGetTickets, cfg.Username)
	utils.Log.Debug().Msgf("Searching tickets with following JQL %s", jql)
	issuesFromLastWorkDay, _, err := i.JiraClient.SearchIssue(jql, &jira.SearchOptions{Expand: "changelog"})
	if err != nil {
		utils.Log.Debug().Err(err).Msg("Can't get issues from search")
		return nil, err
	}
	utils.Log.Debug().Msgf("Got issues from search with jql '%s', issus: %v", jql, issuesFromLastWorkDay)
	issuesThatRealWasInWork := make([]*jira.Issue, 0)
	for idx, issue := range issuesFromLastWorkDay {
		if i.Filters.FilterIssuesByProject(cfg, &issue) && i.Filters.FilterIssueByActivity(cfg, &issue) {
			issuesThatRealWasInWork = append(issuesThatRealWasInWork, &issuesFromLastWorkDay[idx])
		}
	}
	return issuesThatRealWasInWork, nil
}

func (i *Implementator) FromStrKeysToIssues(ctx context.Context, issueKeys []string) ([]*jira.Issue, error) {
	var issues []*jira.Issue
	for _, key := range issueKeys {
		issue, _, err := i.JiraClient.GetIssue(key, &jira.GetQueryOptions{
			Expand: "changelog",
		})
		if err != nil {
			utils.Log.Debug().Err(err).Msgf("Can't get issue %s", key)
			return nil, fmt.Errorf("can't get issue %s. error: %w", key, err)
		}
		utils.Log.Debug().Msgf("Got issue %v", issue)
		issues = append(issues, issue)
	}
	utils.Log.Debug().Msgf("Got following issues: %v", issues)
	return issues, nil
}
