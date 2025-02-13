package main

import (
	"log"

	logrus "github.com/sirupsen/logrus"
)

func loggingwithlog() {
	log.Println("This is my first log message with log package")
	log.Fatal("This is the fatal log message")
}

func loggingwithlogrus() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.WithFields(logrus.Fields{"Status Code": 200}).Info("Successful")
	logrus.WithFields(logrus.Fields{"Status Code": 404}).Warn("Not Found")
	logrus.Fatal("This is the fatal log message")
}

func main() {
	loggingwithlogrus()
	// loggingwithlog()
}
