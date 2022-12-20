package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/chen-keinan/k8s-node-info/pkg/collector"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "node-info",
	Short: "k8s-Node-Info extract file system info from cluster Node",
	Long:  `A tool which provide a way to extract k8s info which is not accessible via apiserver from node cluster based on pre-define commands`,
	RunE: func(cmd *cobra.Command, args []string) error {
		shellCmd := collector.NewShellCmd()
		nodeType, err := shellCmd.FindNodeType()
		if err != nil {
			return err
		}
		infoCollectorMap, err := collector.LoadConfig()
		if err != nil {
			return err
		}
		for _, infoCollector := range infoCollectorMap {
			nodeInfo := make(map[string]interface{})
			for _, ci := range infoCollector.Collectors {
				if ci.NodeType != nodeType && nodeType != collector.MasterNode {
					continue
				}
				output, err := shellCmd.Execute(ci.Audit)
				if err != nil {
					fmt.Print(err)
				}
				if len(output) == 0 {
					continue
				}
				outputParts := strings.Split(output, ",")
				nodeInfo[ci.Name] = outputParts
			}
			nodeData := collector.Node{
				APIVersion: collector.Version,
				Kind:       collector.Kind,
				Type:       nodeType,
				Info:       nodeInfo,
			}
			data, err := json.Marshal(nodeData)
			if err != nil {
				return err
			}
			fmt.Println(string(data))
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
