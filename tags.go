package materialize

import (
	"sort"
	"strings"
)

// Tags provides tags information.
type Tags map[string]struct{}

func newTags(tags []string) Tags {
	m := map[string]struct{}{}
	for _, t := range tags {
		m[t] = struct{}{}
	}
	return m
}

func (tags Tags) score(other Tags) int {
	pos := 1
	for t := range other {
		if _, ok := tags[t]; ok {
			pos++
		}
	}
	if pos > 99 {
		pos = 99
	}

	neg := 0
	for t := range tags {
		if _, ok := other[t]; !ok {
			neg++
		}
	}
	if neg > 99 {
		neg = 99
	}

	return pos*100 - neg
}

var tagEscape = strings.NewReplacer(" ", `\ `, `\`, `\\`)

func (tags Tags) joinKeys() string {
	if len(tags) == 0 {
		return ""
	}
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, tagEscape.Replace(k))
	}
	sort.Strings(keys)
	return strings.Join(keys, " ")
}
