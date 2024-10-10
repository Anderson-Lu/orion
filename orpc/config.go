package orpc

import "github.com/Anderson-Lu/orion/pkg/logger"

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
	RateLimit       []*RateLimitConfig   `yaml:"RateLimit" json:"RateLimit" toml:"RateLimit"`
}

type RateLimitConfig struct {
	Method        string `yaml:"Method" json:"Method" toml:"Method"`
	Capacity      int    `yaml:"Capacity" json:"Capacity" toml:"Capacity"`
	RatePerSecond int    `yaml:"RatePerSecond" json:"RatePerSecond" toml:"RatePerSecond"`
}
