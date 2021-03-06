package cmd

import (
	"fmt"
	"os"

	"github.com/hihoak/auto-standup/internal"
	"github.com/hihoak/auto-standup/internal/clients/jirer"
	"github.com/hihoak/auto-standup/internal/filters"
	"github.com/hihoak/auto-standup/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	impl       *internal.Implementator
	logLevel   string
	configPath string

	username string
	password string

	rootCmd = &cobra.Command{
		Use:              "standup [--log-level 'info']",
		Short:            "cli for generating standup message",
		TraverseChildren: true,
	}
)

// Execute - ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initUtils, initConfig, initImplementator)
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "setting log level for cli. Options: 'debug', 'info'")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config-path", "c", "", "path to your config file in YAML format (default is $HOME/.standup.yaml)")

	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "your username in jira (hint: you can create config in $HOME/.standup.yaml)")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "your password in jira (hint: you can create config in $HOME/.standup.yaml)")
}

func initUtils() {
	utils.Log = utils.NewLogger(logLevel)
	utils.Log.Debug().Msg("Successfully initialize logger...")
}

func initConfig() {
	if username != "" && password != "" {
		utils.Cfg = &utils.Config{
			Username: username,
			Password: password,
		}
		return
	}
	var err error
	if configPath == "" {
		homeDirectory, _ := os.LookupEnv("HOME")
		configPath = homeDirectory + "/.standup.yaml"
	}
	utils.Cfg, err = utils.NewConfig(configPath, addEstimatedTime, addLogTime)
	if err != nil {
		utils.Log.Fatal().Err(err).Msg("failed to init config. Create it in $HOME/.standup.yaml or supply flags")
		return
	}
}

func initImplementator() {
	utils.Log.Debug().Msg("Initializing Jira client...")
	jiraClient, err := jirer.New(utils.Cfg.Username, utils.Cfg.Password)
	if err != nil {
		utils.Log.Fatal().Err(err).Msg("got error while creating Jira client")
	}
	utils.Log.Debug().Msg("Successfully initialize Jira client!")

	impl = internal.NewImplementator(jiraClient, &filters.Filters{})
}
