package cmd

import (
	"fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "geminiz",
	Short: "A CLI tool for interacting with Google's Gemini AI model",
	Long: `Geminiz is a command-line interface tool that allows you to interact 
with Google's Gemini AI model. You can use it to generate text responses 
based on your prompts and manage your API key configuration.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	greeting := figure.NewColorFigure("Geminiz", "roman", "green", true)
	greeting.Print()
	color.Cyan("\nWelcome to Geminiz, a CLI tool powered by Generative AI models 🚀")
	color.Cyan("Author: Fajar Dwi Utomo 🖥️")
	color.Cyan("Use this tool to generate AI-powered responses effortlessly\n")
	color.Yellow("Type 'geminiz prompt' to get started\n")
}

func initConfig() {
	if err := loadConfig(); err != nil {
		if !isSetKeyCommand() {
			color.Red("Error: %v", err)
			color.Yellow("Please set your API key using: geminiz set key <YOUR-API-KEY>")
			os.Exit(1)
		}
	}
}

func isSetKeyCommand() bool {
	if len(os.Args) >= 3 {
		return os.Args[1] == "set" && os.Args[2] == "key"
	}
	return false
}
