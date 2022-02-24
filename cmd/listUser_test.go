package cmd

import (
	"net/url"
	"strconv"
	"testing"
)

func TestAppend(t *testing.T) {
	a := make([]int, 0)
	a = append(a, 1, 2, 3, 4, 5)
	b := make([]int, 0)
	b = append(b, 6, 7, 8, 9, 99)
	t.Log(a)
	t.Log(b)
	c := append(a, b...)
	t.Log(c)
	map2 := make(map[string]string)
	map2["a"] = "hello"
	t.Log(map2)

}

func TestParseUrl(t *testing.T) {
	urlStr := "http://localhost:8080/api/v1/sysadmin/users?pageSize=5&pageNumber=1&orderBy=creationDate%20desc&tempAccountsOnly=false"
	u, _ := url.Parse(urlStr)

	t.Log(u)
	if u != nil {

		m, _ := url.ParseQuery(u.RawQuery)
		pageNumber, _ := strconv.Atoi(m["pageNumber"][0])
		assertEqualInt(t, 1, pageNumber)
	}

}
