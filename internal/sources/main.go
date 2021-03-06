package sources

import (
	"crypto/tls"
	"fmt"

	harborclient "github.com/cryptk/kubernetes-mimic/internal/sources/harbor"
	k8sclient "github.com/cryptk/kubernetes-mimic/internal/sources/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Certificator provides a source to fetch TLS certificates.
type Certificator interface {
	Certificates() (tls.Certificate, error)
}

// Mirrorer ensures that a source can provide a mirror configuration.
type Mirrorer interface {
	Mirrors() map[string]string
}

// WatchingMirrorer is a source that can receive updates to it's mirror configuration without a restart of the application.
type WatchingMirrorer interface {
	Mirrorer
	WatchMirrors(func()) error
	Stop() error
}

// Sources maintains one or more sources of image mirrors.
type Sources struct {
	mirrors      map[string]Mirrorer
	watchers     map[string]WatchingMirrorer
	certificator Certificator

	newMirrorsCB func(map[string]string)
}

// New configures a new mimic Source engine.
func New() (sources *Sources, err error) {
	sources = &Sources{
		mirrors:  make(map[string]Mirrorer),
		watchers: make(map[string]WatchingMirrorer),
	}

	k8sEnabled, err := validateConfigSet("kubernetes_enabled", []string{"kubernetes_certsecret", "kubernetes_configmap"})
	if err != nil {
		return nil, err
	}

	if k8sEnabled {
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

	harborEnabled, err := validateConfigSet("harbor_enabled", []string{"harbor_api_host", "harbor_robot_username", "harbor_robot_password"})
	if err != nil {
		return nil, err
	}

	if harborEnabled {
		harbor, err := harborclient.New()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Harbor Client: %w", err)
		}

		sources.AddMirrorer("harbor", harbor)
	}

	return sources, nil
}

// AddMirrorer adds a new mirror source to the sources engine.
func (sources *Sources) AddMirrorer(name string, source Mirrorer) {
	sources.mirrors[name] = source
}

// AddWatcher adds a new watcher to the sources engine.
func (sources *Sources) AddWatcher(name string, source WatchingMirrorer) {
	sources.watchers[name] = source
}

// Certificates retrieves TLS certificates from the configured certificate source.
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

// Mirrors gathers the map of all mirrors from all configured sources.
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

// Watch enables the source watch to be notified of mirror configuration changes without restarting the application.
func (sources *Sources) Watch(cb func(map[string]string)) {
	sources.newMirrorsCB = cb
	for name, source := range sources.watchers {
		if err := source.WatchMirrors(sources.updateWebhookMirrors); err != nil {
			log.WithError(err).Errorf("failed to initialize watch on mirror source %s", name)
		}

		log.Infof("Watch on source %s started", name)
	}
}

// Stop instructs all WatchingMirrorers to stop their watch processes, usually in preparation for an application halt.
func (sources *Sources) Stop() {
	for name, source := range sources.watchers {
		log.Infof("Stopping watch on source %s", name)

		if err := source.Stop(); err != nil {
			log.WithError(err).Errorf("Failure while stopping watch on source %s", name)
		}
	}
}
