package collector

// Config represents collector configuration.
type Config struct {
	DSN            string
	DefaultService string
	MaxConcurrency int
}
