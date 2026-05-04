package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/inmaildev/inmail-cli/internal/output"
	"github.com/spf13/cobra"
)

const maxHistoryTokens = 4000

type replMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func newReplCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "repl",
		Short: "Start an interactive REPL session",
		Long: `Start an interactive command-line REPL for the InMail API.

Maintains conversation history across commands within the session.
History is trimmed when token limit (~4000 tokens) is approached.

Type 'help' to see available commands. Type 'exit' or press Ctrl+D to quit.

Agent-friendly flags:
  --non-interactive   Read commands from stdin, one per line (no prompts)
  --output json       Emit all output as JSON

Example (agent use):
  echo "messages list --limit 5" | inmail repl --non-interactive --output json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagNonInteractive {
				return runNonInteractiveRepl()
			}
			return runInteractiveRepl()
		},
	}
}

func runInteractiveRepl() error {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[36minmail\033[0m › ",
		HistoryFile:     os.ExpandEnv("$HOME/.inmail/repl_history"),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return fmt.Errorf("init readline: %w", err)
	}
	defer rl.Close()

	var history []replMessage
	output.PrintHeader("InMail REPL — type 'help' for commands, 'exit' to quit")
	fmt.Println()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			output.PrintInfo("Goodbye")
			break
		}

		result, err := dispatchReplCommand(line)
		history = appendHistory(history, "user", line)
		if err != nil {
			output.PrintError(err.Error())
			history = appendHistory(history, "assistant", "error: "+err.Error())
		} else {
			history = appendHistory(history, "assistant", result)
			fmt.Println(result)
		}
		history = trimHistory(history, maxHistoryTokens)
	}
	return nil
}

func runNonInteractiveRepl() error {
	var history []replMessage
	scanner := newLineScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		history = appendHistory(history, "user", line)

		result, err := dispatchReplCommand(line)
		if err != nil {
			if output.Current == output.FormatJSON {
				b, _ := json.Marshal(map[string]string{"error": err.Error(), "input": line})
				fmt.Println(string(b))
			} else {
				fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			}
			history = appendHistory(history, "assistant", "error: "+err.Error())
		} else {
			fmt.Println(result)
			history = appendHistory(history, "assistant", result)
		}
		history = trimHistory(history, maxHistoryTokens)
	}
	return scanner.Err()
}

func dispatchReplCommand(line string) (string, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", nil
	}

	switch parts[0] {
	case "help":
		return replHelp(), nil
	case "history":
		return "History is maintained internally per session.", nil
	case "version":
		return "inmail CLI v" + Version, nil
	}

	rootCmd.SetArgs(parts)
	var buf strings.Builder
	rootCmd.SetOut(&buf)

	if err := rootCmd.Execute(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func replHelp() string {
	return `Available commands:
  users list [--limit N] [--cursor C]
  users get <email>
  users create --email E [--name N] [--password P]
  users update <email> [--name N] [--password P]
  users delete <email>

  messages list [--limit N] [--account E]
  messages get <id>
  messages delete <id>

  accounts list
  accounts get <id>
  accounts create --email E
  accounts update <id> --email E
  accounts delete <id>

  config get
  config create [--attachments] [--catch-all]
  config update [--attachments=false] [--catch-all]

  stats [--days N]

  send --to E --subject S [--text T] [--html H] [--from F]

  send-logs list [--limit N]
  send-logs get <id>

  version
  help
  exit`
}

func appendHistory(history []replMessage, role, content string) []replMessage {
	return append(history, replMessage{Role: role, Content: content})
}

func trimHistory(history []replMessage, maxTokens int) []replMessage {
	total := 0
	for _, m := range history {
		total += len(m.Content) / 4
	}
	for total > maxTokens && len(history) > 1 {
		total -= len(history[0].Content) / 4
		history = history[1:]
	}
	return history
}

type lineScanner struct {
	r    *os.File
	line strings.Builder
	done bool
	err  error
	text string
}

func newLineScanner(f *os.File) *lineScanner {
	return &lineScanner{r: f}
}

func (s *lineScanner) Scan() bool {
	if s.done {
		return false
	}
	s.line.Reset()
	buf := make([]byte, 1)
	for {
		n, err := s.r.Read(buf)
		if n > 0 {
			if buf[0] == '\n' {
				s.text = s.line.String()
				return true
			}
			s.line.WriteByte(buf[0])
		}
		if err != nil {
			if s.line.Len() > 0 {
				s.text = s.line.String()
				s.done = true
				return true
			}
			s.done = true
			return false
		}
	}
}

func (s *lineScanner) Text() string { return s.text }
func (s *lineScanner) Err() error   { return s.err }
