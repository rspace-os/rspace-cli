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
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

type addDocArgs struct {
	ParentfolderArg string
	NameArg         string
	Tags            string
	ContentFile     string
	Content         string
	FormId          string
	InputData       string
	InputDataFormat string
}

var addDocArgV = addDocArgs{}

// addDocumentCmd represents the createNotebook command
var addDocumentCmd = &cobra.Command{
	Use:   "addDocument",
	Short: "Creates a new Basic or Structured document with optional tags and content",
	Long: `Create a new document, with an optional name and parent folder.
If a file is a file of HTML content, it is loaded verbatim, otherwise, plain text files are wrapped in 'pre'
tags to preserve formatting.

You can also create an structured (multi-field) document by passing the 'formId'.
This creates an empty document. You can also create 1 or more documents from a CSV
file with the following characteristics:
- 1st row is a header row and is ignored
- Each row will supply data to create an RSpace document
- The column order should match the fields in the Form definition
- The total number of columns should match the total number of fields in the Form
	`,
	Example: `
// create a new document with tags and HTML content
rspace eln  addDocument --name doc1 --tags tag1,tag2 --contentFile textToPutInDoc.html

// create a doc with tags and  plain-text content, which will be wrapped in '<pre>' tag
rspace eln  addDocument --name doc1 --tags tag1,tag2 --contentFile textToPutInDoc.txt

// create a doc using verbatim text
rspace eln  addDocument --name doc1  --content "some content"

// create an empty, structured (multi-field) document
rspace eln addDocument --name myDoc --formId FM2

// create a multi-field document with data in a CSV file:
rspace eln addDocument --name myDoc --formId FM2 --input data.csv

`,
	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()
		err := doAddDocRun(addDocArgV, context, context.WebClient)
		if err != nil {
			exitWithErr(err)
		}
	},
}

type DocClient interface {
	NewBasicDocumentWithContent(name, tags, content string) (*rspace.Document, error)

	NewDocumentWithContent(post *rspace.DocumentPost) (*rspace.Document, error)
}

func doAddDocRun(addDocArgV addDocArgs, context *Context, docClient DocClient) error {
	var created *rspace.Document
	var err error
	// is basic document if no form is supplied,
	// we make a basic document
	createdDocs := make([]*rspace.DocumentInfo, 0)
	if len(addDocArgV.FormId) == 0 {
		content := getContent(addDocArgV)
		created, err = docClient.NewBasicDocumentWithContent(addDocArgV.NameArg,
			addDocArgV.Tags, content)
		createdDocs = append(createdDocs, created.DocumentInfo)
	} else {
		// we make a structured document
		createdDocs, err = readDocContentFromFile(addDocArgV, docClient)
	}

	if err != nil {
		return err
	}
	docList := rspace.DocumentList{}
	// deref pointers to fit into DocumentList
	toList := make([]rspace.DocumentInfo, 0)
	for _, v := range createdDocs {
		toList = append(toList, *v)
	}
	docList.Documents = toList
	var dlf = DocListFormatter{&docList}
	context.writeResult(&dlf)
	return nil
}

// called when formId is not set.
func readDocContentFromFile(addDocArgV addDocArgs, docClient DocClient) ([]*rspace.DocumentInfo, error) {
	createdDocs := make([]*rspace.DocumentInfo, 0)

	// else is form, we add content if there is any
	// TODO implement this, use CSV data as an example.
	formId, err := idFromGlobalId(addDocArgV.FormId)
	if err != nil {
		return nil, err
	}
	parentFolderId := 0
	if addDocArgV.ParentfolderArg != "" {
		parentFolderId, err = idFromGlobalId(addDocArgV.ParentfolderArg)
		if err != nil {
			return nil, err
		}
	}

	var toPost = rspace.DocumentPost{}
	toPost.Name = addDocArgV.NameArg
	toPost.Tags = addDocArgV.Tags
	toPost.ParentFolderId = parentFolderId
	toPost.FormID = rspace.FormId{formId}
	if len(addDocArgV.InputData) > 0 {
		f, err := os.Open(addDocArgV.InputData)
		if err != nil {
			return nil, err
		}
		csvIn, err := readCsvFile(f)
		if err != nil {
			return nil, err
		}
		err = validateCsvInput(csvIn)
		if err != nil {
			return nil, err
		}
		// iterate of each row in file
		for i, v := range csvIn {
			// ignore 1st line - header
			if i == 0 {
				continue
			}
			messageStdErr(fmt.Sprintf("%d of %d", i, len(csvIn)-1))
			var content []rspace.FieldContent = make([]rspace.FieldContent, 0)
			// convert each column into a field
			for _, v2 := range v {
				content = append(content, rspace.FieldContent{Content: v2})
			}
			toPost.Fields = content
			// add suffix to name if > 1 document being created
			if i > 1 {
				toPost.Name = fmt.Sprintf("%s-%d", toPost.Name, i)
			}
			doc, err := docClient.NewDocumentWithContent(&toPost)
			if err != nil {
				messageStdErr(fmt.Sprintf("Could not create document from data in row %d", i))
				continue
			}
			createdDocs = append(createdDocs, doc.DocumentInfo)
		}
	} else {
		doc, err := docClient.NewDocumentWithContent(&toPost)
		if err != nil {
			exitWithErr(err)
		}
		createdDocs = append(createdDocs, doc.DocumentInfo)

	}
	return createdDocs, nil
}

// validates that csv input is suitable for creating documents from
func validateCsvInput(csvIn [][]string) error {
	if len(csvIn) <= 1 {
		return errors.New(`There must be at least 2 lines in the csv file - 
				1 for headers and at least 1 row of data`)
	}
	row1Len := len(csvIn[0])
	// all should have same number of columns
	for i, v := range csvIn {
		if len(v) != row1Len {
			return errors.New(fmt.Sprintf("Discrepancy in row length for row %d. Expected %d but was %d",
				i, row1Len, len(v)))
		}
	}
	return nil
}

// reads csv input from a source
func readCsvFile(csvIn io.Reader) ([][]string, error) {

	csvReader := csv.NewReader(csvIn)
	result, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getContent(addDocArgV addDocArgs) string {
	if len(addDocArgV.Content) > 0 {
		return addDocArgV.Content
	} else if len(addDocArgV.ContentFile) > 0 {
		bytes, err := ioutil.ReadFile(addDocArgV.ContentFile)
		if err != nil {
			exitWithErr(err)
		}
		lowerCaseFile := strings.ToLower(addDocArgV.ContentFile)
		if strings.HasSuffix(lowerCaseFile, "html") ||
			strings.HasSuffix(lowerCaseFile, "htm") {
			return string(bytes)
		} else {
			return "<pre>" + string(bytes) + "</pre>"
		}
	}
	return ""
}

func init() {
	elnCmd.AddCommand(addDocumentCmd)
	addDocumentCmd.Flags().StringVar(&addDocArgV.NameArg, "name", "", "A name for the document")
	addDocumentCmd.Flags().StringVar(&addDocArgV.ParentfolderArg, "folder", "", "An id for the folder that will contain the new document")
	addDocumentCmd.Flags().StringVar(&addDocArgV.Tags, "tags", "", "One or more tags, comma separated")
	addDocumentCmd.Flags().StringVar(&addDocArgV.ContentFile, "file", "", "A file of text or HTML content to put in a basic document")
	addDocumentCmd.Flags().StringVar(&addDocArgV.Content, "content", "", "Text or HTML content to put in a basic document")
	addDocumentCmd.Flags().StringVar(&addDocArgV.FormId, "formId", "", "Id for a form")
	addDocumentCmd.Flags().StringVar(&addDocArgV.InputData, "input", "", "File of input data in CSV format for adding field Data to structured documents")
}
