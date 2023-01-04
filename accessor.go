package godiff

import (
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type nullifierKind int

const (
	nullifierKindUnspecified = nullifierKind(iota)
	nullifierKindInt
	nullifierKindString
	nullifierKindFloat64
	nullifierKindFloat32
)

type accessor struct {
	pos      int
	deref    bool
	normType reflect.Type //if from,to data type is  different (i.e. int64 vs uint64), norm type is used to reconcile the types
	xType    *xunsafe.Type
	nullifierKind
	*xunsafe.Field
}

func (d *accessor) normalize(value interface{}) (interface{}, error) {
	if value == nil {
		return value, nil
	}
	if d.deref {
		value = d.xType.Deref(value)
	}

	if d.normType != nil {
		ptr := xunsafe.AsPointer(value)
		switch d.normType.Kind() {
		case reflect.Int:
			return *(*int)(ptr), nil
		case reflect.Int32:
			return *(*int32)(ptr), nil
		case reflect.Int64:
			return *(*int64)(ptr), nil
		case reflect.Int16:
			return *(*int16)(ptr), nil
		case reflect.Uint8:
			return *(*uint8)(ptr), nil
		default:
			return nil, fmt.Errorf("unsupported norm type: %T", value)
		}
	}
	return value, nil
}

func (d *accessor) Value(ptr unsafe.Pointer) (value interface{}, err error) {
	if d.IsNil(ptr) {
		return nil, nil
	}
	value = d.Field.Value(ptr)
	if value, err = d.normalize(value); err != nil {
		return nil, err
	}
	return value, nil
}

func (d *accessor) IsNilOrEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	return d.nullifyIfNeeded(value) == nil
}

func (d *accessor) nullifyIfNeeded(value interface{}) interface{} {
	if d.nullifierKind == nullifierKindUnspecified {
		return value
	}
	ptr := xunsafe.AsPointer(value)
	switch d.nullifierKind {
	case nullifierKindInt:
		if *(*int)(ptr) == 0 {
			return nil
		}
	case nullifierKindString:
		if *(*string)(ptr) == "" {
			return nil
		}
	case nullifierKindFloat32:
		if *(*float32)(ptr) == 0 {
			return nil
		}
	case nullifierKindFloat64:
		if *(*float64)(ptr) == 0 {
			return nil
		}
	}
	return value
}

func newAccessor(pos int, field *xunsafe.Field, tag *Tag) accessor {
	result := accessor{
		pos:           pos,
		Field:         field,
		deref:         field.Type.Kind() == reflect.Ptr && (structType(field.Type.Elem()) == nil || isTimeType(field.Type.Elem())),
		nullifierKind: getNullifierKind(tag, field),
	}
	if result.deref {
		result.xType = xunsafe.NewType(field.Type.Elem())
	}
	return result
}

func getNullifierKind(tag *Tag, field *xunsafe.Field) nullifierKind {
	if tag.NullifyEmpty != nil && *tag.NullifyEmpty {
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
			return nullifierKindInt
		case reflect.Float64:
			return nullifierKindFloat64
		case reflect.Float32:
			return nullifierKindFloat32
		case reflect.String:
			return nullifierKindString
		}
	}
	return nullifierKindUnspecified
}
