package shp

// Option funcs can be passed to reading operations.
type Option func(*config)

// PointPrecision sets the precision of coordinates.
func PointPrecision(p uint) Option {
	return func(c *config) {
		c.precision = &p
	}
}

// Config for shp parsing.
type config struct {
	precision *uint
}

func defaultConfig() config {
	return config{}
}
