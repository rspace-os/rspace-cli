package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/richarda23/rspace-client-go/rspace"
	//"errors"
)

const (
	DISPLAY_TIMESTAMP_WIDTH = 16
)

func exitWithStdErrMsg(message string) {
	messageStdErr(message)
	os.Exit(1)
}
func exitWithErr(err error) {
	exitWithStdErrMsg(err.Error())
}
func messageStdErr(message string) {
	fmt.Fprintln(os.Stderr, message)
}

func message(writer io.Writer, message string) {
	fmt.Fprintln(writer, message)
}

func configurePagination() rspace.RecordListingConfig {
	cfg := rspace.NewRecordListingConfig()
	if len(sortOrderArg) > 0 && validateArrayContains(validSortOrders, []string{sortOrderArg}) {
		cfg.SortOrder = sortOrderArg
	}
	if len(orderByArg) > 0 && validateArrayContains(validRecordOrders, []string{orderByArg}) {
		cfg.OrderBy = orderByArg
	}
	// name sort is asc by default
	if orderByArg == "name" && len(sortOrderArg) == 0 {
		cfg.SortOrder = "asc"
	}
	if pageSizeArg > 0 {
		cfg.PageSize = pageSizeArg
	}
	return cfg
}

func idsFromGlobalIds(globalIds []string) []int {
	ids := make([]int, 0)
	for _, v := range globalIds {
		id, err := idFromGlobalId(v)
		if err != nil || id == 0 {
			messageStdErr(fmt.Sprintf("%s is not a valid identifier, skipping", v))
		} else {
			ids = append(ids, id)
		}
	}
	return ids
}

// idFromGlobalId matches either globalId string or a numeric id, returning
// the numeric id or an Error if input string cannot be parsed
// If the globalId is empty string, returns 0.
func idFromGlobalId(globalId string) (int, error) {
	if globalId == "" {
		return 0, nil
	}
	v, _ := strconv.Atoi(globalId)
	if v > 0 {
		return v, nil
	}

	match, err := regexp.MatchString("^[A-Z]{2}\\d+", globalId)
	if match {
		globalId = globalId[2:] //handle globalId
		v, _ := strconv.Atoi(globalId)
		return v, nil
	} else {
		return 0, err
	}
}
func prettyMarshal(anything interface{}) string {
	bytes, _ := json.MarshalIndent(anything, "", "\t")
	return string(bytes)
}

type columnDef struct {
	Title string
	Width int
}

type TableResult struct {
	Headers []columnDef
	Content [][]string
}

type ResultListFormatter interface {
	ToJson() string
	ToTable() *TableResult
	ToQuiet() []identifiable
}

func printTable(ctx *Context, table *TableResult) {
	printTableHeaders(ctx, table.Headers)
	printContent(ctx, table)
}
func printCsv(ctx *Context, table *TableResult) {
	writer := csv.NewWriter(ctx.Writer)
	writer.Write(columnDefsToString(table.Headers))
	if err := writer.WriteAll(table.Content); err != nil {
		exitWithErr(err)
	}
}

func columnDefsToString(headers []columnDef) []string {
	rowToPrint := make([]string, 0)
	for _, header := range headers {
		rowToPrint = append(rowToPrint, header.Title)
	}
	return rowToPrint
}

func printContent(ctx *Context, table *TableResult) {
	for _, row := range table.Content {
		rowToPrint := make([]string, 0)
		for i, cell := range row {
			toPrint := abbreviate(cell, table.Headers[i].Width)
			toPrint = fmt.Sprintf("%-[1]*s", table.Headers[i].Width, toPrint)
			rowToPrint = append(rowToPrint, toPrint)
		}
		ctx.write(strings.Join(rowToPrint, "\t"))
	}
}
func printTableHeaders(ctx *Context, headers []columnDef) {
	headersToPrint := make([]string, 0)
	for _, header := range headers {
		toPrint := abbreviate(header.Title, header.Width)
		toPrint = fmt.Sprintf("%-[1]*s", header.Width, toPrint)
		headersToPrint = append(headersToPrint, toPrint)
	}
	ctx.write(strings.Join(headersToPrint, "\t"))
}
func abbreviate(toAbbreviate string, maxLen int) string {
	if maxLen > 3 && len(toAbbreviate) > maxLen {
		toAbbreviate = toAbbreviate[0:(maxLen-2)] + ".."
	}
	return toAbbreviate
}

type identifiable struct {
	Id string
}

//
func printIds(ctx *Context, source []identifiable) {
	for _, item := range source {
		ctx.write(item.Id)
	}
}

// gets the length of the longest name in the result list
func getMaxNameLength(results []rspace.BasicInfo) int {
	var maxPossible float64 = 25
	var currLongest float64 = 0
	for _, res := range results {
		if nameLen := float64(len(res.GetName())); nameLen > currLongest {
			currLongest = math.Min(maxPossible, nameLen)
		}
	}
	return int(currLongest)
}

func globalIdListToIntList(slice []string) []int {
	results := make([]int, 0)
	for _, v := range slice {
		id, _ := idFromGlobalId(v)
		if id != 0 {
			results = append(results, id)
		}
	}
	return results
}

func maxInt(ints ...int) int {
	max := ints[0]
	for _, i := range ints {
		if i > max {
			max = i
		}
	}
	return max
}
