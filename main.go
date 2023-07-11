package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/sclevine/agouti"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
			file := frame.File[len(path.Dir(os.Args[0]))+1:]
			line := frame.Line
			return "", fmt.Sprintf("%s:%d", file, line)
		},
	})
}

func main() {
	const homepage = "https://%s.zetfix-online.net/"
	delay := time.Second

	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"start-maximized",
			"disable-infobars",
			"--no-sandbox",
			"--app",
			"--disable-3d-apis",
		}),
	)

	if err := driver.Start(); err != nil {
		logrus.Panicf("Invalid driver starting: %s", err)
	}
	defer func() {
		if err := driver.Stop(); err != nil {
			logrus.Panicf("Invalid driver stoping: %s", err)
		}
	}()

	page, err := driver.NewPage()
	if err != nil {
		logrus.Panicf("Invalid page: %s", err)
	}

	go func() {
		if err := page.Navigate(fmt.Sprintf(homepage, strings.ToLower(time.Now().Format("_2Jan")))); err != nil {
			logrus.Panicf("Invalid navigation: %s", err)
		}
	}()

	for {
		if count, _ := page.WindowCount(); count == 0 {
			break
		}

		time.Sleep(delay)
	}
}
