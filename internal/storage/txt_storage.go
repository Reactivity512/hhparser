package storage

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func saveTXT(stats Statistics, directory string) error {
	file, err := os.Create(fmt.Sprintf("%s/stats_%s.txt",
		directory,
		stats.Date.Format("2006-01-02")))
	if err != nil {
		return err
	}
	defer file.Close()

	w := tabwriter.NewWriter(file, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "СТАТИСТИКА ВАКАНСИЙ\n")
	fmt.Fprintf(w, "Дата: %s\n\n", stats.Date.Format("02.01.2006"))

	fmt.Fprint(w, "Технология\t")
	for _, city := range stats.Cities {
		fmt.Fprintf(w, "%s\t", city.Name)
	}
	fmt.Fprintln(w, "ВСЕГО")

	fmt.Fprint(w, "----------\t")
	for range stats.Cities {
		fmt.Fprint(w, "----------\t")
	}
	fmt.Fprintln(w, "----------")

	for _, tech := range stats.Technologies {
		fmt.Fprintf(w, "%s\t", tech.Name)
		for _, city := range stats.Cities {
			fmt.Fprintf(w, "%d\t", city.Vacancies[tech.Name])
		}
		fmt.Fprintf(w, "%d\t", stats.Summary[tech.Name])
		fmt.Fprintln(w)
	}

	return w.Flush()
}
