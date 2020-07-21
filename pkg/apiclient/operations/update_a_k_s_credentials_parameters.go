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

	"github.com/appvia/kore/pkg/apiclient/models"
)

// NewUpdateAKSCredentialsParams creates a new UpdateAKSCredentialsParams object
// with the default values initialized.
func NewUpdateAKSCredentialsParams() *UpdateAKSCredentialsParams {
	var ()
	return &UpdateAKSCredentialsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateAKSCredentialsParamsWithTimeout creates a new UpdateAKSCredentialsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUpdateAKSCredentialsParamsWithTimeout(timeout time.Duration) *UpdateAKSCredentialsParams {
	var ()
	return &UpdateAKSCredentialsParams{

		timeout: timeout,
	}
}

// NewUpdateAKSCredentialsParamsWithContext creates a new UpdateAKSCredentialsParams object
// with the default values initialized, and the ability to set a context for a request
func NewUpdateAKSCredentialsParamsWithContext(ctx context.Context) *UpdateAKSCredentialsParams {
	var ()
	return &UpdateAKSCredentialsParams{

		Context: ctx,
	}
}

// NewUpdateAKSCredentialsParamsWithHTTPClient creates a new UpdateAKSCredentialsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUpdateAKSCredentialsParamsWithHTTPClient(client *http.Client) *UpdateAKSCredentialsParams {
	var ()
	return &UpdateAKSCredentialsParams{
		HTTPClient: client,
	}
}

/*UpdateAKSCredentialsParams contains all the parameters to send to the API endpoint
for the update a k s credentials operation typically these are written to a http.Request
*/
type UpdateAKSCredentialsParams struct {

	/*Body
	  The definition for AKS Credentials

	*/
	Body *models.V1alpha1AKSCredentials
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

// WithTimeout adds the timeout to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) WithTimeout(timeout time.Duration) *UpdateAKSCredentialsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) WithContext(ctx context.Context) *UpdateAKSCredentialsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) WithHTTPClient(client *http.Client) *UpdateAKSCredentialsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) WithBody(body *models.V1alpha1AKSCredentials) *UpdateAKSCredentialsParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) SetBody(body *models.V1alpha1AKSCredentials) {
	o.Body = body
}

// WithName adds the name to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) WithName(name string) *UpdateAKSCredentialsParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) SetName(name string) {
	o.Name = name
}

// WithTeam adds the team to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) WithTeam(team string) *UpdateAKSCredentialsParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the update a k s credentials params
func (o *UpdateAKSCredentialsParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateAKSCredentialsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
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