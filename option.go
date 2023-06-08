package godiff

//ConfigOption represents an option
type ConfigOption func(config *Config)

type Options struct {
	setMarker    bool
	shallow      bool
	nullifyEmpty *bool
	depth        int
}

func (o *Options) decDepth() {
	o.depth--
}
func (o *Options) Apply(options []Option) {
	for _, item := range options {
		item(o)
	}
}

type Option func(c *Options)

//WithTagName updated config tag
func WithTagName(name string) ConfigOption {
	return func(config *Config) {
		config.TagName = name
	}
}

//WithTag updated config tag
func WithTag(tag *Tag) ConfigOption {
	return func(config *Config) {
		config.tag = tag
	}
}

func WithSetMarker(f bool) Option {
	return func(options *Options) {
		options.setMarker = f
	}
}

func WithShallow(f bool) Option {
	return func(options *Options) {
		options.shallow = f
	}
}

//NullifyEmpty updated config option
func NullifyEmpty(flag bool) ConfigOption {
	return func(options *Config) {
		options.NullifyEmpty = &flag
	}
}

//WithConfig updated config tag
func WithConfig(cfg *Config) ConfigOption {
	return func(config *Config) {
		*config = *cfg
	}
}

//WithRegistry updated config with registry
func WithRegistry(registry *Registry) ConfigOption {
	return func(config *Config) {
		config.registry = registry
	}
}
