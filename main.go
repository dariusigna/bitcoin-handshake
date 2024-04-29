package main

import (
	"github.com/dariusigna/bitcoin-handshake/config"
	"github.com/dariusigna/bitcoin-handshake/server"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load("test.env")
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	var cfg config.Server
	envconfig.MustProcess("", &cfg)
	switch cfg.LogLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	server, err := server.New(cfg.Network, "Darius Awesome Node", cfg.TargetNodeAddress)
	if err != nil {
		logrus.Panic("failed to create the server:", err)
	}

	go server.Start()
	<-quit

	server.Stop()
	server.WaitForShutdown()
}
