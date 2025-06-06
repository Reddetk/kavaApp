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
		return nil, err
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
	cfg.Kafka.GroupID = os.ExpandEnv(cfg.Kafka.GroupID)
	for i := range cfg.Kafka.Brokers {
		cfg.Kafka.Brokers[i] = os.ExpandEnv(cfg.Kafka.Brokers[i])
	}
	cfg.Kafka.Topic = os.ExpandEnv(cfg.Kafka.Topic)
	return &cfg, nil
}

// Config holds all configuration sections.
type Config struct {
	Server              ServerConfig              `yaml:"server"`
	Logger              LoggerConfig              `yaml:"logger"`
	Database            DatabaseConfig            `yaml:"database"`
	Kafka               KafkaConfig               `yaml:"kafka"`
	DataPreprocessing   DataPreprocessingConfig   `yaml:"data_preprocessing"`
	Segmentation        SegmentationConfig        `yaml:"segmentation"`
	SurvivalAnalysis    SurvivalAnalysisConfig    `yaml:"survival_analysis"`
	StateTransition     StateTransitionConfig     `yaml:"state_transition"`
	RetentionPrediction RetentionPredictionConfig `yaml:"retention_prediction"`
	ClvUpdater          ClvUpdaterConfig          `yaml:"clv_updater"`
	ReportingAPI        ReportingAPIConfig        `yaml:"reporting_api"`
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

// KafkaConfig holds Kafka connection parameters
type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
	GroupID string   `yaml:"group_id"`
}

// DataPreprocessingConfig holds settings specific to data preprocessing.
type DataPreprocessingConfig struct {
	DataCleaner      DataCleanerConfig      `yaml:"data_cleaner"`
	MetricCalculator MetricCalculatorConfig `yaml:"metric_calculator"`
}

// DataCleanerConfig holds settings for handling outliers and missing values.
type DataCleanerConfig struct {
	OutlierMethod         string `yaml:"outlier_method"`
	MissingValuesStrategy string `yaml:"missing_values_strategy"`
}

// MetricCalculatorConfig holds metrics calculation parameters.
type MetricCalculatorConfig struct {
	TBP TBPConfig `yaml:"tbp"`
	RFM RFMConfig `yaml:"rfm"`
}

// TBPConfig holds settings for calculating Time Between Purchases.
type TBPConfig struct {
	WindowDays int `yaml:"window_days"`
}

// RFMConfig holds weights for Recency, Frequency, and Monetary calculations.
type RFMConfig struct {
	RecencyWeight   float64 `yaml:"recency_weight"`
	FrequencyWeight float64 `yaml:"frequency_weight"`
	MonetaryWeight  float64 `yaml:"monetary_weight"`
}

// SegmentationConfig holds settings for segmentation strategies.
type SegmentationConfig struct {
	RFMClustering      ClusteringConfig `yaml:"rfm_clustering"`
	BehaviorClustering ClusteringConfig `yaml:"behavior_clustering"`
	BatchSize          int              `yaml:"segmentation_batch_size"`
}

// ClusteringConfig can be reused for different clustering methods.
// Some fields may be used only for particular types.
type ClusteringConfig struct {
	Algorithm     string  `yaml:"algorithm"`
	Clusters      int     `yaml:"clusters"`       // For RFM clustering (e.g., KMeans)
	MaxIterations int     `yaml:"max_iterations"` // For RFM clustering
	RandomSeed    int     `yaml:"random_seed"`    // For RFM clustering
	EPS           float64 `yaml:"eps"`            // For behavior clustering (e.g., DBSCAN)
	MinSamples    int     `yaml:"min_samples"`    // For behavior clustering (e.g., DBSCAN)
}

// SurvivalAnalysisConfig holds settings for survival models.
type SurvivalAnalysisConfig struct {
	SurvivalModel SurvivalModelConfig `yaml:"survival_model"`
}

// SurvivalModelConfig holds the type of model and its parameters.
type SurvivalModelConfig struct {
	Type       string                 `yaml:"type"`
	Parameters map[string]interface{} `yaml:"parameters"`
}

// StateTransitionConfig holds settings for state transition modeling.
type StateTransitionConfig struct {
	TransitionMatrix TransitionMatrixConfig `yaml:"transition_matrix"`
}

// TransitionMatrixConfig holds parameters for generating a transition matrix.
type TransitionMatrixConfig struct {
	Smoothing float64 `yaml:"smoothing"`
	Threshold float64 `yaml:"threshold"`
}

// RetentionPredictionConfig holds settings for retention prediction.
type RetentionPredictionConfig struct {
	RetentionCalculator RetentionCalculatorConfig `yaml:"retention_calculator"`
}

// RetentionCalculatorConfig holds parameters for the retention calculator.
type RetentionCalculatorConfig struct {
	Model     string  `yaml:"model"`
	Threshold float64 `yaml:"threshold"`
}

// ClvUpdaterConfig holds settings for the CLV updater.
type ClvUpdaterConfig struct {
	ClvCalculator ClvCalculatorConfig `yaml:"clv_calculator"`
}

// ClvCalculatorConfig holds parameters for calculating CLV.
type ClvCalculatorConfig struct {
	DiscountRate   float64 `yaml:"discount_rate"`
	ForecastPeriod int     `yaml:"forecast_period"`
}

// ReportingAPIConfig holds settings for the reporting API.
type ReportingAPIConfig struct {
	Port     string `yaml:"port"`
	BasePath string `yaml:"base_path"`
}
