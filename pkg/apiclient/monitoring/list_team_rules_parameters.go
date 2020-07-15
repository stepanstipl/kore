// Code generated by go-swagger; DO NOT EDIT.

package monitoring

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

// NewListTeamRulesParams creates a new ListTeamRulesParams object
// with the default values initialized.
func NewListTeamRulesParams() *ListTeamRulesParams {
	var ()
	return &ListTeamRulesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewListTeamRulesParamsWithTimeout creates a new ListTeamRulesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewListTeamRulesParamsWithTimeout(timeout time.Duration) *ListTeamRulesParams {
	var ()
	return &ListTeamRulesParams{

		timeout: timeout,
	}
}

// NewListTeamRulesParamsWithContext creates a new ListTeamRulesParams object
// with the default values initialized, and the ability to set a context for a request
func NewListTeamRulesParamsWithContext(ctx context.Context) *ListTeamRulesParams {
	var ()
	return &ListTeamRulesParams{

		Context: ctx,
	}
}

// NewListTeamRulesParamsWithHTTPClient creates a new ListTeamRulesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewListTeamRulesParamsWithHTTPClient(client *http.Client) *ListTeamRulesParams {
	var ()
	return &ListTeamRulesParams{
		HTTPClient: client,
	}
}

/*ListTeamRulesParams contains all the parameters to send to the API endpoint
for the list team rules operation typically these are written to a http.Request
*/
type ListTeamRulesParams struct {

	/*History
	  The number of historical records to retrieve

	*/
	History *string
	/*Status
	  The alert to filter the results by

	*/
	Status *string
	/*Team
	  Is the name of the team the alerts reside

	*/
	Team string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the list team rules params
func (o *ListTeamRulesParams) WithTimeout(timeout time.Duration) *ListTeamRulesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list team rules params
func (o *ListTeamRulesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list team rules params
func (o *ListTeamRulesParams) WithContext(ctx context.Context) *ListTeamRulesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list team rules params
func (o *ListTeamRulesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list team rules params
func (o *ListTeamRulesParams) WithHTTPClient(client *http.Client) *ListTeamRulesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list team rules params
func (o *ListTeamRulesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithHistory adds the history to the list team rules params
func (o *ListTeamRulesParams) WithHistory(history *string) *ListTeamRulesParams {
	o.SetHistory(history)
	return o
}

// SetHistory adds the history to the list team rules params
func (o *ListTeamRulesParams) SetHistory(history *string) {
	o.History = history
}

// WithStatus adds the status to the list team rules params
func (o *ListTeamRulesParams) WithStatus(status *string) *ListTeamRulesParams {
	o.SetStatus(status)
	return o
}

// SetStatus adds the status to the list team rules params
func (o *ListTeamRulesParams) SetStatus(status *string) {
	o.Status = status
}

// WithTeam adds the team to the list team rules params
func (o *ListTeamRulesParams) WithTeam(team string) *ListTeamRulesParams {
	o.SetTeam(team)
	return o
}

// SetTeam adds the team to the list team rules params
func (o *ListTeamRulesParams) SetTeam(team string) {
	o.Team = team
}

// WriteToRequest writes these params to a swagger request
func (o *ListTeamRulesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.History != nil {

		// query param history
		var qrHistory string
		if o.History != nil {
			qrHistory = *o.History
		}
		qHistory := qrHistory
		if qHistory != "" {
			if err := r.SetQueryParam("history", qHistory); err != nil {
				return err
			}
		}

	}

	if o.Status != nil {

		// query param status
		var qrStatus string
		if o.Status != nil {
			qrStatus = *o.Status
		}
		qStatus := qrStatus
		if qStatus != "" {
			if err := r.SetQueryParam("status", qStatus); err != nil {
				return err
			}
		}

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