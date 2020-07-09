/*
Copyright © 2020 Gleicon Moraes <gleicon@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nazare-cli",
	Short: "nazare-cli is the command line version of nazare",
	Long: `

Nazare is a set of libraries and a server for probabilistic counters and sets.

This command line embeds the libraries and provide a way to store and fetch this data locally, without the need of a server

nazare-cli serverless local HLLCounters and Sets
Use -b </path/to/databasename.db> to persist w/ badgedb (default name nazare.db in the current dir)

HyperLogLog based counters:
Add to a hll counter: nazare-cli -c -a <countername> <item>
Estimate counter size: nazare-cli -c -e <countername>

Cuckoo filter based sets:
Add to a set: nazare-cli -s -a <setname> <item>
Check if an item belongs to a set: nazare-cli -s -i <setname> <item>
Remove an item from a set: nazare-cli -s -r <setname> <item>
Estimate set cardinality: nazare-cli -s -c <setname>
That's it, no way to get an item from a set, cuckoo filter stores hashes and signal it an item was 'seen'

K/V handling:
Set a key: nazare-cli -k -s <keyname> <value>
Get the value of a key: nazare-cli -k -g <keyname>
Delete a key: nazare-cli -k -d <keyname>

`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nazare-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
