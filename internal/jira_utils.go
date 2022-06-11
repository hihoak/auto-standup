package internal

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/hihoak/auto-standup/pkg/utils"
	"strings"
)

// DoneIssuesToReport - ...
func (i *Implementator) DoneIssuesToReport(cfg *utils.Config, issues []*jira.Issue, addLogTime bool) string {
	strIssues := ""
	totalLogTime := 0
	for _, issue := range issues {
		strIssues += fmt.Sprintf("* [%s](%s) - %s", issue.Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", issue.Key), issue.Fields.Summary)
		if addLogTime {
			issueLogTime := i.JiraClient.GetIssueLogTimeForTheLastWorkDay(cfg, *issue)
			totalLogTime += issueLogTime
			strLoggedTime := i.ConvertSecToJiraFormat(issueLogTime)
			if strLoggedTime == "" {
				strLoggedTime = "no time"
			}
			strIssues += fmt.Sprintf(" [log: %s]", strLoggedTime)
		}
		strIssues += "\n"
	}
	if addLogTime {
		strTotalLogTime := i.ConvertSecToJiraFormat(totalLogTime)
		if strTotalLogTime != "" {
			strIssues += fmt.Sprintf("*Суммарно залогировано времени: %s*", strTotalLogTime)
		}
	}
	return strIssues
}

// TodoIssuesToReport - ...
func (i *Implementator) TodoIssuesToReport(issues []*jira.Issue, addRemainingTime bool) string {
	strIssues := ""
	totalTime := 0
	for _, issue := range issues {
		strIssues += fmt.Sprintf("* [%s](%s) - %s", issue.Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", issue.Key), issue.Fields.Summary)
		if addRemainingTime {
			totalTime += issue.Fields.TimeEstimate
			estimateTime := i.ConvertSecToJiraFormat(issue.Fields.TimeEstimate)
			if estimateTime == "" {
				estimateTime = "no estimate"
			}
			strIssues += fmt.Sprintf(" [%s]", estimateTime)
		}
		strIssues += "\n"
	}
	if addRemainingTime {
		totalEstimateTime := i.ConvertSecToJiraFormat(totalTime)
		if totalEstimateTime != "" {
			strIssues += fmt.Sprintf("*Суммарно запланировано времени: %s*", totalEstimateTime)
		}
	}
	return strIssues
}

func (i *Implementator) ConvertSecToJiraFormat(sec int) string {
	weeks := sec / 60 / 60 / 24 / 7
	sec -= weeks * (60 * 60 * 24 * 7)
	days := sec / 60 / 60 / 24
	sec -= days * (60 * 60 * 24)
	hours := sec / 60 / 60
	sec -= hours * (60 * 60)
	minutes := sec / 60
	sec -= minutes * 60

	resultString := ""
	if weeks != 0 {
		resultString += fmt.Sprintf("%dw", weeks)
	}
	if days != 0 {
		resultString += fmt.Sprintf(" %dd", days)
	}
	if hours != 0 {
		resultString += fmt.Sprintf(" %dh", hours)
	}
	if minutes != 0 {
		resultString += fmt.Sprintf(" %dm", minutes)
	}
	if sec != 0 {resultString += fmt.Sprintf(" %ds", sec)
	}
	return strings.TrimSpace(resultString)
}

// GetIssuesFromLastWorkDay - ...
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
		if i.Filters.FilterIssuesByProject(cfg, issue) && i.Filters.FilterIssueByActivity(cfg, issue) {
			issuesThatRealWasInWork = append(issuesThatRealWasInWork, &issuesFromLastWorkDay[idx])
		}
	}
	return issuesThatRealWasInWork, nil
}

// FromStrKeysToIssues - ...
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
