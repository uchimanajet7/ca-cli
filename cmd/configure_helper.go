package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
)

var loadedConfig *config

type config struct {
	Endpoint string
	Profiles map[string]profile
}

type profile struct {
	APIKey   string
	Endpoint string
}

func getConfig() (*config, error) {
	if loadedConfig != nil {
		return loadedConfig, nil
	}

	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	config, err := loadConfig(path)
	if err != nil {
		return nil, err
	}

	loadedConfig = config

	return loadedConfig, nil
}

func getProfile(keyName string) (*profile, error) {
	configVal, err := getConfig()
	if err != nil {
		return nil, err
	}

	// set default
	if keyName == "" {
		keyName = getProfileName()
	}

	profileVal, ok := configVal.Profiles[keyName]
	if !ok {
		return nil, fmt.Errorf("The specified \"%s\" profile did not exist", keyName)
	}

	// set default
	if profileVal.Endpoint == "" {
		profileVal.Endpoint = configVal.Endpoint
	}

	return &profileVal, nil
}

func getConfigDir() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(dir, ".ca")

	return configDir, nil
}

func getConfigFileName() string {
	return "config"
}

func getProfileName() string {
	return "default"
}

func getEndpoint() string {
	return "https://manager.cloudautomator.com/api/v1"
}

func getConfigPath() (string, error) {
	dir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, getConfigFileName())

	return path, nil
}

func maskedAPIKey(apiKey string) string {
	prefix := apiKey[:3]
	suffix := apiKey[len(apiKey)-3:]

	length := len(apiKey) - (len(prefix) + len(suffix))
	masked := strings.Repeat("*", length)

	return prefix + masked + suffix
}

func drawConfigTable(keyName string, config *config) {
	tabelData := make([][]string, len(config.Profiles))

	// set table data
	for key, value := range config.Profiles {
		drawAPIKey := maskedAPIKey(value.APIKey)
		drawEndpoint := value.Endpoint

		if drawEndpoint == "" || drawEndpoint == config.Endpoint {
			drawEndpoint = "default"
		}

		if keyName != "" {
			if key == keyName {
				tabelData = append(tabelData, []string{key, drawAPIKey, drawEndpoint})
				break
			}
		} else {
			tabelData = append(tabelData, []string{key, drawAPIKey, drawEndpoint})
		}
	}

	if len(tabelData) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Profile", "API Key", "Endpoint"})
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.AppendBulk(tabelData)
		table.Render()
	}
}

func loadConfig(path string) (*config, error) {
	var config config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func saveConfig(path string, config *config) error {
	dir, _ := getConfigDir()

	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	writer := bufio.NewWriter(f)
	if err := toml.NewEncoder(writer).Encode(config); err != nil {
		return err
	}

	return nil
}
