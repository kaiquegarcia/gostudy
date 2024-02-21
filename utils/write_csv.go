package utils

import (
	"encoding/csv"
	"os"
)

func WriteCSV(filename string, records [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	cw := csv.NewWriter(file)
	return cw.WriteAll(records)
}
