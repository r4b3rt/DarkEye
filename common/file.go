package common

import (
	"bufio"
	"os"
	"strings"
)

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
