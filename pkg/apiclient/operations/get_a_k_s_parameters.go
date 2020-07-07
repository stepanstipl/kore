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
)

// NewGetAKSParams creates a new GetAKSParams object
// with the default values initialized.
func NewGetAKSParams() *GetAKSParams {
	var ()
	return &GetAKSParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAKSParamsWithTimeout creates a new GetAKSParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAKSParamsWithTimeout(timeout time.Duration) *GetAKSParams {
	var ()
	return &GetAKSParams{

		timeout: timeout,
	}
}

// NewGetAKSParamsWithContext creates a new GetAKSParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAKSParamsWithContext(ctx context.Context) *GetAKSParams {
	var ()
	return &GetAKSParams{

		Context: ctx,
	}
}

// NewGetAKSParamsWithHTTPClient creates a new GetAKSParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetAKSParamsWithHTTPClient(client *http.Client) *GetAKSParams {
	var ()
	return &GetAKSParams{
		HTTPClient: client,
	}
}

/*GetAKSParams contains all the parameters to send to the API endpoint
for the get a k s operation typically these are written to a http.Request
*/
type GetAKSParams struct {

	/*Name
	  Is name the of the AKS cluster you are acting upon

	*/
	Name string
	/*Team
	  Is the name of the team you are acting within

	*/
	Team string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get a k s params
func (o *GetAKSParams) WithTimeout(timeout time.Duration) *GetAKSParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get a k s params
func (o *GetAKSParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get a k s params
func (o *GetAKSParams) WithContext(ctx context.Context) *GetAKSParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get a k s params
func (o *GetAKSParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get a k s params
func (o *GetAKSParams) WithHTTPClient(client *http.Client) *GetAKSParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get a k s params
func (o *GetAKSParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the get a k s params
func (o *GetAKSParams) WithName(name string) *GetAKSParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the get a k s params
func (o *GetAKSParams) SetName(name string) {
	o.Name = name
}

// WithTeam adds the team to the get a k s params
func (o *GetAKSParams) WithTeam(team string) *GetAKSParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the get a k s params
func (o *GetAKSParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *GetAKSParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	// path param team
	if err := r.SetPathParam("team", o.Team); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
