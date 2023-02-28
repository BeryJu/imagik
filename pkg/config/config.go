package config

import (
	"encoding/base64"
	"os"

	"github.com/gorilla/securecookie"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Listen    string `yaml:"listen"`
	LogFormat string `yaml:"logFormat"`
	Debug     bool   `yaml:"debug"`

	SecretKeyString string `yaml:"secretKey"`
	SecretKey       []byte

	AuthDriver       string                      `yaml:"authDriver"`
	AuthStaticConfig *AuthenticationStaticConfig `yaml:"authStaticConfig"`
	AuthOIDCConfig   *AuthenticationOIDCConfig   `yaml:"authOIDCConfig"`

	MetricsDriver         string                `yaml:"metricsDriver"`
	MetricsInfluxDBConfig MetricsInfluxDBConfig `yaml:"metricsInfluxDBConfig"`

	StorageDriver      string              `yaml:"storageDriver"`
	StorageLocalConfig *StorageLocalConfig `yaml:"storageLocalConfig"`
	StorageS3Config    *StorageS3Config    `yaml:"storageS3Config"`
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

type StorageLocalConfig struct {
	Root string `yaml:"root"`
}
type StorageS3Config struct {
	Bucket     string `yaml:"bucket"`
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"accessKey"`
	SecretKey  string `yaml:"secretKey"`
	Region     string `yaml:"region"`
	UsePresign bool   `yaml:"usePresign"`
}

type MetricsInfluxDBConfig struct {
	URL    string `yaml:"url"`
	Token  string `yaml:"token"`
	Org    string `yaml:"org"`
	Bucket string `yaml:"bucket"`
}

func DefaultConfig() {
	C = Config{
		Listen:        "localhost:8000",
		LogFormat:     "plain",
		StorageDriver: "local",
		StorageLocalConfig: &StorageLocalConfig{
			Root: "./root",
		},
		StorageS3Config: &StorageS3Config{
			UsePresign: true,
		},
		AuthDriver:      "null",
		MetricsDriver:   "prometheus",
		SecretKeyString: "",
		Debug:           false,
	}
}

func LoadConfig(path string) error {
	raw, err := os.ReadFile(path)
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
