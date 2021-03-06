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

// NewInvitationSubmitParams creates a new InvitationSubmitParams object
// with the default values initialized.
func NewInvitationSubmitParams() *InvitationSubmitParams {
	var ()
	return &InvitationSubmitParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewInvitationSubmitParamsWithTimeout creates a new InvitationSubmitParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewInvitationSubmitParamsWithTimeout(timeout time.Duration) *InvitationSubmitParams {
	var ()
	return &InvitationSubmitParams{

		timeout: timeout,
	}
}

// NewInvitationSubmitParamsWithContext creates a new InvitationSubmitParams object
// with the default values initialized, and the ability to set a context for a request
func NewInvitationSubmitParamsWithContext(ctx context.Context) *InvitationSubmitParams {
	var ()
	return &InvitationSubmitParams{

		Context: ctx,
	}
}

// NewInvitationSubmitParamsWithHTTPClient creates a new InvitationSubmitParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewInvitationSubmitParamsWithHTTPClient(client *http.Client) *InvitationSubmitParams {
	var ()
	return &InvitationSubmitParams{
		HTTPClient: client,
	}
}

/*InvitationSubmitParams contains all the parameters to send to the API endpoint
for the invitation submit operation typically these are written to a http.Request
*/
type InvitationSubmitParams struct {

	/*Token
	  The generated base64 invitation token which was provided from the team

	*/
	Token string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the invitation submit params
func (o *InvitationSubmitParams) WithTimeout(timeout time.Duration) *InvitationSubmitParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the invitation submit params
func (o *InvitationSubmitParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the invitation submit params
func (o *InvitationSubmitParams) WithContext(ctx context.Context) *InvitationSubmitParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the invitation submit params
func (o *InvitationSubmitParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the invitation submit params
func (o *InvitationSubmitParams) WithHTTPClient(client *http.Client) *InvitationSubmitParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the invitation submit params
func (o *InvitationSubmitParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithToken adds the token to the invitation submit params
func (o *InvitationSubmitParams) WithToken(token string) *InvitationSubmitParams {
	o.SetToken(token)
	return o
}

// SetToken adds the token to the invitation submit params
func (o *InvitationSubmitParams) SetToken(token string) {
	o.Token = token
}

// WriteToRequest writes these params to a swagger request
func (o *InvitationSubmitParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param token
	if err := r.SetPathParam("token", o.Token); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
