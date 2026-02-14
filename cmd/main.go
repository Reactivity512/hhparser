package main

import (
	"fmt"
	"hhparser/internal/config"
	"hhparser/internal/hhparser"
	"hhparser/internal/storage"
	"log"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	fmt.Printf("Старт парсинга: %s\n", startTime.Format("15:04:05"))

	vacancy := hhparser.GetAllVacancy(hhparser.NewParserConfig(cfg))
	if err := storage.SaveStatistics(vacancy, storage.NewStorageConfig(cfg)); err != nil {
		log.Fatal(err)
	}

	endTime := time.Now()
	fmt.Printf("Завершено: %s\n", endTime.Format("15:04:05"))
	fmt.Printf("Длительность: %v\n", endTime.Sub(startTime))
}
