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

// ListRulesReader is a Reader for the ListRules structure.
type ListRulesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListRulesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListRulesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewListRulesBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewListRulesUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewListRulesForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewListRulesInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewListRulesOK creates a ListRulesOK with default headers values
func NewListRulesOK() *ListRulesOK {
	return &ListRulesOK{}
}

/*ListRulesOK handles this case with default header values.

Listing of the rules in kore
*/
type ListRulesOK struct {
	Payload *models.V1beta1AlertRuleList
}

func (o *ListRulesOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/monitoring/rules][%d] listRulesOK  %+v", 200, o.Payload)
}

func (o *ListRulesOK) GetPayload() *models.V1beta1AlertRuleList {
	return o.Payload
}

func (o *ListRulesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1AlertRuleList)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListRulesBadRequest creates a ListRulesBadRequest with default headers values
func NewListRulesBadRequest() *ListRulesBadRequest {
	return &ListRulesBadRequest{}
}

/*ListRulesBadRequest handles this case with default header values.

Validation error of supplied parameters/body
*/
type ListRulesBadRequest struct {
	Payload *models.ValidationError
}

func (o *ListRulesBadRequest) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/monitoring/rules][%d] listRulesBadRequest  %+v", 400, o.Payload)
}

func (o *ListRulesBadRequest) GetPayload() *models.ValidationError {
	return o.Payload
}

func (o *ListRulesBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListRulesUnauthorized creates a ListRulesUnauthorized with default headers values
func NewListRulesUnauthorized() *ListRulesUnauthorized {
	return &ListRulesUnauthorized{}
}

/*ListRulesUnauthorized handles this case with default header values.

If not authenticated
*/
type ListRulesUnauthorized struct {
}

func (o *ListRulesUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/monitoring/rules][%d] listRulesUnauthorized ", 401)
}

func (o *ListRulesUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewListRulesForbidden creates a ListRulesForbidden with default headers values
func NewListRulesForbidden() *ListRulesForbidden {
	return &ListRulesForbidden{}
}

/*ListRulesForbidden handles this case with default header values.

If authenticated but not authorized
*/
type ListRulesForbidden struct {
}

func (o *ListRulesForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/monitoring/rules][%d] listRulesForbidden ", 403)
}

func (o *ListRulesForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewListRulesInternalServerError creates a ListRulesInternalServerError with default headers values
func NewListRulesInternalServerError() *ListRulesInternalServerError {
	return &ListRulesInternalServerError{}
}

/*ListRulesInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type ListRulesInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *ListRulesInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/monitoring/rules][%d] listRulesInternalServerError  %+v", 500, o.Payload)
}

func (o *ListRulesInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *ListRulesInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}