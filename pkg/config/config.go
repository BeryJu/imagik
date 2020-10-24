package config

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Listen           string                      `yaml:"listen"`
	LogFormat        string                      `yaml:"logFormat"`
	RootDir          string                      `yaml:"rootDir"`
	SecretKey        string                      `yaml:"secretKey"`
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
		SecretKey:  "",
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
	if C.SecretKey == "" {
		log.Warning("No CSRF Key has been set, defaulting to a random key. You should set 'secretKey' in the settings to a 32-byte, base64 encoded string to fix this.")
		C.SecretKey = base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	}
	return nil
}

var C Config
