package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/chen-keinan/k8s-node-info/pkg/collector"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "node-info",
	Short: "k8s-Node-Info extract file system info from cluster Node",
	Long:  `A tool which provide a way to extract file system info from node cluster based on pre-define commands`,
	RunE: func(cmd *cobra.Command, args []string) error {
		infoCollectorMap, err := collector.LoadConfig()
		if err != nil {
			return err
		}
		for _, infoCollector := range infoCollectorMap {
			for _, ci := range infoCollector.Collectors {
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				cm := exec.Command("sh", "-c",ci.Audit)
				cm.Stdout = &stdout
				cm.Stderr = &stderr
				err := cm.Run()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print(stdout.String())
			}
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
