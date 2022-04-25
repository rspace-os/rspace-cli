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
	"io/ioutil"
	"strings"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

type addUserArgs struct {
	UsernameArg, FNameArg, LNameArg, EmailArg, RoleArg, AffiliationArg, PasswordFileArg string
}

var userArgs = addUserArgs{}

// addUserCmd represents the addUser command
var addUserCmd = &cobra.Command{
	Use:   "addUser",
	Short: "Adds a new user account",
	Long: `Requires sysadmin permission. Add a new user account rm, equires username and email
	at a minimum. Supply an initial password in a file.
	`,
	Example: ` 
	addUser --username newusername --email someone@somwhere.com --role user|pi| --pwdfile passwordfile
	`,
	Run: func(cmd *cobra.Command, args []string) {
		userPost := validateFlags()
		ctx := initialiseContext()
		user, err := ctx.WebClient.UserNew(userPost)
		if err != nil {
			exitWithErr(err)
		} else {
			ctx.write(prettyMarshal(user))
		}
	},
}

func validateFlags() *rspace.UserPost {
	pwd, err := ioutil.ReadFile(userArgs.PasswordFileArg)
	if err != nil {
		exitWithStdErrMsg("No password file supplied. Please put user password in a file and use the 'pwdfile' argument")
	}

	builder := &rspace.UserPostBuilder{}
	builder.Password(string(pwd)).Username(userArgs.UsernameArg)
	builder.Email(rspace.Email(userArgs.EmailArg))
	builder.FirstName(userArgs.FNameArg).LastName(userArgs.LNameArg)
	builder.Affiliation(userArgs.AffiliationArg)
	post, e := builder.Role(getRoleForArg(userArgs.RoleArg)).Build()
	if e != nil {
		exitWithErr(e)
	}
	return post
}

// pi or user by default
func getRoleForArg(arg string) rspace.UserRoleType {
	if strings.ToLower(arg) == "pi" {
		return rspace.Pi
	} else {
		return rspace.User
	}
}

func init() {
	elnCmd.AddCommand(addUserCmd)
	addUserCmd.Flags().StringVar(&userArgs.UsernameArg, "username", "", "username")
	addUserCmd.Flags().StringVar(&userArgs.FNameArg, "first", "Unknown", "First name")
	addUserCmd.Flags().StringVar(&userArgs.LNameArg, "last", "Unknown", "Last name")
	addUserCmd.Flags().StringVar(&userArgs.EmailArg, "email", "", "Valid email address")
	addUserCmd.Flags().StringVar(&userArgs.RoleArg, "role", "user", "Role, either 'pi' or 'user'")
	addUserCmd.Flags().StringVar(&userArgs.AffiliationArg, "affiliation", "unknown", "Affiliation (Community only)")

	addUserCmd.Flags().StringVar(&userArgs.PasswordFileArg, "pwdfile", "", "a file containing the password")
}
