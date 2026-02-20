package storage

import (
	"fmt"
	"hhparser/internal/config"
	"hhparser/internal/hhparser"
	"os"
	"time"
)

type StorageConfig struct {
	Cities       []config.CityConfig
	Technologies []config.TechnologyConfig
	DataDir      string
}

type Statistics struct {
	Date         time.Time                 `json:"date"`
	Technologies []config.TechnologyConfig `json:"technologiesConfig"`
	Cities       []CityStatistics          `json:"cities"`
	Summary      map[string]int            `json:"summary"`
}

type CityStatistics struct {
	Name      string         `json:"name"`
	Code      int            `json:"code"`
	Vacancies map[string]int `json:"vacancies"`
	Total     int            `json:"total"`
}

func NewStorageConfig(cfg *config.Config) StorageConfig {
	return StorageConfig{
		Cities:       cfg.Cities,
		Technologies: cfg.Technologies,
		DataDir:      cfg.Output.Directory,
	}
}

func SaveStatistics(vacancies []*hhparser.Vacancy, cfg StorageConfig) error {
	stats := collectStatistics(vacancies, cfg)

	if err := ensureDir(cfg.DataDir); err != nil {
		return err
	}

	if err := saveJSON(stats, cfg.DataDir); err != nil {
		return err
	}

	if err := saveTXT(stats, cfg.DataDir); err != nil {
		return err
	}

	return nil
}

func ensureDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}
		return nil
	}

	if !info.IsDir() {
		return fmt.Errorf("path '%s' exists but is not a directory", dir)
	}
	return nil
}

func collectStatistics(vacancies []*hhparser.Vacancy, cfg StorageConfig) Statistics {
	stats := Statistics{
		Date:         time.Now(),
		Technologies: cfg.Technologies,
		Summary:      make(map[string]int),
	}

	for _, city := range cfg.Cities {
		cityStat := CityStatistics{
			Name:      city.Name,
			Code:      city.Code,
			Vacancies: make(map[string]int),
		}

		for _, tech := range cfg.Technologies {
			count := hhparser.GetkeyWordByNameAndCountry(vacancies, tech.Name, city.Code).Count
			cityStat.Vacancies[tech.Name] = count
			cityStat.Total += count
			stats.Summary[tech.Name] += count
		}

		stats.Cities = append(stats.Cities, cityStat)
	}

	return stats
}
