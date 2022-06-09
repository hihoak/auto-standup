package filters

import (
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/hihoak/auto-standup/pkg/utils"
)

//go:generate mockgen -destination "./mocks/mock_filterer.go" -source "./filters.go" -package "mocks" Filterers

// Filterers - ...
type Filterers interface {
	FilterIssueByActivity(cfg *utils.Config, issue jira.Issue) bool
	FilterIssuesByProject(cfg *utils.Config, issue jira.Issue) bool
}

// Filters - ...
type Filters struct{}

// FilterIssueByActivity - ...
func (f *Filters) FilterIssueByActivity(cfg *utils.Config, issue jira.Issue) bool {
	utils.Log.Debug().Msgf("Start FilterIssueByActivity for issue: %v", issue)
	if issue.Changelog == nil {
		utils.Log.Debug().Msgf("Changelog is empty")
		return false
	}
	var validHistory *jira.ChangelogHistory
	utils.Log.Debug().Msgf("Start checking issues '%s' following histories: %v", issue.Key, issue.Changelog.Histories)
	for i := len(issue.Changelog.Histories) - 1; i >= 0; i-- {
		history := issue.Changelog.Histories[i]
		utils.Log.Debug().Msgf("Start checking issues '%s' history that created '%s' author: %v", issue.Key, history.Created, history.Author)
		valid := false
		for _, eligibleUsersForUpdate := range cfg.EligibleUsersHistories {
			if strings.EqualFold(history.Author.Name, eligibleUsersForUpdate) {
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

// FilterIssuesByProject - ...
func (f *Filters) FilterIssuesByProject(cfg *utils.Config, issue jira.Issue) bool {
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
