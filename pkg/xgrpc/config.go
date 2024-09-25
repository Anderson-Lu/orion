package xgrpc

import "github.com/uit/pkg/logger"

type GRPCConfig struct {
	Enable bool   `default:"true"`
	Port   uint32 `default:"8081"`
}

type HTTPConfig struct {
	Enable bool   `default:"true"`
	Port   uint32 `default:"8080"`
}

type PProfConfig struct {
	Enable bool   `default:"true"`
	Port   uint32 `default:"8082"`
}

type Config struct {
	GRPC          *GRPCConfig
	HTTP          *HTTPConfig
	PProf         *PProfConfig
	FrameLogger   *logger.LoggerConfig
	AccessLogger  *logger.LoggerConfig
	ServiceLogger *logger.LoggerConfig
	PanicLogger   *logger.LoggerConfig
}
