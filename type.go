package godiff

import (
	"reflect"
	"time"
)

func structType(p reflect.Type) reflect.Type {
	switch p.Kind() {
	case reflect.Ptr:
		return structType(p.Elem())
	case reflect.Struct:
		return p
	}
	return nil
}

func sliceType(p reflect.Type) reflect.Type {
	switch p.Kind() {
	case reflect.Ptr:
		return sliceType(p.Elem())
	case reflect.Slice:
		return p
	}
	return nil
}

func interfaceType(p reflect.Type) reflect.Type {
	switch p.Kind() {
	case reflect.Ptr:
		return sliceType(p.Elem())
	case reflect.Interface:
		return p
	}
	return nil
}

func mapType(p reflect.Type) reflect.Type {
	switch p.Kind() {
	case reflect.Ptr:
		return mapType(p.Elem())
	case reflect.Map:
		return p
	}
	return nil
}

var timeType = reflect.TypeOf(time.Time{})
var stringsType = reflect.TypeOf([]string{})
var stringMapType = reflect.TypeOf(map[string]interface{}{})

func isTimeType(p reflect.Type) bool {
	sType := structType(p)
	if sType == nil {
		return false
	}
	return timeType == sType
}
