package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/noahgorstein/stardog-go/internal/config"
	tuiv2 "github.com/noahgorstein/stardog-go/internal/tuiV2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "sd-security",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.OnInitialize(initConfig)

		viper.SetDefault("endpoint", "http://localhost:5820")
		viper.SetDefault("username", "admin")
		viper.SetDefault("password", "admin")

		config := config.Config{
			Endpoint: viper.GetString("endpoint"),
			Username: viper.GetString("username"),
			Password: viper.GetString("password"),
		}

		// bubble := tui.New(config)
		// p := tea.NewProgram(bubble, tea.WithAltScreen())
		bubble := tuiv2.New(config)
		p := tea.NewProgram(bubble, tea.WithAltScreen())

		if err := p.Start(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stardog.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".stardog" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".stardog")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

}
