package xgrpc

import "github.com/uit/pkg/logger"

type GRPCConfig struct {
	Enable       bool   `default:"true"`
	Port         uint32 `default:"8081"`
	WithInSecure bool   `default:"true"`
}

type HTTPConfig struct {
	Enable bool   `default:"true"`
	Port   uint32 `default:"8080"`
}

type Config struct {
	GRPC          *GRPCConfig
	HTTP          *HTTPConfig
	FrameLogger   *logger.LoggerConfig
	AccessLogger  *logger.LoggerConfig
	ServiceLogger *logger.LoggerConfig
}
