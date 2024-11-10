package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/nullswan/llama-hackaton/internal/chat"
	"github.com/nullswan/llama-hackaton/internal/code"
	"github.com/nullswan/llama-hackaton/internal/llama"
	"github.com/nullswan/llama-hackaton/internal/logger"
	"github.com/nullswan/llama-hackaton/internal/tools"

	"github.com/spf13/cobra"
)

var targetModel string

var rootCmd = &cobra.Command{
	Use:   "nomi [flags] [arguments]",
	Short: "Llama hackathon project",
	Run:   runApp,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func runApp(_ *cobra.Command, _ []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("Sig received, quitting...")
		cancel()
	}()

	selector := tools.NewSelector()
	toolsLogger := tools.NewLogger(
		true,
	)

	// Initialize Providers
	logger := logger.Init()

	conversation := chat.NewStackedConversation()

	inputHandler := tools.NewInputHandler(
		logger,
	)

	ttjProvider, err := initJSONProviders(
		targetModel,
	)
	if err != nil {
		fmt.Printf("Error initializing providers: %v\n", err)
		return
	}
	defer ttjProvider.Close()

	ttjBackend := tools.NewTextToJSONBackend(
		ttjProvider,
		logger,
	)

	err = interpreter(
		ctx,
		selector,
		toolsLogger,
		ttjBackend,
		inputHandler,
		conversation,
	)
	if err != nil {
		fmt.Printf("Error starting interpreter: %v\n", err)
		return
	}
}

// initJSONProviders initializes the text-to-json provider.
func initJSONProviders(
	targetModel string,
) (*llama.TextToJSONProvider, error) {
	backend, err := llama.LoadTextToJSONProvider(
		targetModel,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error loading text-to-text provider: %w",
			err,
		)
	}

	return backend, nil
}

const executionErrorLimit = 3

func interpreter(
	ctx context.Context,
	selector tools.Selector,
	logger tools.Logger,
	textToJSON tools.TextToJSONBackend,
	inputHandler tools.InputHandler,
	conversation *chat.Conversation,
) error {
	logger.Info("Starting console usecase")

	systemPrompt, err := getConsoleInstruction(
		runtime.GOOS,
	)
	if err != nil {
		return fmt.Errorf("failed to get console instruction: %w", err)
	}

	conversation.AddMessage(
		chat.NewMessage(
			chat.RoleSystem,
			systemPrompt,
		),
	)

	req, err := inputHandler.Read(ctx, ">>> ")
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	conversation.AddMessage(
		chat.NewMessage(
			chat.RoleUser,
			req,
		),
	)

	errorRetries := 0
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())
		default:
			// Handle too many errors
			if errorRetries > executionErrorLimit {
				fmt.Println("Too many errors, how can I help you?")
				resp, err := inputHandler.Read(ctx, ">>> ")
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}

				conversation.AddMessage(
					chat.NewMessage(
						chat.RoleUser,
						resp,
					),
				)

				errorRetries = 0
				continue
			}

			logger.Debug(
				"Calling Llama backend...",
			)
			resp, err := textToJSON.Do(ctx, conversation)
			if err != nil {
				return fmt.Errorf(
					"interpreter: error generating completion: %w",
					err,
				)
			}

			conversation.AddMessage(
				chat.NewMessage(
					chat.RoleAssistant,
					resp,
				),
			)

			var consoleResp consoleResponse
			if err := json.Unmarshal([]byte(resp), &consoleResp); err != nil {
				return fmt.Errorf("failed to unmarshal response: %w", err)
			}

			logger.Debug(
				"Received console response: " + resp,
			)

			switch consoleResp.Action {
			case consoleActionCode:
				// Sanitize code, add code block if necessary
				if consoleResp.Language != "" && consoleResp.Code != "" &&
					!strings.HasPrefix(consoleResp.Code, "```") {
					consoleResp.Code = "```" + consoleResp.Language + "\n" + consoleResp.Code + "\n```"
				}

				result := code.InterpretCodeBlocks(consoleResp.Code)

				if len(result) == 0 {
					logger.Info("No code blocks found")
					continue
				}

				containsError := true
				for _, r := range result {
					fmt.Printf(
						"Received (%d): %s\n%s\n",
						r.ExitCode,
						r.Stdout,
						r.Stderr,
					)

					if r.ExitCode == 0 {
						containsError = false
						break
					}
				}

				formattedResult := code.FormatExecutionResultForLLM(result)
				conversation.AddMessage(
					chat.NewMessage(
						chat.RoleAssistant,
						formattedResult,
					),
				)

				if containsError {
					logger.Info("Code execution failed")
					errorRetries++
					continue
				} else {
					logger.Info("Code execution succeeded")
					errorRetries = 0
					if !selector.SelectBool(
						"Do you want to continue?",
						false,
					) {
						return nil
					}

					req, err := inputHandler.Read(ctx, ">>> ")
					if err != nil {
						return fmt.Errorf("failed to read input: %w", err)
					}

					conversation.AddMessage(
						chat.NewMessage(
							chat.RoleUser,
							req,
						),
					)
				}
			case consoleActionAsk:
				fmt.Println(consoleResp.Question)
				req, err := inputHandler.Read(
					ctx,
					">>> ",
				)
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}

				conversation.AddMessage(
					chat.NewMessage(
						chat.RoleUser,
						req,
					),
				)

				// ask memory
				continue
			}
		}
	}
}

type consoleResponse struct {
	Action   consoleAction `json:"action"`
	Question string        `json:"question"`
	Language string        `json:"language"`
	Code     string        `json:"code"`
}

type consoleAction string

const (
	consoleActionCode consoleAction = "code"
	consoleActionAsk  consoleAction = "ask"
)

const instructionConsoleLinux = `You are running on a Linux machine. Assist the user in achieving their goal by clarifying any unclear steps, and return the appropriate action in JSON format — either asking for more clarification ('ask') or providing executable code ('code').

If generating code (action = code), follow these guidelines:
- Specify whether the code is for Bash or Python.
- Provide the code as an executable string under the code key.
- Ensure scripts are easy to understand, executable directly without edits, and output results to stdout only.

# Steps

1. **Identify User's Goal**:
   - If the goal is unclear, prompt the user with specific follow-up questions that help to proceed. Make the questions as precise as possible to gather the required information efficiently.
2. **Select Solution Type**:
   - When enough information is provided, decide whether the solution requires Bash or Python.
   - Choose the simplest option that satisfies the user's goal.
3. **Generate Script**:
   - Write an executable Bash or Python script that the user can run directly.
   - The script should operate without requiring interaction (e.g., prompts or saving to files).
   - Minimize complexity to improve understandability.
4. **Format the Response**:
   - Structure your output as a JSON object for consistency and clarity.

# Output Format

Your response should be a JSON object with the following keys:

- "action": Indicates if more clarification is needed ('ask') or if a code solution is being provided ('code').
  - action='ask': Include an additional "question" key that contains a specific question for the user to clarify missing requirements.
  - action='code': Include additional keys:
    - "language": Either 'bash' or 'python' to denote the script type.
    - "code": A single executable string containing the script.

# Examples

**Example 1 (Unclear Goal):**

User's request: "I need to copy data between directories."

**JSON Output:**
{
  "action": "ask",
  "question": "Could you please clarify the source and destination directories for copying the data? Should subdirectories be included as well?"
}

**Example 2 (Clear Goal with Code Solution):**

User's request: "List all the active network connections on this machine."

**JSON Output:**
{
  "action": "code",
  "language": "bash",
  "code": "netstat -tuln"
}

# Notes

- If the user request involves manipulating data (text processing, calculations) involving logic best handled in Python, prefer a Python solution.
- Prefer Bash for basic file operations or system commands.
- Output scripts should always produce straightforward results on stdout and should not create or modify files.
- Avoid overcomplicating follow-up questions—be direct in what information is needed for efficient clarification.`

const instructionConsoleMacOS = `You are running on a macOS machine. Assist the user in achieving their goal by clarifying any unclear steps and returning a corresponding action in JSON format, either for further clarification (action=ask) or providing directly executable code (action=code). Ensure that generated scripts are easy to understand, follow the previously outlined instructions, and meet the specifications outlined below.

If using action=code, specify the appropriate coding language, osascript and provide the executable code as a string under the code key. Scripts should be straightforward, executable as-is without additional editing, and output results to stdout, avoiding file storage or dialogs.

# Steps
1. Identify the user's goal. If it is unclear, prompt the user with specific follow-up questions to proceed with the implementation.
2. When you have all the details necessary, determine whether it requires osascript code. Use the simplest option that meets the requirements.
3. Generate the code directly executable from the terminal by the user.
4. Format your response in JSON.

# Output Format
- JSON object with keys:
  - action: Either 'ask' or 'code'.
    - 'ask': Used for clarifying additional details from the user if required before providing a solution.
    - 'code': Used for returning a ready-to-use script.
  - If action='code':
    - language: 'osascript'.
    - code: The script as a single string that can be copy-pasted for immediate execution.

# Examples

Example 1 (Clarification Needed)
Input: The user has asked 'help automate a task' without specifying details.
Output:
{
  'action': 'ask',
  'question': 'Could you please provide more details about the type of task you want to automate, such as opening an application, interacting with system settings, or something else?'
}

Example 2 (Provided Script)
Input: What is the title of my latest email ?
Output:
{
	'action': 'code',
	'language': 'osascript',
	'code': 'tell application "Mail" \n\tset latestMail to first message of inbox \n\tif latestMail is not missing value then \n\t\tset emailSubject to subject of latestMail \n\t\treturn "Subject: " & emailSubject \n\telse \n\t\treturn "No emails found." \n\tend if \nend tell'
}

Example 3 (Provided Script)
Input: Open perplexity
Output:
{
	'action': 'code',
	'language': 'osascript',
	'code': 'tell application "Google Chrome"
    activate
    open location "https://www.perplexity.ai/search?q=how+powerful+it+is+to+interact+with+computer+using+ai"
end tell
'
}

Example 4 (Provided Script)
Input: google doc
Output:
{
	'action': 'code',
	'language': 'osascript',
	'code': 'tell application "System Events"
    open location "https://docs.google.com/document/create"
    delay 5
end tell

tell application "System Events"
    delay 2
    keystroke "The Evolution of AI and Computing Experience"
    keystroke return
    keystroke return

    set articleText to "Artificial Intelligence (AI) has significantly transformed the computing experience in recent years. From enhancing productivity to reshaping industries, AI plays a crucial role in today's technological advancements. "
    
    repeat with i from 1 to (count of characters in articleText)
        keystroke (character i of articleText)
        delay 0.05
    end repeat
end tell
'
}


# Notes
- Begin by determining whether the user has provided sufficient details. Lack of specificity should result in a follow-up question (action=ask).
- Preference should be given to solutions that are simplest in implementation and easy to comprehend.
- Always ensure outputs are directed to the terminal and do not require additional user intervention.
- Avoid using GUI features that require manual clicks, approvals, or dialogs.`

func getConsoleInstruction(osName string) (string, error) {
	switch osName {
	case "linux":
		return instructionConsoleLinux, nil
	case "darwin":
		return instructionConsoleMacOS, nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", osName)
	}
}
