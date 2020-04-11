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
	"github.com/spf13/cobra"
	"rspace"
	"strconv"
	"os"
)
var foldername string
var parentFolder string
// createNotebookCmd represents the createNotebook command
var createFolderCmd = &cobra.Command{
	Use:   "createFolder",
	Short: "Creates a new Folder",
	Long: `Create a new Folder, with an optional name and parent folder ID
	  create-folder --name foldername --folder FL1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		outputFormat = outputFmt(outputFormatArg)
		validateOutputFormatExit (outputFormat)
		post := rspace.FolderPost{IsNotebook:false,}
		doCreateFolder(foldername, parentFolder, post)
	},
}

func doCreateFolder (foldername string, parentFolder string, post rspace.FolderPost) {
		if len(foldername) > 0 {
			post.Name=foldername
		}
		if len(parentFolder) > 0 {
			id,err :=strconv.Atoi(parentFolder)
			if err != nil {
				fmt.Println("Please supply a numeric folder id for the parent folder")
				os.Exit(1)
			}
			post.ParentFolderId=id
		}
		webClient:=setup()
		got, err := webClient.FolderNew(&post)
		if err != nil {
			exitWithErr(err)
		}
		if outputFormat.isJson()  {
			fmt.Println(prettyMarshal(got))
		} else if outputFormat.isTab() || outputFormat.isCsv() {
			folderToTable(got)
		} else if outputFormat.isQuiet(){
			fmt.Println(got.Id)
		} else {
			fmt.Println("?" + string(outputFormat))
		}
}
func folderToTable(folder *rspace.Folder) {
	headers := []columnDef {columnDef{"Id",8}, columnDef{"GlobalId",10},  columnDef{"Name", 25}, columnDef{"ParentFolderId",15}, columnDef{"Created",24}}
	data := []string {strconv.Itoa(folder.Id),folder.GlobalId, folder.Name, strconv.Itoa(folder.ParentFolderId), folder.Created}
	rows := make([][]string, 0)
	rows = append(rows, data)
	if outputFormat.isCsv() {
		printCsv(headers, rows)
	} else {
		printTable(headers, rows)
	}
}
func init() {
	elnCmd.AddCommand(createFolderCmd)
	 createFolderCmd.Flags().StringVarP(&foldername, "name", "n", "", "A name for the folder")
	 createFolderCmd.Flags().StringVarP(&parentFolder, "folder", "p", "", "An id for the folder that will contain the new folder")
}
