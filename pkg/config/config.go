package config

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Listen    string `yaml:"listen"`
	LogFormat string `yaml:"logFormat"`
	RootDir   string `yaml:"rootDir"`
	Debug     bool   `yaml:"debug"`

	SecretKeyString string `yaml:"secretKey"`
	SecretKey       []byte

	AuthDriver       string                      `yaml:"authDriver"`
	AuthStaticConfig *AuthenticationStaticConfig `yaml:"authStaticConfig"`
	AuthOIDCConfig   *AuthenticationOIDCConfig   `yaml:"authOIDCConfig"`

	MetricsDriver         string                `yaml:"metricsDriver"`
	MetricsInfluxDBConfig MetricsInfluxDBConfig `yaml:"metricsInfluxDBConfig"`
}

type AuthenticationStaticConfig struct {
	Tokens map[string]string `yaml:"tokens"`
}
type AuthenticationOIDCConfig struct {
	URL          string `yaml:"url"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
	Redirect     string `yaml:"redirect"`
	Provider     string `yaml:"provider"`
}

type MetricsInfluxDBConfig struct {
	URL    string `yaml:"url"`
	Token  string `yaml:"token"`
	Org    string `yaml:"org"`
	Bucket string `yaml:"bucket"`
}

func DefaultConfig() {
	C = Config{
		Listen:          "localhost:8000",
		LogFormat:       "plain",
		RootDir:         "./root",
		AuthDriver:      "null",
		MetricsDriver:   "prometheus",
		SecretKeyString: "",
		Debug:           false,
	}
}

func LoadConfig(path string) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "Failed to load config file")
	}
	rawExpanded := os.ExpandEnv(string(raw))
	err = yaml.Unmarshal([]byte(rawExpanded), &C)
	if err != nil {
		return errors.Wrap(err, "Failed to parse YAML")
	}
	if C.SecretKeyString == "" {
		log.Warning("No Secret Key has been set, defaulting to a random key. You should set 'secretKey' in the settings to a 32-byte, base64 encoded string to fix this.")
		C.SecretKeyString = base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	}
	// Always assume the secret key is base64 encoded, so we parse it here
	secretKey, err := base64.StdEncoding.DecodeString(C.SecretKeyString)
	if err != nil {
		log.Warning("Failed to parse Secret Key as base64, defaulting to random key.")
		C.SecretKey = securecookie.GenerateRandomKey(32)
	} else {
		C.SecretKey = secretKey
	}
	return nil
}

var C Config

func CleanURL(raw string) string {
	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}
	return filepath.Join(C.RootDir, filepath.FromSlash(path.Clean("/"+raw)))
}
