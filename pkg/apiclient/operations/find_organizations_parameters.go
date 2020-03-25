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

// NewFindOrganizationsParams creates a new FindOrganizationsParams object
// with the default values initialized.
func NewFindOrganizationsParams() *FindOrganizationsParams {
	var ()
	return &FindOrganizationsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewFindOrganizationsParamsWithTimeout creates a new FindOrganizationsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewFindOrganizationsParamsWithTimeout(timeout time.Duration) *FindOrganizationsParams {
	var ()
	return &FindOrganizationsParams{

		timeout: timeout,
	}
}

// NewFindOrganizationsParamsWithContext creates a new FindOrganizationsParams object
// with the default values initialized, and the ability to set a context for a request
func NewFindOrganizationsParamsWithContext(ctx context.Context) *FindOrganizationsParams {
	var ()
	return &FindOrganizationsParams{

		Context: ctx,
	}
}

// NewFindOrganizationsParamsWithHTTPClient creates a new FindOrganizationsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewFindOrganizationsParamsWithHTTPClient(client *http.Client) *FindOrganizationsParams {
	var ()
	return &FindOrganizationsParams{
		HTTPClient: client,
	}
}

/*FindOrganizationsParams contains all the parameters to send to the API endpoint
for the find organizations operation typically these are written to a http.Request
*/
type FindOrganizationsParams struct {

	/*Team
	  Is the name of the team you are acting within

	*/
	Team string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the find organizations params
func (o *FindOrganizationsParams) WithTimeout(timeout time.Duration) *FindOrganizationsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the find organizations params
func (o *FindOrganizationsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the find organizations params
func (o *FindOrganizationsParams) WithContext(ctx context.Context) *FindOrganizationsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the find organizations params
func (o *FindOrganizationsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the find organizations params
func (o *FindOrganizationsParams) WithHTTPClient(client *http.Client) *FindOrganizationsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the find organizations params
func (o *FindOrganizationsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithTeam adds the team to the find organizations params
func (o *FindOrganizationsParams) WithTeam(team string) *FindOrganizationsParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the find organizations params
func (o *FindOrganizationsParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *FindOrganizationsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param team
	if err := r.SetPathParam("team", o.Team); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}