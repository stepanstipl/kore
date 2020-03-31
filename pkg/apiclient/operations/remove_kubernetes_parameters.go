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

// NewRemoveKubernetesParams creates a new RemoveKubernetesParams object
// with the default values initialized.
func NewRemoveKubernetesParams() *RemoveKubernetesParams {
	var ()
	return &RemoveKubernetesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewRemoveKubernetesParamsWithTimeout creates a new RemoveKubernetesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewRemoveKubernetesParamsWithTimeout(timeout time.Duration) *RemoveKubernetesParams {
	var ()
	return &RemoveKubernetesParams{

		timeout: timeout,
	}
}

// NewRemoveKubernetesParamsWithContext creates a new RemoveKubernetesParams object
// with the default values initialized, and the ability to set a context for a request
func NewRemoveKubernetesParamsWithContext(ctx context.Context) *RemoveKubernetesParams {
	var ()
	return &RemoveKubernetesParams{

		Context: ctx,
	}
}

// NewRemoveKubernetesParamsWithHTTPClient creates a new RemoveKubernetesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewRemoveKubernetesParamsWithHTTPClient(client *http.Client) *RemoveKubernetesParams {
	var ()
	return &RemoveKubernetesParams{
		HTTPClient: client,
	}
}

/*RemoveKubernetesParams contains all the parameters to send to the API endpoint
for the remove kubernetes operation typically these are written to a http.Request
*/
type RemoveKubernetesParams struct {

	/*Name
	  Is name the of the GKE cluster you are acting upon

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

// WithTimeout adds the timeout to the remove kubernetes params
func (o *RemoveKubernetesParams) WithTimeout(timeout time.Duration) *RemoveKubernetesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the remove kubernetes params
func (o *RemoveKubernetesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the remove kubernetes params
func (o *RemoveKubernetesParams) WithContext(ctx context.Context) *RemoveKubernetesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the remove kubernetes params
func (o *RemoveKubernetesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the remove kubernetes params
func (o *RemoveKubernetesParams) WithHTTPClient(client *http.Client) *RemoveKubernetesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the remove kubernetes params
func (o *RemoveKubernetesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the remove kubernetes params
func (o *RemoveKubernetesParams) WithName(name string) *RemoveKubernetesParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the remove kubernetes params
func (o *RemoveKubernetesParams) SetName(name string) {
	o.Name = name
}

// WithTeam adds the team to the remove kubernetes params
func (o *RemoveKubernetesParams) WithTeam(team string) *RemoveKubernetesParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the remove kubernetes params
func (o *RemoveKubernetesParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *RemoveKubernetesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
