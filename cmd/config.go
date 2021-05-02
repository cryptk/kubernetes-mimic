package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func configDefaults() {
	viper.SetEnvPrefix("mimic")
	viper.AutomaticEnv()

	// Kubernetes source configuration options
	viper.SetDefault("kubernetes_enabled", true)
	viper.SetDefault("kubernetes_namespace", "") // Leaving this blank will autodiscover from the Kubernetes environment
	viper.SetDefault("kubernetes_certsecret", "mimic-certs")
	viper.SetDefault("kubernetes_configmap", "mimic-mirrors")
	viper.SetDefault("kubernetes_watch", true)

	// Harbor source configuration options
	viper.SetDefault("enable_harbor", false)
	viper.SetDefault("harbor_api_host", "")
	viper.SetDefault("harbor_harbor_registryurl", "") // Leaving this blank will autodiscover from the Harbor API
	viper.SetDefault("harbor_robot_username", "")
	viper.SetDefault("harbor_robot_password", "")

	// Certificate source
	viper.SetDefault("certificate_source", "kubernetes")

	// Web server options
	viper.SetDefault("listenport", 8443)
	viper.SetDefault("listenhost", "0.0.0.0")

	// Generic options
	viper.SetDefault("loglevel", "info")
	viper.SetDefault("logformat", "text")
	viper.SetDefault("watchmirrors", true)

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
