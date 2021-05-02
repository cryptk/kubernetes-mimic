[![Go Report Card](https://goreportcard.com/badge/github.com/cryptk/kubernetes-mimic)](https://goreportcard.com/report/github.com/cryptk/kubernetes-mimic) [![License](https://img.shields.io/github/license/cryptk/kubernetes-mimic)](LICENSE)

# kubernetes-mimic
Kubernetes Mimic is a Mutating Webhook that will watch for pod creation and update events in a Kubernetes cluster and automatically adjust their container images to pull from an image mirror as opposed to upstream servers.

It aims to make using an internal image mirror simple and hassle-free.  It can even automatically auto-discover configured repository mirrors from [Harbor](https://goharbor.io/).

This project is still in it's early stages, and as such, documentation is less than ideal.

## Integrations

Currently Mimic can only integrate with [Harbor](https://goharbor.io/) for autodiscovery of [Proxy Cache](https://goharbor.io/docs/2.1.0/administration/configure-proxy-cache/) projects.  When this integration is enabled, Mimic will watch for pods being created with an Image that is pulled from a source that is also available as a public Proxy Cache in Harbor and will update the Image source as necessary to pull the image from the Harbor cache instead.

There are plans to also support Artifactory.  Any other desired integrations should be requested by [opening an issue](cryptk/kubernetes-mimic/issues)

## Image building

Mimic can be built into a docker image using all of the normal techniques.  Assuming you are wanting a Linux AMD64 Docker image, you can build it with the following command from within the base of the repository.

`docker build -t mimic:latest .`

## Deployment

Currently the deployment is manual, and there are example manifests in the [manifests folder](deploy/manifests).  As the project matures, the deployment of Mimic will be handled via Helm (GH-14).

The process is as follows:

1. Create a Kubernetes Namespace to deploy Mimic into

`kubectl apply -f ./deploy/manifests/namespace`

2. Generate SSL certificates used for communication between the kubernetes API layer and the webhook:

`./deploy/scripts/webhook-create-signed-cert.sh --service mimic --secret mimic-certs --namespace mimic`

3. Add the CA Bundle for the generated certificate to the mutating webhook configuration

`./deploy/scripts/webhook-patch-ca-bundle.sh ./deploy/manifests/templates/mutatingwebhookconfiguration.yaml ./deploy/manifests/mutatingwebhookconfiguration-cabundle.yaml`

4. Deploy the rest of the required resources

`kubectl apply -f ./deploy/manifests`

## Configuration

Mimic accepts it's configuration via environment variables.

| Variable | Default | Description |
|----------|---------|-------------|
| MIMIC_LISTENPORT | 8443 | What port should the Mimic API server listen on |
| MIMIC_LISTENHOST | "0.0.0.0" | What host should the Mimic API server listen on |
| MIMIC_LOGLEVEL | "info" | What level should mimic log at.  Valid options are trace, debug, info, warn, error, fatal and panic |
| MIMIC_LOGFORMAT | "text" | What format should the logs be rendered in.  Valid options are text, json |
| MIMIC_CERTIFICATE_SOURCE | kubernetes | Where to load TLS certificates from.  Currently the only valid option is "kubernetes" which will load the TLS certificates from a kubernetes secret |
| MIMIC_WATCHMIRRORS | true | Should sources be watched for updates and new mirrors automatically.  Sources that support watching can also be toggled individually |
| MIMIC_KUBERNETES_ENABLED | true | Should the Kubernetes integration be enabled |
| MIMIC_KUBERNETES_NAMESPACE | "" | What Namespace should Mimic look for it's resources in.  If this is not specified, Mimic will attempt to autodiscover what namespace it is in automatically |
| MIMIC_KUBERNETES_CERTSECRET | "mimic-certs" | The name of the Kubernetes Secret that holds the TLS certificates for the webhook server |
| MIMIC_KUBERNETES_CONFIGMAP | "mimic-mirrors" | The name of the Kubernetes ConfigMap that holds the mirror configuration.  Please see the [example configmap](./deploy/manifests/configmap.yaml) |
| MIMIC_KUBERNETES_WATCH | true | Should Mimic watch the ConfigMap to automatically pull in changes as opposed to requiring an application restart to load new changes |
| MIMIC_HARBOR_ENABLED | false | Should Mimic attempt to auto-discover docker mirrors configured within a Harbor installation |
| MIMIC_HARBOR_API_HOST | "" | Hostname that Mimic should use for communications with the Harbor API |
| MIMIC_HARBOR_REGISTRYURL | "" | Hostname that Harbor serves it's repository mirrors from.  If this is left blank, Mimic will attempt to autodiscover this from the Harbor API |
| MIMIC_HARBOR_ROBOT_USERNAME | "" | Robot account username from Harbor.  Needed to autodiscover the Registry URL from the Harbor API |
| MIMIC_HARBOR_ROBOT_PASSWORD | "" | Robot account password from Harbor.  Needed to autodiscover the Registry URL from the Harbor API |
