package godiff

import (
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type entry struct {
	index int
	value interface{}
}

type indexer struct {
	field *xunsafe.Field
}

func (i *indexer) indexBy(xSlice *xunsafe.Slice, ptr unsafe.Pointer, by string) map[interface{}]*entry {
	if by == "" || by == "." {
		return i.indexPrimitive(xSlice, ptr)
	}
	elemType := xSlice.Type.Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%s not yet supported", elemType.String()))
	}
	i.field = xunsafe.FieldByName(elemType, by)
	return i.indexByField(xSlice, ptr)
}

func (i *indexer) indexByField(xSlice *xunsafe.Slice, ptr unsafe.Pointer) map[interface{}]*entry {
	var result = make(map[interface{}]*entry)
	l := xSlice.Len(ptr)
	for j := 0; j < l; j++ {
		value := xSlice.ValueAt(ptr, j)
		ptr := xunsafe.AsPointer(value)
		key := i.field.Value(ptr)
		result[key] = &entry{index: j, value: value}
	}
	return result
}

func (i *indexer) indexPrimitive(xSlice *xunsafe.Slice, ptr unsafe.Pointer) map[interface{}]*entry {
	var result = make(map[interface{}]*entry)
	l := xSlice.Len(ptr)
	for i := 0; i < l; i++ {
		value := xSlice.ValueAt(ptr, i)
		result[value] = &entry{index: i, value: value}
	}
	return result
}
