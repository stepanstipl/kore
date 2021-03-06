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

// NewUpdateAWSOrganizationParams creates a new UpdateAWSOrganizationParams object
// with the default values initialized.
func NewUpdateAWSOrganizationParams() *UpdateAWSOrganizationParams {
	var ()
	return &UpdateAWSOrganizationParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateAWSOrganizationParamsWithTimeout creates a new UpdateAWSOrganizationParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUpdateAWSOrganizationParamsWithTimeout(timeout time.Duration) *UpdateAWSOrganizationParams {
	var ()
	return &UpdateAWSOrganizationParams{

		timeout: timeout,
	}
}

// NewUpdateAWSOrganizationParamsWithContext creates a new UpdateAWSOrganizationParams object
// with the default values initialized, and the ability to set a context for a request
func NewUpdateAWSOrganizationParamsWithContext(ctx context.Context) *UpdateAWSOrganizationParams {
	var ()
	return &UpdateAWSOrganizationParams{

		Context: ctx,
	}
}

// NewUpdateAWSOrganizationParamsWithHTTPClient creates a new UpdateAWSOrganizationParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUpdateAWSOrganizationParamsWithHTTPClient(client *http.Client) *UpdateAWSOrganizationParams {
	var ()
	return &UpdateAWSOrganizationParams{
		HTTPClient: client,
	}
}

/*UpdateAWSOrganizationParams contains all the parameters to send to the API endpoint
for the update a w s organization operation typically these are written to a http.Request
*/
type UpdateAWSOrganizationParams struct {

	/*Body
	  The definition for AWS organization

	*/
	Body *models.V1alpha1AWSOrganization
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

// WithTimeout adds the timeout to the update a w s organization params
func (o *UpdateAWSOrganizationParams) WithTimeout(timeout time.Duration) *UpdateAWSOrganizationParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update a w s organization params
func (o *UpdateAWSOrganizationParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update a w s organization params
func (o *UpdateAWSOrganizationParams) WithContext(ctx context.Context) *UpdateAWSOrganizationParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update a w s organization params
func (o *UpdateAWSOrganizationParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update a w s organization params
func (o *UpdateAWSOrganizationParams) WithHTTPClient(client *http.Client) *UpdateAWSOrganizationParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update a w s organization params
func (o *UpdateAWSOrganizationParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update a w s organization params
func (o *UpdateAWSOrganizationParams) WithBody(body *models.V1alpha1AWSOrganization) *UpdateAWSOrganizationParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update a w s organization params
func (o *UpdateAWSOrganizationParams) SetBody(body *models.V1alpha1AWSOrganization) {
	o.Body = body
}

// WithName adds the name to the update a w s organization params
func (o *UpdateAWSOrganizationParams) WithName(name string) *UpdateAWSOrganizationParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the update a w s organization params
func (o *UpdateAWSOrganizationParams) SetName(name string) {
	o.Name = name
}

// WithTeam adds the team to the update a w s organization params
func (o *UpdateAWSOrganizationParams) WithTeam(team string) *UpdateAWSOrganizationParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the update a w s organization params
func (o *UpdateAWSOrganizationParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateAWSOrganizationParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
