// Code generated by go-swagger; DO NOT EDIT.

package korefeatures

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

// NewGetFeatureParams creates a new GetFeatureParams object
// with the default values initialized.
func NewGetFeatureParams() *GetFeatureParams {
	var ()
	return &GetFeatureParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetFeatureParamsWithTimeout creates a new GetFeatureParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetFeatureParamsWithTimeout(timeout time.Duration) *GetFeatureParams {
	var ()
	return &GetFeatureParams{

		timeout: timeout,
	}
}

// NewGetFeatureParamsWithContext creates a new GetFeatureParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetFeatureParamsWithContext(ctx context.Context) *GetFeatureParams {
	var ()
	return &GetFeatureParams{

		Context: ctx,
	}
}

// NewGetFeatureParamsWithHTTPClient creates a new GetFeatureParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetFeatureParamsWithHTTPClient(client *http.Client) *GetFeatureParams {
	var ()
	return &GetFeatureParams{
		HTTPClient: client,
	}
}

/*GetFeatureParams contains all the parameters to send to the API endpoint
for the get feature operation typically these are written to a http.Request
*/
type GetFeatureParams struct {

	/*Name
	  The name of the feature you wish to retrieve

	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get feature params
func (o *GetFeatureParams) WithTimeout(timeout time.Duration) *GetFeatureParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get feature params
func (o *GetFeatureParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get feature params
func (o *GetFeatureParams) WithContext(ctx context.Context) *GetFeatureParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get feature params
func (o *GetFeatureParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get feature params
func (o *GetFeatureParams) WithHTTPClient(client *http.Client) *GetFeatureParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get feature params
func (o *GetFeatureParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the get feature params
func (o *GetFeatureParams) WithName(name string) *GetFeatureParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get feature params
func (o *GetFeatureParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *GetFeatureParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}