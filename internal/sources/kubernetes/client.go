package kubernetes

import (
	"context"
	"crypto/tls"
	"fmt"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Client handles interactions with the kubernetes API.
type Client struct {
	namespace        string
	certSecret       string
	mirrorsConfigMap string

	clientset    *kubernetes.Clientset
	certificates *tls.Certificate
	mirrors      *v1.ConfigMap
}

// New configures a new Kubernetes Client.
func New(namespace string, certSecretName string, mirrorsConfigMap string) (*Client, error) {
	if namespace == "" {
		namespace = getNamespace()
		log.WithField("namespace", namespace).Debug("Namespace not specified, Namespace discovered from environment")
	}

	client := &Client{
		namespace:        namespace,
		certSecret:       certSecretName,
		mirrorsConfigMap: mirrorsConfigMap,
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to generate InClusterConfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	client.clientset = clientset

	log.Info("Kubernetes ClientSet initialized")

	mirrors, err := client.fetchMirrorsConfig()
	if err != nil {
		log.WithError(err).Error("Failed to fetch Mirrors Config from Kubernetes API")
	}

	client.mirrors = mirrors

	log.Info("Mirrors Config retrieved from Kubernetes ConfigMap")

	return client, nil
}

// Certificates fetches TLS certificates from a kubernetes secret.
func (client *Client) Certificates() (tls.Certificate, error) {
	if client.certificates != nil {
		return *client.certificates, nil
	}

	certs, err := client.clientset.CoreV1().Secrets(client.namespace).Get(context.TODO(),
		client.certSecret,
		metav1.GetOptions{})
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to retrieve TLS certificates from Kubernetes API: %w", err)
	}

	pair, err := tls.X509KeyPair(certs.Data["cert.pem"], certs.Data["key.pem"])
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to generate X509 Key Pair from fetched certificates: %w", err)
	}

	client.certificates = &pair

	return *client.certificates, nil
}

// Mirrors returns all mirrors which have been retrieved from the Kubernetes API.
func (client *Client) Mirrors() map[string]string {
	return client.mirrors.Data
}

func (client *Client) fetchMirrorsConfig() (*v1.ConfigMap, error) {
	configmap, err := client.clientset.CoreV1().ConfigMaps(client.namespace).Get(context.TODO(),
		client.mirrorsConfigMap,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve configmap from Kubernetes API: %w", err)
	}

	return configmap, nil
}

// WatchMirrors initiates a watch on the Kubernetes ConfigMap to be notified of changes without needing an application restart.
func (client *Client) WatchMirrors(cb func()) error {
	selector := fields.OneTermEqualSelector("metadata.name", client.mirrorsConfigMap)
	listOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}

	watcher, err := client.clientset.CoreV1().ConfigMaps(client.namespace).Watch(context.TODO(), listOptions)
	if err != nil {
		return fmt.Errorf("failed to create watch on kubernetes configmap: %w", err)
	}

	// Handle the events that come in from Kubernetes in a goroutine.
	go func() {
		for event := range watcher.ResultChan() {
			configmap, ok := event.Object.(*v1.ConfigMap)
			if !ok {
				log.Warn("Received an event for something other than a ConfigMap")
			}

			log.WithFields(log.Fields{
				"eventType": event.Type,
				"configmap": configmap.Name,
				"namespace": configmap.Namespace,
			}).Debug("Event received")

			client.mirrors = configmap

			cb()
		}
	}()

	return nil
}

// Stop gracefully ends the watch, typically in preparation for an application halt.
func (client *Client) Stop() error {
	return nil
}
