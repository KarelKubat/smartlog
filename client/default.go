package client

var DefaultClient *Client

func Info(msg string) error {
	return DefaultClient.Info(msg)
}

func Infof(format string, args ...interface{}) error {
	return DefaultClient.Infof(format, args...)
}

func Warn(msg string) error {
	return DefaultClient.Warn(msg)
}

func Warnf(format string, args ...interface{}) error {
	return DefaultClient.Warnf(format, args...)
}

func Fatal(msg string) error {
	return DefaultClient.Fatal(msg)
}

func Fatalf(format string, args ...interface{}) error {
	return DefaultClient.Fatalf(format, args...)
}
