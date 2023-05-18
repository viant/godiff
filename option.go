package godiff

//Option represents an option
type Option func(config *Config)

//WithTagName updated config tag
func WithTagName(name string) Option {
	return func(config *Config) {
		config.TagName = name
	}
}

//WithTag updated config tag
func WithTag(tag *Tag) Option {
	return func(config *Config) {
		config.tag = tag
	}
}

func WithPresence(f bool) Option {
	return func(config *Config) {
		config.withPresence = f
	}
}

//NullifyEmpty updated config option
func NullifyEmpty(flag bool) Option {
	return func(config *Config) {
		config.NullifyEmpty = &flag
	}
}

//WithConfig updated config tag
func WithConfig(cfg *Config) Option {
	return func(config *Config) {
		*config = *cfg
	}
}

//WithRegistry updated config with registry
func WithRegistry(registry *Registry) Option {
	return func(config *Config) {
		config.registry = registry
	}
}
