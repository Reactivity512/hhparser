package storage

import (
	"encoding/json"
	"fmt"
	"hhparser/internal/config"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveJSON_Success(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir := t.TempDir()

	// Подготавливаем тестовые данные
	stats := Statistics{
		Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
		Technologies: []config.TechnologyConfig{
			{Name: "Golang", Search: "Golang", Enabled: true},
			{Name: "Python", Search: "Python", Enabled: true},
			{Name: "Java", Search: "Java", Enabled: true},
		},
		Cities: []CityStatistics{
			{
				Name: "MOSCOW",
				Code: 1,
				Vacancies: map[string]int{
					"Golang": 306,
					"Python": 3181,
					"Java":   815,
				},
				Total: 306 + 3181 + 815,
			},
			{
				Name: "KRASNODAR",
				Code: 2,
				Vacancies: map[string]int{
					"Golang": 4,
					"Python": 72,
					"Java":   29,
				},
				Total: 4 + 72 + 29,
			},
		},
		Summary: map[string]int{
			"Golang": 306 + 4,
			"Python": 3181 + 72,
			"Java":   815 + 29,
			"total":  306 + 3181 + 815 + 4 + 72 + 29,
		},
	}

	// Вызываем тестируемую функцию
	err := saveJSON(stats, tempDir)

	// Проверяем что нет ошибки
	assert.NoError(t, err)

	// Проверяем что файл создан с правильным именем
	expectedPath := filepath.Join(tempDir, "2026-02-14.json")
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err, "File should exist")

	// Проверяем содержимое файла
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	var savedStats Statistics
	err = json.Unmarshal(data, &savedStats)
	require.NoError(t, err)

	// Проверяем дату
	assert.Equal(t, stats.Date.Format("2006-01-02"), savedStats.Date.Format("2006-01-02"))

	// Проверяем технологии
	assert.Equal(t, len(stats.Technologies), len(savedStats.Technologies))
	assert.Equal(t, stats.Technologies[0].Name, savedStats.Technologies[0].Name)

	// Проверяем города
	assert.Equal(t, len(stats.Cities), len(savedStats.Cities))
	assert.Equal(t, stats.Cities[0].Name, savedStats.Cities[0].Name)
	assert.Equal(t, stats.Cities[0].Total, savedStats.Cities[0].Total)
	assert.Equal(t, stats.Cities[0].Vacancies["Golang"], savedStats.Cities[0].Vacancies["Golang"])
	assert.Equal(t, stats.Cities[1].Vacancies["Python"], savedStats.Cities[1].Vacancies["Python"])

	// Проверяем суммарную статистику
	assert.Equal(t, stats.Summary["Golang"], savedStats.Summary["Golang"])
	assert.Equal(t, stats.Summary["total"], savedStats.Summary["total"])
}

func TestSaveJSON_EmptyStats(t *testing.T) {
	tempDir := t.TempDir()

	// Пустая статистика
	stats := Statistics{
		Date:         time.Now(),
		Technologies: []config.TechnologyConfig{},
		Cities:       []CityStatistics{},
		Summary:      map[string]int{},
	}

	err := saveJSON(stats, tempDir)

	assert.NoError(t, err)

	// Проверяем что файл создан
	expectedPath := filepath.Join(tempDir, time.Now().Format("2006-01-02")+".json")
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	var savedStats Statistics
	err = json.Unmarshal(data, &savedStats)
	require.NoError(t, err)

	assert.Empty(t, savedStats.Technologies)
	assert.Empty(t, savedStats.Cities)
	assert.Empty(t, savedStats.Summary)
}

func TestSaveJSON_WithDisabledTechnologies(t *testing.T) {
	tempDir := t.TempDir()

	// Технологии с enabled = false
	stats := Statistics{
		Date: time.Now(),
		Technologies: []config.TechnologyConfig{
			{Name: "Golang", Search: "Golang", Enabled: true},
			{Name: "Python", Search: "Python", Enabled: false}, // отключена
			{Name: "Java", Search: "Java", Enabled: true},
		},
		Cities: []CityStatistics{
			{
				Name: "MOSCOW",
				Vacancies: map[string]int{
					"Golang": 306,
					"Python": 3181, // но данные есть!
					"Java":   815,
				},
			},
		},
	}

	err := saveJSON(stats, tempDir)
	assert.NoError(t, err)

	// Проверяем что все данные сохранились, включая отключенные технологии
	expectedPath := filepath.Join(tempDir, time.Now().Format("2006-01-02")+".json")
	data, _ := os.ReadFile(expectedPath)

	assert.Contains(t, string(data), "Python") // данные по Python должны быть
	assert.Contains(t, string(data), "false")  // enabled=false тоже должно быть
}

// Таблица тестов
func TestSaveJSON_TableDriven(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		stats     Statistics
		wantError bool
		checkFunc func(*testing.T, string, Statistics)
	}{
		{
			name: "полные данные",
			stats: Statistics{
				Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
				Technologies: []config.TechnologyConfig{
					{Name: "Golang", Enabled: true},
				},
				Cities: []CityStatistics{
					{
						Name:      "MOSCOW",
						Vacancies: map[string]int{"Golang": 306},
						Total:     306,
					},
				},
				Summary: map[string]int{"Golang": 306, "total": 306},
			},
			wantError: false,
		},
		{
			name: "нет городов",
			stats: Statistics{
				Date:         time.Now(),
				Technologies: []config.TechnologyConfig{{Name: "Golang"}},
				Cities:       []CityStatistics{},
				Summary:      map[string]int{},
			},
			wantError: false,
		},
		{
			name: "нет технологий",
			stats: Statistics{
				Date:         time.Now(),
				Technologies: []config.TechnologyConfig{},
				Cities:       []CityStatistics{{Name: "MOSCOW"}},
				Summary:      map[string]int{},
			},
			wantError: false,
		},
		{
			name: "nil карты",
			stats: Statistics{
				Date:         time.Now(),
				Technologies: nil,
				Cities:       nil,
				Summary:      nil,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := saveJSON(tt.stats, tempDir)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				filename := tt.stats.Date.Format("2006-01-02") + ".json"
				if tt.stats.Date.IsZero() {
					filename = "0001-01-01.json"
				}
				expectedPath := filepath.Join(tempDir, filename)

				_, err := os.Stat(expectedPath)
				assert.NoError(t, err, "Файл %s должен существовать", filename)
			}
		})
	}
}

// Тест на некорректную директорию
func TestSaveJSON_InvalidDirectory(t *testing.T) {
	stats := Statistics{
		Date:         time.Now(),
		Technologies: []config.TechnologyConfig{{Name: "Golang"}},
	}

	// Пытаемся записать в несуществующую директорию
	err := saveJSON(stats, "/nonexistent/directory/path")

	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err), "Ожидалась ошибка 'файл не найден'")
}

// Тест на большие данные
func TestSaveJSON_LargeData(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем тест с большими данными в коротком режиме")
	}

	tempDir := t.TempDir()

	// Создаем много технологий
	technologies := make([]config.TechnologyConfig, 100)
	for i := 0; i < 100; i++ {
		technologies[i] = config.TechnologyConfig{
			Name:    fmt.Sprintf("Tech%d", i),
			Search:  fmt.Sprintf("search%d", i),
			Enabled: i%2 == 0,
		}
	}

	// Создаем много городов
	cities := make([]CityStatistics, 50)
	summary := make(map[string]int)

	for i := 0; i < 50; i++ {
		cityName := fmt.Sprintf("CITY%d", i)
		vacancies := make(map[string]int)
		total := 0

		for j := 0; j < 100; j++ {
			techName := fmt.Sprintf("Tech%d", j)
			count := i * j * 10
			vacancies[techName] = count
			total += count
			summary[techName] += count
		}

		cities[i] = CityStatistics{
			Name:      cityName,
			Code:      i,
			Vacancies: vacancies,
			Total:     total,
		}
	}

	stats := Statistics{
		Date:         time.Now(),
		Technologies: technologies,
		Cities:       cities,
		Summary:      summary,
	}

	err := saveJSON(stats, tempDir)
	assert.NoError(t, err)

	// Проверяем размер файла (должен быть разумным)
	expectedPath := filepath.Join(tempDir, time.Now().Format("2006-01-02")+".json")
	info, err := os.Stat(expectedPath)
	require.NoError(t, err)

	t.Logf("Размер файла с большими данными: %d байт", info.Size())
	assert.Less(t, info.Size(), int64(10*1024*1024), "Файл слишком большой (>10MB)")
}

// Тест на конкурентную запись
func TestSaveJSON_ConcurrentWrites(t *testing.T) {
	tempDir := t.TempDir()

	concurrency := 10
	errChan := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			stats := Statistics{
				Date: time.Now().Add(time.Duration(id) * time.Hour), // разное время
				Technologies: []config.TechnologyConfig{
					{Name: fmt.Sprintf("Golang%d", id)},
				},
				Cities: []CityStatistics{
					{
						Name: "MOSCOW",
						Vacancies: map[string]int{
							fmt.Sprintf("Golang%d", id): id * 100,
						},
					},
				},
			}
			errChan <- saveJSON(stats, tempDir)
		}(i)
	}

	for i := 0; i < concurrency; i++ {
		err := <-errChan
		assert.NoError(t, err, "Конкурентная запись не должна вызывать ошибок")
	}
}

// Интеграционный тест: сохранили -> загрузили -> сравнили
func TestSaveAndLoadJSON_Integration(t *testing.T) {
	tempDir := t.TempDir()

	// Оригинальные данные
	original := Statistics{
		Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
		Technologies: []config.TechnologyConfig{
			{Name: "Golang", Search: "Golang", Enabled: true},
		},
		Cities: []CityStatistics{
			{
				Name: "MOSCOW",
				Code: 1,
				Vacancies: map[string]int{
					"Golang": 306,
				},
				Total: 306,
			},
		},
		Summary: map[string]int{
			"Golang": 306,
			"total":  306,
		},
	}

	// Сохраняем
	err := saveJSON(original, tempDir)
	require.NoError(t, err)

	// Загружаем обратно
	filePath := filepath.Join(tempDir, "2026-02-14.json")
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)

	var loaded Statistics
	err = json.Unmarshal(data, &loaded)
	require.NoError(t, err)

	// Сравниваем
	assert.Equal(t, original.Date.Format("2006-01-02"), loaded.Date.Format("2006-01-02"))
	assert.Equal(t, len(original.Technologies), len(loaded.Technologies))
	assert.Equal(t, original.Technologies[0].Name, loaded.Technologies[0].Name)
	assert.Equal(t, original.Technologies[0].Enabled, loaded.Technologies[0].Enabled)
	assert.Equal(t, original.Cities[0].Name, loaded.Cities[0].Name)
	assert.Equal(t, original.Cities[0].Total, loaded.Cities[0].Total)
	assert.Equal(t, original.Summary["Golang"], loaded.Summary["Golang"])
}

func BenchmarkSaveJSON(b *testing.B) {
	tempDir := b.TempDir()

	stats := Statistics{
		Date: time.Now(),
		Technologies: []config.TechnologyConfig{
			{Name: "Golang", Enabled: true},
			{Name: "Python", Enabled: true},
			{Name: "Java", Enabled: true},
			{Name: "C++", Enabled: true},
		},
		Cities: []CityStatistics{
			{
				Name: "MOSCOW",
				Vacancies: map[string]int{
					"Golang": 306,
					"Python": 3181,
					"Java":   815,
					"C++":    783,
				},
				Total: 306 + 3181 + 815 + 783,
			},
			{
				Name: "KRASNODAR",
				Vacancies: map[string]int{
					"Golang": 4,
					"Python": 72,
					"Java":   29,
					"C++":    15,
				},
				Total: 4 + 72 + 29 + 15,
			},
		},
		Summary: map[string]int{
			"Golang": 306 + 4,
			"Python": 3181 + 72,
			"Java":   815 + 29,
			"C++":    783 + 15,
			"total":  306 + 3181 + 815 + 783 + 4 + 72 + 29 + 15,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := saveJSON(stats, tempDir)
		if err != nil {
			b.Error(err)
		}
	}
}
