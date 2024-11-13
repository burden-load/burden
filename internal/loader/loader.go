package loader

import (
	"burden/pkg/model"
	"encoding/json"
	"fmt"
	"os"
)

func LoadCollection(filePath string) ([]model.Request, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var collection model.PostmanCollection
	err = json.Unmarshal(data, &collection)
	if err != nil {
		return nil, fmt.Errorf("error loading collection: %w", err)
	}

	// Преобразуем данные коллекции в массив `model.Request`
	var requests []model.Request
	for _, item := range collection.Item {
		for _, subItem := range item.Item {
			var body string
			if subItem.Request.Body != nil { // Проверяем, что Body не nil
				body = subItem.Request.Body.Raw
			}
			requests = append(requests, model.Request{
				Method: subItem.Request.Method,
				URL:    subItem.Request.URL,
				Body:   body,
			})
		}
	}

	return requests, nil
}
