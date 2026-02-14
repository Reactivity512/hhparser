package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

func saveJSON(stats Statistics, directory string) error {
	file, err := os.Create(fmt.Sprintf("%s/%s.json",
		directory,
		stats.Date.Format("2006-01-02")))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(stats)
}
