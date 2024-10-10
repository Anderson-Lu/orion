package main

import (
	"errors"
	"os"
	"strings"
)

type OrionGenerator struct {
	err    error
	output string
}

func (o *OrionGenerator) Check(args []string) *OrionGenerator {

	cl.Log("Checking ...")

	if len(args) == 0 {
		o.err = errors.New("output path not specified")
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
	return o
}

func (o *OrionGenerator) Excute() error {

	if o.err != nil {
		return o.err
	}

	if err := o.CreateFolder(o.output+"/cmd", func() (name, content string) {
		return o.output + "/cmd/main.go", ""
	}); err != nil {
		cl.Log("excute err: " + err.Error())
		return err
	}

	if err := o.CreateFolder(o.output + "/proto"); err != nil {
		cl.Log("excute err: " + err.Error())
		return err
	}

	if err := o.CreateFolder(o.output + "/service"); err != nil {
		cl.Log("excute err: " + err.Error())
		return err
	}

	if err := o.CreateFolder(o.output + "/config"); err != nil {
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
	fs.WriteString(content)
	return nil
}
