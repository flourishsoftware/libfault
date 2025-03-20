package libfault

import "fmt"

// New creates a new basic fault error.
func (c *Config) New(message string, skipFramesDelta int, w ...Wrapper) error {
	f := &fundamental{
		msg:      message,
		location: c.getLocation(skipFramesDelta),
		config:   c,
	}

	var err error = f
	for _, fn := range w {
		err = fn(err)
	}

	return err
}

// New is a package-level convenience function that creates a default Config and calls New.
func New(message string, w ...Wrapper) error {
	config := &Config{}
	return config.New(message, 0, w...)
}

// Newf includes formatting specifiers.
func (c *Config) Newf(message string, skipFramesDelta int, va ...any) error {
	f := &fundamental{
		msg:      fmt.Sprintf(message, va...),
		location: c.getLocation(skipFramesDelta),
		config:   c,
	}
	return f
}

// Newf is a package-level convenience function that creates a default Config and calls Newf.
func Newf(message string, va ...any) error {
	config := &Config{}
	return config.Newf(message, 0, va...)
}

type fundamental struct {
	msg      string
	location string
	config   *Config
}

func (f *fundamental) Error() string {
	return f.msg
}
