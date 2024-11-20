package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Config holds the application configuration
type Config struct {
	APIKey string
	Model  string
}

var (
	config = &Config{
		Model: "gemini-1.5-flash",
	}
)

var setKeyCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration for Geminiz",
}

var setKeyApiCmd = &cobra.Command{
	Use:   "key [API_KEY]",
	Short: "Set the GEMINI_API_KEY",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetKey,
}

func init() {
	rootCmd.AddCommand(setKeyCmd)
	setKeyCmd.AddCommand(setKeyApiCmd)
}

func runSetKey(cmd *cobra.Command, args []string) error {
	apiKey := args[0]
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	if err := saveAPIKey(apiKey); err != nil {
		return fmt.Errorf("failed to save API key: %w", err)
	}

	color.Green("API key successfully set!")
	return nil
}

func saveAPIKey(apiKey string) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("unable to determine config directory: %w", err)
	}

	configDir = filepath.Join(configDir, "geminiz")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config")
	file, err := os.Create(configFile)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "GEMINI_API_KEY=%s\n", apiKey); err != nil {
		return fmt.Errorf("failed to write API key: %w", err)
	}

	return nil
}

func loadConfig() error {
	if envKey := os.Getenv("GEMINI_API_KEY"); envKey != "" {
		config.APIKey = envKey
		return nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("unable to determine config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "geminiz", "config")
	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("API key not found. Use 'geminiz set key <GEMINI_API_KEY>' to set it")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "GEMINI_API_KEY=") {
			config.APIKey = strings.TrimPrefix(line, "GEMINI_API_KEY=")
			return nil
		}
	}

	return fmt.Errorf("API key not found in config file")
}
