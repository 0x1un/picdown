package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

func CreateDirectory(dir string) {
	if dir == "" {
		return
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0644)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}
