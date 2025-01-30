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
	"strings"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

type shareArgs struct {
	permArg         string
	userIdsArg      string
	groupIdsArg     string
	targetFolderArg int
}

var shareArgsa shareArgs

// listDocumentsCmd represents the listDocuments command
var shareCmd = &cobra.Command{
	Use:   "share",
	Args:  cobra.MinimumNArgs(1),
	Short: "Shares one or more documents or notebooks with one or more users or groups",
	Long: `Share documents or notebooks with groups and individual users
	Documents/notebooks can be specified by plain IDs (e.g. 12345) or by Global Ids (e.g. SD12345)
	`,
	Example: `
// share a document and notebook with edit permission with a group into a designated folder
rspace eln share SD12345 NB23456 --groups 122345--permission edit --folder 7689

// Share a document with users in my group with read permission
rspace eln share 12345 --users 122345,45678


	`,

	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()
		if len(shareArgsa.groupIdsArg) == 0 && len(shareArgsa.userIdsArg) == 0 {
			exitWithStdErrMsg("You must specify either >= 1 group id (--groups) or >=1  user id (--users) to share with")
		}
		doShare(context, args, &shareArgsa)
	},
}

func doShare(ctx *Context, args []string, shareArgs *shareArgs) {
	post := rspace.SharePost{}
	var ids []int = idsFromGlobalIds(args)
	if len(ids) == 0 {
		exitWithStdErrMsg("No valid items to share, exiting")
	}

	post.ItemsToShare = ids
	var uIds, gIds []int
	if len(shareArgsa.groupIdsArg) > 0 {
		groupIdSs := strings.Split(shareArgsa.groupIdsArg, ",")
		gIds = idsFromGlobalIds(groupIdSs)
	}
	if len(shareArgsa.userIdsArg) > 0 {
		userIdSs := strings.Split(shareArgsa.userIdsArg, ",")
		uIds = idsFromGlobalIds(userIdSs)
	}
	if len(uIds) == 0 && len(gIds) == 0 {
		exitWithStdErrMsg("No valid user or group Ids to share with")
	}
	perm := shareArgsa.permArg
	if ok := validateArrayContains([]string{"read", "edit"}, []string{perm}); !ok {
		exitWithStdErrMsg(fmt.Sprintf("%s is not a valid permission", perm))
	}
	//TODO remove info logging from client, support target folder, don't send 0 as a value
	uPosts := make([]rspace.UserShare, 0)
	if len(uIds) > 0 {
		for _, id := range uIds {
			uPosts = append(uPosts, rspace.UserShare{id, perm})
		}
	}
	post.Users = uPosts
	gPosts := make([]rspace.GroupShare, 0)

	if len(gIds) > 0 {
		for _, id := range gIds {
			gPosts = append(gPosts, rspace.GroupShare{Id: id, Permission: perm})
		}
	}
	if len(gPosts) == 1 {
		gPosts[0].SharedFolderId = shareArgsa.targetFolderArg
	}
	post.Groups = gPosts
	var err error
	shares, err := ctx.WebClient.Share(&post)
	if err != nil {
		exitWithErr(err)
	}
	formatter := &ShareInfoListFormatter{shares}
	ctx.writeResult(formatter)
}

type ShareInfoListFormatter struct {
	shareList *rspace.ShareInfoList
}

func (fs *ShareInfoListFormatter) ToJson() string {
	return prettyMarshal(fs.shareList)
}

func (ds *ShareInfoListFormatter) ToQuiet() []identifiable {
	rows := make([]identifiable, 0)
	for _, res := range ds.shareList.ShareInfos {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}

func (ds *ShareInfoListFormatter) ToTable() *TableResult {
	results := ds.shareList.ShareInfos

	headers := []columnDef{columnDef{"Id", 8}, columnDef{"ItemId", 10}, columnDef{"ItemName", 25},
		columnDef{"SharedWith", 25}, columnDef{"Permission", 10}}

	rows := make([][]string, 0)
	for _, res := range results {
		data := []string{strconv.Itoa(res.Id), strconv.Itoa(res.ItemId), res.ItemName,
			res.TargetType, res.Permission}
		rows = append(rows, data)
	}
	return &TableResult{headers, rows}

}
func toIdentifiableShareInfo(results []*rspace.ShareResult) []identifiable {
	rows := make([]identifiable, 0)

	for _, res := range results {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}
func init() {
	elnCmd.AddCommand(shareCmd)
	shareCmd.Flags().StringVar(&shareArgsa.groupIdsArg, "groups", "",
		"Comma separated list of group Ids to share with")
	shareCmd.Flags().StringVar(&shareArgsa.userIdsArg, "users", "",
		"Comma separated list of user Ids to share with")
	shareCmd.Flags().StringVar(&shareArgsa.permArg, "permission", "read",
		"Permission - 'read' or 'edit'")
	shareCmd.Flags().IntVar(&shareArgsa.targetFolderArg, "folder", 0,
		"Target folder id (group sharing only)")
}
