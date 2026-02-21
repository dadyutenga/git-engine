package cli

import (
	"os"

	"github.com/dadyutenga/git-engine/internal/infrastructure/ssh"
	"gopkg.in/yaml.v3"
)

// Config represents the CLI configuration.
type Config struct {
	SSH ssh.Config `yaml:"ssh"`
}

// LoadConfig reads YAML configuration from path.
func LoadConfig(path string) (Config, error) {
	cfg := Config{}
	content, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return cfg, err
	}
	if cfg.SSH.Port == 0 {
		cfg.SSH.Port = 22
	}
	return cfg, nil
}
