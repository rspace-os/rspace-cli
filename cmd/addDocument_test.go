package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/richarda23/rspace-client-go/rspace"
)

func TestGetContent(t *testing.T) {
	// plain text content
	args := addDocArgs{}
	args.Content = "abcdefg"
	content := getContent(args)

	if content != "abcdefg" {
		t.Fatalf("unexpected content")
	}

	// text file
	args.Content = ""
	args.ContentFile = "testData/textContent.txt"
	content = getContent(args)
	if !strings.Contains(content, "<pre>") {
		t.Fatalf("expected content should be wrapped in <pre> tag but was %s", content)
	}

	//html file
	args.ContentFile = "testData/textContent.html"
	content = getContent(args)
	if !strings.Contains(content, "<p> some html </p>") {
		t.Fatalf("expected verbatim html but was '%s'", content)
	}
}

func TestAddDocumentOK(t *testing.T) {
	ctx := Context{}
	outWriter := bytes.NewBufferString("")
	ctx.Writer = outWriter
	// output only gets ID field, so we don't have to make a fully populated DocumentInfo
	ctx.Format = outputFmt("quiet")
	// supply all required fields
	doAddDocRun(addDocArgs{Content: "hello", NameArg: "name", Tags: "tags"}, &ctx, &DocumentAddOK{})
	assertOKDoc(outWriter, t)

	// content is optional
	doAddDocRun(addDocArgs{NameArg: "name", Tags: "tags"}, &ctx, &DocumentAddOK{})
	assertOKDoc(outWriter, t)

	// tags is optional
	doAddDocRun(addDocArgs{NameArg: "name"}, &ctx, &DocumentAddOK{})
	assertOKDoc(outWriter, t)

	// name is optional
	doAddDocRun(addDocArgs{}, &ctx, &DocumentAddOK{})
	assertOKDoc(outWriter, t)
}

func TestReadCsv(t *testing.T) {
	dataFile := "testData/ExperimentData.csv"
	ioReader, err := os.Open(dataFile)
	if err != nil {
		t.Fatalf("unexpected error opening test-data file")
	}
	data, _ := readCsvFile(ioReader)
	if len(data) != 3 {
		t.Fatalf("Expected 3 rows")
	}
	if data[0][0] != "Date" {
		t.Fatalf("1st row should be header row")
	}
}

func TestValidateCsvInput(t *testing.T) {
	tooFewRows := `name,desc,otherFieldName`
	r := strings.NewReader(tooFewRows)
	data, err := readCsvFile(r)
	fmt.Println(data)
	err = validateCsvInput(data)
	if err == nil {
		t.Fatalf("Should fail - too few rows")
	}

	misMatchedRows := `f1,f2,f3
	d1,d2,d3
	d4,d5
	d6,d7,d8
	`
	r = strings.NewReader(misMatchedRows)
	data, err = readCsvFile(r)
	err = validateCsvInput(data)
	if err == nil {
		t.Fatalf("Should fail - mismatched rows")
	}

}
func assertOKDoc(outWriter io.Reader, t *testing.T) {
	outString, _ := ioutil.ReadAll(outWriter)
	// quiet ouput only outputs the ID
	if string(outString) != "1234\n" {
		t.Fatalf("Expected quiet output of id 1234 but was '%s'", string(outString))
	}
}

// stub for returning a created document
type DocumentAddOK struct{}

// simulates a correct response
func (ds *DocumentAddOK) NewBasicDocumentWithContent(name, tags, content string) (*rspace.Document, error) {
	namable := &rspace.IdentifiableNamable{Id: 1234}
	rc := rspace.DocumentInfo{IdentifiableNamable: namable}
	return &rspace.Document{DocumentInfo: &rc}, nil
}

func (ds *DocumentAddOK) NewDocumentWithContent(post *rspace.DocumentPost) (*rspace.Document, error) {
	namable := &rspace.IdentifiableNamable{Id: 1234}
	rc := rspace.DocumentInfo{IdentifiableNamable: namable}
	return &rspace.Document{DocumentInfo: &rc}, nil
}
