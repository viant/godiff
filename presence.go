package godiff

import (
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type PresenceProvider struct {
	Holder        *xunsafe.Field
	Fields        []*xunsafe.Field
	IdentityIndex int
}

//Has returns true if value on holder field with index has been set
func (p *PresenceProvider) Has(ptr unsafe.Pointer, index int) bool {
	hasPtr := p.Holder.ValuePointer(ptr)
	if index >= len(p.Fields) || p.Fields[index] == nil {
		return false
	}
	return p.Fields[index].Bool(hasPtr)
}

//IsFieldSet returns true if field has been set
func (p *PresenceProvider) IsFieldSet(ptr unsafe.Pointer, index int) bool {
	if p == nil || p.Holder == nil {
		return true //we do not have field presence provider so we assume all fields are set
	}
	return p.Has(ptr, index)
}

func (p *PresenceProvider) Init(filedPos map[string]int) error {
	if p.Holder == nil || len(filedPos) == 0 {
		return nil
	}
	if holder := p.Holder; holder != nil {
		p.Fields = make([]*xunsafe.Field, len(filedPos))
		holderType := holder.Type
		if holderType.Kind() == reflect.Ptr {
			holderType = holderType.Elem()
		}
		for i := 0; i < holderType.NumField(); i++ {
			presentField := holderType.Field(i)
			pos, ok := filedPos[presentField.Name]
			if !ok {
				return fmt.Errorf("failed to match presence field %v", presentField.Name)
			}
			p.Fields[pos] = xunsafe.NewField(presentField)
		}
	}

	return nil
}
