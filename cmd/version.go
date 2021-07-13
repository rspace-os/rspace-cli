package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// this variable is searched/replaced by build.sh
// if this variable is renamed, the script should be updated
//go:embed version.txt
var rsVersion string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of RSpace CLI",
	Long:  `All software has versions. This is RSpace's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("This is RSpace CLI version %s", rsVersion)
	},
}
