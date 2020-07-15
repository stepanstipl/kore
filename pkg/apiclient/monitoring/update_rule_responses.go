// Code generated by go-swagger; DO NOT EDIT.

package monitoring

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/appvia/kore/pkg/apiclient/models"
)

// UpdateRuleReader is a Reader for the UpdateRule structure.
type UpdateRuleReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateRuleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateRuleOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdateRuleBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewUpdateRuleUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewUpdateRuleForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateRuleInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUpdateRuleOK creates a UpdateRuleOK with default headers values
func NewUpdateRuleOK() *UpdateRuleOK {
	return &UpdateRuleOK{}
}

/*UpdateRuleOK handles this case with default header values.

The rule has been deleted
*/
type UpdateRuleOK struct {
	Payload *models.V1beta1AlertRule
}

func (o *UpdateRuleOK) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/monitoring/rules/{group}/{version}/{kind}/{namespace}/{resource}/{name}][%d] updateRuleOK  %+v", 200, o.Payload)
}

func (o *UpdateRuleOK) GetPayload() *models.V1beta1AlertRule {
	return o.Payload
}

func (o *UpdateRuleOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1AlertRule)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateRuleBadRequest creates a UpdateRuleBadRequest with default headers values
func NewUpdateRuleBadRequest() *UpdateRuleBadRequest {
	return &UpdateRuleBadRequest{}
}

/*UpdateRuleBadRequest handles this case with default header values.

Validation error of supplied parameters/body
*/
type UpdateRuleBadRequest struct {
	Payload *models.ValidationError
}

func (o *UpdateRuleBadRequest) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/monitoring/rules/{group}/{version}/{kind}/{namespace}/{resource}/{name}][%d] updateRuleBadRequest  %+v", 400, o.Payload)
}

func (o *UpdateRuleBadRequest) GetPayload() *models.ValidationError {
	return o.Payload
}

func (o *UpdateRuleBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateRuleUnauthorized creates a UpdateRuleUnauthorized with default headers values
func NewUpdateRuleUnauthorized() *UpdateRuleUnauthorized {
	return &UpdateRuleUnauthorized{}
}

/*UpdateRuleUnauthorized handles this case with default header values.

If not authenticated
*/
type UpdateRuleUnauthorized struct {
}

func (o *UpdateRuleUnauthorized) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/monitoring/rules/{group}/{version}/{kind}/{namespace}/{resource}/{name}][%d] updateRuleUnauthorized ", 401)
}

func (o *UpdateRuleUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateRuleForbidden creates a UpdateRuleForbidden with default headers values
func NewUpdateRuleForbidden() *UpdateRuleForbidden {
	return &UpdateRuleForbidden{}
}

/*UpdateRuleForbidden handles this case with default header values.

If authenticated but not authorized
*/
type UpdateRuleForbidden struct {
}

func (o *UpdateRuleForbidden) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/monitoring/rules/{group}/{version}/{kind}/{namespace}/{resource}/{name}][%d] updateRuleForbidden ", 403)
}

func (o *UpdateRuleForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateRuleInternalServerError creates a UpdateRuleInternalServerError with default headers values
func NewUpdateRuleInternalServerError() *UpdateRuleInternalServerError {
	return &UpdateRuleInternalServerError{}
}

/*UpdateRuleInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type UpdateRuleInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *UpdateRuleInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/monitoring/rules/{group}/{version}/{kind}/{namespace}/{resource}/{name}][%d] updateRuleInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateRuleInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *UpdateRuleInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}