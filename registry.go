package godiff

import (
	"reflect"
	"sync"
)

//Registry represents differ registry
type Registry struct {
	sync.RWMutex
	differs map[reflect.Type]map[reflect.Type]*Differ
}

func (r *Registry) Get(from, to reflect.Type, tag *Tag, options ...Option) (*Differ, error) {
	if tag != nil && (tag.PairSeparator != "" || tag.ItemSeparator != "") {
		return New(from, to, WithRegistry(r), WithTag(tag))
	}
	fromDiffers := r.getFromDiffers(from)
	r.RWMutex.RLock()
	differ, ok := fromDiffers[to]
	r.RWMutex.RUnlock()
	if ok {
		return differ, nil
	}
	var err error
	options = append(options, WithRegistry(r), WithTag(tag))
	if differ, err = New(from, to, options...); err != nil {
		return nil, err
	}
	r.RWMutex.Lock()
	fromDiffers[to] = differ
	r.RWMutex.Unlock()
	return differ, nil
}

func (r *Registry) getFromDiffers(from reflect.Type) map[reflect.Type]*Differ {
	r.RWMutex.RLock()
	fromDiffers, ok := r.differs[from]
	r.RWMutex.RUnlock()
	if ok {
		return fromDiffers
	}
	r.RWMutex.Lock()
	fromDiffers = make(map[reflect.Type]*Differ)
	r.differs[from] = make(map[reflect.Type]*Differ)
	r.RWMutex.Unlock()
	return fromDiffers
}

func NewRegistry() *Registry {
	return &Registry{differs: map[reflect.Type]map[reflect.Type]*Differ{}}
}
