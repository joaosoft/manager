package gomanager

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// GoMockOption ...
type GoMockOption func(gomock *GoMock)

// Reconfigure ...
func (gomock *GoMock) Reconfigure(options ...GoMockOption) {
	for _, option := range options {
		option(gomock)
	}
}

// WithPath ...
func WithPath(path string) GoMockOption {
	return func(gomock *GoMock) {
		if path != "" {
			if !strings.HasSuffix(path, "/") {
				path += "/"
			}
			global["path"] = path
		}
	}
}

// WithRunInBackground ...
func WithRunInBackground(background bool) GoMockOption {
	return func(gomock *GoMock) {
		gomock.background = background
	}
}

// WithConfigurationFile ...
func WithConfigurationFile(file string) GoMockOption {
	return func(gomock *GoMock) {
		config := &Configurations{}
		if _, err := readFile(file, config); err != nil {
			panic(err)
		}
		gomock.Reconfigure(
			WithSqlConfiguration(&config.Connections.SqlConfig),
			WithRedisConfiguration(&config.Connections.RedisConfig),
			WithNsqConfiguration(&config.Connections.NsqConfig))
	}
}

// WithRedisConfiguration ...
func WithRedisConfiguration(config *RedisConfig) GoMockOption {
	return func(gomock *GoMock) {
		global["redis"] = config
	}
}

// WithSqlConfiguration ...
func WithSqlConfiguration(config *SqlConfig) GoMockOption {
	return func(gomock *GoMock) {
		global["sql"] = config
	}
}

// WithNsqConfiguration ...
func WithNsqConfiguration(config *NsqConfig) GoMockOption {
	return func(gomock *GoMock) {
		global["nsq"] = config
	}
}

// WithLogLevel ...
func WithLogLevel(level logrus.Level) GoMockOption {
	return func(gomock *GoMock) {
		logrus.SetLevel(level)
	}
}

// WithConfigurations ...
func WithConfigurations(config *Configurations) GoMockOption {
	return func(gomock *GoMock) {
		gomock.Reconfigure(
			WithSqlConfiguration(&config.Connections.SqlConfig),
			WithRedisConfiguration(&config.Connections.RedisConfig),
			WithNsqConfiguration(&config.Connections.NsqConfig))
	}
}
