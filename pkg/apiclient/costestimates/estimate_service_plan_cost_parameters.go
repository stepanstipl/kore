// Code generated by go-swagger; DO NOT EDIT.

package costestimates

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

// NewEstimateServicePlanCostParams creates a new EstimateServicePlanCostParams object
// with the default values initialized.
func NewEstimateServicePlanCostParams() *EstimateServicePlanCostParams {
	var ()
	return &EstimateServicePlanCostParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewEstimateServicePlanCostParamsWithTimeout creates a new EstimateServicePlanCostParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewEstimateServicePlanCostParamsWithTimeout(timeout time.Duration) *EstimateServicePlanCostParams {
	var ()
	return &EstimateServicePlanCostParams{

		timeout: timeout,
	}
}

// NewEstimateServicePlanCostParamsWithContext creates a new EstimateServicePlanCostParams object
// with the default values initialized, and the ability to set a context for a request
func NewEstimateServicePlanCostParamsWithContext(ctx context.Context) *EstimateServicePlanCostParams {
	var ()
	return &EstimateServicePlanCostParams{

		Context: ctx,
	}
}

// NewEstimateServicePlanCostParamsWithHTTPClient creates a new EstimateServicePlanCostParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewEstimateServicePlanCostParamsWithHTTPClient(client *http.Client) *EstimateServicePlanCostParams {
	var ()
	return &EstimateServicePlanCostParams{
		HTTPClient: client,
	}
}

/*EstimateServicePlanCostParams contains all the parameters to send to the API endpoint
for the estimate service plan cost operation typically these are written to a http.Request
*/
type EstimateServicePlanCostParams struct {

	/*Body
	  The specification for the plan you want estimating

	*/
	Body *models.V1ServicePlan

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the estimate service plan cost params
func (o *EstimateServicePlanCostParams) WithTimeout(timeout time.Duration) *EstimateServicePlanCostParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the estimate service plan cost params
func (o *EstimateServicePlanCostParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the estimate service plan cost params
func (o *EstimateServicePlanCostParams) WithContext(ctx context.Context) *EstimateServicePlanCostParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the estimate service plan cost params
func (o *EstimateServicePlanCostParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the estimate service plan cost params
func (o *EstimateServicePlanCostParams) WithHTTPClient(client *http.Client) *EstimateServicePlanCostParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the estimate service plan cost params
func (o *EstimateServicePlanCostParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the estimate service plan cost params
func (o *EstimateServicePlanCostParams) WithBody(body *models.V1ServicePlan) *EstimateServicePlanCostParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the estimate service plan cost params
func (o *EstimateServicePlanCostParams) SetBody(body *models.V1ServicePlan) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *EstimateServicePlanCostParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
