package cmd

import (
	"testing"

	"github.com/richarda23/rspace-client-go/rspace"
)

func TestIdFromGlobalId(t *testing.T) {
	id, err := idFromGlobalId("not an id")
	if err != nil {
		t.Fatalf("Error should b non-nil")
	}
	id, err = idFromGlobalId("FL12334")
	assertIdEquals(t, id, 12334)

	id, err = idFromGlobalId("12334")
	assertIdEquals(t, id, 12334)

	id, err = idFromGlobalId("")
	assertIdEquals(t, id, 0)
}

func TestAbbreviate(t *testing.T) {
	assertEqualString(t, "abcd..", abbreviate("abcdefg", 6))
	assertEqualString(t, "abcde", abbreviate("abcde", 3))
	assertEqualString(t, "abc", abbreviate("abc", 3))
}

func assertEqualString(t *testing.T, expected, got string) {
	if got != expected {
		t.Fatalf("got %s, expected %s", got, expected)
	}
}

func TestMaxColWidth(t *testing.T) {
	doc1 := rspace.IdentifiableNamable{Name: "abcde"}
	doc2 := rspace.IdentifiableNamable{Name: "abcdefdfkdsfj"}
	docs := make([]rspace.BasicInfo, 0)
	docs = append(docs, doc1, doc2)
	maxLen := getMaxNameLength(docs)
	if maxLen != 13 {
		t.Fatalf("Expected max size 13 but was %d", maxLen)
	}
	/// > 25 absolute max
	docs = append(docs, rspace.IdentifiableNamable{Name: "abcdefdfsdfdsfdsfdsfdsfdsfdsfsdfdsfkdsfj"})
	maxLen = getMaxNameLength(docs)
	if maxLen != 25 {
		t.Fatalf("Expected max size 25 but was %d", maxLen)
	}

}

func assertIdEquals(t *testing.T, got, expected int) {
	if got != expected {
		t.Fatalf("Expected %d but was %d", expected, got)
	}
}
