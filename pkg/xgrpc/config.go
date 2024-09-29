package xgrpc

import "github.com/uit/pkg/logger"

type ServerConfig struct {
	Port              uint32 `default:"8081"`
	EnableGRPCGateway bool   `default:"true"`
}

type PromtheusConfig struct {
	Enable bool   `default:"true"`
	Port   uint32 `default:"8082"`
}

type Config struct {
	Server          *ServerConfig
	PromtheusConfig *PromtheusConfig
	FrameLogger     *logger.LoggerConfig
	AccessLogger    *logger.LoggerConfig
	ServiceLogger   *logger.LoggerConfig
	PanicLogger     *logger.LoggerConfig
}
