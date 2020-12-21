package common

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GenFileName(filename string) string {
	saveFile := filename + "_" + time.Now().Format("2006/1/2 15:04:05")
	saveFile = strings.Replace(saveFile, " ", "_", -1)
	saveFile = strings.Replace(saveFile, ":", "_", -1)
	saveFile = strings.Replace(saveFile, "/", "_", -1) + ".csv"
	saveFile = filepath.Join(BaseDir, saveFile)
	return saveFile
}

func Write2CSV(filename string, jsonObject []byte) (string, error) {
	filename = GenFileName(filename)
	if jsonObject == nil {
		return filename, nil
	}
	w, err := os.Create(filename)
	if err != nil {
		return filename, err
	}
	defer w.Close()
	return filename, Convert(bytes.NewReader(jsonObject), w)
}

func GenDicFromFile(filename string) []string {
	if filename == "" {
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()
	result := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		one := scanner.Text()
		if strings.HasPrefix(one, "#") {
			continue
		}
		one = strings.TrimSpace(one)
		one = strings.Trim(one, "\r\n")
		result = append(result, one)
	}
	return result
}
