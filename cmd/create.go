package cmd

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/hihoak/auto-standup/pkg/utils"
	"github.com/spf13/cobra"
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
	output := "Что вы делали с прошлого опроса?\n"
	doneTasksString := ""
	var doneIssues []*jira.Issue
	var err error
	if len(doneTickets) == 0 {
		utils.Log.Debug().Msg("Flag Done Tickets is empty, starting logic with getting issues from last work day...")
		doneIssues, err = impl.GetIssuesFromLastWorkDay(utils.Cfg)
		if err != nil {
			utils.Log.Fatal().Err(err).Msg("Failed to get tickets from last work day")
		}
	} else {
		utils.Log.Debug().Msg("Done tickets provided via flag, getting them by keys...")
		doneIssues, err = impl.FromStrKeysToIssues(context.TODO(), doneTickets)
		if err != nil {
			utils.Log.Fatal().Err(err).Msg("Failed to get DONE tickets from jira")
		}
	}
	doneTasksString = impl.IssuesToStr(doneIssues)
	if err != nil {
		utils.Log.Fatal().Err(err).Msg("got error while trying to get 'Done' tasks report")
	}
	if doneTasksString == "" {
		output += "Ничего не было сделано"
	} else {
		output += doneTasksString
	}

	output += "\nЧто вы будете делать до следующего опроса?\n"
	toDoIssues, err := impl.FromStrKeysToIssues(context.TODO(), toDoTickets)
	if err != nil {
		utils.Log.Fatal().Err(err).Msg("failed to get TODO issues from jira")
	}
	toDoTasksString := impl.IssuesToStr(toDoIssues)
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
