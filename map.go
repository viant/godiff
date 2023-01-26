package godiff

import (
	"fmt"
	"reflect"
)

type mapDiffer struct {
	config        *Config
	itemDiffer    *Differ
	isInterface   bool
	tag           *Tag
	from, to      reflect.Type
	isStringIface bool
}

func (s *mapDiffer) diff(changeLog *ChangeLog, path *Path, from, to interface{}, changeType ChangeType) error {
	if from == nil && to == nil {
		return nil
	}
	if !s.isStringIface {
		return fmt.Errorf("type: %T not supported yet", from)
	}
	var ok bool
	var fromMap, toMap map[string]interface{}
	if from != nil {
		if fromMap, ok = from.(map[string]interface{}); !ok {
			return fmt.Errorf("invalid from type: %T", from)
		}
	}
	if to != nil {
		if toMap, ok = to.(map[string]interface{}); !ok {
			return fmt.Errorf("invalid to type: %T", from)
		}
	}
	var err error

	if from == nil {
		for k, v := range toMap {
			if err = s.diffIfaceElement(changeLog, path, nil, v, k, ChangeTypeDelete); err != nil {
				return err
			}
		}
	} else if to == nil {
		for k, v := range fromMap {
			if err = s.diffIfaceElement(changeLog, path, v, nil, k, ChangeTypeCreate); err != nil {
				return err
			}
		}
	} else {

		for k, fromItem := range fromMap {
			toItem := toMap[k]
			if err = s.diffIfaceElement(changeLog, path, fromItem, toItem, k, ChangeTypeCreate); err != nil {
				return err
			}
		}

		for k, toItem := range toMap {
			if _, has := fromMap[k]; has {
				continue
			}
			if err = s.diffIfaceElement(changeLog, path, nil, toItem, k, ChangeTypeCreate); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *mapDiffer) diffIfaceElement(changeLog *ChangeLog, path *Path, from, to interface{}, key string, changeType ChangeType) error {
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
	return itemDiffer.diff(changeLog, path.Entry(key), from, to, changeType)
}

func newMapDiffer(from, to reflect.Type, config *Config, tag *Tag) (*mapDiffer, error) {
	isStringIface := from.Key().Kind() == reflect.String && from.Elem().Kind() == reflect.Interface
	return &mapDiffer{config: config, tag: tag, from: from, to: to, isStringIface: isStringIface}, nil
}
