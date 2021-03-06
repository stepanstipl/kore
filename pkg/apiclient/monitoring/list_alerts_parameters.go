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

// NewListAlertsParams creates a new ListAlertsParams object
// with the default values initialized.
func NewListAlertsParams() *ListAlertsParams {
	var ()
	return &ListAlertsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewListAlertsParamsWithTimeout creates a new ListAlertsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewListAlertsParamsWithTimeout(timeout time.Duration) *ListAlertsParams {
	var ()
	return &ListAlertsParams{

		timeout: timeout,
	}
}

// NewListAlertsParamsWithContext creates a new ListAlertsParams object
// with the default values initialized, and the ability to set a context for a request
func NewListAlertsParamsWithContext(ctx context.Context) *ListAlertsParams {
	var ()
	return &ListAlertsParams{

		Context: ctx,
	}
}

// NewListAlertsParamsWithHTTPClient creates a new ListAlertsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewListAlertsParamsWithHTTPClient(client *http.Client) *ListAlertsParams {
	var ()
	return &ListAlertsParams{
		HTTPClient: client,
	}
}

/*ListAlertsParams contains all the parameters to send to the API endpoint
for the list alerts operation typically these are written to a http.Request
*/
type ListAlertsParams struct {

	/*History
	  The number of historical records to retrieve

	*/
	History *string
	/*Label
	  A label to filter the alert by

	*/
	Label *string
	/*Latest
	  Indicates to we only want the latest alert status

	*/
	Latest *string
	/*Status
	  The alert to filter the results by

	*/
	Status *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the list alerts params
func (o *ListAlertsParams) WithTimeout(timeout time.Duration) *ListAlertsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list alerts params
func (o *ListAlertsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list alerts params
func (o *ListAlertsParams) WithContext(ctx context.Context) *ListAlertsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list alerts params
func (o *ListAlertsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list alerts params
func (o *ListAlertsParams) WithHTTPClient(client *http.Client) *ListAlertsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list alerts params
func (o *ListAlertsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithHistory adds the history to the list alerts params
func (o *ListAlertsParams) WithHistory(history *string) *ListAlertsParams {
	o.SetHistory(history)
	return o
}

// SetHistory adds the history to the list alerts params
func (o *ListAlertsParams) SetHistory(history *string) {
	o.History = history
}

// WithLabel adds the label to the list alerts params
func (o *ListAlertsParams) WithLabel(label *string) *ListAlertsParams {
	o.SetLabel(label)
	return o
}

// SetLabel adds the label to the list alerts params
func (o *ListAlertsParams) SetLabel(label *string) {
	o.Label = label
}

// WithLatest adds the latest to the list alerts params
func (o *ListAlertsParams) WithLatest(latest *string) *ListAlertsParams {
	o.SetLatest(latest)
	return o
}

// SetLatest adds the latest to the list alerts params
func (o *ListAlertsParams) SetLatest(latest *string) {
	o.Latest = latest
}

// WithStatus adds the status to the list alerts params
func (o *ListAlertsParams) WithStatus(status *string) *ListAlertsParams {
	o.SetStatus(status)
	return o
}

// SetStatus adds the status to the list alerts params
func (o *ListAlertsParams) SetStatus(status *string) {
	o.Status = status
}

// WriteToRequest writes these params to a swagger request
func (o *ListAlertsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	if o.Label != nil {

		// query param label
		var qrLabel string
		if o.Label != nil {
			qrLabel = *o.Label
		}
		qLabel := qrLabel
		if qLabel != "" {
			if err := r.SetQueryParam("label", qLabel); err != nil {
				return err
			}
		}

	}

	if o.Latest != nil {

		// query param latest
		var qrLatest string
		if o.Latest != nil {
			qrLatest = *o.Latest
		}
		qLatest := qrLatest
		if qLatest != "" {
			if err := r.SetQueryParam("latest", qLatest); err != nil {
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

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
