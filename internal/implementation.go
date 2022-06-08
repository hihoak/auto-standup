package internal

import (
	"github.com/hihoak/auto-standup/internal/clients/jirer"
	"github.com/hihoak/auto-standup/internal/filters"
)

type Implementator struct {
	JiraClient jirer.Jirer
	Filters filters.Filterers
}

func NewImplementator(jiraClient jirer.Jirer, filters filters.Filterers) *Implementator {
	return &Implementator{
		JiraClient: jiraClient,
		Filters: filters,
	}
}
