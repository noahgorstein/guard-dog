package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/noahgorstein/go-stardog/stardog"
	"github.com/noahgorstein/guard-dog/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Credits: https://github.com/carolynvs/stingoftheviper
const (
	defaultConfigFilename = ".guard-dog"
	envPrefix             = "GUARD_DOG"
)

func NewRootCommand() *cobra.Command {
	username := ""
	password := ""
	endpoint := ""

	rootCmd := &cobra.Command{
		Use:   "guard-dog",
		Short: "a TUI to manage users, roles and permissions in Stardog ⭐🐕",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {

			basicAuthTransport := &stardog.BasicAuthTransport{
				Username: username,
				Password: password,
			}

			client, _ := stardog.NewClient(endpoint, basicAuthTransport.Client())

			alive, _, err := client.ServerAdmin.IsAlive(context.Background())
			if err != nil || !*alive {
				log.Fatalf("Server: %s is not alive.", endpoint)
			}

			bubble := tui.New(*client, username)
			p := tea.NewProgram(bubble, tea.WithAltScreen())

			if err := p.Start(); err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}

		},
	}

	rootCmd.Flags().StringVarP(&username, "username", "u", "admin", "username")
	rootCmd.Flags().StringVarP(&password, "password", "p", "admin", "password")
	rootCmd.Flags().StringVarP(&endpoint, "server", "s", "http://localhost:5820", "server")

	return rootCmd
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName(defaultConfigFilename)

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	v.AddConfigPath(home)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			_ = v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
