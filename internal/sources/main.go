package sources

import (
	"crypto/tls"
	"fmt"

	k8sclient "github.com/cryptk/kubernetes-mimic/internal/sources/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Certificator interface {
	Certificates() (tls.Certificate, error)
}

type Mirrorer interface {
	Mirrors() map[string]string
}

type WatchingMirrorer interface {
	Mirrorer
	WatchMirrors(func()) error
	Stop() error
}

type Sources struct {
	mirrors      map[string]Mirrorer
	watchers     map[string]WatchingMirrorer
	certificator Certificator

	newMirrorsCB func(map[string]string)
}

func New() (sources *Sources, err error) {

	sources = &Sources{
		mirrors:  make(map[string]Mirrorer),
		watchers: make(map[string]WatchingMirrorer),
	}

	if validateConfigSet("kubernetes_enabled", []string{"kubernetes_namespace", "kubernetes_certsecret", "kubernetes_configmap"}) {
		k8s, err := k8sclient.New(
			viper.GetString("kubernetes_namespace"),
			viper.GetString("kubernetes_certsecret"),
			viper.GetString("kubernetes_configmap"),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Kubernetes Client: %w", err)
		}
		sources.AddMirrorer("kubernetes", k8s)
		if viper.GetBool("kubernetes_watch") {
			sources.AddWatcher("kubernetes", k8s)
		}
		if viper.GetString("certificate_source") == "kubernetes" {
			sources.certificator = k8s
		}
	}

	return sources, nil
}

func (sources *Sources) AddMirrorer(name string, source Mirrorer) {
	sources.mirrors[name] = source
}

func (sources *Sources) AddWatcher(name string, source WatchingMirrorer) {
	sources.watchers[name] = source
}

func (sources *Sources) Certificates() tls.Certificate {
	if sources.certificator == nil {
		log.Fatal("Unable to retrieve certificates, no certificate source configured")
	}
	certs, err := sources.certificator.Certificates()
	if err != nil {
		log.WithError(err).Fatal("Unable to fetch certificates")
	}
	log.Debug("TLS certificates retrieved")
	return certs
}

func (sources *Sources) Mirrors() (mirrors map[string]string) {
	mirrors = make(map[string]string)
	for _, mirror := range sources.mirrors {
		for k, v := range mirror.Mirrors() {
			mirrors[k] = v
		}
	}
	return mirrors
}

func (sources *Sources) updateWebhookMirrors() {
	sources.newMirrorsCB(sources.Mirrors())
}

func (sources *Sources) Watch(cb func(map[string]string)) {
	sources.newMirrorsCB = cb
	for name, source := range sources.watchers {
		if err := source.WatchMirrors(sources.updateWebhookMirrors); err != nil {
			log.WithError(err).Errorf("failed to initialize watch on mirror source %s", name)
		}
		log.Infof("Watch on source %s started", name)
	}
}

func (sources *Sources) Stop() {
	for name, source := range sources.watchers {
		log.Infof("Stopping watch on source %s", name)
		source.Stop()
	}
}
