/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"strconv"
	"github.com/spf13/cobra"
)


// listDocumentsCmd represents the listDocuments command
var unshareCmd = &cobra.Command{
    	Use:   "unshare",
	Args:  cobra.ExactArgs(1),
	Short: "Unshares a document or notebook",
	Long: `Unshare documents or notebooks 
	Documents/notebooks can be specified by the 'id' attribute of a ShareInfo response
	`,
	Example: `
// share a document
rspace eln unshare 1245184

	`,

	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()
		if sharedID, err := strconv.Atoi(args[0]); err == nil {
			doUnshare(context, sharedID)  
		} else {
			exitWithStdErrMsg(fmt.Sprintf("%s. %s is not a valid ID", err.Error(),args[0]))	
		}
		
	},
}

func doUnshare(ctx *Context, sharedID int) {

	_, err := ctx.WebClient.Unshare(sharedID)
	if err != nil {
		exitWithErr(err)
	}

}

func init() {
	elnCmd.AddCommand(unshareCmd)
}
