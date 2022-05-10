/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/currycan/supkube/core"
)

// ipvsCmd represents the ipvs command
var ipvsCmd = &cobra.Command{
	Use:   "ipvs",
	Short: "create or care local ipvs lb",
	Run: func(cmd *cobra.Command, args []string) {
		core.Ipvs.VsAndRsCare()
	},
}

func init() {
	rootCmd.AddCommand(ipvsCmd)

	// Here you will define your flags and configuration settings.
	ipvsCmd.Flags().BoolVar(&core.Ipvs.RunOnce, "run-once", false, "is run once mode")
	ipvsCmd.Flags().BoolVarP(&core.Ipvs.Clean, "clean", "c", true, " clean Vip ipvs rule before join node, if Vip has no ipvs rule do nothing.")
	ipvsCmd.Flags().StringVar(&core.Ipvs.VirtualServer, "vs", "", "virturl server like 10.54.0.2:6443")
	ipvsCmd.Flags().StringSliceVar(&core.Ipvs.RealServer, "rs", []string{}, "virturl server like 192.168.0.2:6443")

	ipvsCmd.Flags().StringVar(&core.Ipvs.HealthPath, "health-path", "/healthz", "health check path")
	ipvsCmd.Flags().StringVar(&core.Ipvs.HealthSchem, "health-schem", "https", "health check schem")
	ipvsCmd.Flags().Int32Var(&core.Ipvs.Interval, "interval", 5, "health check interval, unit is sec.")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ipvsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ipvsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
