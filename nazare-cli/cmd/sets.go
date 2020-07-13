package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// setsCmd represents the sets command
var setsCmd = &cobra.Command{
	Use:   "sets",
	Short: "Cuckoo filter based sets",
	Long: `Cuckoo filter is a probabilistic structure that can tell if an item belongs 
	to a large set using way less memory than holding the entire set 
	but with a trade-off on false positives (they are possible). 
	False negatives (never seen that item) are not possible.
	
	Operations available: Add, Test if item is present, Remove item, Estimate cardinality

	- Add to a set: nazare-cli sets -a <setname> <item>
	- Check if an item belongs to a set: nazare-cli sets -i <setname> <item>
	- Remove an item from a set: nazare-cli sets -r <setname> <item>
	- Estimate set cardinality: nazare-cli sets -e <setname>

`,
	RunE: setsHandler,
}

func init() {
	rootCmd.AddCommand(setsCmd)
	setsCmd.Flags().BoolP("add", "a", false, "<setname> <item> - adds <item> to set <setname>")
	setsCmd.Flags().BoolP("ismember", "i", false, "<setname> <item> - check if <item> belongs to set <setname>")
	setsCmd.Flags().BoolP("remove", "r", false, "<setname> <item> - removes <item> from <setname>")
	setsCmd.Flags().BoolP("estimate", "e", false, "<setname> estimate set cardinality (size)")
}

func setsHandler(cmd *cobra.Command, args []string) error {
	var isMember, removed bool
	var err error

	addFlag, _ := cmd.Flags().GetBool("add")
	isMemberFlag, _ := cmd.Flags().GetBool("ismember")
	removeFlag, _ := cmd.Flags().GetBool("remove")
	estimateFlag, _ := cmd.Flags().GetBool("estimate")

	if !addFlag && !isMemberFlag && !removeFlag && !estimateFlag {
		return errors.New("Invalid command: no flag given")
	}

	if addFlag {
		if len(args) < 2 {
			return errors.New("Invalid parameters, Add requires <setname> <item>")
		}
		if err = ldb.localSets.SAdd([]byte(args[0]), []byte(args[1])); err != nil {
			return errors.Unwrap(fmt.Errorf("Error adding to set: %w", err))
		}
	}

	if isMemberFlag {
		if len(args) < 2 {
			return errors.New("Invalid parameters, isMember requires <setname> <item>")
		}
		if isMember, err = ldb.localSets.SisMember([]byte(args[0]), []byte(args[1])); err != nil {
			return errors.Unwrap(fmt.Errorf("Error checking isMember: %w", err))
		}
		fmt.Fprintf(os.Stdout, "%t\n", isMember)

	}

	if removeFlag {
		if len(args) < 2 {
			return errors.New("Invalid parameters, remove requires <setname> <item>")
		}
		if removed, err = ldb.localSets.SRem([]byte(args[0]), []byte(args[1])); err != nil {
			return errors.Unwrap(fmt.Errorf("Error removing from set: %w", err))
		}
		fmt.Fprintf(os.Stdout, "%t\n", removed)

	}

	if estimateFlag {
		var res uint
		if len(args) < 1 {
			return errors.New("Invalid parameters, Estimate requires <countername>")
		}
		if res, err = ldb.localSets.SCard([]byte(args[0])); err != nil {
			return errors.Unwrap(fmt.Errorf("Error checking estimate from set: %w", err))
		}
		fmt.Fprintf(os.Stdout, "%d\n", res)
	}
	return nil
}
