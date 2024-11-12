package loader

import (
	"burden/pkg/model"
	"encoding/json"
	"io"
	"log"
	"os"
)

func LoadCollection(filePath string) ([]model.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var requests []model.Request
	if err := json.Unmarshal(byteValue, &requests); err != nil {
		return nil, err
	}

	log.Printf("Успешно загружено %d запросов из коллекции", len(requests))

	return requests, nil
}
