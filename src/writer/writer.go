package writer

import (
	"os"
)

func WriteFile(path string, content []byte) error {

	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(content)

	return err
}

func CreateDirectory(path string) error {
	return os.Mkdir(path, os.ModePerm)
}
