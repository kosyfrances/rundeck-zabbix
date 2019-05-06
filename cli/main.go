package main

import (
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/kosyfrances/rundeck-zabbix/cli/cmd"
)

func main() {
	logFile := path.Join("/tmp", "rundeck-zabbix.log")

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()
	mw := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(mw)

	// execute cli
	cmd.Execute()
}
