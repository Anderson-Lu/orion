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
	Registry        *Registry            `yaml:"Registry" json:"Registry" toml:"Registry"`
	PromtheusConfig *PromtheusConfig     `yaml:"PromtheusConfig" json:"PromtheusConfig" toml:"PromtheusConfig"`
	FrameLogger     *logger.LoggerConfig `yaml:"FrameLogger" json:"FrameLogger" toml:"FrameLogger"`
	AccessLogger    *logger.LoggerConfig `yaml:"AccessLogger" json:"AccessLogger" toml:"AccessLogger"`
	ServiceLogger   *logger.LoggerConfig `yaml:"ServiceLogger" json:"ServiceLogger" toml:"ServiceLogger"`
	PanicLogger     *logger.LoggerConfig `yaml:"PanicLogger" json:"PanicLogger" toml:"PanicLogger"`
	RateLimit       []*ratelimit.Config  `yaml:"RateLimit" json:"RateLimit" toml:"RateLimit"`
	Tracing         *TracingConfig       `yaml:"Tracing" json:"Tracing" toml:"Tracing"`
}

type Registry struct {
	Service string `yaml:"Service" json:"Service" toml:"Service"`
	IP      string `yaml:"IP" json:"IP" toml:"IP"`
	Port    uint32 `yaml:"Port" json:"Port" toml:"Port"`
}

type TracingConfig struct {
	Address     string `yaml:"Address" json:"Address" toml:"Address"`
	Namespace   string `yaml:"Namespace" json:"Namespace" toml:"Namespace"`
	ServiceName string `yaml:"ServiceName" json:"ServiceName" toml:"ServiceName"`
	InstanceId  string `yaml:"InstanceId" json:"InstanceId" toml:"InstanceId"`
	Env         string `yaml:"Env" json:"Env" toml:"Env"`
}
