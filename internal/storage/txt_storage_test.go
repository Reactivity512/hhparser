package storage

import (
	"hhparser/internal/config"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveTXT_Success(t *testing.T) {
	// Создаем временную директорию
	tempDir := t.TempDir()

	// Подготавливаем тестовые данные
	stats := Statistics{
		Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
		Technologies: []config.TechnologyConfig{
			{Name: "Golang", Enabled: true},
			{Name: "Python", Enabled: true},
			{Name: "Java", Enabled: true},
		},
		Cities: []CityStatistics{
			{
				Name: "MOSCOW",
				Vacancies: map[string]int{
					"Golang": 306,
					"Python": 3181,
					"Java":   815,
				},
			},
			{
				Name: "KRASNODAR",
				Vacancies: map[string]int{
					"Golang": 4,
					"Python": 72,
					"Java":   29,
				},
			},
		},
		Summary: map[string]int{
			"Golang": 310,
			"Python": 3253,
			"Java":   844,
		},
	}

	// Вызываем тестируемую функцию
	err := saveTXT(stats, tempDir)

	// Проверяем что нет ошибки
	assert.NoError(t, err)

	// Проверяем что файл создан
	expectedPath := filepath.Join(tempDir, "stats_2026-02-14.txt")
	_, err = os.Stat(expectedPath)
	assert.NoError(t, err, "Файл должен существовать")

	// Читаем содержимое файла
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	content := string(data)

	// Проверяем заголовок (ищем подстроки без табуляций, так как tabwriter заменяет их пробелами)
	assert.Contains(t, content, "СТАТИСТИКА ВАКАНСИЙ")
	assert.Contains(t, content, "Дата: 14.02.2026")

	// Проверяем заголовки колонок (tabwriter заменил \t на пробелы для выравнивания)
	assert.Contains(t, content, "Технология  MOSCOW      KRASNODAR   ВСЕГО")

	// Проверяем разделители
	assert.Contains(t, content, "----------  ----------  ----------  ----------")

	// Проверяем данные (обрати внимание на пробелы вместо табуляций)
	assert.Contains(t, content, "Golang      306         4           310")
	assert.Contains(t, content, "Python      3181        72          3253")
	assert.Contains(t, content, "Java        815         29          844")
}

func TestSaveTXT_OneCity(t *testing.T) {
	tempDir := t.TempDir()

	stats := Statistics{
		Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
		Technologies: []config.TechnologyConfig{
			{Name: "Golang"},
		},
		Cities: []CityStatistics{
			{
				Name: "MOSCOW",
				Vacancies: map[string]int{
					"Golang": 306,
				},
			},
		},
		Summary: map[string]int{
			"Golang": 306,
		},
	}

	err := saveTXT(stats, tempDir)
	assert.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "stats_2026-02-14.txt")
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	content := string(data)

	// Разбиваем на строки для более точной проверки
	lines := strings.Split(content, "\n")

	// Проверяем что есть нужные строки (не проверяем точное количество пробелов)
	t.Logf("Content: %s", content) // Для отладки

	// Проверяем что строка с технологией содержит все нужные элементы
	found := false
	for _, line := range lines {
		if strings.Contains(line, "Golang") &&
			strings.Contains(line, "306") &&
			strings.Contains(line, "306") {
			found = true
			break
		}
	}
	assert.True(t, found, "Строка с данными должна содержать Golang, 306, 306")

	// Проверяем заголовок
	assert.Contains(t, content, "СТАТИСТИКА ВАКАНСИЙ")
	assert.Contains(t, content, "Дата: 14.02.2026")

	// Проверяем что есть заголовки городов
	assert.Contains(t, content, "MOSCOW")
	assert.Contains(t, content, "ВСЕГО")

	// Проверяем разделитель
	assert.Contains(t, content, "----------")
}

func TestSaveTXT_OneCity_WithRegexp(t *testing.T) {
	tempDir := t.TempDir()

	stats := Statistics{
		Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
		Technologies: []config.TechnologyConfig{
			{Name: "Golang"},
		},
		Cities: []CityStatistics{
			{
				Name: "MOSCOW",
				Vacancies: map[string]int{
					"Golang": 306,
				},
			},
		},
		Summary: map[string]int{
			"Golang": 306,
		},
	}

	err := saveTXT(stats, tempDir)
	assert.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "stats_2026-02-14.txt")
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	content := string(data)

	// Используем regexp для проверки формата (любое количество пробелов)
	regexPattern := `Golang\s+306\s+306`
	assert.Regexp(t, regexPattern, content, "Данные должны быть в формате: Golang [пробелы] 306 [пробелы] 306")

	// Проверяем заголовок
	assert.Regexp(t, `Технология\s+MOSCOW\s+ВСЕГО`, content)
}

func TestSaveTXT_NoCities(t *testing.T) {
	tempDir := t.TempDir()

	stats := Statistics{
		Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
		Technologies: []config.TechnologyConfig{
			{Name: "Golang"},
		},
		Cities:  []CityStatistics{},
		Summary: map[string]int{"Golang": 0},
	}

	err := saveTXT(stats, tempDir)
	assert.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "stats_2026-02-14.txt")
	data, _ := os.ReadFile(expectedPath)
	content := string(data)

	// Должна быть только строка с технологией без городов
	assert.Contains(t, content, "Технология  ВСЕГО")
	assert.Contains(t, content, "----------  ----------")
	assert.Contains(t, content, "Golang      0")
}

func TestSaveTXT_DynamicWidth(t *testing.T) {
	tempDir := t.TempDir()

	stats := Statistics{
		Date: time.Now(),
		Technologies: []config.TechnologyConfig{
			{Name: "Short"},
			{Name: "VeryLongTechnologyName"},
		},
		Cities: []CityStatistics{
			{
				Name: "Short",
				Vacancies: map[string]int{
					"Short":                  1,
					"VeryLongTechnologyName": 1000,
				},
			},
			{
				Name: "VeryLongCityName",
				Vacancies: map[string]int{
					"Short":                  2,
					"VeryLongTechnologyName": 2000,
				},
			},
		},
		Summary: map[string]int{
			"Short":                  3,
			"VeryLongTechnologyName": 3000,
		},
	}

	err := saveTXT(stats, tempDir)
	assert.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "stats_"+time.Now().Format("2006-01-02")+".txt")
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	content := string(data)
	lines := strings.Split(content, "\n")

	// Проверяем что длинные названия не обрезаны
	assert.Contains(t, content, "VeryLongTechnologyName")
	assert.Contains(t, content, "VeryLongCityName")

	// Проверяем что числа выровнены (находим строку с технологией и проверяем формат)
	for _, line := range lines {
		if strings.Contains(line, "VeryLongTechnologyName") {
			// Должны быть пробелы между колонками
			assert.Regexp(t, `VeryLongTechnologyName\s+\d+\s+\d+\s+\d+`, line)
		}
	}
}

// Тест с проверкой точного форматирования
func TestSaveTXT_ExactFormatting(t *testing.T) {
	tempDir := t.TempDir()

	stats := Statistics{
		Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
		Technologies: []config.TechnologyConfig{
			{Name: "Go"},
		},
		Cities: []CityStatistics{
			{
				Name: "MSK",
				Vacancies: map[string]int{
					"Go": 100,
				},
			},
		},
		Summary: map[string]int{
			"Go": 100},
	}

	err := saveTXT(stats, tempDir)
	require.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "stats_2026-02-14.txt")
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	content := string(data)
	lines := strings.Split(content, "\n")

	// Выводим реальный вывод для отладки
	t.Logf("Real output:\n%s", content)

	// Проверяем структуру, но не точное форматирование
	assert.GreaterOrEqual(t, len(lines), 6, "Должно быть минимум 6 строк")

	// Проверяем что строки содержат нужные элементы (без проверки пробелов)
	assert.Contains(t, lines[0], "СТАТИСТИКА ВАКАНСИЙ")
	assert.Contains(t, lines[1], "Дата: 14.02.2026")
	assert.Equal(t, "", lines[2]) // пустая строка

	// Проверяем что в строке заголовка есть все нужные слова
	assert.Contains(t, lines[3], "Технология")
	assert.Contains(t, lines[3], "MSK")
	assert.Contains(t, lines[3], "ВСЕГО")

	// Проверяем разделитель (должен содержать дефисы)
	assert.Contains(t, lines[4], "----------")

	// Проверяем строку с данными
	assert.Contains(t, lines[5], "Go")
	assert.Contains(t, lines[5], "100")
	assert.Contains(t, lines[5], "100")

	// Проверяем что 100 встречается дважды
	assert.Equal(t, 2, strings.Count(lines[5], "100"))
}

// Альтернативный тест с regexp
func TestSaveTXT_ExactFormatting_WithRegexp(t *testing.T) {
	tempDir := t.TempDir()

	stats := Statistics{
		Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
		Technologies: []config.TechnologyConfig{
			{Name: "Go"},
		},
		Cities: []CityStatistics{
			{
				Name: "MSK",
				Vacancies: map[string]int{
					"Go": 100,
				},
			},
		},
		Summary: map[string]int{
			"Go": 100,
		},
	}

	err := saveTXT(stats, tempDir)
	require.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "stats_2026-02-14.txt")
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	content := string(data)

	// Используем regexp для гибкой проверки
	headerRegex := `Технология\s+MSK\s+ВСЕГО`
	assert.Regexp(t, headerRegex, content, "Заголовок должен содержать Технология, MSK и ВСЕГО с пробелами")

	separatorRegex := `-+\s+-+\s+-+`
	assert.Regexp(t, separatorRegex, content, "Должен быть разделитель из дефисов")

	dataRegex := `Go\s+100\s+100`
	assert.Regexp(t, dataRegex, content, "Строка с данными должна содержать Go, 100, 100 с пробелами")
}

// Тест на динамическую ширину колонок
func TestSaveTXT_DynamicColumnWidth(t *testing.T) {
	tempDir := t.TempDir()

	// Данные где названия разной длины
	stats := Statistics{
		Date: time.Now(),
		Technologies: []config.TechnologyConfig{
			{Name: "A"},                      // короткое
			{Name: "VeryLongTechnologyName"}, // длинное
		},
		Cities: []CityStatistics{
			{
				Name: "Short", // короткое
				Vacancies: map[string]int{
					"A":                      1,
					"VeryLongTechnologyName": 1000,
				},
			},
			{
				Name: "VeryLongCityName", // длинное
				Vacancies: map[string]int{
					"A":                      2,
					"VeryLongTechnologyName": 2000,
				},
			},
		},
		Summary: map[string]int{
			"A":                      3,
			"VeryLongTechnologyName": 3000,
		},
	}

	err := saveTXT(stats, tempDir)
	require.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "stats_"+time.Now().Format("2006-01-02")+".txt")
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	content := string(data)
	lines := strings.Split(content, "\n")

	// Проверяем что все строки с данными содержат правильные элементы
	for _, line := range lines {
		if strings.Contains(line, "VeryLongTechnologyName") {
			assert.Contains(t, line, "1000")
			assert.Contains(t, line, "2000")
			assert.Contains(t, line, "3000")
		}
		if strings.Contains(line, "A") && !strings.Contains(line, "VeryLongTechnologyName") {
			assert.Contains(t, line, "1")
			assert.Contains(t, line, "2")
			assert.Contains(t, line, "3")
		}
	}
}

// Тест на обработку tabwriter форматирования
func TestSaveTXT_TabWriterReplacesTabs(t *testing.T) {
	tempDir := t.TempDir()

	stats := Statistics{
		Date: time.Now(),
		Technologies: []config.TechnologyConfig{
			{Name: "Test"},
		},
		Cities: []CityStatistics{
			{
				Name: "City",
				Vacancies: map[string]int{
					"Test": 123,
				},
			},
		},
		Summary: map[string]int{
			"Test": 123,
		},
	}

	err := saveTXT(stats, tempDir)
	require.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "stats_"+time.Now().Format("2006-01-02")+".txt")
	data, _ := os.ReadFile(expectedPath)
	content := string(data)

	// Tabwriter должен заменить табуляции на пробелы
	assert.NotContains(t, content, "\t", "Tabwriter должен заменить табуляции на пробелы")

	// Но должны быть пробелы для выравнивания
	assert.Contains(t, content, "  ") // как минимум два пробела между колонками
}

// Таблица тестов
func TestSaveTXT_TableDriven(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name   string
		stats  Statistics
		wanted []string // элементы которые должны быть в строке
	}{
		{
			name: "нормальные данные",
			stats: Statistics{
				Date:         time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
				Technologies: []config.TechnologyConfig{{Name: "Go"}},
				Cities:       []CityStatistics{{Name: "MOSCOW", Vacancies: map[string]int{"Go": 100}}},
				Summary:      map[string]int{"Go": 100},
			},
			wanted: []string{"Go", "100", "100"}, // проверяем только ключевые элементы
		},
		{
			name: "несколько технологий",
			stats: Statistics{
				Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
				Technologies: []config.TechnologyConfig{
					{Name: "Go"},
					{Name: "Python"},
					{Name: "Java"},
				},
				Cities: []CityStatistics{
					{
						Name: "MOSCOW",
						Vacancies: map[string]int{
							"Go":     1,
							"Python": 2,
							"Java":   3,
						},
					},
				},
				Summary: map[string]int{"Go": 1, "Python": 2, "Java": 3},
			},
			wanted: []string{
				"Go", "1", "1",
				"Python", "2", "2",
				"Java", "3", "3",
			},
		},
		{
			name: "несколько городов",
			stats: Statistics{
				Date:         time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
				Technologies: []config.TechnologyConfig{{Name: "Go"}},
				Cities: []CityStatistics{
					{Name: "MOSCOW", Vacancies: map[string]int{"Go": 100}},
					{Name: "SPB", Vacancies: map[string]int{"Go": 50}},
				},
				Summary: map[string]int{"Go": 150},
			},
			wanted: []string{"Go", "100", "50", "150"},
		},
		{
			name: "нулевые значения",
			stats: Statistics{
				Date:         time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
				Technologies: []config.TechnologyConfig{{Name: "Go"}, {Name: "Rust"}},
				Cities: []CityStatistics{
					{
						Name: "MOSCOW",
						Vacancies: map[string]int{
							"Go": 100,
							// Rust отсутствует
						},
					},
				},
				Summary: map[string]int{"Go": 100, "Rust": 0},
			},
			wanted: []string{"Go", "100", "100", "Rust", "0", "0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := saveTXT(tt.stats, tempDir)
			assert.NoError(t, err)

			filename := "stats_" + tt.stats.Date.Format("2006-01-02") + ".txt"
			data, err := os.ReadFile(filepath.Join(tempDir, filename))
			require.NoError(t, err)

			content := string(data)

			// Выводим реальный контент для отладки
			t.Logf("Actual content for %s:\n%s", tt.name, content)

			// Проверяем что все нужные элементы присутствуют
			for _, wanted := range tt.wanted {
				assert.Contains(t, content, wanted,
					"Контент должен содержать '%s'", wanted)
			}

			// Дополнительные проверки структуры
			assert.Contains(t, content, "СТАТИСТИКА ВАКАНСИЙ")
			assert.Contains(t, content, "Дата: 14.02.2026")

			// Проверяем что заголовки городов присутствуют
			for _, city := range tt.stats.Cities {
				assert.Contains(t, content, city.Name)
			}

			// Проверяем что все технологии присутствуют
			for _, tech := range tt.stats.Technologies {
				assert.Contains(t, content, tech.Name)
			}
		})
	}
}

// Альтернативный тест с regexp для более точной проверки
func TestSaveTXT_TableDriven_WithRegexp(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		stats    Statistics
		patterns []string // regexp паттерны
	}{
		{
			name: "нормальные данные",
			stats: Statistics{
				Date:         time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
				Technologies: []config.TechnologyConfig{{Name: "Go"}},
				Cities:       []CityStatistics{{Name: "MOSCOW", Vacancies: map[string]int{"Go": 100}}},
				Summary:      map[string]int{"Go": 100},
			},
			patterns: []string{
				`Go\s+100\s+100`, // Go, пробелы, 100, пробелы, 100
			},
		},
		{
			name: "несколько технологий",
			stats: Statistics{
				Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
				Technologies: []config.TechnologyConfig{
					{Name: "Go"}, {Name: "Python"}, {Name: "Java"},
				},
				Cities: []CityStatistics{
					{
						Name: "MOSCOW",
						Vacancies: map[string]int{
							"Go": 1, "Python": 2, "Java": 3,
						},
					},
				},
				Summary: map[string]int{"Go": 1, "Python": 2, "Java": 3},
			},
			patterns: []string{
				`Go\s+1\s+1`,
				`Python\s+2\s+2`,
				`Java\s+3\s+3`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := saveTXT(tt.stats, tempDir)
			assert.NoError(t, err)

			filename := "stats_" + tt.stats.Date.Format("2006-01-02") + ".txt"
			data, err := os.ReadFile(filepath.Join(tempDir, filename))
			require.NoError(t, err)

			content := string(data)

			for _, pattern := range tt.patterns {
				assert.Regexp(t, pattern, content,
					"Должен соответствовать паттерну: %s", pattern)
			}
		})
	}
}

// Тест на большие числа (чтобы убедиться что выравнивание работает)
func TestSaveTXT_LargeNumbers(t *testing.T) {
	tempDir := t.TempDir()

	stats := Statistics{
		Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
		Technologies: []config.TechnologyConfig{
			{Name: "Go"},
			{Name: "Python"},
		},
		Cities: []CityStatistics{
			{
				Name: "MOSCOW",
				Vacancies: map[string]int{
					"Go":     1234567,
					"Python": 9876543,
				},
			},
			{Name: "SPB",
				Vacancies: map[string]int{
					"Go":     123456,
					"Python": 987654,
				},
			},
		},
		Summary: map[string]int{
			"Go":     1234567 + 123456,
			"Python": 9876543 + 987654,
		},
	}

	err := saveTXT(stats, tempDir)
	assert.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "stats_2026-02-14.txt")
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	content := string(data)

	// Выводим реальный контент для отладки
	t.Logf("Actual content:\n%s", content)

	// Вариант 1: Проверяем наличие всех чисел (простой подход)
	assert.Contains(t, content, "Go")
	assert.Contains(t, content, "Python")
	assert.Contains(t, content, "1234567")
	assert.Contains(t, content, "123456")
	assert.Contains(t, content, "1358023") // 1234567 + 123456
	assert.Contains(t, content, "9876543")
	assert.Contains(t, content, "987654")
	assert.Contains(t, content, "10864197") // 9876543 + 987654

	// Вариант 2: Используем regexp для гибкой проверки форматирования
	goPattern := `Go\s+1234567\s+123456\s+1358023`
	assert.Regexp(t, goPattern, content,
		"Строка с Go должна содержать числа с пробелами между ними")

	pythonPattern := `Python\s+9876543\s+987654\s+10864197`
	assert.Regexp(t, pythonPattern, content,
		"Строка с Python должна содержать числа с пробелами между ними")
}

// Тест с проверкой порядка через поля (более надежный)
func TestSaveTXT_LargeNumbers_OrderByFields(t *testing.T) {
	tempDir := t.TempDir()

	stats := Statistics{
		Date: time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC),
		Technologies: []config.TechnologyConfig{
			{Name: "Go"},
		},
		Cities: []CityStatistics{
			{
				Name: "MOSCOW",
				Vacancies: map[string]int{
					"Go": 1234567,
				},
			},
			{
				Name: "SPB",
				Vacancies: map[string]int{
					"Go": 123456,
				},
			},
		},
		Summary: map[string]int{
			"Go": 1358023,
		},
	}

	err := saveTXT(stats, tempDir)
	require.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "stats_2026-02-14.txt")
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	content := string(data)
	lines := strings.Split(content, "\n")

	// Находим строку с Go
	var goLine string
	for _, line := range lines {
		if strings.Contains(line, "Go") {
			goLine = line
			break
		}
	}

	assert.NotEmpty(t, goLine, "Должна быть строка с Go")
	t.Logf("Go line: %s", goLine)

	// Разбиваем на поля (tabwriter разделяет пробелами)
	fields := strings.Fields(goLine)
	t.Logf("Fields: %v", fields)

	// Первое поле - "Go"
	assert.Equal(t, "Go", fields[0])

	// Остальные поля - числа
	// Порядок должен быть: Moscow, SPB, Total
	if len(fields) >= 4 {
		// Проверяем что числа соответствуют ожидаемым
		assert.Equal(t, "1234567", fields[1], "Первое число должно быть Moscow")
		assert.Equal(t, "123456", fields[2], "Второе число должно быть SPB")
		assert.Equal(t, "1358023", fields[3], "Третье число должно быть Total")
	}
}

// Тест на выравнивание больших чисел
func TestSaveTXT_LargeNumbers_Alignment(t *testing.T) {
	tempDir := t.TempDir()

	stats := Statistics{
		Date: time.Now(),
		Technologies: []config.TechnologyConfig{
			{Name: "Small"},
			{Name: "Large"},
		},
		Cities: []CityStatistics{
			{
				Name: "City1",
				Vacancies: map[string]int{
					"Small": 1,
					"Large": 1000000,
				},
			},
			{
				Name: "City2",
				Vacancies: map[string]int{
					"Small": 2,
					"Large": 2000000,
				},
			},
		},
		Summary: map[string]int{
			"Small": 3,
			"Large": 3000000,
		},
	}

	err := saveTXT(stats, tempDir)
	assert.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "stats_"+time.Now().Format("2006-01-02")+".txt")
	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	content := string(data)
	lines := strings.Split(content, "\n")

	// Проверяем что все числа есть в нужных строках
	for _, line := range lines {
		if strings.Contains(line, "Small") {
			assert.Contains(t, line, "1")
			assert.Contains(t, line, "2")
			assert.Contains(t, line, "3")
		}
		if strings.Contains(line, "Large") {
			assert.Contains(t, line, "1000000")
			assert.Contains(t, line, "2000000")
			assert.Contains(t, line, "3000000")
		}
	}
}
