package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Server      Server      `yaml:"server"`
	Postgres    Postgres    `yaml:"postgres"`
	FaceScanAPI FaceScanAPI `yaml:"face_scan_api"`
}

type Server struct {
	Host             string `yaml:"server_host"`
	Port             string `yaml:"server_port"`
	AuthorizationKey string `yaml:"server_authorization_key"`
}

type Postgres struct {
	Host         string `yaml:"postgres_host"`
	Port         int    `yaml:"postgres_port"`
	User         string `yaml:"postgres_user"`
	Password     string `yaml:"postgres_password"`
	DatabaseName string `yaml:"postgres_database_name"`
	SslMode      string `yaml:"postgres_ssl_mode"`
}
type FaceScanAPI struct {
	URL           string `yaml:"url"`
	Authorization string `yaml:"face_scan_api_authorization_key"`
	MimeType      string `yaml:"face_scan_api_mime_type"`
}

func LoadConfig() (cfg *Config, err error) {
	data, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		err = fmt.Errorf("os.Open(...): %w", err)
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		err = fmt.Errorf("yaml.Unmarshal(...): %w", err)
		return nil, err
	}
	return cfg, nil
}
