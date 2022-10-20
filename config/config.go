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
	ListenAddr   string `json:"listen_addr"`
	ListenPort   uint16 `json:"listen_port"`
	UseSSL       bool   `json:"use_ssl"`
	KeyFilePath  string `json:"key_file_path"`
	CertFilePath string `json:"cert_file_path"`
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
	}
}
