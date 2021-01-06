package common

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

//GenFileName add comment
func GenFileName(filename string) string {
	saveFile := filename + "_" + time.Now().Format("2006/1/2 15:04:05") + ".csv"
	re := regexp.MustCompile(" |/|:")
	saveFile = re.ReplaceAllString(saveFile, "_")
	return filepath.Join(BaseDir, saveFile)
}

//Write2CSV add comment
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
	_, _ = w.WriteString("\xEF\xBB\xBF")
	return filename, Convert(bytes.NewReader(jsonObject), w)
}

//GenDicFromFile add comment
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
		one = TrimLR.ReplaceAllString(one, "")
		result = append(result, one)
	}
	return result
}
