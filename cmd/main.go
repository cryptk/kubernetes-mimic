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

	k8s "github.com/cryptk/kubernetes-mimic/internal/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	initConfig()

	kclient, err := k8s.NewClient(
		viper.GetString("namespace"),
		viper.GetString("certSecretName"),
		viper.GetString("mirrorsConfigMapName"),
		viper.GetBool("watchmirrorsconfig"),
	)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize the kubernetes client")
	}

	x509KeyPair := kclient.GetCertificates()
	mirrorsConfig := kclient.GetMirrorsConfig()

	whsrvr := &WebhookServer{
		mirrorsConfig: mirrorsConfig,
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", viper.GetString("listenhost"), viper.GetInt("listenport")),
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{x509KeyPair},
				MinVersion:   tls.VersionTLS12,
			},
		},
	}

	kclient.SetMirrorsCallback(whsrvr.updateMirrorsConfig)

	log.Info("Webhook Server initialized")

	http.HandleFunc("/ping", pong)
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

func pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}
