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

// ListTeamRulesReader is a Reader for the ListTeamRules structure.
type ListTeamRulesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListTeamRulesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListTeamRulesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewListTeamRulesBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewListTeamRulesUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewListTeamRulesForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewListTeamRulesInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewListTeamRulesOK creates a ListTeamRulesOK with default headers values
func NewListTeamRulesOK() *ListTeamRulesOK {
	return &ListTeamRulesOK{}
}

/*ListTeamRulesOK handles this case with default header values.

A list of the rules
*/
type ListTeamRulesOK struct {
	Payload *models.V1beta1AlertList
}

func (o *ListTeamRulesOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/monitoring/teams/{team}/rules][%d] listTeamRulesOK  %+v", 200, o.Payload)
}

func (o *ListTeamRulesOK) GetPayload() *models.V1beta1AlertList {
	return o.Payload
}

func (o *ListTeamRulesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1AlertList)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListTeamRulesBadRequest creates a ListTeamRulesBadRequest with default headers values
func NewListTeamRulesBadRequest() *ListTeamRulesBadRequest {
	return &ListTeamRulesBadRequest{}
}

/*ListTeamRulesBadRequest handles this case with default header values.

Validation error of supplied parameters/body
*/
type ListTeamRulesBadRequest struct {
	Payload *models.ValidationError
}

func (o *ListTeamRulesBadRequest) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/monitoring/teams/{team}/rules][%d] listTeamRulesBadRequest  %+v", 400, o.Payload)
}

func (o *ListTeamRulesBadRequest) GetPayload() *models.ValidationError {
	return o.Payload
}

func (o *ListTeamRulesBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListTeamRulesUnauthorized creates a ListTeamRulesUnauthorized with default headers values
func NewListTeamRulesUnauthorized() *ListTeamRulesUnauthorized {
	return &ListTeamRulesUnauthorized{}
}

/*ListTeamRulesUnauthorized handles this case with default header values.

If not authenticated
*/
type ListTeamRulesUnauthorized struct {
}

func (o *ListTeamRulesUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/monitoring/teams/{team}/rules][%d] listTeamRulesUnauthorized ", 401)
}

func (o *ListTeamRulesUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewListTeamRulesForbidden creates a ListTeamRulesForbidden with default headers values
func NewListTeamRulesForbidden() *ListTeamRulesForbidden {
	return &ListTeamRulesForbidden{}
}

/*ListTeamRulesForbidden handles this case with default header values.

If authenticated but not authorized
*/
type ListTeamRulesForbidden struct {
}

func (o *ListTeamRulesForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/monitoring/teams/{team}/rules][%d] listTeamRulesForbidden ", 403)
}

func (o *ListTeamRulesForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewListTeamRulesInternalServerError creates a ListTeamRulesInternalServerError with default headers values
func NewListTeamRulesInternalServerError() *ListTeamRulesInternalServerError {
	return &ListTeamRulesInternalServerError{}
}

/*ListTeamRulesInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type ListTeamRulesInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *ListTeamRulesInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/monitoring/teams/{team}/rules][%d] listTeamRulesInternalServerError  %+v", 500, o.Payload)
}

func (o *ListTeamRulesInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *ListTeamRulesInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}