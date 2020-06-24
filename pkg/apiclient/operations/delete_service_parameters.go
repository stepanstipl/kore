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

// NewDeleteServiceParams creates a new DeleteServiceParams object
// with the default values initialized.
func NewDeleteServiceParams() *DeleteServiceParams {
	var ()
	return &DeleteServiceParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteServiceParamsWithTimeout creates a new DeleteServiceParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeleteServiceParamsWithTimeout(timeout time.Duration) *DeleteServiceParams {
	var ()
	return &DeleteServiceParams{

		timeout: timeout,
	}
}

// NewDeleteServiceParamsWithContext creates a new DeleteServiceParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeleteServiceParamsWithContext(ctx context.Context) *DeleteServiceParams {
	var ()
	return &DeleteServiceParams{

		Context: ctx,
	}
}

// NewDeleteServiceParamsWithHTTPClient creates a new DeleteServiceParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeleteServiceParamsWithHTTPClient(client *http.Client) *DeleteServiceParams {
	var ()
	return &DeleteServiceParams{
		HTTPClient: client,
	}
}

/*DeleteServiceParams contains all the parameters to send to the API endpoint
for the delete service operation typically these are written to a http.Request
*/
type DeleteServiceParams struct {

	/*Cascade
	  If true then all objects owned by this object will be deleted too.

	*/
	Cascade *bool
	/*Name
	  Is the name of the service

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

// WithTimeout adds the timeout to the delete service params
func (o *DeleteServiceParams) WithTimeout(timeout time.Duration) *DeleteServiceParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete service params
func (o *DeleteServiceParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete service params
func (o *DeleteServiceParams) WithContext(ctx context.Context) *DeleteServiceParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete service params
func (o *DeleteServiceParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete service params
func (o *DeleteServiceParams) WithHTTPClient(client *http.Client) *DeleteServiceParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete service params
func (o *DeleteServiceParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCascade adds the cascade to the delete service params
func (o *DeleteServiceParams) WithCascade(cascade *bool) *DeleteServiceParams {
	o.SetCascade(cascade)
	return o
}

// SetCascade adds the cascade to the delete service params
func (o *DeleteServiceParams) SetCascade(cascade *bool) {
	o.Cascade = cascade
}

// WithName adds the name to the delete service params
func (o *DeleteServiceParams) WithName(name string) *DeleteServiceParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the delete service params
func (o *DeleteServiceParams) SetName(name string) {
	o.Name = name
}

// WithTeam adds the team to the delete service params
func (o *DeleteServiceParams) WithTeam(team string) *DeleteServiceParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the delete service params
func (o *DeleteServiceParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteServiceParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
