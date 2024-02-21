package utils

import (
	"encoding/csv"
	"os"
)

func ReadCSV(filename string) ([][]string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	cr := csv.NewReader(file)
	return cr.ReadAll()
}
