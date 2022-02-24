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
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

// listDocumentsCmd represents the listDocuments command
var listUsersCmd = &cobra.Command{
	Use:   "listUsers",
	Short: "Lists users - requires sysadmin role!",
	Long: `List users, sorted or paginated.

Please note that currently users are ordered by account creation date only (default is most recent first),
but there is no limit on the page size, so you can retrieve all users in one query, then
filter or sort using standard Unix utilities.

	`,
	Example: `
// find out how many users you have:
rspace eln listUsers --all -f json | jq '.TotalHits'

// or just use 'grep' if you don't have 'jq' installed:
rspace eln listUsers --all -f json | grep TotalHits

// now you can get all users in one go and sort by any column e.g. by username:
rspace eln listUsers --all | sort -k2
	`,

	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()
		cfg := configurePagination()
		cfg.OrderBy = "creationDate"
		doListusers(context, cfg)
	},
}

var allArg bool = false

func getNextPageIndex(userList *rspace.UserList) int {
	links := userList.Links
	for _, link := range links {
		if link.Rel == "next" {
			u, _ := url.Parse(link.Link)
			if u != nil {
				m, _ := url.ParseQuery(u.RawQuery)
				pageNumber, _ := strconv.Atoi(m["pageNumber"][0])
				return pageNumber

			}
		}
	}
	return 0
}

func doListusers(ctx *Context, cfg rspace.RecordListingConfig) {
	var usersList *rspace.UserList
	var err error
	usersList, err = ctx.WebClient.Users(time.Time{}, time.Time{}, cfg)
	if err != nil {
		exitWithErr(err)
	}
	var next_link = getNextPageIndex(usersList)
	users := usersList.Users

	// get all pages and aggregate results
	if allArg {
		for next_link != 0 {

			cfg.PageNumber = next_link
			usersList, err = ctx.WebClient.Users(time.Time{}, time.Time{}, cfg)
			if err != nil {
				exitWithErr(err)
			}
			users = append(users, usersList.Users...)
			next_link = getNextPageIndex(usersList)
			//formatter := &UserListFormatter{usersList}
			//ctx.writeResult(formatter)
		}
	}
	allUserList := rspace.UserList{Users: users, TotalHits: usersList.TotalHits,
		Links: usersList.Links, PageNumber: 0}
	formatter := &UserListFormatter{&allUserList}
	ctx.writeResult(formatter)
}

type UserListFormatter struct {
	*rspace.UserList
}

func (fs *UserListFormatter) ToJson() string {
	return prettyMarshal(fs.UserList)
}

func (ds *UserListFormatter) ToQuiet() []identifiable {
	rows := make([]identifiable, 0)
	for _, res := range ds.UserList.Users {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}

func (ds *UserListFormatter) ToTable() *TableResult {
	results := ds.UserList.Users

	sort.SliceStable(results,
		func(i, j int) bool { return len(results[j].Email) < len(results[i].Email) })
	maxEmailLen := len(results[0].Email)
	sort.SliceStable(results,
		func(i, j int) bool { return len(results[j].Username) < len(results[i].Username) })
	maxUnameLen := len(results[0].Username)
	sort.SliceStable(results,
		func(i, j int) bool { return len(results[j].FirstName) < len(results[i].FirstName) })
	maxFnameLen := maxInt(9, len(results[0].FirstName))
	sort.SliceStable(results,
		func(i, j int) bool { return len(results[j].LastName) < len(results[i].LastName) })
	maxLnameLen := maxInt(9, len(results[0].LastName))
	headers := []columnDef{columnDef{"Id", 8}, columnDef{"Username", maxUnameLen},
		columnDef{"email", maxEmailLen}, columnDef{"FirstName", maxFnameLen},
		columnDef{"LastName", maxLnameLen}}

	rows := make([][]string, 0)
	for _, res := range results {
		data := []string{strconv.Itoa(res.Id), res.Username, res.Email, res.FirstName,
			res.LastName}
		rows = append(rows, data)
	}
	return &TableResult{headers, rows}

}
func toIdentifiableUser(results []*rspace.UserInfo) []identifiable {
	rows := make([]identifiable, 0)

	for _, res := range results {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}
func init() {
	elnCmd.AddCommand(listUsersCmd)
	initPaginationFromArgs(listUsersCmd)
	listUsersCmd.PersistentFlags().BoolVar(&allArg, "all", false, "Gets all users")
}
