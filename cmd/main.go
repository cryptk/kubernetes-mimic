package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kapi "github.com/cryptk/kubernetes-mimic/internal/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	initConfig()

	// kubernetesAPI := kapi.NewKubernetesApi(viper.GetString("namespace"), viper.GetString("certSecretName"), viper.GetString("mirrorsConfigMapName"))

	kubernetesAPI := kapi.KubernetesAPI{
		Namespace:            viper.GetString("namespace"),
		CertSecretName:       viper.GetString("certSecretName"),
		MirrorsConfigMapName: viper.GetString("mirrorsConfigMapName"),
	}
	kubernetesAPI.Start()

	x509KeyPair := kubernetesAPI.GetCertificates()
	mirrorsConfig := kubernetesAPI.GetMirrorsConfig()

	if viper.GetBool("watchmirrorsconfig") {
		err := kubernetesAPI.StartWatchMirrorsConfig()
		if err != nil {
			log.WithError(err).Error("Failed to watch configmaps")
		}
	}

	whsrvr := &WebhookServer{
		mirrorsConfig: mirrorsConfig,
		server: &http.Server{
			Addr:      fmt.Sprintf("%s:%d", viper.GetString("listenhost"), viper.GetInt("listenport")),
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{x509KeyPair}},
			// Handler:           nil,
			// ReadTimeout:       0,
			// ReadHeaderTimeout: 0,
			// WriteTimeout:      0,
			// IdleTimeout:       0,
			// MaxHeaderBytes:    0,
			// TLSNextProto:      map[string]func(*http.Server, *tls.Conn, http.Handler){},
			// ConnState: func(net.Conn, http.ConnState) {
			// },
			// ErrorLog: &log.Logger{},
			// BaseContext: func(net.Listener) context.Context {
			// },
			// ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			// },
		},
	}

	kubernetesAPI.NewMirrorsCallback = whsrvr.updateMirrorsConfig

	log.Info("Webhook Server initialized")

	http.HandleFunc("/ping", pong)
	http.HandleFunc("/mutate", whsrvr.serve)

	go func() {
		if err := whsrvr.server.ListenAndServeTLS("", ""); err != nil {
			if err != http.ErrServerClosed {
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
	err := whsrvr.server.Shutdown(context.Background())
	if err != nil {
		log.WithError(err).Panic("Web server failed to shut down... lets PANIC instead!")
	}
}

func pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}
