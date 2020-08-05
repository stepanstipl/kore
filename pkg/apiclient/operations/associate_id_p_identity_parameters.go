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

// NewAssociateIDPIdentityParams creates a new AssociateIDPIdentityParams object
// with the default values initialized.
func NewAssociateIDPIdentityParams() *AssociateIDPIdentityParams {
	var ()
	return &AssociateIDPIdentityParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewAssociateIDPIdentityParamsWithTimeout creates a new AssociateIDPIdentityParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewAssociateIDPIdentityParamsWithTimeout(timeout time.Duration) *AssociateIDPIdentityParams {
	var ()
	return &AssociateIDPIdentityParams{

		timeout: timeout,
	}
}

// NewAssociateIDPIdentityParamsWithContext creates a new AssociateIDPIdentityParams object
// with the default values initialized, and the ability to set a context for a request
func NewAssociateIDPIdentityParamsWithContext(ctx context.Context) *AssociateIDPIdentityParams {
	var ()
	return &AssociateIDPIdentityParams{

		Context: ctx,
	}
}

// NewAssociateIDPIdentityParamsWithHTTPClient creates a new AssociateIDPIdentityParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewAssociateIDPIdentityParamsWithHTTPClient(client *http.Client) *AssociateIDPIdentityParams {
	var ()
	return &AssociateIDPIdentityParams{
		HTTPClient: client,
	}
}

/*AssociateIDPIdentityParams contains all the parameters to send to the API endpoint
for the associate ID p identity operation typically these are written to a http.Request
*/
type AssociateIDPIdentityParams struct {

	/*Body*/
	Body *models.V1UpdateIDPIdentity
	/*User
	  The name of the user you wish to retrieve identities for

	*/
	User string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the associate ID p identity params
func (o *AssociateIDPIdentityParams) WithTimeout(timeout time.Duration) *AssociateIDPIdentityParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the associate ID p identity params
func (o *AssociateIDPIdentityParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the associate ID p identity params
func (o *AssociateIDPIdentityParams) WithContext(ctx context.Context) *AssociateIDPIdentityParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the associate ID p identity params
func (o *AssociateIDPIdentityParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the associate ID p identity params
func (o *AssociateIDPIdentityParams) WithHTTPClient(client *http.Client) *AssociateIDPIdentityParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the associate ID p identity params
func (o *AssociateIDPIdentityParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the associate ID p identity params
func (o *AssociateIDPIdentityParams) WithBody(body *models.V1UpdateIDPIdentity) *AssociateIDPIdentityParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the associate ID p identity params
func (o *AssociateIDPIdentityParams) SetBody(body *models.V1UpdateIDPIdentity) {
	o.Body = body
}

// WithUser adds the user to the associate ID p identity params
func (o *AssociateIDPIdentityParams) WithUser(user string) *AssociateIDPIdentityParams {
	o.SetUser(user)
	return o
}

// SetUser adds the user to the associate ID p identity params
func (o *AssociateIDPIdentityParams) SetUser(user string) {
	o.User = user
}

// WriteToRequest writes these params to a swagger request
func (o *AssociateIDPIdentityParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
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