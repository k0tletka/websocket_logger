package config

import (
	"errors"
	"github.com/BurntSushi/toml"
	"os"
)

var (
	ErrFileNotExists      = errors.New("invalid file, file must exist")
	ErrInvalidHistorySize = errors.New("invalid history size provided. It must be above zero")
	ErrSSLNoKeyOrCert     = errors.New("use_ssl enabled, but key or cert are absent")
)

type RootConfig struct {
	LogLocation string `toml:"log_location"`
	HistorySize int    `toml:"history_size"`

	HTTPConfig HTTPServerConfiguration `toml:"http"`
}

type HTTPServerConfiguration struct {
	ListenAddr     string          `toml:"listen_addr"`
	ListenPort     uint16          `toml:"listen_port"`
	UseSSL         bool            `toml:"use_ssl"`
	KeyFilePath    string          `toml:"key_file_path"`
	CertFilePath   string          `toml:"cert_file_path"`
	BasicAuthUsers []BasicAuthUser `toml:"ba_user"`
}

type BasicAuthUser struct {
	Name       string `toml:"name"`
	Base64Hash string `toml:"hash"`
}

func (r *RootConfig) validateFields() error {
	if _, err := os.Stat(r.LogLocation); os.IsNotExist(err) {
		return ErrFileNotExists
	}

	if r.HistorySize < 0 {
		return ErrInvalidHistorySize
	}

	if r.HTTPConfig.UseSSL && (r.HTTPConfig.KeyFilePath == "" || r.HTTPConfig.CertFilePath == "") {
		return ErrSSLNoKeyOrCert
	}

	return nil
}

func GetConfiguration() (*RootConfig, error) {
	config := getDefaultConfiguration()

	configLocation, ok := os.LookupEnv("CONFIGFILE")
	if !ok {
		configLocation = "config.toml"
	}

	_, err := toml.DecodeFile(configLocation, config)
	if err != nil {
		return nil, err
	}

	return config, config.validateFields()
}

func getDefaultConfiguration() *RootConfig {
	return &RootConfig{
		LogLocation: "/var/log/messages.log",
		HistorySize: 1000,
		HTTPConfig: HTTPServerConfiguration{
			ListenAddr: "127.0.0.1",
			ListenPort: 80,
			UseSSL:     false,
		},
	}
}
