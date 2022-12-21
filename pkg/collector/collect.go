package collector

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func CollectNodeData(cmd *cobra.Command, args []string) error {
	shellCmd := NewShellCmd()
	nodeType, err := shellCmd.FindNodeType()
	if err != nil {
		return err
	}
	infoCollectorMap, err := LoadConfig()
	if err != nil {
		return err
	}
	for _, infoCollector := range infoCollectorMap {
		nodeInfo := make(map[string]interface{})
		for _, ci := range infoCollector.Collectors {
			if ci.NodeType != nodeType && nodeType != MasterNode {
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
			filterdParts := make([]string, 0)
			for _, part := range outputParts {
				if len(part) == 0 || part == "[^\"]\\S*'" {
					continue
				}
				filterdParts = append(filterdParts, part)
			}
			if len(filterdParts) > 0 {
				nodeInfo[ci.Key] = filterdParts
			}
		}
		nodeData := Node{
			APIVersion: Version,
			Kind:       Kind,
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
}
