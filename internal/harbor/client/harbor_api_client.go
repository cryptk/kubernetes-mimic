// Code generated by go-swagger; DO NOT EDIT.

package client

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/artifact"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/auditlog"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/gc"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/icon"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/ping"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/preheat"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/project"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/replication"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/repository"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/retention"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/robot"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/robotv1"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/scan"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/scan_all"
	"github.com/cryptk/kubernetes-mimic/internal/harbor/client/systeminfo"
)

// Default harbor API HTTP client.
var Default = NewHTTPClient(nil)

const (
	// DefaultHost is the default Host
	// found in Meta (info) section of spec file
	DefaultHost string = "localhost"
	// DefaultBasePath is the default BasePath
	// found in Meta (info) section of spec file
	DefaultBasePath string = "/api/v2.0"
)

// DefaultSchemes are the default schemes found in Meta (info) section of spec file
var DefaultSchemes = []string{"http", "https"}

// NewHTTPClient creates a new harbor API HTTP client.
func NewHTTPClient(formats strfmt.Registry) *HarborAPI {
	return NewHTTPClientWithConfig(formats, nil)
}

// NewHTTPClientWithConfig creates a new harbor API HTTP client,
// using a customizable transport config.
func NewHTTPClientWithConfig(formats strfmt.Registry, cfg *TransportConfig) *HarborAPI {
	// ensure nullable parameters have default
	if cfg == nil {
		cfg = DefaultTransportConfig()
	}

	// create transport and client
	transport := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	return New(transport, formats)
}

// New creates a new harbor API client
func New(transport runtime.ClientTransport, formats strfmt.Registry) *HarborAPI {
	// ensure nullable parameters have default
	if formats == nil {
		formats = strfmt.Default
	}

	cli := new(HarborAPI)
	cli.Transport = transport
	cli.Artifact = artifact.New(transport, formats)
	cli.Auditlog = auditlog.New(transport, formats)
	cli.Gc = gc.New(transport, formats)
	cli.Icon = icon.New(transport, formats)
	cli.Ping = ping.New(transport, formats)
	cli.Preheat = preheat.New(transport, formats)
	cli.Project = project.New(transport, formats)
	cli.Replication = replication.New(transport, formats)
	cli.Repository = repository.New(transport, formats)
	cli.Retention = retention.New(transport, formats)
	cli.Robot = robot.New(transport, formats)
	cli.Robotv1 = robotv1.New(transport, formats)
	cli.Scan = scan.New(transport, formats)
	cli.ScanAll = scan_all.New(transport, formats)
	cli.Systeminfo = systeminfo.New(transport, formats)
	return cli
}

// DefaultTransportConfig creates a TransportConfig with the
// default settings taken from the meta section of the spec file.
func DefaultTransportConfig() *TransportConfig {
	return &TransportConfig{
		Host:     DefaultHost,
		BasePath: DefaultBasePath,
		Schemes:  DefaultSchemes,
	}
}

// TransportConfig contains the transport related info,
// found in the meta section of the spec file.
type TransportConfig struct {
	Host     string
	BasePath string
	Schemes  []string
}

// WithHost overrides the default host,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithHost(host string) *TransportConfig {
	cfg.Host = host
	return cfg
}

// WithBasePath overrides the default basePath,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithBasePath(basePath string) *TransportConfig {
	cfg.BasePath = basePath
	return cfg
}

// WithSchemes overrides the default schemes,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithSchemes(schemes []string) *TransportConfig {
	cfg.Schemes = schemes
	return cfg
}

// HarborAPI is a client for harbor API
type HarborAPI struct {
	Artifact artifact.ClientService

	Auditlog auditlog.ClientService

	Gc gc.ClientService

	Icon icon.ClientService

	Ping ping.ClientService

	Preheat preheat.ClientService

	Project project.ClientService

	Replication replication.ClientService

	Repository repository.ClientService

	Retention retention.ClientService

	Robot robot.ClientService

	Robotv1 robotv1.ClientService

	Scan scan.ClientService

	ScanAll scan_all.ClientService

	Systeminfo systeminfo.ClientService

	Transport runtime.ClientTransport
}

// SetTransport changes the transport on the client and all its subresources
func (c *HarborAPI) SetTransport(transport runtime.ClientTransport) {
	c.Transport = transport
	c.Artifact.SetTransport(transport)
	c.Auditlog.SetTransport(transport)
	c.Gc.SetTransport(transport)
	c.Icon.SetTransport(transport)
	c.Ping.SetTransport(transport)
	c.Preheat.SetTransport(transport)
	c.Project.SetTransport(transport)
	c.Replication.SetTransport(transport)
	c.Repository.SetTransport(transport)
	c.Retention.SetTransport(transport)
	c.Robot.SetTransport(transport)
	c.Robotv1.SetTransport(transport)
	c.Scan.SetTransport(transport)
	c.ScanAll.SetTransport(transport)
	c.Systeminfo.SetTransport(transport)
}
