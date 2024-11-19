package cmd

import (
	"bufio"
	"context"
	"fmt"
	"geminiz/config"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var (
	GEMINI_API_KEY = config.GetEnv("GEMINI_API_KEY")
	GEMINI_MODEL   = config.GetEnv("GEMINI_MODEL")
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Generate text using generative AI models",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for {
			color.Yellow("Enter your prompts one by one. Type 'done' to finish and generate, or 'q' to exit.\n")
			var prompts []string
			scanner := bufio.NewScanner(os.Stdin)

			for {
				fmt.Print("Enter your prompt: ")
				if !scanner.Scan() {
					break
				}
				userInput := strings.TrimSpace(scanner.Text())

				if userInput == "q" {
					color.Green("Exiting Geminiz. Goodbye!")
					return
				}

				if userInput == "done" {
					if len(prompts) == 0 {
						color.Red("No prompts entered. Exiting.")
						return
					}
					break
				}

				if userInput == "" {
					color.Red("Prompt cannot be empty. Please try again.")
					continue
				}

				prompts = append(prompts, userInput)
				color.Cyan("Prompt added: %s", userInput)
			}

			color.Green("\nGenerating response for all prompts...\n")
			getResponse(prompts)

			fmt.Print("\nDo you want to continue? (Type 'q' to quit, or press Enter to continue): ")
			if scanner.Scan() {
				choice := strings.TrimSpace(scanner.Text())
				if choice == "q" {
					color.Green("Exiting Geminiz. Goodbye!")
					return
				}
			}
		}
	},
}

func getResponse(prompts []string) {
	promptArgs := strings.Join(prompts, ". ")
	ctx := context.Background()

	color.Green("Your full prompt: %s", promptArgs)
	client, err := genai.NewClient(ctx, option.WithAPIKey(GEMINI_API_KEY))
	if err != nil {
		color.Red("Error initializing Generative AI client: %v", err)
		os.Exit(1)
	}
	defer client.Close()

	model := client.GenerativeModel(GEMINI_MODEL)
	result, err := model.GenerateContent(ctx, genai.Text(promptArgs))
	if err != nil {
		color.Red("Error generating response: %v", err)
		os.Exit(1)
	}

	response := result.Candidates[0].Content.Parts[0]
	respFormat := strings.Replace(fmt.Sprintf("%s", response), "*", "", -1)
	color.Yellow("\nGenerated Response:\n\n")
	color.White(respFormat)
}

func init() {
	greeting := figure.NewColorFigure("Geminiz", "roman", "green", true)
	greeting.Print()
	color.Cyan("\nWelcome to Geminiz, a CLI tool powered by Generative AI models 🚀")
	color.Cyan("Author by Fajar Dwi Utomo 🖥️")
	color.Cyan("Use this tool to generate AI-powered responses effortlessly\n")
	color.Yellow("Type 'geminiz prompt' to get started\n")

	rootCmd.AddCommand(promptCmd)
	promptCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}
