package orpc

import (
	"github.com/Anderson-Lu/orion/pkg/logger"
	"github.com/Anderson-Lu/orion/pkg/ratelimit"
)

type ServerConfig struct {
	Port              uint32 `default:"8081" yaml:"Port" json:"Port" toml:"Port"`
	EnableGRPCGateway bool   `default:"true" yaml:"EnableGRPCGateway" json:"EnableGRPCGateway" toml:"EnableGRPCGateway"`
}

type PromtheusConfig struct {
	Enable bool   `default:"true" yaml:"Enable" json:"Enable" toml:"Enable"`
	Port   uint32 `default:"8082" yaml:"Port" json:"Port" toml:"Port"`
}

type Config struct {
	Server          *ServerConfig        `yaml:"Server" json:"Server" toml:"Server"`
	PromtheusConfig *PromtheusConfig     `yaml:"PromtheusConfig" json:"PromtheusConfig" toml:"PromtheusConfig"`
	FrameLogger     *logger.LoggerConfig `yaml:"FrameLogger" json:"FrameLogger" toml:"FrameLogger"`
	AccessLogger    *logger.LoggerConfig `yaml:"AccessLogger" json:"AccessLogger" toml:"AccessLogger"`
	ServiceLogger   *logger.LoggerConfig `yaml:"ServiceLogger" json:"ServiceLogger" toml:"ServiceLogger"`
	PanicLogger     *logger.LoggerConfig `yaml:"PanicLogger" json:"PanicLogger" toml:"PanicLogger"`
	RateLimit       []*ratelimit.Config  `yaml:"RateLimit" json:"RateLimit" toml:"RateLimit"`
}
