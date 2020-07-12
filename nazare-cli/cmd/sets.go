package cmd

import (
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
}

func init() {
	rootCmd.AddCommand(setsCmd)
	setsCmd.Flags().StringP("add", "a", "", "<setname> <item> - adds <item> to set <setname>")
	setsCmd.Flags().StringP("ismember", "i", "", "<setname> <item> - check if <item> belongs to set <setname>")
	setsCmd.Flags().StringP("remove", "r", "", "<setname> <item> - removes <item> from <setname>")
	setsCmd.Flags().StringP("estimate", "e", "", "<setname> estimate set cardinality (size)")
}
