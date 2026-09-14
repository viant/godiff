package godiff

import (
	"reflect"
)

type (
	ifaceDiffer struct {
		config *Config
		tag    *Tag
	}
)

func (d *ifaceDiffer) diff(changeLog *ChangeLog, path *Path, from, to interface{}, changeType ChangeType, options *Options) error {
	from, fromType := d.value(from)
	to, toType := d.value(to)
	if from == nil && to == nil {
		switch changeType {
		case ChangeTypeCreate:
			changeLog.AddCreate(path, nil)
		case ChangeTypeDelete:
			changeLog.AddDelete(path, nil)
		}
		return nil
	}
	if fromType == nil {
		fromType = toType
	}
	if toType == nil {
		toType = fromType
	}
	differ, err := d.config.registry.Get(fromType, toType, d.tag)
	if err != nil {
		return err
	}
	return differ.diff(changeLog, path, from, to, changeType, options)
}

// value keeps concrete type authority while safely unwrapping pointer values.
func (d *ifaceDiffer) value(value interface{}) (interface{}, reflect.Type) {
	if value == nil {
		return nil, nil
	}
	v := reflect.ValueOf(value)
	typ := v.Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if v.IsNil() {
			for typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			return nil, typ
		}
		v = v.Elem()
	}
	return v.Interface(), typ
}

func newIfaceDiffer(config *Config, tag *Tag) (*ifaceDiffer, error) {
	ret := &ifaceDiffer{config: config, tag: tag}
	return ret, nil
}
