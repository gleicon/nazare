package cmd

import (
	"fmt"
	"os"

	"github.com/gleicon/nazare/datalayer"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var ldb *datalayer.LocalDB

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nazare-cli",
	Short: "nazare-cli is the command line version of nazare",
	Long: `

Nazare is a set of libraries and a server for probabilistic counters and sets.

This command line embeds the libraries and provide a way to store and fetch this data locally, without the need of a server

nazare-cli serverless local HLLCounters and Sets
Use -b </path/to/databasename.db> to persist w/ badgedb (default name nazare.db in the current dir)

`,
}

/*
Execute the first command init
*/
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nazare-cli.yaml)")

	dbPath := rootCmd.Flags().StringP("database", "b", "nazare.db", "full database path and name. defaults to nazare.db at local dir")
	ldb = datalayer.NewLocalDB()
	ldb.Start(*dbPath)
}

// initConfig reads in config file and ENV variables if set.
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

		// Search config in home directory with name ".nazare-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".nazare-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
