package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// CollectNodeData run spec audit command and output it result data
func CollectNodeData(cmd *cobra.Command) error {
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
		outputFormat := cmd.Flag("output").Value.String()
		err := printOutput(nodeData, outputFormat, os.Stdout)
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

func printOutput(nodeData Node, output string, writer io.Writer) error {
	switch output {
	case "json":
		data, err := json.Marshal(nodeData)
		if err != nil {
			return err
		}
		fmt.Fprint(writer, string(data))
	case "table":
		data := make([][]string, 0)
		for key, ndata := range nodeData.Info {
			results, ok := ndata.([]string)
			if !ok {
				return fmt.Errorf("no data found")
			}
			if len(results) > 0 {
				joinedResults := join(results...)

				data = append(data, []string{key, joinedResults})
			}

		}
		table := tablewriter.NewWriter(writer)
		table.SetHeader([]string{"Key", "Value"})
		table.SetBorder(false) // Set Border to false
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
	}
	return nil
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(fmt.Sprintf(",%s", strings.TrimSpace(str)))
	}
	return strings.TrimPrefix(sb.String(), ",")
}
