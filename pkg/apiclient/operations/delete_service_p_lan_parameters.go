// Code generated by go-swagger; DO NOT EDIT.

package operations

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
	"github.com/go-openapi/swag"
)

// NewDeleteServicePLanParams creates a new DeleteServicePLanParams object
// with the default values initialized.
func NewDeleteServicePLanParams() *DeleteServicePLanParams {
	var ()
	return &DeleteServicePLanParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteServicePLanParamsWithTimeout creates a new DeleteServicePLanParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeleteServicePLanParamsWithTimeout(timeout time.Duration) *DeleteServicePLanParams {
	var ()
	return &DeleteServicePLanParams{

		timeout: timeout,
	}
}

// NewDeleteServicePLanParamsWithContext creates a new DeleteServicePLanParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeleteServicePLanParamsWithContext(ctx context.Context) *DeleteServicePLanParams {
	var ()
	return &DeleteServicePLanParams{

		Context: ctx,
	}
}

// NewDeleteServicePLanParamsWithHTTPClient creates a new DeleteServicePLanParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeleteServicePLanParamsWithHTTPClient(client *http.Client) *DeleteServicePLanParams {
	var ()
	return &DeleteServicePLanParams{
		HTTPClient: client,
	}
}

/*DeleteServicePLanParams contains all the parameters to send to the API endpoint
for the delete service p lan operation typically these are written to a http.Request
*/
type DeleteServicePLanParams struct {

	/*Cascade
	  If true then all objects owned by this object will be deleted too.

	*/
	Cascade *bool
	/*Name
	  The name of the service plan you wish to delete

	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the delete service p lan params
func (o *DeleteServicePLanParams) WithTimeout(timeout time.Duration) *DeleteServicePLanParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete service p lan params
func (o *DeleteServicePLanParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete service p lan params
func (o *DeleteServicePLanParams) WithContext(ctx context.Context) *DeleteServicePLanParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete service p lan params
func (o *DeleteServicePLanParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete service p lan params
func (o *DeleteServicePLanParams) WithHTTPClient(client *http.Client) *DeleteServicePLanParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete service p lan params
func (o *DeleteServicePLanParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCascade adds the cascade to the delete service p lan params
func (o *DeleteServicePLanParams) WithCascade(cascade *bool) *DeleteServicePLanParams {
	o.SetCascade(cascade)
	return o
}

// SetCascade adds the cascade to the delete service p lan params
func (o *DeleteServicePLanParams) SetCascade(cascade *bool) {
	o.Cascade = cascade
}

// WithName adds the name to the delete service p lan params
func (o *DeleteServicePLanParams) WithName(name string) *DeleteServicePLanParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the delete service p lan params
func (o *DeleteServicePLanParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteServicePLanParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Cascade != nil {

		// query param cascade
		var qrCascade bool
		if o.Cascade != nil {
			qrCascade = *o.Cascade
		}
		qCascade := swag.FormatBool(qrCascade)
		if qCascade != "" {
			if err := r.SetQueryParam("cascade", qCascade); err != nil {
				return err
			}
		}

	}

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
