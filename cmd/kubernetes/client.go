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

// type KubernetesAPI interface {
// 	GetCertificates() tls.Certificate
// 	GetMirrorsConfig(mirrorsConfigMapName string, namespace string)
// 	StartWatchMirrorsConfig() error
// }

type KubernetesAPI struct {
	Namespace            string
	CertSecretName       string
	MirrorsConfigMapName string
	NewMirrorsCallback   func(*map[string]string)

	clientset        *kubernetes.Clientset
	certificates     tls.Certificate
	mirrorsConfigmap *v1.ConfigMap
}

func (kapi *KubernetesAPI) Start() {
	if kapi.Namespace == "" {
		kapi.Namespace = getNamespace()
		log.WithField("namespace", kapi.Namespace).Debug("Namespace not specified, Namespace discovered from environment")
	}

	// 	kapi := kubernetesAPI{
	// 		namespace:            namespace,
	// 		certSecretName:       certSecretName,
	// 		mirrorsConfigMapName: mirrorsConfigMapName,
	// 	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.WithError(err).Fatal("Unable to initialize InClusterConfig for Kubernetes API Client")
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.WithError(err).Fatal("Failed to create Kubernetes ClientSet")
	}

	kapi.clientset = clientset
	log.Info("Kubernetes ClientSet initialized")

	certs, err := kapi.fetchCertificates()
	if err != nil {
		log.WithError(err).Fatal("Unable to fetch SSL Certificates from Kubernetes API")
	}
	log.Info("SSL Certificates retrieved from Kubernetes Secret")
	kapi.certificates = certs

	mirrorsConfigMap, err := kapi.fetchMirrorsConfig()
	if err != nil {
		log.WithError(err).Error("Failed to fetch Mirrors Config from Kubernetes API")
	}
	log.Info("Mirrors Config retrieved from Kubernetes ConfigMap")
	kapi.mirrorsConfigmap = mirrorsConfigMap

}

func (kapi *KubernetesAPI) GetCertificates() tls.Certificate {
	return kapi.certificates
}

func (kapi *KubernetesAPI) fetchCertificates() (tls.Certificate, error) {

	certs, err := kapi.clientset.CoreV1().Secrets(kapi.Namespace).Get(context.TODO(), kapi.CertSecretName, metav1.GetOptions{})
	if err != nil {
		return tls.Certificate{}, err
	}

	pair, err := tls.X509KeyPair(certs.Data["cert.pem"], certs.Data["key.pem"])
	if err != nil {
		return tls.Certificate{}, err
	}

	return pair, nil
}

func (kapi *KubernetesAPI) GetMirrorsConfig() *map[string]string {
	return &kapi.mirrorsConfigmap.Data
}

func (kapi *KubernetesAPI) fetchMirrorsConfig() (*v1.ConfigMap, error) {
	configmap, err := kapi.clientset.CoreV1().ConfigMaps(kapi.Namespace).Get(context.TODO(), kapi.MirrorsConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return configmap, nil
}

func (kapi *KubernetesAPI) StartWatchMirrorsConfig() error {
	selector := fields.OneTermEqualSelector("metadata.name", kapi.MirrorsConfigMapName)
	listOptions := metav1.ListOptions{
		FieldSelector: selector.String(),
	}
	watcher, err := kapi.clientset.CoreV1().ConfigMaps(kapi.Namespace).Watch(context.TODO(), listOptions)
	if err != nil {
		return err
	}
	go kapi.watchMirrorsConfig(watcher.ResultChan())
	return nil
}

func (kapi *KubernetesAPI) watchMirrorsConfig(c <-chan watch.Event) {
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
		kapi.mirrorsConfigmap = configmap

		kapi.NewMirrorsCallback(&kapi.mirrorsConfigmap.Data)
	}
}
