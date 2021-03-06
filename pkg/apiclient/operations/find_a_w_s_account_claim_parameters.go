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

// NewFindAWSAccountClaimParams creates a new FindAWSAccountClaimParams object
// with the default values initialized.
func NewFindAWSAccountClaimParams() *FindAWSAccountClaimParams {
	var ()
	return &FindAWSAccountClaimParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewFindAWSAccountClaimParamsWithTimeout creates a new FindAWSAccountClaimParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewFindAWSAccountClaimParamsWithTimeout(timeout time.Duration) *FindAWSAccountClaimParams {
	var ()
	return &FindAWSAccountClaimParams{

		timeout: timeout,
	}
}

// NewFindAWSAccountClaimParamsWithContext creates a new FindAWSAccountClaimParams object
// with the default values initialized, and the ability to set a context for a request
func NewFindAWSAccountClaimParamsWithContext(ctx context.Context) *FindAWSAccountClaimParams {
	var ()
	return &FindAWSAccountClaimParams{

		Context: ctx,
	}
}

// NewFindAWSAccountClaimParamsWithHTTPClient creates a new FindAWSAccountClaimParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewFindAWSAccountClaimParamsWithHTTPClient(client *http.Client) *FindAWSAccountClaimParams {
	var ()
	return &FindAWSAccountClaimParams{
		HTTPClient: client,
	}
}

/*FindAWSAccountClaimParams contains all the parameters to send to the API endpoint
for the find a w s account claim operation typically these are written to a http.Request
*/
type FindAWSAccountClaimParams struct {

	/*Name
	  Is name the of the resource you are acting on

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

// WithTimeout adds the timeout to the find a w s account claim params
func (o *FindAWSAccountClaimParams) WithTimeout(timeout time.Duration) *FindAWSAccountClaimParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the find a w s account claim params
func (o *FindAWSAccountClaimParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the find a w s account claim params
func (o *FindAWSAccountClaimParams) WithContext(ctx context.Context) *FindAWSAccountClaimParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the find a w s account claim params
func (o *FindAWSAccountClaimParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the find a w s account claim params
func (o *FindAWSAccountClaimParams) WithHTTPClient(client *http.Client) *FindAWSAccountClaimParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the find a w s account claim params
func (o *FindAWSAccountClaimParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the find a w s account claim params
func (o *FindAWSAccountClaimParams) WithName(name string) *FindAWSAccountClaimParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the find a w s account claim params
func (o *FindAWSAccountClaimParams) SetName(name string) {
	o.Name = name
}

// WithTeam adds the team to the find a w s account claim params
func (o *FindAWSAccountClaimParams) WithTeam(team string) *FindAWSAccountClaimParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the find a w s account claim params
func (o *FindAWSAccountClaimParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *FindAWSAccountClaimParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
