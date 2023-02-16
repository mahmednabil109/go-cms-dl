package main

import (
	"bytes"
	"io"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Config struct holds general configuration for the downloader
type Config struct {
	Goroutines  int    `yaml:"routines_count"`
	DownlaodDir string `yaml:"download_dir"`
	CMS_URL     string `yaml:"cms_url"`
	DotFilePath string `yaml:"dotfile_path"`
}

// Parse reads the config file `config.yml` near the executable
func Parse(log *zap.Logger) (Config, error) {
	file, err := os.Open("./config.yml")
	if err != nil {
		log.Panic("failed to open config file", zap.Error(err))
	}
	defer file.Close()

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	if err != nil {
		log.Panic("failed to read config file content", zap.Error(err))
	}

	var cfg Config
	err = yaml.Unmarshal(buffer.Bytes(), &cfg)
	if err != nil {
		log.Panic("faild to parse config file", zap.Error(err))
	}
	return cfg, nil
}
