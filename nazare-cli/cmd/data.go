package cmd

import (
	"errors"
	"fmt"
	"os"

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
	RunE: dataHandler,
}

func init() {
	rootCmd.AddCommand(dataCmd)
	dataCmd.Flags().BoolP("set", "s", false, "<key> <value> - sets <key> to <value>")
	dataCmd.Flags().BoolP("get", "g", false, "<key> - gets value for <key>")
	dataCmd.Flags().BoolP("del", "d", false, "deletes <key>")
}

func dataHandler(cmd *cobra.Command, args []string) error {
	var removed bool
	var value []byte
	var err error

	setFlag, _ := cmd.Flags().GetBool("set")
	getFlag, _ := cmd.Flags().GetBool("get")
	deleteFlag, _ := cmd.Flags().GetBool("del")

	if !setFlag && !getFlag && !deleteFlag {
		return errors.New("Invalid command: no flag given")
	}

	if setFlag {
		if len(args) < 2 {
			return errors.New("Invalid parameters, Add requires <key> <value>")
		}
		if err = ldb.LocalDatastorage.Add([]byte(args[0]), []byte(args[1])); err != nil {
			return errors.Unwrap(fmt.Errorf("Error adding to db: %w", err))
		}
	}

	if getFlag {
		if len(args) < 1 {
			return errors.New("Invalid parameters, Get requires <key>")
		}
		if value, err = ldb.LocalDatastorage.Get([]byte(args[0])); err != nil {
			return errors.Unwrap(fmt.Errorf("Error getting value from db: %w", err))
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(value))

	}

	if deleteFlag {
		if len(args) < 1 {
			return errors.New("Invalid parameters, delete requires <key>")
		}
		if removed, err = ldb.LocalDatastorage.Delete([]byte(args[0])); err != nil {
			return errors.Unwrap(fmt.Errorf("Error deleting from db: %w", err))
		}
		fmt.Fprintf(os.Stdout, "%t\n", removed)

	}

	return nil
}
