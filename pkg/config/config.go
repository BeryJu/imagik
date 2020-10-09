package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Listen           string                      `yaml:"listen"`
	LogFormat        string                      `yaml:"logFormat"`
	RootDir          string                      `yaml:"rootDir"`
	AuthDriver       string                      `yaml:"authDriver"`
	AuthStaticConfig *AuthenticationStaticConfig `yaml:"authStaticConfig"`
	AuthOIDCConfig   *AuthenticationOIDCConfig   `yaml:"authOIDCConfig"`
}

type AuthenticationStaticConfig struct {
	Tokens map[string]string `yaml:"tokens"`
}
type AuthenticationOIDCConfig struct {
	URL          string `yaml:"url"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
	Redirect     string `yaml:"redirect"`
}

func DefaultConfig() {
	C = Config{
		Listen:     "localhost:8000",
		LogFormat:  "plain",
		RootDir:    "./root",
		AuthDriver: "null",
	}
}

func LoadConfig(path string) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "Failed to load config file")
	}
	err = yaml.Unmarshal(raw, &C)
	if err != nil {
		return errors.Wrap(err, "Failed to parse YAML")
	}
	fmt.Printf("%+v\n", &C)
	return nil
}

var C Config
