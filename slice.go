package godiff

import (
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type sliceDiffer struct {
	config      *Config
	itemDiffer  *Differ
	isInterface bool
	tag         *Tag
	fromSlice   *xunsafe.Slice
	fromIndexer indexer
	toSlice     *xunsafe.Slice
	toIndexer   indexer
}

func (s *sliceDiffer) diff(changeLog *ChangeLog, path *Path, from, to interface{}, changeType ChangeType, options *Options) error {
	if s.isInterface {
		return s.diffIfacedSlice(changeLog, path, from, to, changeType, options)

	}
	return s.diffTypedSlice(changeLog, path, from, to, changeType, options)
}

func (s *sliceDiffer) diffTypedSlice(changeLog *ChangeLog, path *Path, from interface{}, to interface{}, changeType ChangeType, options *Options) error {
	fromPtr := xunsafe.AsPointer(from)
	toPtr := xunsafe.AsPointer(to)

	if s.tag.Sort {
		if s.itemDiffer == nil {
			from = sortPrimitive(from)
			fromPtr = xunsafe.AsPointer(from)
			to = sortPrimitive(to)
			toPtr = xunsafe.AsPointer(to)
		}

	}
	changeType = ChangeTypeUpdate
	var fromLen, toLen int
	if from != nil {
		if fromLen = s.fromSlice.Len(fromPtr); fromLen == 0 {
			changeType = ChangeTypeCreate
		}
	}
	if to != nil {
		if toLen = s.toSlice.Len(toPtr); toLen == 0 && fromLen > 0 {
			changeType = ChangeTypeDelete
		}
	}

	if by := s.tag.IndexBy; by != "" && fromLen > 0 && toLen > 0 {
		fromIndex := s.fromIndexer.indexBy(s.fromSlice, fromPtr, by)
		toIndex := s.toIndexer.indexBy(s.toSlice, toPtr, by)
		return s.diffIndexedElement(changeLog, path, fromIndex, toIndex, changeType, options)
	}

	return s.diffSliceElements(changeLog, path, changeType, fromPtr, toPtr, fromLen, toLen, options)
}

func (s *sliceDiffer) diffSliceElements(changeLog *ChangeLog, path *Path, changeType ChangeType, fromPtr, toPtr unsafe.Pointer, fromLen int, toLen int, options *Options) error {
	var repeat int
	if repeat = fromLen; repeat < toLen {
		repeat = toLen
	}
	var err error
	for i := 0; i < repeat; i++ {
		switch changeType {
		case ChangeTypeCreate:
			value := s.toSlice.ValueAt(toPtr, i)
			if s.itemDiffer != nil {
				if err = s.itemDiffer.diff(changeLog, path.Element(i), nil, value, changeType, options); err != nil {
					return err
				}
				continue
			}
			changeLog.AddCreate(path.Element(i), value)

		case ChangeTypeDelete:
			value := s.fromSlice.ValueAt(fromPtr, i)
			if s.itemDiffer != nil {
				if err = s.itemDiffer.diff(changeLog, path.Element(i), value, nil, changeType, options); err != nil {
					return err
				}
				continue
			}
			changeLog.AddCreate(path.Element(i), value)
		case ChangeTypeUpdate:
			if fromLen <= i {
				value := s.toSlice.ValueAt(toPtr, i)
				if s.itemDiffer != nil {
					if err = s.itemDiffer.diff(changeLog, path.Element(i), nil, value, ChangeTypeCreate, options); err != nil {
						return err
					}
					continue
				}
				changeLog.AddCreate(path.Element(i), value)
				continue
			} else if toLen <= i {
				value := s.fromSlice.ValueAt(fromPtr, i)
				if s.itemDiffer != nil {
					if err = s.itemDiffer.diff(changeLog, path.Element(i), value, nil, ChangeTypeDelete, options); err != nil {
						return err
					}
					continue
				}
				changeLog.AddCreate(path.Element(i), value)
				continue
			}

			fromItem := s.fromSlice.ValueAt(fromPtr, i)
			toItem := s.toSlice.ValueAt(toPtr, i)
			if s.itemDiffer != nil {
				if err = s.itemDiffer.diff(changeLog, path.Element(i), fromItem, toItem, ChangeTypeUpdate, options); err != nil {
					return err
				}
				continue
			}
			if !matches(fromItem, toItem) {
				changeLog.AddUpdate(path.Element(i), fromItem, toItem)
			}
		}
	}
	return nil
}

func (s *sliceDiffer) diffIndexedElement(changeLog *ChangeLog, path *Path, fromIndex map[interface{}]*entry, toIndex map[interface{}]*entry, changeType ChangeType, options *Options) error {
	for k := range fromIndex {
		fromValue := fromIndex[k]
		toValue, ok := toIndex[k]
		if !ok {
			changeLog.AddDelete(path.Element(fromValue.index), fromValue.value)
			continue
		}
		if s.itemDiffer != nil {
			if err := s.itemDiffer.diff(changeLog, path, fromValue.value, toValue.value, ChangeTypeUpdate, options); err != nil {
				return err
			}
			continue
		}
		if !matches(fromValue.value, toValue.value) {
			changeLog.AddUpdate(path.Element(fromValue.index), fromValue.value, toValue.value)
		}
	}

	for k := range toIndex {
		if _, ok := fromIndex[k]; ok {
			continue
		}
		toValue := toIndex[k]
		changeLog.AddCreate(path.Element(toValue.index), toValue.value)
	}
	return nil
}

func (s *sliceDiffer) diffIfacedSlice(changeLog *ChangeLog, path *Path, from interface{}, to interface{}, changeType ChangeType, options *Options) error {
	var repeat int

	changeType = ChangeTypeUpdate
	fromPtr := xunsafe.AsPointer(from)
	var fromLen = -1
	if from != nil {
		fromLen = s.fromSlice.Len(fromPtr)
	} else {
		changeType = ChangeTypeCreate
	}

	toLen := -1
	toPtr := xunsafe.AsPointer(to)
	if to != nil {
		toLen = s.fromSlice.Len(toPtr)
	} else {
		changeType = ChangeTypeDelete
	}

	if repeat = fromLen; repeat < toLen {
		repeat = toLen
	}

	var err error
	for i := 0; i < repeat; i++ {
		switch changeType {
		case ChangeTypeCreate:
			value := s.toSlice.ValueAt(toPtr, i)
			if err = s.diffIfaceElement(changeLog, path, nil, value, i, changeType, options); err != nil {
				return err
			}
			continue

		case ChangeTypeDelete:
			value := s.fromSlice.ValueAt(fromPtr, i)
			if err = s.diffIfaceElement(changeLog, path, value, nil, i, changeType, options); err != nil {
				return err
			}
		case ChangeTypeUpdate:

			if i < fromLen && i >= toLen {
				value := s.toSlice.ValueAt(toPtr, i)
				if err = s.diffIfaceElement(changeLog, path, nil, value, i, ChangeTypeCreate, options); err != nil {
					return err
				}
				continue
			}

			if i < toLen && i >= fromLen {
				value := s.toSlice.ValueAt(toPtr, i)
				if err = s.diffIfaceElement(changeLog, path, value, nil, i, ChangeTypeDelete, options); err != nil {
					return err
				}
				continue
			}

			fromItem := s.fromSlice.ValueAt(fromPtr, i)
			toItem := s.toSlice.ValueAt(toPtr, i)
			if err = s.diffIfaceElement(changeLog, path, fromItem, toItem, i, ChangeTypeUpdate, options); err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

func (s *sliceDiffer) diffIfaceElement(changeLog *ChangeLog, path *Path, from, to interface{}, index int, changeType ChangeType, options *Options) error {
	var fromValue, toValue reflect.Value
	if from == nil && to == nil {
		return nil
	}
	if to != nil {
		toValue = reflect.ValueOf(to)
		fromValue = toValue
	} else if from != nil {
		fromValue = reflect.ValueOf(from)
		toValue = fromValue
	} else {
		toValue = reflect.ValueOf(to)
		fromValue = reflect.ValueOf(from)
	}

	if fromValue.Kind() == reflect.Ptr {
		fromValue = fromValue.Elem()
		from = fromValue.Interface()
	}
	if toValue.Kind() == reflect.Ptr {
		toValue = toValue.Elem()
		to = toValue.Elem()
	}

	itemDiffer, err := s.config.registry.Get(fromValue.Type(), toValue.Type(), s.tag)
	if err != nil {
		return err
	}
	return itemDiffer.diff(changeLog, path.Element(index), from, to, changeType, options)
}

func newSliceDiffer(from, to reflect.Type, config *Config, tag *Tag) (*sliceDiffer, error) {
	if tag == nil {
		tag = &Tag{}
	}
	result := &sliceDiffer{
		config:    config,
		fromSlice: xunsafe.NewSlice(from),
		tag:       tag,
	}

	result.toSlice = result.fromSlice
	if from != to {
		result.toSlice = xunsafe.NewSlice(to)
	}

	if interfaceType(result.toSlice.Type.Elem()) != nil || interfaceType(result.fromSlice.Type.Elem()) != nil {
		result.isInterface = true
		return result, nil
	}

	fromElem := structType(from.Elem())
	toElem := structType(to.Elem())
	if fromElem != nil {
		differ, err := newStructDiffer(fromElem, toElem, config)
		if err != nil {
			return nil, err
		}
		result.itemDiffer = &Differ{structDiffer: differ}
	}
	return result, nil
}
