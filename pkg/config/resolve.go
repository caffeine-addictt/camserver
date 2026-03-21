package config

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/lattesec/log"
)

const (
	DefaultFileName      = "camserver"
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

// LoadConfig looks at the following paths and attempts
// to locate the first usable config file
//
//   - `-C` flag
//   - $pwd
//   - $HOME/.camserver.yaml
//   - $XDG_CONFIG_HOME/camserver/camserver.yaml
//
// Returns cfg, path used and any errors
func LoadConfig(customPath string) (*Config, string, error) {
	if customPath != "" {
		cfg, err := loadConfigFile(customPath)
		if err == nil {
			log.Info().WithMeta("scope", "cfg").Msgf("loaded '%s'", customPath).Send()
			return cfg, customPath, nil
		}
		log.Error().WithMeta("scope", "cfg").Msgf("failed to load '%s': %v", customPath, err).Send()
	}

	searchLocations := []string{}

	pwd, err := os.Getwd()
	if err != nil {
		log.Error().WithMeta("scope", "cfg").Msgf("failed to get pwd: %v", err).Send()
	} else {
		searchLocations = append(searchLocations, pwd)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error().WithMeta("scope", "cfg").Msgf("failed to get user home dir: %v", err).Send()
	} else {
		searchLocations = append(searchLocations, homeDir)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Error().WithMeta("scope", "cfg").Msgf("failed to get user config dir: %v", err).Send()
	} else {
		searchLocations = append(searchLocations, filepath.Join(configDir, DefaultConfigDirName))
	}

	for _, dir := range searchLocations {
		log.Debug().WithMeta("scope", "cfg").Msgf("trying '%s'", dir).Send()
		errs := make(map[string]error, 0)

		for _, possibleFilename := range []string{
			"." + DefaultFileName + ".yml",
			"." + DefaultFileName + ".yaml",
			DefaultFileName + ".yml",
			DefaultFileName + ".yaml",
		} {
			location := filepath.Join(dir, possibleFilename)

			cfg, err := loadConfigFile(location)
			if err == nil {
				log.Info().WithMeta("scope", "cfg").Msgf("loaded '%s'", location).Send()
				return cfg, location, nil
			}
			if !errors.Is(err, os.ErrNotExist) {
				errs[location] = err
			}
		}

		if len(errs) > 0 {
			log.Warn().WithMeta("scope", "cfg").Msgf("failed to load: %v", errs).Send()
		}
	}

	return nil, "", errors.New("could not load")
}
