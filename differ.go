package godiff

import (
	"fmt"
	"github.com/viant/parsly/matcher/option"
	"github.com/viant/parsly/splitter"
	"reflect"
	"strings"
)

//Differ represents a differ
type Differ struct {
	config  *Config
	decoder func(value interface{}) interface{}
	*structDiffer
	*mapDiffer
	*sliceDiffer
	*ifaceDiffer
}

//Diff creates change log based on comparison from and to values
func (d *Differ) Diff(from, to interface{}) *ChangeLog {
	changeLog := &ChangeLog{}
	root := &Path{}
	fieldChangeType := discoverChangeType(from, to)
	var err error

	err = d.diff(changeLog, root, from, to, fieldChangeType)
	if err != nil {
		changeLog.AddError(root, err)
	}
	return changeLog
}

func (d *Differ) diff(changeLog *ChangeLog, aPath *Path, from, to interface{}, fieldChangeType ChangeType) error {
	var err error

	if d.decoder != nil {
		if from != nil {
			from = d.decoder(from)
		}
		if to != nil {
			to = d.decoder(to)
		}
	}

	if d.structDiffer != nil {
		err = d.structDiffer.diff(changeLog, aPath, from, to, fieldChangeType)
	} else if d.sliceDiffer != nil {
		err = d.sliceDiffer.diff(changeLog, aPath, from, to, fieldChangeType)
	} else if d.ifaceDiffer != nil {
		err = d.ifaceDiffer.diff(changeLog, aPath, from, to, fieldChangeType)
	} else if d.mapDiffer != nil {
		err = d.mapDiffer.diff(changeLog, aPath, from, to, fieldChangeType)
	} else {
		if !matches(from, to) {
			switch {
			case fieldChangeType == ChangeTypeDelete && to == nil:
				changeLog.AddDelete(aPath, from)
			case fieldChangeType == ChangeTypeCreate && from == nil:
				changeLog.AddCreate(aPath, to)
			default:
				changeLog.AddUpdate(aPath, from, to)
			}
		}
	}
	return err
}

func discoverChangeType(from interface{}, to interface{}) ChangeType {
	fieldChangeType := ChangeTypeUpdate
	if from == nil {
		fieldChangeType = ChangeTypeCreate
	} else if to == nil {
		fieldChangeType = ChangeTypeDelete
	}
	return fieldChangeType
}

func (d *Differ) decodedSliceDiff() (*Differ, error) {
	var err error
	if d.sliceDiffer, err = newSliceDiffer(stringsType, stringsType, d.config, d.config.tag); err != nil {
		return nil, err
	}
	tag := d.config.tag
	itemSplitter := splitter.New(strings.Split(tag.ItemSeparator, "|"), option.NewCase(false))
	d.decoder = func(value interface{}) interface{} {
		text, _ := value.(string)
		text = tag.removeWhitespace(text)
		var ret = make([]string, 0)
		items := itemSplitter.Split(text)
		for _, item := range items {
			ret = append(ret, item)
		}
		return &ret
	}
	return d, nil
}

func (d *Differ) decodedMapDiff() (*Differ, error) {
	var err error
	tag := d.config.tag
	clonedTag := tag.clone()
	clonedTag.PairDelimiter = ""
	clonedTag.PairSeparator = ""
	if d.mapDiffer, err = newMapDiffer(stringMapType, stringMapType, d.config, clonedTag); err != nil {
		return nil, err
	}
	itemSplitter := splitter.New(strings.Split(tag.PairDelimiter, "|"), option.NewCase(false))
	d.decoder = func(value interface{}) interface{} {
		text, _ := value.(string)
		var ret = make(map[string]interface{}, 0)
		items := itemSplitter.Split(text)
		for _, item := range items {
			pair := strings.Split(item, tag.PairSeparator)
			pair[0] = tag.removeWhitespace(pair[0])
			pair[1] = tag.removeWhitespace(pair[1])
			if len(pair) == 2 {
				ret[pair[0]] = pair[1]
			}
		}
		return ret
	}
	return d, nil
}

//New creates a differ
func New(from, to reflect.Type, opts ...Option) (*Differ, error) {
	var result = &Differ{config: &Config{}}
	for _, opt := range opts {
		opt(result.config)
	}
	result.config.Init()
	if from == nil {
		from = to
	} else if to == nil {
		to = from
	}
	tag := result.config.tag
	var err error

	switch {
	case structType(from) != nil && structType(to) != nil:
		if result.structDiffer, err = newStructDiffer(from, to, result.config); err != nil {
			return nil, err
		}
		return result, nil
	case sliceType(from) != nil && sliceType(to) != nil:
		if result.sliceDiffer, err = newSliceDiffer(from, to, result.config, result.config.tag); err != nil {
			return nil, err
		}
		return result, nil
	case mapType(from) != nil && mapType(to) != nil:
		if result.mapDiffer, err = newMapDiffer(from, to, result.config, result.config.tag); err != nil {
			return nil, err
		}
		return result, nil
	case interfaceType(from) != nil || interfaceType(to) != nil:
		differ, err := newIfaceDiffer(result.config, result.config.tag)
		if err != nil {
			return nil, err
		}
		return &Differ{config: result.config, ifaceDiffer: differ}, nil
	case from == to:

		if from.Kind() == reflect.String {
			if tag != nil {
				if delimiter := tag.PairDelimiter; delimiter != "" {
					return result.decodedMapDiff()
				} else if itemSeparator := tag.ItemSeparator; itemSeparator != "" {
					return result.decodedSliceDiff()
				}
			}
		}
		return result, nil
	default:
	}
	return nil, fmt.Errorf("unsupported match types: %s, %s", from.String(), to.String())
}
