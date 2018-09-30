package util

import (
	"fmt"
	"os"
	"strings"
)

func Exists(path string) bool {
	// (from: https://stackoverflow.com/a/12518877/2885946)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true;
	} else {
		return false;
	}
}

// (from: https://siongui.github.io/2016/04/23/go-read-yes-no-from-console/)
func Ask4confirm(message string) bool {
	var s string

	fmt.Printf("%s[y/N]: ", message)
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.ToLower(strings.TrimSpace(s))

	return s == "y" || s == "yes"
}

func GetEnv() map[string]string {
	// (from: https://gobyexample.com/environment-variables)
	env := map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env[pair[0]] = pair[1]
	}
	return env
}
