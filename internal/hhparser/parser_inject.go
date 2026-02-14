package hhparser

import (
	"strconv"
	"strings"
)

func injectSearchCounts(content string) (int, error) {
	var start = strings.Index(content, "searchCounts") + 15
	var bracket uint32 = 1

	for end := start; end < len(content); end++ {
		switch content[end] {
		case 123:
			bracket++
		case 125:
			bracket--
		}
		if bracket == 0 {
			result := []byte(content[start:end])

			return injectCount(string(result))
		}
	}

	return 0, nil
}

func injectCount(content string) (int, error) {
	var result string
	var start = strings.Index(content, "\"value\":") + 8

	for end := start; end < len(content); end++ {
		if content[end] == ',' {
			result = content[start:end]
			break
		}
	}

	count, err := strconv.Atoi(result)
	if err != nil {
		return 0, ErrVacancyNotInteger
	}

	return count, nil
}
