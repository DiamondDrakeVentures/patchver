package common

// Logger is a simple interface for logging providers.
type Logger interface {
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)

	Prefix() string
	SetPrefix(prefix string)
}
