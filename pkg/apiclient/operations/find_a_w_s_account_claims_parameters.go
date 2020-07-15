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

// NewFindAWSAccountClaimsParams creates a new FindAWSAccountClaimsParams object
// with the default values initialized.
func NewFindAWSAccountClaimsParams() *FindAWSAccountClaimsParams {
	var ()
	return &FindAWSAccountClaimsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewFindAWSAccountClaimsParamsWithTimeout creates a new FindAWSAccountClaimsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewFindAWSAccountClaimsParamsWithTimeout(timeout time.Duration) *FindAWSAccountClaimsParams {
	var ()
	return &FindAWSAccountClaimsParams{

		timeout: timeout,
	}
}

// NewFindAWSAccountClaimsParamsWithContext creates a new FindAWSAccountClaimsParams object
// with the default values initialized, and the ability to set a context for a request
func NewFindAWSAccountClaimsParamsWithContext(ctx context.Context) *FindAWSAccountClaimsParams {
	var ()
	return &FindAWSAccountClaimsParams{

		Context: ctx,
	}
}

// NewFindAWSAccountClaimsParamsWithHTTPClient creates a new FindAWSAccountClaimsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewFindAWSAccountClaimsParamsWithHTTPClient(client *http.Client) *FindAWSAccountClaimsParams {
	var ()
	return &FindAWSAccountClaimsParams{
		HTTPClient: client,
	}
}

/*FindAWSAccountClaimsParams contains all the parameters to send to the API endpoint
for the find a w s account claims operation typically these are written to a http.Request
*/
type FindAWSAccountClaimsParams struct {

	/*Team
	  Is the name of the team you are acting within

	*/
	Team string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the find a w s account claims params
func (o *FindAWSAccountClaimsParams) WithTimeout(timeout time.Duration) *FindAWSAccountClaimsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the find a w s account claims params
func (o *FindAWSAccountClaimsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the find a w s account claims params
func (o *FindAWSAccountClaimsParams) WithContext(ctx context.Context) *FindAWSAccountClaimsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the find a w s account claims params
func (o *FindAWSAccountClaimsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the find a w s account claims params
func (o *FindAWSAccountClaimsParams) WithHTTPClient(client *http.Client) *FindAWSAccountClaimsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the find a w s account claims params
func (o *FindAWSAccountClaimsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithTeam adds the team to the find a w s account claims params
func (o *FindAWSAccountClaimsParams) WithTeam(team string) *FindAWSAccountClaimsParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the find a w s account claims params
func (o *FindAWSAccountClaimsParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *FindAWSAccountClaimsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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