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

// NewRemoveClusterParams creates a new RemoveClusterParams object
// with the default values initialized.
func NewRemoveClusterParams() *RemoveClusterParams {
	var ()
	return &RemoveClusterParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewRemoveClusterParamsWithTimeout creates a new RemoveClusterParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewRemoveClusterParamsWithTimeout(timeout time.Duration) *RemoveClusterParams {
	var ()
	return &RemoveClusterParams{

		timeout: timeout,
	}
}

// NewRemoveClusterParamsWithContext creates a new RemoveClusterParams object
// with the default values initialized, and the ability to set a context for a request
func NewRemoveClusterParamsWithContext(ctx context.Context) *RemoveClusterParams {
	var ()
	return &RemoveClusterParams{

		Context: ctx,
	}
}

// NewRemoveClusterParamsWithHTTPClient creates a new RemoveClusterParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewRemoveClusterParamsWithHTTPClient(client *http.Client) *RemoveClusterParams {
	var ()
	return &RemoveClusterParams{
		HTTPClient: client,
	}
}

/*RemoveClusterParams contains all the parameters to send to the API endpoint
for the remove cluster operation typically these are written to a http.Request
*/
type RemoveClusterParams struct {

	/*Cascade
	  If true then all objects owned by this object will be deleted too.

	*/
	Cascade *bool
	/*Name
	  Is the name of the cluster

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

// WithTimeout adds the timeout to the remove cluster params
func (o *RemoveClusterParams) WithTimeout(timeout time.Duration) *RemoveClusterParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the remove cluster params
func (o *RemoveClusterParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the remove cluster params
func (o *RemoveClusterParams) WithContext(ctx context.Context) *RemoveClusterParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the remove cluster params
func (o *RemoveClusterParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the remove cluster params
func (o *RemoveClusterParams) WithHTTPClient(client *http.Client) *RemoveClusterParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the remove cluster params
func (o *RemoveClusterParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCascade adds the cascade to the remove cluster params
func (o *RemoveClusterParams) WithCascade(cascade *bool) *RemoveClusterParams {
	o.SetCascade(cascade)
	return o
}

// SetCascade adds the cascade to the remove cluster params
func (o *RemoveClusterParams) SetCascade(cascade *bool) {
	o.Cascade = cascade
}

// WithName adds the name to the remove cluster params
func (o *RemoveClusterParams) WithName(name string) *RemoveClusterParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the remove cluster params
func (o *RemoveClusterParams) SetName(name string) {
	o.Name = name
}

// WithTeam adds the team to the remove cluster params
func (o *RemoveClusterParams) WithTeam(team string) *RemoveClusterParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the remove cluster params
func (o *RemoveClusterParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *RemoveClusterParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// path param team
	if err := r.SetPathParam("team", o.Team); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
