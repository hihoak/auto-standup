package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"standup-cli/pkg/utils"
)

var (
	logLevel string
	configPath string

	username string
	password string

	rootCmd = &cobra.Command{
		Use:   "standup [--log-level 'info']",
		Short: "cli for generating standup message",
		TraverseChildren: true,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initUtils, initConfig)
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
	utils.Cfg, err = utils.NewConfig(configPath)
	if err != nil {
		utils.Log.Fatal().Err(err).Msg("failed to init config. Create it in $HOME/.standup.yaml or supply flags")
		return
	}
}
