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

func splitTags(s string) Tags {
	return newTags(split(s))
}

func TestTags_score(t *testing.T) {
	check := func(tags, query string, expectedScore int) {
		t.Helper()
		sc := splitTags(tags).score(splitTags(query))
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

func TestTags_joinKeys(t *testing.T) {
	check := func(tags []string, exp string) {
		t.Helper()
		k := newTags(tags).joinKeys()
		if k != exp {
			t.Errorf("joinKeys not match: %s (expected=%s)", k, exp)
		}
	}

	check([]string{}, "")
	check([]string{"foo", "bar"}, "bar foo")
	check([]string{"foo", "bar", "baz"}, "bar baz foo")
	check([]string{"f o", `b\r`, "baz"}, `b\\r baz f\ o`)
}
