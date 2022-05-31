package cmd

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/spf13/cobra"
	"auto-standup/internal/clients/jirer"
	"auto-standup/pkg/utils"
)

var (
	doneTickets []string
	toDoTickets []string

	createCmd = &cobra.Command{
		Use: "create --log-level \"info\" --done \"RE-1000,RE-2000\" --todo \"RE-3000,RE-4000\"",
		Short: "creating standup message",
		Run: create,
	}
)

func init() {
	createCmd.PersistentFlags().StringSliceVarP(&doneTickets, "done", "d", []string{}, "enter tasks that was done before previous survey")
	createCmd.PersistentFlags().StringSliceVarP(&toDoTickets, "to-do", "t", []string{}, "enter tasks that you plan to do before the next survey")
	rootCmd.AddCommand(createCmd)
}

func create(_ *cobra.Command, _ []string) {
	utils.Log.Debug().Msg("Initializing Jira client...")
	jiraClient, err := jirer.NewJiraClient(utils.Cfg.Username, utils.Cfg.Password)
	if err != nil {
		utils.Log.Fatal().Err(err).Msg("got error while creating Jira client")
	}
	utils.Log.Debug().Msg("Successfully initialize Jira client!")

	output := "Что вы делали с прошлого опроса?\n"
	doneTasksString := ""
	var doneIssues []*jira.Issue
	if len(doneTickets) == 0 {
		utils.Log.Debug().Msg("Flag Done Tickets is empty, starting logic with getting issues from last work day...")
		doneIssues, err = jiraClient.GetIssuesFromLastWorkDay(utils.Cfg)
		if err != nil {
			utils.Log.Fatal().Err(err).Msg("Failed to get tickets from last work day")
		}
	} else {
		utils.Log.Debug().Msg("Done tickets provided via flag, getting them by keys...")
		doneIssues, err = jiraClient.FromStrKeysToIssues(doneTickets)
		if err != nil {
			utils.Log.Fatal().Err(err).Msg("Failed to get DONE tickets from jira")
		}
	}
	doneTasksString = jiraClient.IssuesToStr(doneIssues)
	if err != nil {
		utils.Log.Fatal().Err(err).Msg("got error while trying to get 'Done' tasks report")
	}
	if doneTasksString == "" {
		output += "Ничего не было сделано"
	} else {
		output += doneTasksString
	}

	output += "\nЧто вы будете делать до следующего опроса?\n"
	toDoIssues, err := jiraClient.FromStrKeysToIssues(toDoTickets)
	if err != nil {
		utils.Log.Fatal().Err(err).Msg("failed to get TODO issues from jira")
	}
	toDoTasksString := jiraClient.IssuesToStr(toDoIssues)
	if err != nil {
		utils.Log.Fatal().Err(err).Msg("got error while trying to get 'ToDo' tasks report")
	}
	if toDoTasksString == "" {
		output += "Ничего не планирую делать"
	} else {
		output += toDoTasksString
	}

	fmt.Println(output)
}
