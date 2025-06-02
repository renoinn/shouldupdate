package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config stores application names as keys and their versions as values.
type Config map[string]string

var configFile string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Warning: Error getting user home directory: %v. Using current directory for config file.", err)
		configFile = "versions.toml" // Fallback to current directory
	} else {
		configFile = filepath.Join(homeDir, ".config", "shouldupdate", "versions.toml")
	}
}

// loadConfig loads the configuration from the configFile.
// If the file doesn't exist, it returns an empty Config.
func loadConfig() (Config, error) {
	config := make(Config)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// This log message is for debug/startup; user feedback is handled by command handlers.
		log.Printf("Info: Config file '%s' not found. A new one will be created upon adding an application.", configFile)
		return config, nil // Return empty config, it will be saved on first 'add'
	}

	// Read the file content
	data, err := os.ReadFile(configFile)
	if err != nil {
		// Log the error for debugging, but return a user-friendly one.
		log.Printf("Debug: Error reading config file %s: %v", configFile, err)
		return nil, fmt.Errorf("could not read config file '%s': %w", configFile, err)
	}

	// Decode the TOML data
	if err := toml.Unmarshal(data, &config); err != nil {
		// Log the error for debugging.
		log.Printf("Debug: Error unmarshalling TOML from %s: %v", configFile, err)
		return nil, fmt.Errorf("could not parse config file '%s' (TOML format error): %w", configFile, err)
	}

	// log.Printf("%sConfig loaded successfully from %s%s\n", colorGreen, configFile, colorReset) // Less verbose, success is implicit.
	return config, nil
}

// saveConfig saves the configuration to the configFile.
func saveConfig(config Config) error {
	// Marshal the config map to TOML []byte
	data, err := toml.Marshal(config)
	if err != nil {
		log.Printf("Debug: Error marshalling config to TOML: %v", err)
		return fmt.Errorf("could not format configuration for saving: %w", err)
	}

	// Ensure the directory structure exists
	dirPath := filepath.Dir(configFile)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Printf("Debug: Error creating directory structure %s: %v", dirPath, err)
		return fmt.Errorf("could not create config directory '%s': %w", dirPath, err)
	}

	// Write the TOML data to the file
	// os.WriteFile truncates the file if it exists, and creates it if it doesn't.
	err = os.WriteFile(configFile, data, 0644) // 0644 are standard file permissions
	if err != nil {
		log.Printf("Debug: Error writing config to file %s: %v", configFile, err)
		return fmt.Errorf("could not write configuration to file '%s': %w", configFile, err)
	}

	// log.Printf("%sConfig saved successfully to %s%s\n", colorGreen, configFile, colorReset) // Less verbose, success is implicit.
	return nil
}
