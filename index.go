package godiff

import (
	"github.com/viant/xunsafe"
	"unsafe"
)

type entry struct {
	index int
	value interface{}
}

type indexer struct{}

func (i *indexer) indexBy(xSlice *xunsafe.Slice, ptr unsafe.Pointer, by string) map[interface{}]*entry {
	if by == "" || by == "." {
		return i.indexPrimitive(xSlice, ptr)
	}
	panic("not yet supported")
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
