package commonflags

// Log contains options for Logrus
type Log struct {
	// LogLevel specifies logrus logging level
	Level string `long:"log-level" env:"LOG_LEVEL" default:"info" choice:"debug" choice:"info" choice:"notice" choice:"warning" choice:"error" choice:"critical" description:"Log level"`

	// LogFormat specifies logrus logging format
	Format string `long:"log-format" env:"LOG_FORMAT" default:"text" choice:"text" choice:"json" description:"Log format"`
}
