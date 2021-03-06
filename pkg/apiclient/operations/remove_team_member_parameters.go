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

// NewRemoveTeamMemberParams creates a new RemoveTeamMemberParams object
// with the default values initialized.
func NewRemoveTeamMemberParams() *RemoveTeamMemberParams {
	var ()
	return &RemoveTeamMemberParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewRemoveTeamMemberParamsWithTimeout creates a new RemoveTeamMemberParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewRemoveTeamMemberParamsWithTimeout(timeout time.Duration) *RemoveTeamMemberParams {
	var ()
	return &RemoveTeamMemberParams{

		timeout: timeout,
	}
}

// NewRemoveTeamMemberParamsWithContext creates a new RemoveTeamMemberParams object
// with the default values initialized, and the ability to set a context for a request
func NewRemoveTeamMemberParamsWithContext(ctx context.Context) *RemoveTeamMemberParams {
	var ()
	return &RemoveTeamMemberParams{

		Context: ctx,
	}
}

// NewRemoveTeamMemberParamsWithHTTPClient creates a new RemoveTeamMemberParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewRemoveTeamMemberParamsWithHTTPClient(client *http.Client) *RemoveTeamMemberParams {
	var ()
	return &RemoveTeamMemberParams{
		HTTPClient: client,
	}
}

/*RemoveTeamMemberParams contains all the parameters to send to the API endpoint
for the remove team member operation typically these are written to a http.Request
*/
type RemoveTeamMemberParams struct {

	/*Team
	  Is the name of the team you are acting within

	*/
	Team string
	/*User
	  Is the user you are removing from the team

	*/
	User string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the remove team member params
func (o *RemoveTeamMemberParams) WithTimeout(timeout time.Duration) *RemoveTeamMemberParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the remove team member params
func (o *RemoveTeamMemberParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the remove team member params
func (o *RemoveTeamMemberParams) WithContext(ctx context.Context) *RemoveTeamMemberParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the remove team member params
func (o *RemoveTeamMemberParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the remove team member params
func (o *RemoveTeamMemberParams) WithHTTPClient(client *http.Client) *RemoveTeamMemberParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the remove team member params
func (o *RemoveTeamMemberParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithTeam adds the team to the remove team member params
func (o *RemoveTeamMemberParams) WithTeam(team string) *RemoveTeamMemberParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the remove team member params
func (o *RemoveTeamMemberParams) SetTeam(team string) {
	o.Team = team
}

// WithUser adds the user to the remove team member params
func (o *RemoveTeamMemberParams) WithUser(user string) *RemoveTeamMemberParams {
	o.SetUser(user)
	return o
}

// SetUser adds the user to the remove team member params
func (o *RemoveTeamMemberParams) SetUser(user string) {
	o.User = user
}

// WriteToRequest writes these params to a swagger request
func (o *RemoveTeamMemberParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param team
	if err := r.SetPathParam("team", o.Team); err != nil {
		return err
	}

	// path param user
	if err := r.SetPathParam("user", o.User); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
