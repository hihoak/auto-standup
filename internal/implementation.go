package internal

import (
	"context"

	"github.com/andygrunwald/go-jira"
	"github.com/hihoak/auto-standup/internal/clients/jirer"
	"github.com/hihoak/auto-standup/internal/filters"
	"github.com/hihoak/auto-standup/pkg/utils"
)

//go:generate mockgen -destination "./mocks/implementer_mock.go" -source "./implementation.go" -package "mocks" Implementer

// Implementer - ...
type Implementer interface {
	// IssuesToStr - ...
	IssuesToStr(issues []*jira.Issue) string
	// GetIssuesFromLastWorkDay - ...
	GetIssuesFromLastWorkDay(cfg *utils.Config) ([]*jira.Issue, error)
	// FromStrKeysToIssues - ...
	FromStrKeysToIssues(ctx context.Context, issueKeys []string) ([]*jira.Issue, error)
}

// Implementator - ...
type Implementator struct {
	JiraClient jirer.Jirer
	Filters    filters.Filterers
}

// NewImplementator - ...
func NewImplementator(jiraClient jirer.Jirer, filters filters.Filterers) *Implementator {
	return &Implementator{
		JiraClient: jiraClient,
		Filters:    filters,
	}
}
