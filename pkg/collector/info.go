package collector

import (
	"embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

const configFolder = "config"

//go:embed config
var config embed.FS

func LoadConfig() (map[string]*SpecInfo, error) {
	dirEntries, err := config.ReadDir(configFolder)
	if err != nil {
		return nil, err
	}
	specInfoMap := make(map[string]*SpecInfo)
	for _, entry := range dirEntries {
		fContent, err := config.ReadFile(fmt.Sprintf("%s/%s", configFolder, entry.Name()))
		if err != nil {
			return nil, err
		}
		si, err := getSpecInfo(string(fContent))
		if err != nil {
			return nil, err
		}
		specInfoMap[si.Name] = si
	}
	return specInfoMap, nil
}

type SpecInfo struct {
	Version    string      `yaml:"version"`
	Name       string      `yaml:"name"`
	Title      string      `yaml:"title"`
	Collectors []Collector `yaml:"collectors"`
}

type Collector struct {
	Name  string `yaml:"name"`
	Audit string `yaml:"audit"`
}

func getSpecInfo(info string) (*SpecInfo, error) {
	var specInfo SpecInfo
	err := yaml.Unmarshal([]byte(info), &specInfo)
	if err != nil {
		return nil, err
	}
	return &specInfo, nil
}
