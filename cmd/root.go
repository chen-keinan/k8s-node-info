package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "node-info",
	Short: "k8s-Node-Info extract file system info from cluster Node",
	Long:  `A tool which provide a way to extract file system info fro node cluster based on pre-define commands`,
	Run: func(cmd *cobra.Command, args []string) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cm := exec.Command("stat", "-c", "%U:%G", "/etc/kubernetes/manifests/kube-controller-manager.yaml")
		cm.Stdout = &stdout
		cm.Stderr = &stderr
		err := cm.Run()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print(stdout.String())
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
