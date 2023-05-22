package godiff

import (
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
)

type (
	structDiffer struct {
		config           *Config
		from             *xunsafe.Struct
		to               *xunsafe.Struct
		fromType         reflect.Type
		toType           reflect.Type
		fields           []*field
		presenceProvider PresenceProvider
	}
)

func (s *structDiffer) diff(changeLog *ChangeLog, path *Path, from, to interface{}, changeType ChangeType, options *Options) error {
	fromPtr := xunsafe.AsPointer(from)
	toPtr := xunsafe.AsPointer(to)
	var err error
	var fromValue, toValue interface{}

	for _, field := range s.fields {
		if fromValue, err = field.from.Value(fromPtr); err != nil {
			changeLog.AddError(path.Field(field.name), err)
			continue
		}
		if toValue, err = field.to.Value(toPtr); err != nil {
			changeLog.AddError(path.Field(field.name), err)
			continue
		}
		if options.presence {
			hasPtr := toPtr
			if s.presenceProvider.Holder != nil {
				if s.presenceProvider.Holder.IsNil(toPtr) {
					hasPtr = fromPtr
				}
			}
			if !s.presenceProvider.IsFieldSet(hasPtr, int(field.to.Index)) {
				continue //skip diff, to/dest is not set
			}
			if field.tag != nil && field.tag.Presence {
				continue //do not compare presence tag
			}
		}

		if fromValue == nil && toValue == nil {
			continue
		}
		if field.differ != nil {
			diffChangeType := ChangeTypeUpdate
			if field.Kind == reflect.Slice {
				if fromValue != nil {
					fromValue = field.from.Addr(fromPtr)
				} else {
					diffChangeType = ChangeTypeCreate
				}
				if toValue != nil {
					toValue = field.to.Addr(toPtr)
					diffChangeType = ChangeTypeDelete
				}
			}
			if fromValue == nil {

			}
			if err = field.differ.diff(changeLog, path.Field(field.name), fromValue, toValue, diffChangeType, options); err != nil {
				return err
			}
			continue
		}

		switch changeType {
		case ChangeTypeCreate:
			if !field.to.IsNil(toPtr) {
				changeLog.AddCreate(path.Field(field.name), toValue)
				continue
			}
		case ChangeTypeDelete:
			if !field.from.IsNil(fromPtr) {
				changeLog.AddDelete(path.Field(field.name), fromValue)
				continue
			}
		}
		if !matches(fromValue, toValue) {
			changeLog.AddUpdate(path.Field(field.name), fromValue, toValue)
		}
	}
	return nil
}

func (s *structDiffer) matchFields() error {
	var fields = make([]*field, 0, len(s.from.Fields))
	typesMatches := s.to == s.from
	matcher := matcher{}
	if !typesMatches {
		matcher.build(s.to, s.config)
	}

	var filedPos = map[string]int{}
	for i := range s.from.Fields {
		fromField := &s.from.Fields[i]
		filedPos[fromField.Name] = int(fromField.Index)
		tag, err := ParseTag(fromField.Tag.Get(s.config.TagName))
		if err != nil {
			return err
		}
		if tag.Presence {
			s.presenceProvider.Holder = fromField
			continue
		}
		if tag.Ignore {
			continue
		}
		tag.init(s.config)

		fromAccessor := newAccessor(i, fromField, tag)
		toAccessor := fromAccessor
		if !typesMatches {
			if match := matcher.match(fromField.Name); match != nil {
				toAccessor = *match
			} else {
				continue
			}
		}
		aField := newField(fromField, fromAccessor, toAccessor, tag)
		fields = append(fields, aField)

		switch aField.Kind {
		case reflect.Map:
			differ, err := newMapDiffer(aField.from.Type, aField.to.Type, s.config, aField.tag)
			if err != nil {
				return err
			}
			aField.differ = &Differ{config: s.config, mapDiffer: differ}
		case reflect.Struct:
			if aField.from.Type == s.fromType {
				aField.differ = &Differ{config: s.config, structDiffer: s}
				continue
			}
			differ, err := newStructDiffer(aField.from.Type, aField.to.Type, s.config)
			if err != nil {
				return err
			}
			aField.differ = &Differ{config: s.config, structDiffer: differ}
		case reflect.Slice:
			differ, err := newSliceDiffer(aField.from.Type, aField.to.Type, s.config, aField.tag)
			if err != nil {
				return err
			}
			aField.differ = &Differ{config: s.config, sliceDiffer: differ}
		case reflect.Interface:
			differ, err := newIfaceDiffer(s.config, aField.tag)
			if err != nil {
				return err
			}
			aField.differ = &Differ{config: s.config, ifaceDiffer: differ}
		}
		if aField.tag != nil && aField.tag.decodable() && aField.differ == nil {
			if aField.differ, err = New(aField.from.Type, aField.to.Type, WithConfig(s.config), WithTag(aField.tag)); err != nil {
				return err
			}
		}
	}
	if s.presenceProvider.Holder != nil {
		s.presenceProvider.Init(filedPos)
	}
	s.fields = fields
	return nil
}

func newStructDiffer(from, to reflect.Type, config *Config) (*structDiffer, error) {
	var result = structDiffer{config: config, fromType: from, toType: timeType}

	fromType := structType(from)
	if fromType == nil {
		return nil, fmt.Errorf("invalid 'from' struct type: %s", from.String())
	}
	toType := structType(to)
	if toType == nil {
		return nil, fmt.Errorf("invalid 'to' struct type: %s", to.String())
	}
	result.from = xunsafe.NewStruct(fromType)
	result.to = result.from
	if toType != fromType {
		result.to = xunsafe.NewStruct(toType)
	}
	if err := result.matchFields(); err != nil {
		return nil, err
	}

	return &result, nil
}
