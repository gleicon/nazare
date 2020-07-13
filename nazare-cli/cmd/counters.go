package cmd

import (
	"errors"
	"fmt"
	"os"

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
	RunE: countersHandler,
}

func init() {
	rootCmd.AddCommand(countersCmd)

	countersCmd.Flags().BoolP("add", "a", false, "<countername> <item> - adds <item> to a counter <countername>")
	countersCmd.Flags().BoolP("estimate", "e", false, "<countername> estimate counter cardinality")

}

func countersHandler(cmd *cobra.Command, args []string) error {
	var err error
	addFlag, _ := cmd.Flags().GetBool("add")
	estimateFlag, _ := cmd.Flags().GetBool("estimate")
	if !addFlag && !estimateFlag {
		return errors.New("Invalid command: no flag given")
	}
	if addFlag {
		if len(args) < 2 {
			return errors.New("Invalid parameters, Add requires <countername> <item>")
		}
		if err = ldb.localCounters.IncrementCounter([]byte(args[0]), []byte(args[1])); err != nil {
			return errors.Unwrap(fmt.Errorf("Error adding to counter: %w", err))
		}
	}

	if estimateFlag {
		var res uint64
		if len(args) < 1 {
			return errors.New("Invalid parameters, Estimate requires <countername>")
		}
		if res, err = ldb.localCounters.RetrieveCounterEstimate([]byte(args[0])); err != nil {
			return errors.Unwrap(fmt.Errorf("Error checking estimate from counter: %w", err))
		}
		fmt.Fprintf(os.Stdout, "%d\n", res)
	}
	return nil
}
