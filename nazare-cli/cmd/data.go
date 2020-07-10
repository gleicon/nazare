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
