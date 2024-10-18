package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sekiseigumi/dattebayo/internal/tui"
	"github.com/sekiseigumi/dattebayo/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Version: "v0.0.1",
		Use:     "dattebayo",
		Long:    "Dattebayo helps you do 127.0.0.1 things. Unleashes Super-Ultra-Big Ball Rasengan if the superuser does it. Believe it!",
		Example: "dattebayo",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "init" {
				return nil
			}

			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			viper.SetEnvPrefix("DATTEBAYO")
			viper.SetConfigType("yaml")
			viper.AutomaticEnv()

			userHomeDir, _ := os.UserHomeDir()
			dattebayoConfigDir := filepath.Join(userHomeDir, ".dattebayo")

			if _, err := os.Stat(dattebayoConfigDir); os.IsNotExist(err) {
				return fmt.Errorf("dattebayo not initialized. run `dattebayo init` to initialize")
			}

			configFile, _ := cmd.Flags().GetString("config")
			if configFile != "" {
				// Use config file from flag
				viper.SetConfigFile(configFile)
			} else {
				// Default to $HOME/.dattebayo/config.yaml
				viper.AddConfigPath(dattebayoConfigDir)
				viper.SetConfigName("config")
			}

			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("unable to read config: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var config shared.Config

			if err := viper.Unmarshal(&config); err != nil {
				return fmt.Errorf("unable to unmarshal config: %w", err)
			}

			if config.Superuser.Username == "" || config.Superuser.Password == "" {
				return fmt.Errorf("config file corrupted. missing superuser credentials. remove config file and run `dattebayo init` to initialize")
			}

			tuiInstance := tui.NewTUI(config)
			if _, err := tuiInstance.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.AddCommand(initialize())

	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.dattebayo/config.yaml)")

	return rootCmd.ExecuteContext(context.Background())
}
