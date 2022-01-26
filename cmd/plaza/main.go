package main

import (
	"github.com/kjubybot/nightwave-plaza-bot/internal/plaza"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := plaza.Init(); err != nil {
		logrus.Fatal(err)
	}

	if err := plaza.Run(); err != nil {
		logrus.Fatal(err)
	}
}
