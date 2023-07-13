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
	const hostname = "https://%s.zetfix-online.net/"
	delay := time.Second

	homepage := fmt.Sprintf(hostname, strings.ToLower(time.Now().Format("_2Jan")))
	logrus.Infof("Homepage: %s", homepage)
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"start-maximized",
			"disable-infobars",
			"--no-sandbox",
			"--app",
			"--disable-3d-apis",
		}),
	)
	logrus.Info("*** Driver: configured ***")

	logrus.Info("Driver: starting...")
	if err := driver.Start(); err != nil {
		logrus.Panicf("Invalid driver starting: %s", err)
	}
	defer func() {
		logrus.Info("Driver: stopping...")
		if err := driver.Stop(); err != nil {
			logrus.Panicf("Invalid driver stopping: %s", err)
		}
		logrus.Info("*** Driver: stopped ***")
	}()
	logrus.Info("*** Driver: started ***")

	logrus.Info("Page: creation...")
	page, err := driver.NewPage()
	if err != nil {
		logrus.Panicf("Invalid page: %s", err)
	}
	logrus.Info("*** Page: created ***")

	go func() {
		logrus.Info("Page: loading...")
		if err := page.Navigate(homepage); err != nil {
			logrus.Panicf("Invalid navigation: %s", err)
		}
		logrus.Info("*** Page: loaded ***")
	}()

	for {
		if count, _ := page.WindowCount(); count == 0 {
			logrus.Info("*** All windows closed ***")
			break
		}

		time.Sleep(delay)
	}
}
