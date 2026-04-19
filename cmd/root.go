package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/geqo/togram/internal/config"
	"github.com/geqo/togram/internal/detect"
	"github.com/geqo/togram/internal/telegram"
	"github.com/spf13/cobra"
)

var (
	flagChat  string
	flagToken string
)

var rootCmd = &cobra.Command{
	Use:   "togram [flags] [file]",
	Short: "Send messages and files to Telegram",
	Example: `  echo "hello" | togram
  cat photo.jpg | togram
  togram report.pdf
  togram -c @mychat --token 123:ABC video.mp4`,
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&flagChat, "chat", "c", "", "chat ID or @username")
	rootCmd.Flags().StringVar(&flagToken, "token", "", "bot token")
	rootCmd.AddCommand(completionCmd)
}

var completionCmd = &cobra.Command{
	Use:    "completion [bash|zsh|fish]",
	Short:  "Generate shell completion script",
	Hidden: true,
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		}
		return nil
	},
}

func run(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	token := first(flagToken, cfg.Token)
	if token == "" {
		return fmt.Errorf("bot token required: use --token or set token in %s", config.Path)
	}

	chatID := first(flagChat, cfg.ChatID)
	if chatID == "" {
		return fmt.Errorf("chat ID required: use -c or set chat in %s", config.Path)
	}

	client := telegram.New(token)

	if len(args) == 1 {
		return sendFile(client, chatID, args[0])
	}
	return sendStdin(client, chatID)
}

func sendFile(client *telegram.Client, chatID, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	t, r, err := detect.FromReader(f)
	if err != nil {
		return err
	}

	// prefer extension-based detection for files (more reliable than magic bytes)
	if ext := detect.FromFilename(path); ext != detect.TypeDocument {
		t = ext
	}

	return client.Send(chatID, t, r, filepath.Base(path))
}

func sendStdin(client *telegram.Client, chatID string) error {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return fmt.Errorf("nothing to send: pipe input or provide a file argument")
	}

	t, r, err := detect.FromReader(os.Stdin)
	if err != nil {
		return err
	}

	return client.Send(chatID, t, r, "")
}

func first(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
