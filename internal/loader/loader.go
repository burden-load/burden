package loader

import (
	"os"
)

func LoadCollection(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
}
