package collector

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CollectNodeData run spec audit command and output it result data
func CollectNodeData() error {
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
			filterdParts := filterAuditResults(output)
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
		err := printOutput(nodeData)
		if err != nil {
			return err
		}
	}
	return nil
}

func filterAuditResults(output string) []string {
	if len(output) == 0 {
		return []string{}
	}
	outputParts := strings.Split(output, ",")
	filterdParts := make([]string, 0)
	for _, part := range outputParts {
		if len(part) == 0 || part == "[^\"]\\S*'" {
			continue
		}
		filterdParts = append(filterdParts, part)
	}
	return filterdParts
}

func printOutput(nodeData Node) error {
	data, err := json.Marshal(nodeData)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
