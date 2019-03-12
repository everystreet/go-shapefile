package shapefile

import (
	"github.com/mercatormaps/go-shapefile/cpg"
)

type Option func(*Config)

func CharacterEncoding(enc cpg.CharacterEncoding) Option {
	return func(c *Config) {
		c.charEnc = enc
	}
}

func FilterFields(names ...string) Option {
	return func(c *Config) {
		c.fields = names
	}
}

type Config struct {
	charEnc cpg.CharacterEncoding
	fields  []string
}

func (c *Config) CharacterEncoding() cpg.CharacterEncoding {
	return c.charEnc
}

func (c *Config) FilteredFields() []string {
	return c.fields
}

func defaultConfig() *Config {
	return &Config{}
}
