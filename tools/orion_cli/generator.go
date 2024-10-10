package main

import (
	"errors"
	"os"
	"strings"
)

type OrionGenerator struct {
	err    error
	output string
	module string
}

func (o *OrionGenerator) Check(args []string) *OrionGenerator {

	cl.Log("Checking ...")

	if len(args) != 2 || args[0] == "" || args[1] == "" {
		o.err = errors.New("need output folder and module name")
		return o
	}
	c, e := os.Stat(args[0])
	if e != nil && !strings.Contains(e.Error(), "no such file or directory") {
		o.err = errors.New("output path check error:" + e.Error())
		return o
	}
	if c != nil && c.IsDir() {
		o.err = errors.New("output path existed")
		return o
	}

	cl.Log("Create folder: " + args[0])

	if err := os.MkdirAll(args[0], os.ModePerm); err != nil {
		o.err = errors.New("output path init error:" + err.Error())
		return o
	}
	o.output = args[0]
	o.module = args[1]
	return o
}

func (o *OrionGenerator) Excute() error {

	if o.err != nil {
		return o.err
	}

	if err := o.CreateFolder(o.output+"/cmd", func() (name, content string) {
		return o.output + "/cmd/main.go", _tpl_main
	}); err != nil {
		cl.Log("excute err: " + err.Error())
		return err
	}

	if err := o.CreateFolder(o.output + "/proto"); err != nil {
		cl.Log("excute err: " + err.Error())
		return err
	}

	if err := o.CreateFolder(o.output+"/service", func() (name, content string) {
		return o.output + "/service/service.go", _tpl_service
	}); err != nil {
		cl.Log("excute err: " + err.Error())
		return err
	}

	if err := o.CreateFolder(o.output+"/config", func() (name, content string) {
		return o.output + "/config/config.go", _tpl_config
	}); err != nil {
		cl.Log("excute err: " + err.Error())
		return err
	}

	if err := o.CreateFile(o.output+"/go.mod", _tpl_gomod); err != nil {
		cl.Log("excute err: " + err.Error())
		return err
	}

	return nil
}

func (o *OrionGenerator) CreateFolder(dir string, files ...func() (name string, content string)) error {
	cl.Log("create dir: " + dir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return errors.New("create dir error:" + err.Error())
	}
	for _, gfunc := range files {
		fileName, fileContent := gfunc()
		if err := o.CreateFile(fileName, fileContent); err != nil {
			return errors.New("create file error:" + err.Error())
		}
	}
	return nil
}

func (o *OrionGenerator) CreateFile(name string, content string) error {
	cl.Log("create file: " + name)
	fs, err := os.Create(name)
	if err != nil {
		return errors.New("create dir error:" + err.Error())
	}
	defer fs.Close()
	content = strings.ReplaceAll(content, "${module}", o.module)
	fs.WriteString(content)
	return nil
}

var (
	_tpl_main = `
package main

import (
	"log"

<<<<<<< HEAD
	"github.com/orion/orpc"
	_ "github.com/orion/orpc/build"
	"github.com/orion/pkg/logger"

	// "github.com/orion/example/orion_server/proto_go/proto/todo"
=======
	"github.com/Anderson-Lu/orion/orpc"
	_ "github.com/Anderson-Lu/orion/orpc/build"
	"github.com/Anderson-Lu/orion/pkg/logger"

	// "github.com/Anderson-Lu/orion/example/orion_server/proto_go/proto/todo"
>>>>>>> dev_0_0_2
	"${module}/service"
)

func main() {

	
	c := &orpc.Config{
		Server:          &orpc.ServerConfig{Port: 8080, EnableGRPCGateway: true},
		PromtheusConfig: &orpc.PromtheusConfig{Enable: true, Port: 9092},
		FrameLogger:     &logger.LoggerConfig{Path: "../log/frame.log", LogLevel: "info"},
		AccessLogger:    &logger.LoggerConfig{Path: "../log/access.log"},
		ServiceLogger:   &logger.LoggerConfig{Path: "../log/service.log"},
		PanicLogger:     &logger.LoggerConfig{Path: "../log/panic.log"},
	}

	handler, _ := service.NewService(c)
	server, err := orpc.New(
		orpc.WithConfigFile("../config/config.toml"),
		// TODO
		// orpc.WithGRPCHandler(handler, &todo.UitTodo_ServiceDesc),
		// orpc.WithGrpcGatewayEndpointFunc(todo.RegisterUitTodoHandlerFromEndpoint),
		orpc.WithFlags(),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}	
`

	_tpl_service = `
package service

import (
	"${module}/config"
<<<<<<< HEAD
	"github.com/orion/pkg/logger"
=======
	"github.com/Anderson-Lu/orion/pkg/logger"
>>>>>>> dev_0_0_2
)

func NewService(c *config.Config) (*Service, error) {

	lg, err := logger.NewLogger(c.ServiceLogger)
	if err != nil {
		return nil, err
	}

	return &Service{
		c: c,
		l: lg,
	}, nil
}

type Service struct {
	c *config.Config
	l *logger.Logger
}
`

	_tpl_config = `
package config

<<<<<<< HEAD
import "github.com/orion/orpc"
=======
import "github.com/Anderson-Lu/orion/orpc"
>>>>>>> dev_0_0_2

type Config struct {
	orpc.Config
	// your own configs
}
`

	_tpl_gomod = `
module ${module}

go 1.21

toolchain go1.22.1
`
)
