// Code generated by go-swagger; DO NOT EDIT.

package retention

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetRentenitionMetadataParams creates a new GetRentenitionMetadataParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetRentenitionMetadataParams() *GetRentenitionMetadataParams {
	return &GetRentenitionMetadataParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetRentenitionMetadataParamsWithTimeout creates a new GetRentenitionMetadataParams object
// with the ability to set a timeout on a request.
func NewGetRentenitionMetadataParamsWithTimeout(timeout time.Duration) *GetRentenitionMetadataParams {
	return &GetRentenitionMetadataParams{
		timeout: timeout,
	}
}

// NewGetRentenitionMetadataParamsWithContext creates a new GetRentenitionMetadataParams object
// with the ability to set a context for a request.
func NewGetRentenitionMetadataParamsWithContext(ctx context.Context) *GetRentenitionMetadataParams {
	return &GetRentenitionMetadataParams{
		Context: ctx,
	}
}

// NewGetRentenitionMetadataParamsWithHTTPClient creates a new GetRentenitionMetadataParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetRentenitionMetadataParamsWithHTTPClient(client *http.Client) *GetRentenitionMetadataParams {
	return &GetRentenitionMetadataParams{
		HTTPClient: client,
	}
}

/* GetRentenitionMetadataParams contains all the parameters to send to the API endpoint
   for the get rentenition metadata operation.

   Typically these are written to a http.Request.
*/
type GetRentenitionMetadataParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get rentenition metadata params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetRentenitionMetadataParams) WithDefaults() *GetRentenitionMetadataParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get rentenition metadata params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetRentenitionMetadataParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get rentenition metadata params
func (o *GetRentenitionMetadataParams) WithTimeout(timeout time.Duration) *GetRentenitionMetadataParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get rentenition metadata params
func (o *GetRentenitionMetadataParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get rentenition metadata params
func (o *GetRentenitionMetadataParams) WithContext(ctx context.Context) *GetRentenitionMetadataParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get rentenition metadata params
func (o *GetRentenitionMetadataParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get rentenition metadata params
func (o *GetRentenitionMetadataParams) WithHTTPClient(client *http.Client) *GetRentenitionMetadataParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get rentenition metadata params
func (o *GetRentenitionMetadataParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetRentenitionMetadataParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
