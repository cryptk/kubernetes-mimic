package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetEnvPrefix("mimiccfg")
	viper.AutomaticEnv()

	viper.SetDefault("listenport", 8443)
	viper.SetDefault("listenhost", "0.0.0.0")
	viper.SetDefault("certSecretName", "mimic-certs")
	viper.SetDefault("mirrorsConfigMapName", "mimic-mirrors")
	viper.SetDefault("watchmirrorsconfig", true)
	viper.SetDefault("namespace", "")
	viper.SetDefault("loglevel", "info")
	viper.SetDefault("logformat", "text")

	switch viper.GetString("logformat") {
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}

	switch viper.GetString("loglevel") {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.SetOutput(os.Stdout)
}
