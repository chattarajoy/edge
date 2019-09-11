package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "uniProxy",
		Short: "uniProxy is a proxy to various services in Qubole",
		Long:  "uniProxy or universal Proxy handles proxies of various protocols",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "path to config file")
}

func initConfig() {

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".uniProxy")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Print("Using config file: ", viper.ConfigFileUsed())
	}
}
