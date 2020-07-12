package cmd

import (
	"github.com/spf13/cobra"
)

// countersCmd represents the counters command
var countersCmd = &cobra.Command{
	Use:   "counters",
	Short: "HyperLogLog based counters",
	Long: `HyperLogLog is an algorithm to estimate cadinality for large quantities of data without consuming
	 large memory for its elements. Essentially a counter with error deviation.

	 Operations available: Add to the set, estimate cardinality (count)
	- Add to a hll counter: nazare-cli counters -a <countername> <item>
	- Estimate counter size: nazare-cli counters -e <countername>
`,
}

func init() {
	rootCmd.AddCommand(countersCmd)

	countersCmd.Flags().StringP("add", "a", "", "<countername> <item> - adds <item> to a counter <countername>")
	countersCmd.Flags().StringP("estimate", "e", "", "<countername> estimate counter cardinality")

}
