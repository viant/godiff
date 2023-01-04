package godiff

import (
	"github.com/viant/xunsafe"
	"reflect"
	"strings"
)

type (
	field struct {
		name   string
		from   accessor
		to     accessor
		Kind   reflect.Kind
		tag    *Tag
		differ *Differ
	}

	matcher struct {
		index map[string]*accessor
	}
)

func newField(fromField *xunsafe.Field, fromAccessor accessor, toAccessor accessor, tag *Tag) *field {
	aField := &field{
		name:   fromField.Name,
		from:   fromAccessor,
		to:     toAccessor,
		tag:    tag,
		differ: nil,
	}
	if tag.Name != "" {
		aField.name = tag.Name
	}
	if structType(fromField.Type) != nil && !isTimeType(fromField.Type) {
		aField.Kind = reflect.Struct
	} else if sliceType(fromField.Type) != nil {
		aField.Kind = reflect.Slice
	} else if interfaceType(fromField.Type) != nil {
		aField.Kind = reflect.Interface
	} else if mapType(fromField.Type) != nil {
		aField.Kind = reflect.Map
	}
	return aField
}

func (m *matcher) build(xStruct *xunsafe.Struct, config *Config) {
	m.index = make(map[string]*accessor, 3*len(xStruct.Fields))
	for i := range xStruct.Fields {
		xField := &xStruct.Fields[i]
		tag, _ := ParseTag(string(xField.Tag))
		tag.init(config)
		fieldAccessor := newAccessor(i, xField, tag)
		m.index[xField.Name] = &fieldAccessor
		m.index[strings.ToLower(xField.Name)] = &fieldAccessor
		m.index[m.normKey(xField.Name)] = &fieldAccessor
	}
}

func (m *matcher) match(name string) *accessor {
	if result, ok := m.index[name]; ok {
		return result
	}
	if result, ok := m.index[strings.ToLower(name)]; ok {
		return result
	}
	if result, ok := m.index[m.normKey(name)]; ok {
		return result
	}
	return nil
}

func (m *matcher) normKey(key string) string {
	return strings.ReplaceAll(strings.ToLower(key), "_", "")
}
