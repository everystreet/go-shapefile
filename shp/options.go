package shp

// Option funcs can be passed to NewScanner().
type Option func(*Config)

// PointPrecision sets the precision of coordinates.
func PointPrecision(p uint) Option {
	return func(c *Config) {
		c.precision = &p
	}
}

// Config for shp parsing.
type Config struct {
	precision *uint
}

func defaultConfig() *Config {
	return &Config{}
}
