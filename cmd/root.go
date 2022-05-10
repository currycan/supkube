/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/currycan/supkube/cert"
	"github.com/currycan/supkube/core"
	"github.com/currycan/supkube/pkg/logger"
)

var (
	cfgFile string
	Info    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "supkube",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.supkube/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&Info, "info", false, "logger ture for Info, false for Debug")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home := cert.GetUserHomeDir()
	logFile := fmt.Sprintf("%s/.supkube/supkube.log", home)
	if !core.FileExist(home + "/.supkube") {
		err := os.MkdirAll(home+"/.supkube", os.ModePerm)
		if err != nil {
			fmt.Println("create default supkube config dir failed, please create it by your self mkdir -p /root/.supkube && touch /root/.supkube/config.yaml")
		}
	}
	if Info {
		logger.Cfg(5, logFile)
	} else {
		logger.Cfg(6, logFile)
	}
}
