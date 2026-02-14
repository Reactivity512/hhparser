# Vacancy Parser for hh.ru

[![CI](https://github.com/Reactivity512/hhparser/actions/workflows/main.yml/badge.svg)](https://github.com/Reactivity512/hhparser/actions/workflows/main.yml)
[![Coverage](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fraw.githubusercontent.com%2FReactivity512%2Fhhparser%2Fmain%2Fcoverage.json&query=%24.coverage&label=Coverage&color=brightgreen)](https://github.com/Reactivity512/hhparser)
[![Go Report Card](https://goreportcard.com/badge/github.com/Reactivity512/hhparser)](https://goreportcard.com/report/github.com/Reactivity512/hhparser)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Парсер вакансий с hh.ru для сбора статистики по технологиям и городам. Собирает данные о количестве вакансий для различных языков программирования и технологий.

## 📋 Содержание
- [Возможности](#возможности)
- [Структура проекта](#структура-проекта)
- [Установка](#установка)
- [Конфигурация](#конфигурация)
- [Использование](#использование)
- [Тестирование](#тестирование)
- [CI/CD](#cicd)
- [Участие в разработке](#участие-в-разработке)

## ✨ Возможности

- Парсинг вакансий по технологиям:
  - Языки: C++, Go, Java, C#, Python, PHP, JavaScript ...
  - Фреймворки: Node.js, Spring, Django, Laravel ...
  - Роли: DevOps, Team Lead ...
- Поддержка нескольких городов (Москва, Краснодар ...)
- Ограничение количества одновременных запросов
- Retry логика при ошибках
- Сохранение в нескольких форматах (JSON, TXT)
- Автоматическое создание структуры директорий

## 📁 Структура проекта
```
hhparser/
├── cmd/
│ └── main.go       # Точка входа
├── internal/
│ ├── config/       # Конфигурация
│ │ ├── config.go
│ ├── hhparser/     # Парсер hh.ru
│ │ ├── hhparser.go
│ │ ├── parser_inject.go
│ │ └── parser_inject_test.go
│ └── storage/      # Сохранение данных
│ ├── storage.go
│ ├── json_storage.go.go
│ ├── json_storage_test.go
│ ├── txt_storage.go
│ └── txt_storage_test.go
├── configs/
│ └── config.yaml  # Конфигурационный файл
├── data/          # Директория для данных
├── golangci-lint/ # Линтер
├── .gitignore
├── go.mod
├── go.sum
├── Taskfile.yml   # Задачи для Task
└── README.md
```

## 🔧 Установка

### Предварительные требования
- Go 1.21 или выше
- Git
- (Опционально) Task

### Клонирование репозитория

```bash
git clone https://github.com/Reactivity512/hhparser.git
cd hhparser
```

### Установка зависимостей
```bash
go mod download
go mod tidy
```

### Установка инструментов разработки
```bash
# Установка golangci-lint
# Windows (choco)
choco install golangci-lint

# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Установка Task (опционально)
go install github.com/go-task/task/v3/cmd/task@latest
```

## ⚙️ Конфигурация
Основной конфигурационный файл `configs/config.yaml`
```yaml
# Города для парсинга (code - код города для hh.ru)
cities:
  - id: 1
    name: "MOSCOW"
    code: 1
    enabled: true
  - id: 2
    name: "KRASNODAR"
    code: 53
    enabled: true

# Технологии для поиска
technologies:
  - name: "Cpp"
    search: "C%2B%2B"
    category: "languages"
    enabled: true
  - name: "Golang"
    search: "Golang"
    category: "languages"
    enabled: true
  - name: "Java"
    search: "Java"
    category: "languages"
    enabled: true
  - name: "CSharp"
    search: "C%23"
    category: "languages"
    enabled: true
  - name: "Python"
    search: "Python"
    category: "languages"
    enabled: true
  - name: "Php"
    search: "Php"
    category: "languages"
    enabled: true
  - name: "Javascript"
    search: "Javascript"
    category: "languages"
    enabled: true
  - name: "Laravel"
    search: "Laravel"
    category: "framework"
    enabled: true
  - name: "Spring"
    search: "Spring"
    category: "framework"
    enabled: true
  - name: "Nodejs"
    search: "Node.js"
    category: "framework"
    enabled: true
  - name: "Django"
    search: "Django"
    category: "framework"
    enabled: true
  - name: "Devops"
    search: "Devops"
    category: "roles"
    enabled: true
  - name: "TeamLead"
    search: "Team+lead"
    category: "roles"
    enabled: true

# Настройки парсера
parser:
  max_goroutines: 4  # Кол-во одновременных соединений с hh.ru, можно больше, но тогда бываю разрывы соединения
  timeout_seconds: 10
  retry_count: 2
  rate_limit_ms: 200
  url_search_vacancies: "https://hh.ru/search/vacancy?text=%s&salary=&ored_clusters=true&area=%d&hhtmFrom=vacancy_search_list&hhtmFromLabel=vacancy_search_line"

# Пути для данных
output:
  format: "json"
  directory: "./data"
  filename_prefix: "vacancies"
```

## 🚀 Использование

### Запуск парсера
```bash
# Простой запуск
go run cmd/main.go
```

### Использование Task (рекомендуется)
```bash
# Список всех задач
task

# Запуск парсера
task run

# Сборка
task build

# Сборка под все платформы
task build-all

# Сборка релизной версии
task build-release

# Установка зависимостей
task deps

# Обновление зависимостей
task deps-update

# Очистка временных файлов
task clean

# Очистка только данных
task clean-data

# Запуск тестов
task test

# Форматирование кода
task fmt

# Запуск линтера
task lint

# Запустить go vet
task vet
```

## 🧪 Тестирование

### Запуск тестов
```bash
# Все тесты
go test ./... -v

# Тесты с покрытием
go test ./... -coverprofile=coverage.out

# Тесты с детектором гонок (требует, чтобы была включена поддержка CGO)
go test -race ./...

# Включение поддержка CGO
export CGO_ENABLED=1 # В Linux/MacOS
set CGO_ENABLED=1   # В Windows (cmd)
$env:CGO_ENABLED="1" # В Windows (PowerShell)

# Конкретный пакет
go test ./internal/storage -v
```

### Линтер
```bash
# Запуск линтера
golangci-lint run ./...

# С автоисправлением
golangci-lint run --fix ./...

# Только быстрые проверки
golangci-lint run --fast ./...
```

## 📊 Результаты

Данные сохраняются в директории `data/` в нескольких форматах:

JSON
```json
{
  "date": "2026-02-15T01:11:59.6621431+03:00",
  "technologiesConfig": [
    {
      "Name": "Golang",
      "Search": "Golang",
      "Category": "languages",
      "Enabled": true
    },
    {
      "Name": "Python",
      "Search": "Python",
      "Category": "languages",
      "Enabled": true
    },
  ],
  "cities": [
    {
      "name": "MOSCOW",
      "code": 1,
      "vacancies": {
        "Golang": 406,
        "Python": 4545,
      },
      "total": 4951
    },
    {
      "name": "KRASNODAR",
      "code": 53,
      "vacancies": {
        "Golang": 5,
        "Python": 74,
      },
      "total": 79
    }
  ],
  "summary": {
    "Golang": 411,
    "Python": 4619,
  }
}
```

TXT
```
СТАТИСТИКА ВАКАНСИЙ
Дата: 15.02.2026

Технология  MOSCOW      KRASNODAR   ВСЕГО
----------  ----------  ----------  ----------  
Golang      406         5           411    
Python      4545        74          4619  
```

## 🔄 CI/CD

GitHub Actions
```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go mod download

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: v1.64.4
      - run: go test -race -coverprofile=coverage.txt -covermode=atomic ./...
      - run: go build -o parser cmd/main.go
      
```

## 🤝 Участие в разработке

1. Форкните репозиторий
2. Создайте ветку для фичи (`git checkout -b feature/amazing-feature`)
3. Запустите тесты и линтер (`task test && task lint`)
4. Сделайте коммит (`git commit -m 'Add amazing feature'`)
5. Запушьте ветку (`git push origin feature/amazing-feature`)
6. Откройте **Pull Reques**

[![CI](https://github.com/Reactivity512/hhparser/actions/workflows/main.yml/badge.svg)](https://github.com/Reactivity512/hhparser/actions/workflows/main.yml)
[![Coverage](https://github.com/Reactivity512/hhparser/badge.svg)](https://coveralls.io/github/Reactivity512/hhparser)
[![Go Report Card](https://goreportcard.com/badge/github.com/Reactivity512/hhparser)](https://goreportcard.com/report/github.com/Reactivity512/hhparser)
