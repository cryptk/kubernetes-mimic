package harbor

import (
	"fmt"
	"net/url"

	apiclient "github.com/cryptk/kubernetes-mimic/internal/harbor/client"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/project"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Client handles interactions with the Harbor API.
type Client struct {
	harborAPI  *apiclient.HarborAPI
	harborAuth runtime.ClientAuthInfoWriter
	// transport *httptransport.Runtime

	registryURL string
	mirrors     map[string]string
}

// New configures a new Harbor Client.
func New() (*Client, error) {
	transport := httptransport.New(viper.GetString("harbor_api_host"), "/api/v2.0", []string{"http"})
	harborClient := apiclient.New(transport, strfmt.Default)
	harborAuth := httptransport.BasicAuth(viper.GetString("harbor_robot_username"), viper.GetString("harbor_robot_password"))

	client := &Client{
		harborAPI:  harborClient,
		harborAuth: harborAuth,
	}

	if err := client.getExternalURL(); err != nil {
		return nil, err
	}

	if err := client.populateMirrors(); err != nil {
		return nil, err
	}

	return client, nil
}

// Mirrors returns all mirrors which have been retrieved from the Kubernetes API.
func (client *Client) Mirrors() map[string]string {
	return client.mirrors
}

func (client *Client) getExternalURL() error {
	if registryURL := viper.GetString("harbor_registryurl"); registryURL != "" {
		log.Infof("Harbor RegistryURL overridden from environment variable, using %s", registryURL)
		client.registryURL = registryURL

		return nil
	}

	sysinfo, err := client.harborAPI.Systeminfo.GetSysteminfo(nil, client.harborAuth)
	if err != nil {
		return fmt.Errorf("failed to fetch system info from Harbor API: %w", err)
	}

	log.Infof("Harbor RegistryURL discovered from Harbor API, using %s", *sysinfo.Payload.RegistryURL)
	client.registryURL = *sysinfo.Payload.RegistryURL

	return nil
}

func (client *Client) populateMirrors() error {
	projects, err := client.harborAPI.Project.ListProjects(nil, nil)

	if err != nil {
		return fmt.Errorf("failure to fetch Projects from Harbor API: %w", err)
	}

	log.Debugf("Found %d public harbor Projects", len(projects.Payload))

	mirrors := make(map[string]string)

	for _, prj := range projects.Payload {
		if prj.RegistryID > 0 {
			params := project.NewGetProjectSummaryParams()
			params.WithProjectNameOrID(prj.Name)

			prjsummary, err := client.harborAPI.Project.GetProjectSummary(params, nil)
			if err != nil {
				return fmt.Errorf("failed to retrieve summary for project %s: %w", prj.Name, err)
			}

			url, err := url.Parse(prjsummary.Payload.Registry.URL)
			if err != nil {
				return fmt.Errorf("failed to parse Harbor Registry URL: %w", err)
			}

			log.Debugf("Discovered Proxy Cache project %s mirroring %s", prj.Name, url.Host)

			mirrors[url.Host] = fmt.Sprintf("%s/%s", client.registryURL, prj.Name)
			if url.Host == "hub.docker.com" {
				mirrors["docker.io"] = fmt.Sprintf("%s/%s", client.registryURL, prj.Name)
			}
		}
	}

	client.mirrors = mirrors

	return nil
}
