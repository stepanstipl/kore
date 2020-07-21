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

// NewDeleteAKSCredentialsParams creates a new DeleteAKSCredentialsParams object
// with the default values initialized.
func NewDeleteAKSCredentialsParams() *DeleteAKSCredentialsParams {
	var ()
	return &DeleteAKSCredentialsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteAKSCredentialsParamsWithTimeout creates a new DeleteAKSCredentialsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeleteAKSCredentialsParamsWithTimeout(timeout time.Duration) *DeleteAKSCredentialsParams {
	var ()
	return &DeleteAKSCredentialsParams{

		timeout: timeout,
	}
}

// NewDeleteAKSCredentialsParamsWithContext creates a new DeleteAKSCredentialsParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeleteAKSCredentialsParamsWithContext(ctx context.Context) *DeleteAKSCredentialsParams {
	var ()
	return &DeleteAKSCredentialsParams{

		Context: ctx,
	}
}

// NewDeleteAKSCredentialsParamsWithHTTPClient creates a new DeleteAKSCredentialsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeleteAKSCredentialsParamsWithHTTPClient(client *http.Client) *DeleteAKSCredentialsParams {
	var ()
	return &DeleteAKSCredentialsParams{
		HTTPClient: client,
	}
}

/*DeleteAKSCredentialsParams contains all the parameters to send to the API endpoint
for the delete a k s credentials operation typically these are written to a http.Request
*/
type DeleteAKSCredentialsParams struct {

	/*Name
	  Is name the of the AKS credentials you are acting upon

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

// WithTimeout adds the timeout to the delete a k s credentials params
func (o *DeleteAKSCredentialsParams) WithTimeout(timeout time.Duration) *DeleteAKSCredentialsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete a k s credentials params
func (o *DeleteAKSCredentialsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete a k s credentials params
func (o *DeleteAKSCredentialsParams) WithContext(ctx context.Context) *DeleteAKSCredentialsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete a k s credentials params
func (o *DeleteAKSCredentialsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete a k s credentials params
func (o *DeleteAKSCredentialsParams) WithHTTPClient(client *http.Client) *DeleteAKSCredentialsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete a k s credentials params
func (o *DeleteAKSCredentialsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the delete a k s credentials params
func (o *DeleteAKSCredentialsParams) WithName(name string) *DeleteAKSCredentialsParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the delete a k s credentials params
func (o *DeleteAKSCredentialsParams) SetName(name string) {
	o.Name = name
}

// WithTeam adds the team to the delete a k s credentials params
func (o *DeleteAKSCredentialsParams) WithTeam(team string) *DeleteAKSCredentialsParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the delete a k s credentials params
func (o *DeleteAKSCredentialsParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteAKSCredentialsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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