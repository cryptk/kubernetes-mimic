package kubernetes

import (
	"context"
	"crypto/tls"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	namespace            string
	certSecretName       string
	mirrorsConfigMapName string
	mirrorsCallback      func(*map[string]string)

	clientset        *kubernetes.Clientset
	certificates     tls.Certificate
	mirrorsConfigmap *v1.ConfigMap
}

func NewClient(namespace string, certSecretName string, mirrorsConfigMapName string) (*Client, error) {
	if namespace == "" {
		namespace = getNamespace()
		log.WithField("namespace", namespace).Debug("Namespace not specified, Namespace discovered from environment")
	}

	client := &Client{
		namespace:            namespace,
		certSecretName:       certSecretName,
		mirrorsConfigMapName: mirrorsConfigMapName,
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	client.clientset = clientset
	log.Info("Kubernetes ClientSet initialized")

	certs, err := client.fetchCertificates()
	if err != nil {
		log.WithError(err).Fatal("Unable to fetch SSL Certificates from Kubernetes API")
	}
	log.Info("SSL Certificates retrieved from Kubernetes Secret")
	client.certificates = certs

	mirrorsConfigMap, err := client.fetchMirrorsConfig()
	if err != nil {
		log.WithError(err).Error("Failed to fetch Mirrors Config from Kubernetes API")
	}
	log.Info("Mirrors Config retrieved from Kubernetes ConfigMap")
	client.mirrorsConfigmap = mirrorsConfigMap

	return client, nil
}

func (client *Client) SetMirrorsCallback(cb func(*map[string]string)) {
	client.mirrorsCallback = cb
}

func (client *Client) GetCertificates() tls.Certificate {
	return client.certificates
}

func (client *Client) fetchCertificates() (tls.Certificate, error) {
	certs, err := client.clientset.CoreV1().Secrets(client.namespace).Get(context.TODO(),
		client.certSecretName,
		metav1.GetOptions{})
	if err != nil {
		return tls.Certificate{}, err
	}

	pair, err := tls.X509KeyPair(certs.Data["cert.pem"], certs.Data["key.pem"])
	if err != nil {
		return tls.Certificate{}, err
	}

	return pair, nil
}

func (client *Client) GetMirrorsConfig() *map[string]string {
	return &client.mirrorsConfigmap.Data
}

func (client *Client) fetchMirrorsConfig() (*v1.ConfigMap, error) {
	configmap, err := client.clientset.CoreV1().ConfigMaps(client.namespace).Get(context.TODO(),
		client.mirrorsConfigMapName,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, err
	}

	return configmap, nil
}

func (client *Client) StartWatchMirrorsConfig() error {
	selector := fields.OneTermEqualSelector("metadata.name", client.mirrorsConfigMapName)
	listOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}
	watcher, err := client.clientset.CoreV1().ConfigMaps(client.namespace).Watch(context.TODO(), listOptions)
	if err != nil {
		return err
	}
	go client.watchMirrorsConfig(watcher.ResultChan())
	return nil
}

func (client *Client) watchMirrorsConfig(c <-chan watch.Event) {
	for event := range c {
		configmap, ok := event.Object.(*v1.ConfigMap)
		if !ok {
			log.Error("Received an event for something other than a ConfigMap")
		}
		log.WithFields(log.Fields{
			"eventType":          event.Type,
			"configmapName":      configmap.Name,
			"configmapNamespace": configmap.Namespace,
		}).Debug("Event received")
		client.mirrorsConfigmap = configmap

		client.mirrorsCallback(&client.mirrorsConfigmap.Data)
	}
}
