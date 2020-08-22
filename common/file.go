package common

import "os"

func SaveFile(content, filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, _ = f.Write([]byte(content+"\n"))
	return nil
}
