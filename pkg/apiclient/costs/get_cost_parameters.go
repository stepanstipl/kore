// Code generated by go-swagger; DO NOT EDIT.

package costs

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

// NewGetCostParams creates a new GetCostParams object
// with the default values initialized.
func NewGetCostParams() *GetCostParams {
	var ()
	return &GetCostParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetCostParamsWithTimeout creates a new GetCostParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetCostParamsWithTimeout(timeout time.Duration) *GetCostParams {
	var ()
	return &GetCostParams{

		timeout: timeout,
	}
}

// NewGetCostParamsWithContext creates a new GetCostParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetCostParamsWithContext(ctx context.Context) *GetCostParams {
	var ()
	return &GetCostParams{

		Context: ctx,
	}
}

// NewGetCostParamsWithHTTPClient creates a new GetCostParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetCostParamsWithHTTPClient(client *http.Client) *GetCostParams {
	var ()
	return &GetCostParams{
		HTTPClient: client,
	}
}

/*GetCostParams contains all the parameters to send to the API endpoint
for the get cost operation typically these are written to a http.Request
*/
type GetCostParams struct {

	/*Name
	  The name of the cost to retrieve

	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get cost params
func (o *GetCostParams) WithTimeout(timeout time.Duration) *GetCostParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get cost params
func (o *GetCostParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get cost params
func (o *GetCostParams) WithContext(ctx context.Context) *GetCostParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get cost params
func (o *GetCostParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get cost params
func (o *GetCostParams) WithHTTPClient(client *http.Client) *GetCostParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get cost params
func (o *GetCostParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the get cost params
func (o *GetCostParams) WithName(name string) *GetCostParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get cost params
func (o *GetCostParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *GetCostParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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