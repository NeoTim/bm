package main

import (
	"fmt"
	"os"
)

// Debug prints debug messages.
func Debug(arg ... interface{}) {
	if EnableDebug {
		fmt.Println(arg...)
	}
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
