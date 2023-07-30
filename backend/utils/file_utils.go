package utils

import (
	"io"
	"os"
)

// saveFile saves the uploaded file to the specified path.
func SaveFile(file io.Reader, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}

	return nil
}
