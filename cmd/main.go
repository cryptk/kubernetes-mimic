package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cryptk/kubernetes-mimic/internal/config"
	"github.com/cryptk/kubernetes-mimic/internal/sources"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	config.Defaults()

	sources, err := sources.New()
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize mirror sources")
	}

	whsrvr := &webhookServer{
		mirrorsConfig: sources.Mirrors(),
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", viper.GetString("listenhost"), viper.GetInt("listenport")),
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{sources.Certificates()},
				MinVersion:   tls.VersionTLS12,
			},
		},
	}

	if viper.GetBool("watchmirrors") {
		sources.Watch(whsrvr.updateMirrors)
	}

	log.Info("Webhook Server initialized")

	http.HandleFunc("/mutate", whsrvr.serve)

	go func() {
		if err := whsrvr.server.ListenAndServeTLS("", ""); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.WithError(err).Fatal("failed to listen and serve webhook server")
			}
		}

		log.WithField("port", viper.GetInt("listenport")).Info("Webhook Server listening")
	}()

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Info("Got OS shutdown signal, shutting down webhook server gracefully...")

	err = whsrvr.server.Shutdown(context.Background())
	if err != nil {
		log.WithError(err).Fatal("web server failed to shut down... lets PANIC instead!")
	}
}
