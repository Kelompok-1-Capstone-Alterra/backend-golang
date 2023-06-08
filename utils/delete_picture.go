package utils

import "os"

func Delete_picture(filePath string) error {
	err := os.Remove("assets/images/"+filePath)
	if err != nil {
		return err
	}
	return nil
}
