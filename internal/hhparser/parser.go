package hhparser

import (
	"errors"
	"fmt"
	"hhparser/internal/config"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	ErrNoConnection      = errors.New("http: Нет коннекта к HH")
	ErrCanNotReadData    = errors.New("io: Не могу прочитать данные с HH")
	ErrVacancyNotInteger = errors.New("strconv: Не могу перевести количество вакансий в число")
	ErrVacancyNotFind    = errors.New("getkeyWordByName: Вакансия не найдена в списке")
)

type ParserConfig struct {
	Cities             []config.CityConfig
	Technologies       []config.TechnologyConfig
	MaxGoroutines      int
	Timeout            time.Duration
	RetryCount         int
	RateLimit          time.Duration
	UrlSearchVacancies string
}

type Vacancy struct {
	Name       string
	SearchName string
	Count      int
	NumCity    int
}

func NewParserConfig(cfg *config.Config) ParserConfig {
	return ParserConfig{
		Cities:             cfg.Cities,
		Technologies:       cfg.Technologies,
		MaxGoroutines:      cfg.Parser.MaxGoroutines,
		Timeout:            cfg.Parser.Timeout,
		RetryCount:         cfg.Parser.RetryCount,
		RateLimit:          cfg.Parser.RateLimit,
		UrlSearchVacancies: cfg.Parser.UrlSearchVacancies,
	}
}

func GetAllVacancy(cfg ParserConfig) []*Vacancy {
	var keyWords = creatingKeywordsFromConfig(cfg)

	var wg sync.WaitGroup
	wg.Add(len(keyWords))

	semaphore := make(chan struct{}, cfg.MaxGoroutines)

	for _, keyWord := range keyWords {
		semaphore <- struct{}{} // Занимаем слот (блокируется, если уже 4 горутины работают)

		go func(kw *Vacancy) {
			defer wg.Done()
			defer func() { <-semaphore }() // Освобождаем слот при завершении

			kw.getCountVacancyFrom(cfg.UrlSearchVacancies, cfg.RetryCount)
		}(keyWord)
	}

	wg.Wait()

	return keyWords
}

func creatingKeywordsFromConfig(cfg ParserConfig) []*Vacancy {
	var vacancies []*Vacancy
	for _, city := range cfg.Cities {
		for _, tech := range cfg.Technologies {
			vacancies = append(vacancies, &Vacancy{
				Name:       tech.Name,
				SearchName: tech.Search,
				NumCity:    city.Code,
			})
		}
	}

	return vacancies
}

func (vacancy *Vacancy) getCountVacancyFrom(url string, maxRetries int) {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		var link = fmt.Sprintf(url, vacancy.SearchName, vacancy.NumCity)

		res, err := http.Get(link)
		if err != nil {
			panic(ErrNoConnection)
		}
		content, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			panic(ErrCanNotReadData)
		}

		var countVac, _ = injectSearchCounts(string(content))
		if countVac > 0 {
			vacancy.Count = countVac
			return
		}
	}
}

func GetkeyWordByNameAndCountry(vacancies []*Vacancy, name string, country int) *Vacancy {
	for _, vacancy := range vacancies {
		if vacancy.Name == name && vacancy.NumCity == country {
			return vacancy
		}
	}

	panic(ErrVacancyNotFind)
}
