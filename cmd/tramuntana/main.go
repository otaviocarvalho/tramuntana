package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/otaviocarvalho/tramuntana/hook"
	"github.com/otaviocarvalho/tramuntana/internal/config"
	"github.com/spf13/cobra"
)

var (
	version   = "v0.1.0"
	cfgPath   string
	cfg       *config.Config
	installHook bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "tramuntana",
		Short: "Bridge Telegram group topics to Claude Code sessions via tmux",
	}

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the Telegram bot",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cfgPath != "" {
				_ = godotenv.Load(cfgPath)
			}
			var err error
			cfg, err = config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Starting tramuntana bot...")
			// Bot wiring comes in Task 07
			return nil
		},
	}
	serveCmd.Flags().StringVar(&cfgPath, "config", "", "path to .env config file")

	hookCmd := &cobra.Command{
		Use:   "hook",
		Short: "Run the Claude Code SessionStart hook",
		RunE: func(cmd *cobra.Command, args []string) error {
			if installHook {
				return hook.Install()
			}
			return hook.Run()
		},
	}
	hookCmd.Flags().BoolVar(&installHook, "install", false, "install hook into Claude Code settings")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("tramuntana %s\n", version)
		},
	}

	rootCmd.AddCommand(serveCmd, hookCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
