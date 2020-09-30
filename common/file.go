package common

import (
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
