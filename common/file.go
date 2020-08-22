package common

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveFile(content, filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, _ = f.Write([]byte(content+"\n"))
	return nil
}

func GenFileName(filename string) string {
	saveFile := time.Now().Format("2006/1/2 15:04:05")
	saveFile = strings.Replace(saveFile, " ", "_", -1)
	saveFile = strings.Replace(saveFile, ":", "_", -1)
	saveFile = strings.Replace(saveFile, "/", "_", -1) + ".csv"
	saveFile = filepath.Join(BaseDir, saveFile)
	return saveFile
}
