package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/lattesec/log"
)

const (
	DefaultFileName      = "camserver.yaml"
	DefaultConfigDirName = "camserver"
)

func loadConfigFile(path string) (*Config, error) {
	f, err := os.OpenFile(filepath.Clean(path), os.O_RDONLY, 0o600)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadConfig(customPath *string) (*Config, error) {
	if customPath != nil {
		cfg, err := loadConfigFile(*customPath)
		if err == nil {
			log.Info().WithMeta("scope", "cfg").Msgf("loaded '%s'", *customPath).Send()
			return cfg, nil
		}
		log.Error().WithMeta("scope", "cfg").Msgf("failed to load '%s': %v", *customPath, err).Send()
	}

	searchLocations := []string{"."}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error().WithMeta("scope", "cfg").Msgf("failed to get user home dir: %v", err).Send()
	} else {
		searchLocations = append(searchLocations, filepath.Join(homeDir, "."+DefaultFileName))
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Error().WithMeta("scope", "cfg").Msgf("failed to get user config dir: %v", err).Send()
	} else {
		searchLocations = append(searchLocations, filepath.Join(configDir, DefaultConfigDirName, DefaultFileName))
	}

	for _, location := range searchLocations {
		cfg, err := loadConfigFile(location)
		if err == nil {
			log.Info().WithMeta("scope", "cfg").Msgf("loaded '%s'", location).Send()
			return cfg, nil
		}
		log.Debug().WithMeta("scope", "cfg").Msgf("failed to load '%s': %v", location, err).Send()
	}

	return &Config{}, nil
}
