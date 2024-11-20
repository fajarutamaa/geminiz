package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Generate text using generative AI models",
	Args:  cobra.NoArgs,
	RunE:  runPrompt,
}

func init() {
	rootCmd.AddCommand(promptCmd)
}

func runPrompt(cmd *cobra.Command, args []string) error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		prompts, err := collectPrompts(scanner)
		if err != nil {
			return err
		}

		if len(prompts) == 0 {
			return nil
		}

		if err := generateResponse(prompts); err != nil {
			color.Red("Error: %v", err)
			continue
		}

		if !continueSession(scanner) {
			return nil
		}
	}
}

func collectPrompts(scanner *bufio.Scanner) ([]string, error) {
	var prompts []string
	color.Yellow("Enter your prompts one by one. Type 'done' to finish and generate, or 'q' to exit.\n")

	for {
		fmt.Print("Enter your prompt: ")
		if !scanner.Scan() {
			return nil, fmt.Errorf("failed to read input")
		}

		input := strings.TrimSpace(scanner.Text())
		switch input {
		case "q":
			color.Green("Exiting Geminiz. Goodbye!")
			return nil, nil
		case "done":
			if len(prompts) == 0 {
				color.Red("No prompts entered.")
				return nil, nil
			}
			return prompts, nil
		case "":
			color.Red("Prompt cannot be empty. Please try again.")
			continue
		default:
			prompts = append(prompts, input)
			color.Cyan("Prompt added: %s", input)
		}
	}
}

func generateResponse(prompts []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	promptText := strings.Join(prompts, ". ")
	color.Green("Your full prompt: %s", promptText)

	client, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return fmt.Errorf("failed to initialize AI client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(config.Model)
	result, err := model.GenerateContent(ctx, genai.Text(promptText))
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("no response generated")
	}

	response := result.Candidates[0].Content.Parts[0]
	respText := strings.ReplaceAll(fmt.Sprintf("%v", response), "*", "")
	color.Yellow("\nGenerated Response:\n\n")
	color.White(respText)
	return nil
}

func continueSession(scanner *bufio.Scanner) bool {
	fmt.Print("\nDo you want to continue? (Type 'q' to quit, or press Enter to continue): ")
	if !scanner.Scan() {
		return false
	}
	return strings.TrimSpace(scanner.Text()) != "q"
}
