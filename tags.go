package materialize

// Tags provides tags information.
type Tags map[string]struct{}

func newTags(tags []string) Tags {
	m := map[string]struct{}{}
	for _, t := range tags {
		m[t] = struct{}{}
	}
	return m
}

func (tags Tags) score(queryTags []string) int {
	pos := 1
	for _, t := range queryTags {
		if _, ok := tags[t]; ok {
			pos++
		}
	}
	if pos > 99 {
		pos = 99
	}

	neg := 0
	m := newTags(queryTags)
	for t := range tags {
		if _, ok := m[t]; !ok {
			neg++
		}
	}
	if neg > 99 {
		neg = 99
	}

	return pos*100 - neg
}
