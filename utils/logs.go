package utils

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func ShowErrorLogs(errData error) {
	if errData != nil {
		file, line := Caller(2)
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.Errorf("%s[%d] %s", filepath.Base(file), line, errData.Error())
	}
}

func ShowInfoLogs(infoWantToLog string) {
	if infoWantToLog != "" {
		file, line := Caller(2)
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.Infof("%s[%d] %s", filepath.Base(file), line, infoWantToLog)
	}
}

func LogFull[T any](object T) string {
	e, err := json.Marshal(object)
	if err != nil {
		logrus.Errorf("Error marshal object %v - Error %v", object, err)
		return ""
	}

	jsonString := strings.ReplaceAll(string(e), "\"", " ")
	return fmt.Sprintf("%v", jsonString)
}

func Caller(level int) (string, int) {
	_, file, line, _ := runtime.Caller(level)
	return file, line
}
