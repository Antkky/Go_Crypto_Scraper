package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// readConfig reads and unmarshals the configuration file.
func ReadConfig(filePath string) ([]ExchangeConfig, error) {
	rawConfig, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var configs []ExchangeConfig
	if err := json.Unmarshal(rawConfig, &configs); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return configs, nil
}
