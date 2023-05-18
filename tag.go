package godiff

import (
	"fmt"
	"strconv"
	"strings"
)

//Tag represents a tag
type Tag struct {
	Name          string
	Presence      bool
	PairSeparator string
	PairDelimiter string
	pairDelimiter []string
	ItemSeparator string

	Whitespace   string
	IndexBy      string
	Sort         bool
	TimeLayout   string
	Precision    *int
	Ignore       bool
	NullifyEmpty *bool
}

func (t *Tag) decodable() bool {
	return t.PairDelimiter != "" || t.ItemSeparator != ""
}

func (t *Tag) init(config *Config) {
	if t.NullifyEmpty == nil {
		t.NullifyEmpty = config.NullifyEmpty
	}
	if t.PairDelimiter != "" {
		t.pairDelimiter = strings.Split(t.PairDelimiter, "|")
	}
	if t.PairDelimiter != "" && t.PairSeparator == "" {
		t.PairSeparator = "="
	}
}

func (t *Tag) removeWhitespace(value string) string {
	if t.Whitespace == "" {
		return value
	}
	value = strings.TrimSpace(value)
	for _, ws := range strings.Split(t.Whitespace, "|") {
		value = strings.ReplaceAll(value, ws, "")
	}
	return value
}

func (t Tag) clone() *Tag {
	clone := t
	return &clone
}

//ParseTag parses tag
func ParseTag(tagString string) (*Tag, error) {
	tag := &Tag{}
	if tagString == "-" {
		tag.Ignore = true
		return tag, nil
	}

	elements := strings.Split(tagString, ",")
	if len(elements) == 0 {
		return tag, nil
	}
	for _, element := range elements {
		if count := strings.Count(element, "$coma"); count > 0 {
			element = strings.Replace(element, "$coma", ",", count)
		}
		nv := strings.Split(element, "=")
		switch len(nv) {
		case 2:
			switch strings.ToLower(strings.TrimSpace(nv[0])) {
			case "presence":
				tag.Presence = true
			case "ignore":
				tag.Ignore = true
			case "name":
				tag.Name = strings.TrimSpace(nv[1])
			case "indexby":
				tag.IndexBy = strings.TrimSpace(nv[1])
			case "timelayout":
				tag.TimeLayout = strings.TrimSpace(nv[1])
			case "precision":
				precision, err := strconv.Atoi(strings.TrimSpace(nv[1]))
				if err != nil {
					return nil, fmt.Errorf("invalid precission: %w, %v", err, nv[1])
				}
				tag.Precision = &precision
			case "whitespace":
				tag.Whitespace = strings.TrimSpace(nv[1])
			case "pairseparator":
				tag.PairSeparator = strings.TrimSpace(nv[1])
			case "pairdelimiter":
				tag.PairDelimiter = strings.TrimSpace(nv[1])
			case "itemseparator":
				tag.ItemSeparator = strings.TrimSpace(nv[1])
			case "sort":
				tag.Sort, _ = strconv.ParseBool(strings.TrimSpace(nv[1]))
			}
			continue
		case 1:
			tag.Name = strings.TrimSpace(element)

		}
	}
	return tag, nil
}
