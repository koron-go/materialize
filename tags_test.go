package materialize

import (
	"strings"
	"testing"
)

func split(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, " ")
}

func TestTagsScore(t *testing.T) {
	check := func(tags, query string, expectedScore int) {
		t.Helper()
		sc := newTags(split(tags)).score(split(query))
		if sc != expectedScore {
			t.Errorf("score not match: %d (expected %d) tags=%q query=%q",
				sc, expectedScore, tags, query)
		}
	}

	check("", "", 100)
	check("foo", "", 99)
	check("foo", "foo", 200)
	check("foo", "foo bar", 200)
	check("foo bar", "foo", 199)
}
