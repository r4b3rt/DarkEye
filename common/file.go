package common

import (
	"bufio"
	"encoding/csv"
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

func CreateCSV(fileName string, cols []string) (*csv.Writer, *os.File, string, error) {
	fileName = GenFileName(fileName)
	f, err := os.Create(fileName)
	if err != nil {
		return nil, nil, "", err
	}
	_, _ = f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)
	_ = w.Write(cols)
	return w, f, fileName, nil
}

func GenArraryFromFile(filename string) []string{
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
