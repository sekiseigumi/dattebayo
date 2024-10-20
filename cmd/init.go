package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func initialize() *cobra.Command {
	init := &cobra.Command{
		Use:     "init",
		Short:   "init dattebayo for the first time.",
		Long:    "creates an admin user and initializes the program for the first use.",
		Example: "dattebayo init",
		Aliases: []string{"i"},
		RunE: func(cmd *cobra.Command, args []string) error {
			currentUser, err := user.Current()
			if err != nil {
				return fmt.Errorf("unable to get current user: %w", err)
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Welcome to Dattebayo! Create a superuser to get started.")
			fmt.Printf("Enter username for superuser [%s]: ", currentUser.Username)
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)
			if username == "" {
				username = currentUser.Username
			}
			fmt.Printf("Enter password for superuser: ")
			passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("unable to read password: %w", err)
			}
			password := strings.TrimSpace(string(passwordBytes))
			fmt.Println()
			if password == "" {
				return fmt.Errorf("password cannot be empty")
			}

			fmt.Println("Creating superuser...")
			userHomeDir, _ := os.UserHomeDir()
			dattebayoConfigDir := filepath.Join(userHomeDir, ".dattebayo")
			configFilePath := filepath.Join(dattebayoConfigDir, "config.yaml")

			if _, err := os.Stat(dattebayoConfigDir); os.IsNotExist(err) {
				if err := os.MkdirAll(dattebayoConfigDir, 0755); err != nil {
					return fmt.Errorf("unable to create dattebayo config directory: %w", err)
				}
			}

			if _, err := os.Stat(configFilePath); err == nil {
				return fmt.Errorf("config file already exists at %s", configFilePath)
			}

			configFile, err := os.Create(configFilePath)
			if err != nil {
				return fmt.Errorf("unable to create config file: %w", err)
			}

			defer configFile.Close()

			configContent := fmt.Sprintf(`# Dattebayo Configuration File
# Superuser credentials
superuser:
  # username of the superuser (change if needed)
  username: %s

  # password of the superuser (change if needed)
  password: %s

# General settings
# timer before all the servers start (in seconds, default is 5)
# start_timer: 5
`, username, password)

			if _, err := configFile.WriteString(configContent); err != nil {
				return fmt.Errorf("unable to write to config file: %w", err)
			}

			fmt.Printf("Superuser created successfully. Config file created at %s\n", configFilePath)
			fmt.Println("Run `dattebayo` to start the program.")
			return nil
		},
	}
	return init
}
