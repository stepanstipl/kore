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

	"github.com/appvia/kore/pkg/apiclient/models"
)

// NewUpdateFeatureParams creates a new UpdateFeatureParams object
// with the default values initialized.
func NewUpdateFeatureParams() *UpdateFeatureParams {
	var ()
	return &UpdateFeatureParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateFeatureParamsWithTimeout creates a new UpdateFeatureParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUpdateFeatureParamsWithTimeout(timeout time.Duration) *UpdateFeatureParams {
	var ()
	return &UpdateFeatureParams{

		timeout: timeout,
	}
}

// NewUpdateFeatureParamsWithContext creates a new UpdateFeatureParams object
// with the default values initialized, and the ability to set a context for a request
func NewUpdateFeatureParamsWithContext(ctx context.Context) *UpdateFeatureParams {
	var ()
	return &UpdateFeatureParams{

		Context: ctx,
	}
}

// NewUpdateFeatureParamsWithHTTPClient creates a new UpdateFeatureParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUpdateFeatureParamsWithHTTPClient(client *http.Client) *UpdateFeatureParams {
	var ()
	return &UpdateFeatureParams{
		HTTPClient: client,
	}
}

/*UpdateFeatureParams contains all the parameters to send to the API endpoint
for the update feature operation typically these are written to a http.Request
*/
type UpdateFeatureParams struct {

	/*Body
	  The specification for the feature you are creating or updating

	*/
	Body *models.V1KoreFeature
	/*Name
	  The name of the feature you wish to create or update

	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the update feature params
func (o *UpdateFeatureParams) WithTimeout(timeout time.Duration) *UpdateFeatureParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update feature params
func (o *UpdateFeatureParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update feature params
func (o *UpdateFeatureParams) WithContext(ctx context.Context) *UpdateFeatureParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update feature params
func (o *UpdateFeatureParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update feature params
func (o *UpdateFeatureParams) WithHTTPClient(client *http.Client) *UpdateFeatureParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update feature params
func (o *UpdateFeatureParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update feature params
func (o *UpdateFeatureParams) WithBody(body *models.V1KoreFeature) *UpdateFeatureParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update feature params
func (o *UpdateFeatureParams) SetBody(body *models.V1KoreFeature) {
	o.Body = body
}

// WithName adds the name to the update feature params
func (o *UpdateFeatureParams) WithName(name string) *UpdateFeatureParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the update feature params
func (o *UpdateFeatureParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateFeatureParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
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
