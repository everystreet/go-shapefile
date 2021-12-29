package dbf

import "golang.org/x/text/encoding"

// Option funcs can be applied to reading operations.
type Option func(*config)

// CharacterDecoder sets the encoding of character field values.
// By default, ASCII is assumed.
func CharacterDecoder(dec *encoding.Decoder) Option {
	return func(c *config) {
		c.decoder = dec
	}
}

// FilterFields allows filtering by field name.
// If this option is used, only these fields will be returned in the Record.
// Without this option, all available fields are returned.
func FilterFields(names ...string) Option {
	return func(c *config) {
		c.fields = names
	}
}

// Config for dbf parsing.
type config struct {
	decoder *encoding.Decoder
	fields  []string
}

func defaultConfig() config {
	return config{
		decoder: encoding.Nop.NewDecoder(),
	}
}
