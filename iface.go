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
	if from == nil && to == nil {
		return nil
	}
	var fromValue, toValue reflect.Value
	var fromStruct, toStruct reflect.Type

	if from != nil {
		fromValue = reflect.ValueOf(from)
		fromStruct = structType(fromValue.Type())
	}

	if to != nil {
		toValue = reflect.ValueOf(to)
		toStruct = structType(toValue.Type())
	}

	if fromStruct != nil && toStruct != nil {
		differ, err := d.config.registry.Get(fromStruct, toStruct, d.tag)
		if err != nil {
			return err
		}
		if fromValue.Kind() == reflect.Ptr {
			from = fromValue.Elem().Interface()
		}
		if toValue.Kind() == reflect.Ptr {
			to = toValue.Elem().Interface()
		}
		return differ.diff(changeLog, path, from, to, changeType, options)
	}

	if from == nil && toStruct != nil {
		differ, err := d.config.registry.Get(toStruct, toStruct, d.tag)
		if err != nil {
			return err
		}
		if toValue.Kind() == reflect.Ptr {
			to = toValue.Elem().Interface()
		}
		return differ.diff(changeLog, path, from, to, ChangeTypeCreate, options)
	}

	if to == nil && fromStruct != nil {
		differ, err := d.config.registry.Get(fromStruct, fromStruct, d.tag)
		if err != nil {
			return err
		}
		if fromValue.Kind() == reflect.Ptr {
			from = fromValue.Elem().Interface()
		}
		return differ.diff(changeLog, path, from, to, ChangeTypeDelete, options)
	}

	return nil
}

func newIfaceDiffer(config *Config, tag *Tag) (*ifaceDiffer, error) {
	ret := &ifaceDiffer{config: config, tag: tag}
	return ret, nil
}
