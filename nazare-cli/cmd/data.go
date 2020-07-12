package cmd

import (
	"github.com/spf13/cobra"
)

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "K/V data storage on badgerdb",
	Long: `Local key/value storage. No data type handling
	
	Operations available: Set, Get and Del(ete)
	- Set a key: nazare-cli data -s <keyname> <value>
	- Get the value of a key: nazare-cli data -g <keyname>
	- Delete a key: nazare-cli data -d <keyname>
	`,
}

func init() {
	rootCmd.AddCommand(dataCmd)
	dataCmd.Flags().StringP("set", "s", "", "<key> <value> - sets <key> to <value>")
	dataCmd.Flags().StringP("get", "g", "", "<key> - gets value for <key>")
	dataCmd.Flags().StringP("del", "d", "", "deletes <key>")
}
