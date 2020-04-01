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

// NewFindEKSsParams creates a new FindEKSsParams object
// with the default values initialized.
func NewFindEKSsParams() *FindEKSsParams {
	var ()
	return &FindEKSsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewFindEKSsParamsWithTimeout creates a new FindEKSsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewFindEKSsParamsWithTimeout(timeout time.Duration) *FindEKSsParams {
	var ()
	return &FindEKSsParams{

		timeout: timeout,
	}
}

// NewFindEKSsParamsWithContext creates a new FindEKSsParams object
// with the default values initialized, and the ability to set a context for a request
func NewFindEKSsParamsWithContext(ctx context.Context) *FindEKSsParams {
	var ()
	return &FindEKSsParams{

		Context: ctx,
	}
}

// NewFindEKSsParamsWithHTTPClient creates a new FindEKSsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewFindEKSsParamsWithHTTPClient(client *http.Client) *FindEKSsParams {
	var ()
	return &FindEKSsParams{
		HTTPClient: client,
	}
}

/*FindEKSsParams contains all the parameters to send to the API endpoint
for the find e k ss operation typically these are written to a http.Request
*/
type FindEKSsParams struct {

	/*Team
	  Is the name of the team you are acting within

	*/
	Team string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the find e k ss params
func (o *FindEKSsParams) WithTimeout(timeout time.Duration) *FindEKSsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the find e k ss params
func (o *FindEKSsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the find e k ss params
func (o *FindEKSsParams) WithContext(ctx context.Context) *FindEKSsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the find e k ss params
func (o *FindEKSsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the find e k ss params
func (o *FindEKSsParams) WithHTTPClient(client *http.Client) *FindEKSsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the find e k ss params
func (o *FindEKSsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithTeam adds the team to the find e k ss params
func (o *FindEKSsParams) WithTeam(team string) *FindEKSsParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the find e k ss params
func (o *FindEKSsParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *FindEKSsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
