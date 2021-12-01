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

func Error(msg string) error {
	return DefaultClient.Error(msg)
}

func Errorf(format string, args ...interface{}) error {
	return DefaultClient.Errorf(format, args...)
}
