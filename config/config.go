package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
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
		err = fmt.Errorf("ioutil.ReadFile(...): %w", err)
		return nil, err
	}

	cfg = &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		err = fmt.Errorf("yaml.Unmarshal(...): %w", err)
		return nil, err
	}

	// Override with environment variables
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}
	if key := os.Getenv("SERVER_AUTHORIZATION_KEY"); key != "" {
		cfg.Server.AuthorizationKey = key
	}

	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		cfg.Postgres.Host = host
	}
	if portStr := os.Getenv("POSTGRES_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Postgres.Port = port
		}
	}
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		cfg.Postgres.User = user
	}
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		cfg.Postgres.Password = password
	}
	if dbName := os.Getenv("POSTGRES_DATABASE_NAME"); dbName != "" {
		cfg.Postgres.DatabaseName = dbName
	}
	if sslMode := os.Getenv("POSTGRES_SSL_MODE"); sslMode != "" {
		cfg.Postgres.SslMode = sslMode
	}

	if url := os.Getenv("FACE_SCAN_API_URL"); url != "" {
		cfg.FaceScanAPI.URL = url
	}
	if authKey := os.Getenv("FACE_SCAN_API_AUTHORIZATION_KEY"); authKey != "" {
		cfg.FaceScanAPI.Authorization = authKey
	}
	if mimeType := os.Getenv("FACE_SCAN_API_MIME_TYPE"); mimeType != "" {
		cfg.FaceScanAPI.MimeType = mimeType
	}

	return cfg, nil
}
