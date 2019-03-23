package dbf

import (
	"github.com/mercatormaps/go-shapefile/cpg"
)

// Option funcs can be parsed to Scanner.Scan().
type Option func(*Config)

// CharacterEncoding sets the encoding of character field values.
// By default, ASCII is assumed.
func CharacterEncoding(enc cpg.CharacterEncoding) Option {
	return func(c *Config) {
		c.charEnc = enc
	}
}

// FilterFields allows filtering by field name.
// If this option is used, only these fields will be returned in the Record.
// Without this option, all available fields are returned.
func FilterFields(names ...string) Option {
	return func(c *Config) {
		c.fields = names
	}
}

// Config for dbf parsing.
type Config struct {
	charEnc cpg.CharacterEncoding
	fields  []string
}

// CharacterEncoding returns the configured encoding.
func (c *Config) CharacterEncoding() cpg.CharacterEncoding {
	return c.charEnc
}

// FilteredFields returns the configured field names.
func (c *Config) FilteredFields() []string {
	return c.fields
}

func defaultConfig() *Config {
	return &Config{
		charEnc: cpg.EncodingASCII,
	}
}
