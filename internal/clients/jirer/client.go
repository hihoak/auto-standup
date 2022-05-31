package jirer

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"standup-cli/pkg/utils"
	"strings"
)

type Jira struct {
	client                 *jira.Client
}

func NewJiraClient(username, password string) (*Jira, error) {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	client, err := jira.NewClient(tp.Client(), "https://jit.ozon.ru")
	if err != nil {
		utils.Log.Debug().Err(err).Msg("Can't initialize Jira client!")
		return nil, fmt.Errorf("can't initialize Jira client. Error: %w", err)
	}

	return &Jira{
		client: client,
	}, nil
}

func (j *Jira) FromStrKeysToIssues(issueKeys []string) ([]*jira.Issue, error) {
	var issues []*jira.Issue
	for _, key := range issueKeys {
		issue, _, err := j.client.Issue.Get(key, &jira.GetQueryOptions{
			Expand: "changelog",
		})
		if err != nil {
			utils.Log.Debug().Err(err).Msgf("Can't get issue %s", key)
			return nil, fmt.Errorf("can't get issue %s. error: %w", key, err)
		}
		utils.Log.Debug().Msgf("Got issue %w", issue)
		issues = append(issues, issue)
	}
	utils.Log.Debug().Msgf("Got following issues: %v", issues)
	return issues, nil
}

func (j *Jira) IssuesToStr(issues []*jira.Issue) string {
	strIssues := ""
	for _, issue := range issues {
		strIssues += fmt.Sprintf("* [%s](%s) - %s\n", issue.Key, fmt.Sprintf("https://jit.ozon.ru/browse/%s", issue.Key), issue.Fields.Summary)
	}
	return strIssues
}

func (j *Jira) GetIssuesFromLastWorkDay(cfg *utils.Config) ([]*jira.Issue, error) {
	jql := fmt.Sprintf("updatedDate >= \"-%dd\" AND assignee = %s", cfg.NumberOfDaysForGetTickets, cfg.Username)
	utils.Log.Debug().Msgf("Searching tickets with following JQL %s", jql)
	issuesFromLastWorkDay, _, err := j.client.Issue.Search(jql, &jira.SearchOptions{Expand: "changelog"})
	if err != nil {
		utils.Log.Debug().Err(err).Msg("Can't get issues from search")
		return nil, err
	}
	utils.Log.Debug().Msgf("Got issues from search with jql '%s', issus: %v", jql, issuesFromLastWorkDay)
	issuesThatRealWasInWork := make([]*jira.Issue, 0)
	for idx, issue := range issuesFromLastWorkDay {
		if j.filterIssuesByProject(cfg, &issue) && j.filterIssueByActivity(cfg, &issue) {
			issuesThatRealWasInWork = append(issuesThatRealWasInWork, &issuesFromLastWorkDay[idx])
		}
	}
	return issuesThatRealWasInWork, nil
}

func (j *Jira) filterIssuesByProject(cfg *utils.Config, issue *jira.Issue) bool {
	utils.Log.Debug().Msgf("Check that issue not in exclude projects: %v", issue)
	for _, excludeProject := range cfg.ExcludeJiraProjects {
		if strings.EqualFold(utils.GetProjectFromIssueKey(issue.Key), excludeProject) {
			utils.Log.Debug().Msgf("Ticket %s in exclude projects", issue.Key)
			return false
		}
	}
	utils.Log.Debug().Msgf("Ticket %s not in exclude projects", issue.Key)
	return true
}

func (j *Jira) filterIssueByActivity(cfg *utils.Config, issue *jira.Issue) bool {
	var validHistory *jira.ChangelogHistory
	utils.Log.Debug().Msgf("Start checking issues '%s' following histories: %v", issue.Key, issue.Changelog.Histories)
	for i := len(issue.Changelog.Histories) - 1; i >= 0; i-- {
		history := issue.Changelog.Histories[i]
		utils.Log.Debug().Msgf("Start checking issues '%s' history that created '%s' author: %v", issue.Key, history.Created, history.Author)
		valid := false
		for _, eligibleUsersForUpdate := range cfg.EligibleUsersHistories {
			if history.Author.Name == eligibleUsersForUpdate {
				validHistory = &history
				utils.Log.Debug().Msgf("For issue '%s' found valid latest history: %v", issue.Key, *validHistory)
				valid = true
				break
			}
		}
		if !valid {
			utils.Log.Debug().Msg("It's not valid history")
		} else {
			break
		}
	}

	if validHistory == nil {
		utils.Log.Debug().Msgf("Not found valid histories for ticket %s", issue.Key)
		return false
	}

	startTime := cfg.CmdStartTime.AddDate(0, 0, -cfg.NumberOfDaysForGetTickets)
	endTime := cfg.CmdStartTime
	historyTimeCreation, err := validHistory.CreatedTime()
	if err != nil {
		utils.Log.Debug().Err(err).Msgf("Can't parse created time %s", validHistory.Created)
		return false
	}

	if startTime.Before(historyTimeCreation) && endTime.After(historyTimeCreation) {
		utils.Log.Debug().Msgf("Issues %s is valid, add it to report", issue.Key)
		return true
	}
	utils.Log.Debug().Msgf("Issues %s is not valid, last activity was before last work day ", issue.Key)
	return false
}
