package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Cities       []CityConfig       `mapstructure:"cities"`
	Technologies []TechnologyConfig `mapstructure:"technologies"`
	Parser       ParserConfig       `mapstructure:"parser"`
	Output       OutputConfig       `mapstructure:"output"`
}

type CityConfig struct {
	ID      int    `mapstructure:"id"`
	Name    string `mapstructure:"name"`
	Code    int    `mapstructure:"code"`
	Enabled bool   `mapstructure:"enabled"`
}

type TechnologyConfig struct {
	Name     string `mapstructure:"name"`
	Search   string `mapstructure:"search"`
	Category string `mapstructure:"category"`
	Enabled  bool   `mapstructure:"enabled"`
}

type ParserConfig struct {
	MaxGoroutines      int    `mapstructure:"max_goroutines"`
	TimeoutSeconds     int    `mapstructure:"timeout_seconds"`
	RetryCount         int    `mapstructure:"retry_count"`
	RateLimitMs        int    `mapstructure:"rate_limit_ms"`
	UrlSearchVacancies string `mapstructure:"url_search_vacancies"`

	// Вычисляемые поля
	Timeout   time.Duration
	RateLimit time.Duration
}

type OutputConfig struct {
	Format         string `mapstructure:"format"`
	Directory      string `mapstructure:"directory"`
	FilenamePrefix string `mapstructure:"filename_prefix"`
}

func Load() (*Config, error) {
	err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	setDefaultValues()

	if err := readConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Вычислить поля
	config.Parser.Timeout = time.Duration(config.Parser.TimeoutSeconds) * time.Second
	config.Parser.RateLimit = time.Duration(config.Parser.RateLimitMs) * time.Millisecond

	config.filterEnabled()

	return &config, nil
}

func loadConfig() error {
	// Получаем директорию исполняемого файла
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Получаем рабочую директорию
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get working directory: %w", err)
	}

	// Ищем конфиг в нескольких местах
	searchPaths := []string{
		// Относительно рабочей директории
		filepath.Join(workDir, "configs"),
		filepath.Join(workDir, "..", "configs"), // если в cmd

		// Относительно исполняемого файла
		filepath.Join(exeDir, "configs"),
		filepath.Join(exeDir, "..", "configs"),

		// Абсолютные пути из env
		os.Getenv("VACANCY_CONFIG_PATH"),

		// Стандартные пути
		"./configs",
		"../configs",
		".",
	}

	// Добавляем все пути в viper
	for _, path := range searchPaths {
		if path != "" {
			viper.AddConfigPath(path)
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	return nil
}

func setDefaultValues() {
	viper.AutomaticEnv()

	viper.SetDefault("parser.max_goroutines", 4)
	viper.SetDefault("parser.timeout_seconds", 10)
	viper.SetDefault("parser.retry_count", 2)
	viper.SetDefault("parser.rate_limit_ms", 200)
	viper.SetDefault("output.format", "json")
}

func readConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found, using defaults and environment variables")
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}
	return nil
}

func (c *Config) filterEnabled() {
	// Оставляем только включенные города
	enabledCities := make([]CityConfig, 0, len(c.Cities))
	for _, city := range c.Cities {
		if city.Enabled {
			enabledCities = append(enabledCities, city)
		}
	}
	c.Cities = enabledCities

	// Оставляем только включенные технологии
	enabledTechs := make([]TechnologyConfig, 0, len(c.Technologies))
	for _, tech := range c.Technologies {
		if tech.Enabled {
			enabledTechs = append(enabledTechs, tech)
		}
	}
	c.Technologies = enabledTechs
}

func (c *Config) Validate() error {
	if len(c.Cities) == 0 {
		return fmt.Errorf("нет включенных городов для парсинга")
	}

	if len(c.Technologies) == 0 {
		return fmt.Errorf("нет включенных технологий для парсинга")
	}

	if c.Parser.MaxGoroutines <= 0 {
		return fmt.Errorf("max_goroutines должен быть > 0")
	}

	return nil
}
