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
	"strconv"
	"github.com/spf13/cobra"
	"github.com/richarda23/rspace-client-go/rspace"
)

var sharesFilterArg string

var listSharesCmd = &cobra.Command{
    	Use:   "listShares",
	Args:  cobra.NoArgs,
	Short: "List shared document or notebook",
	Long: `List all shared documents or notebooks. Optionally filter 
the output by the SharedItem's name, ID or global ID.  
	`,
	Example: `
// List shared items
rspace eln listShares

// List shared items which name contains "test"
rspace eln listShares --filter test

	`,

	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()
		cfg := configurePagination()
		doShareList(context,cfg)
	},
}


func doShareList(ctx *Context, cfg rspace.RecordListingConfig) {
	shareList, err := ctx.WebClient.ShareList(sharesFilterArg, cfg)
	if err != nil {
		exitWithErr(err)
	}
	formatter := &SharedItemListFormatter{shareList}
	ctx.writeResult(formatter)
}

type SharedItemListFormatter struct {
	sharedList *rspace.SharedItemList
}

func (fs *SharedItemListFormatter) ToJson() string {
	return prettyMarshal(fs.sharedList)
}

func (ds *SharedItemListFormatter) ToQuiet() []identifiable {
	rows := make([]identifiable, 0)
	for _, res := range ds.sharedList.Shares {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}

func (ds *SharedItemListFormatter) ToTable() *TableResult {
	results := ds.sharedList.Shares

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


func init() {
	elnCmd.AddCommand(listSharesCmd)
	listSharesCmd.Flags().StringVar(&sharesFilterArg, "filter", "", "Restrict results to SharedItem's name, ID or global ID")

	initPaginationFromArgs(listSharesCmd)
}
