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
