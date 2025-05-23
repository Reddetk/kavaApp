package config

import (
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// LoadConfig reads the YAML file and unmarshals it into the Config struct.
// It also loads database configuration from environment variables if present.
func LoadConfig(path string) (*Config, error) {
	// load vars in .env
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Expand environment variables in the DSN string from the YAML file.
	cfg.Database.DSN = os.ExpandEnv(cfg.Database.DSN)
	return &cfg, nil
}

// Config holds all configuration sections.
type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Logger      LoggerConfig      `yaml:"logger"`
	Database    DatabaseConfig    `yaml:"database"`
	Apriori     AprioriConfig     `yaml:"apriori"`
	ABCAnalysis ABCAnalysisConfig `yaml:"abc_analysis"`
}

// ServerConfig holds the server-related settings.
type ServerConfig struct {
	Address             string `yaml:"address"`
	Port                string `yaml:"port"`
	ReadTimeoutSeconds  int    `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `yaml:"write_timeout_seconds"`
	IdleTimeoutSeconds  int    `yaml:"idle_timeout_seconds"`
}

// LoggerConfig holds logger settings.
type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// DatabaseConfig holds the DB connection parameters.
type DatabaseConfig struct {
	DSN                    string `yaml:"dsn"`
	MaxOpenConns           int    `yaml:"max_open_conns"`
	MaxIdleConns           int    `yaml:"max_idle_conns"`
	ConnMaxLifetimeMinutes int    `yaml:"conn_max_lifetime_minutes"`
}

// AprioriConfig holds settings for the Apriori algorithm.
type AprioriConfig struct {
	DefaultMinSupport    float64 `yaml:"default_min_support"`
	DefaultMinConfidence float64 `yaml:"default_min_confidence"`
	MaxRecommendations   int     `yaml:"max_recommendations"`
}

// ABCAnalysisConfig holds settings for ABC analysis.
type ABCAnalysisConfig struct {
	AThreshold float64 `yaml:"a_threshold"`
	BThreshold float64 `yaml:"b_threshold"`
}
