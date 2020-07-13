/*
 * @Author: domchan
 * @Date: 2018-12-28 15:30:41
 * @Last Modified by: domchan
 * @Last Modified time: 2018-12-28 15:31:43
 */
package cmd

import (
	"flag"
	"github.com/spf13/cobra"
)

// AddFlags adds all command line flags to the given command.
func AddFlags(rootCmd *cobra.Command) {
	flag.CommandLine.VisitAll(func(gf *flag.Flag) {
		rootCmd.PersistentFlags().AddGoFlag(gf)
	})
}
