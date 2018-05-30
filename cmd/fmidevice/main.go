package main

import (
	"fmt"
	"os"

	"github.com/michaljirman/fmidevice"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd *cobra.Command

func main() {
	rootCmd.Execute()
}

func init() {
	rootCmd = &cobra.Command{
		Use:   "fmidevice",
		Short: "iDevice locating tool",
	}
	configPath := os.Getenv("HOME")
	configName := ".fmidevice"

	rootCmd.AddCommand(fmidevice.LocateCmd)
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No configuration file found")
	} else {
		fmt.Printf("Configuration loaded from: %s/%s\n", configPath, configName)
	}
}
