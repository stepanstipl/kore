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

// NewRemoveFeatureParams creates a new RemoveFeatureParams object
// with the default values initialized.
func NewRemoveFeatureParams() *RemoveFeatureParams {
	var ()
	return &RemoveFeatureParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewRemoveFeatureParamsWithTimeout creates a new RemoveFeatureParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewRemoveFeatureParamsWithTimeout(timeout time.Duration) *RemoveFeatureParams {
	var ()
	return &RemoveFeatureParams{

		timeout: timeout,
	}
}

// NewRemoveFeatureParamsWithContext creates a new RemoveFeatureParams object
// with the default values initialized, and the ability to set a context for a request
func NewRemoveFeatureParamsWithContext(ctx context.Context) *RemoveFeatureParams {
	var ()
	return &RemoveFeatureParams{

		Context: ctx,
	}
}

// NewRemoveFeatureParamsWithHTTPClient creates a new RemoveFeatureParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewRemoveFeatureParamsWithHTTPClient(client *http.Client) *RemoveFeatureParams {
	var ()
	return &RemoveFeatureParams{
		HTTPClient: client,
	}
}

/*RemoveFeatureParams contains all the parameters to send to the API endpoint
for the remove feature operation typically these are written to a http.Request
*/
type RemoveFeatureParams struct {

	/*Name
	  The name of the feature you wish to delete

	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the remove feature params
func (o *RemoveFeatureParams) WithTimeout(timeout time.Duration) *RemoveFeatureParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the remove feature params
func (o *RemoveFeatureParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the remove feature params
func (o *RemoveFeatureParams) WithContext(ctx context.Context) *RemoveFeatureParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the remove feature params
func (o *RemoveFeatureParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the remove feature params
func (o *RemoveFeatureParams) WithHTTPClient(client *http.Client) *RemoveFeatureParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the remove feature params
func (o *RemoveFeatureParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the remove feature params
func (o *RemoveFeatureParams) WithName(name string) *RemoveFeatureParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the remove feature params
func (o *RemoveFeatureParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *RemoveFeatureParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
