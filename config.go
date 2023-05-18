package godiff

import "reflect"

//Config represents a config
type Config struct {
	TimeLayout   string
	NullifyEmpty *bool
	TagName      string //diff by default
	StrictMode   bool   //non-strict mode allows string with non-string matches
	tag          *Tag
	registry     *Registry
	withPresence bool
}

//Init init config
func (c *Config) Init() {
	if c.TagName == "" {
		c.TagName = "diff"
	}
	if c.registry == nil {
		c.registry = &Registry{differs: map[reflect.Type]map[reflect.Type]*Differ{}}
	}
}
