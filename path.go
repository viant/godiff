package godiff

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	//PathKindRoot defines root path kind
	PathKindRoot = PathKind(iota)
	//PathKinField defines field path kind
	PathKinField
	//PathKindKey defines key path kind
	PathKindKey
	//PathKindIndex defines index path kind
	PathKindIndex
)

type (
	//PathKind defines patch kind
	PathKind int
	//Path represents an arbitrary data structure path
	Path struct {
		Kind  PathKind    `json:",omitempty"`
		Path  *Path       `json:",omitempty"`
		Name  string      `json:",omitempty"`
		Index int         `json:",omitempty"`
		Key   interface{} `json:",omitempty"`
	}
)

//Field add fields node
func (p *Path) Field(name string) *Path {
	return &Path{Name: name, Kind: PathKinField, Path: p}
}

//Entry adds map entry node
func (p *Path) Entry(name string) *Path {
	return &Path{Key: name, Kind: PathKindKey, Path: p}
}

//Element adds slice element node
func (p *Path) Element(index int) *Path {
	return &Path{Index: index, Kind: PathKindIndex, Path: p}
}

//String stringifies a path
func (p *Path) String() string {
	builder := new(strings.Builder)
	p.stringify(builder)
	return builder.String()
}

func (p *Path) stringify(builder *strings.Builder) {
	if p.Path != nil {
		p.Path.stringify(builder)
	}
	switch p.Kind {
	case PathKindKey:
		builder.WriteByte('[')
		switch actual := p.Key.(type) {
		case string:
			builder.WriteString(actual)
		case int:
			builder.WriteString(strconv.Itoa(actual))
		case int64:
			builder.WriteString(strconv.Itoa(int(actual)))
		default:
			builder.WriteString(fmt.Sprintf("%v", actual))
		}
		builder.WriteByte(']')
	case PathKinField:
		if builder.Len() > 0 {
			builder.WriteByte('.')
		}
		builder.WriteString(p.Name)
	case PathKindIndex:
		builder.WriteByte('[')
		builder.WriteString(strconv.Itoa(p.Index))
		builder.WriteByte(']')
	}
}
