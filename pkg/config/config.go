package config

import (
	"io/ioutil"

	"github.com/BeryJu/gopyazo/pkg/drivers/auth"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	RootDir                     string            `yaml:"root_dir"`
	APIPathPrefix               string            `yaml:"api_path_prefix"`
	AuthenticationDriver        string            `yaml:"authentication_driver"`
	AuthenticationDriverContext map[string]string `yaml:"authentication_driver_context"`

	authDriver auth.AuthDriver
}

func Load(path string) *Configuration {
	c := &Configuration{}
	yamlFile, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(yamlFile), &c)
	if err != nil {
		panic(err)
	}
	c.GetAuth()
	return c
}

func (c *Configuration) GetAuth() auth.AuthDriver {
	if c.authDriver != nil {
		return c.authDriver
	}
	switch c.AuthenticationDriver {
	case "static":
		c.authDriver = &auth.StaticAuth{}
		c.authDriver.Init(c.AuthenticationDriverContext)
	}
	return c.authDriver
}

var Config *Configuration
