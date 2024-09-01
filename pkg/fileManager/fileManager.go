package fileManager

import (
	"fmt"
	"io/ioutil"
	"os"
)

const defaultFilePath = "./files/images/%s.jpeg"

func SaveFile(fileName string, imageData []byte) error {
	file, err := os.Create(fmt.Sprintf(defaultFilePath, fileName))
	if err != nil {
		err = fmt.Errorf("os.Create(...): %w", err)
		return err
	}

	err = ioutil.WriteFile(file.Name(), imageData, 0644)
	if err != nil {
		err = fmt.Errorf("ioutil.WriteFile(...): %w", err)
		return err
	}

	return nil
}

func DeleteFile(fileName string) error {
	err := os.Remove(fmt.Sprintf(defaultFilePath, fileName))
	if err != nil {
		err = fmt.Errorf("os.Remove(...): %w", err)
		return err
	}

	return nil
}
