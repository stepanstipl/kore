// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/appvia/kore/pkg/apiclient/models"
)

// GetServiceCredentialSchemaForPlanReader is a Reader for the GetServiceCredentialSchemaForPlan structure.
type GetServiceCredentialSchemaForPlanReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetServiceCredentialSchemaForPlanReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetServiceCredentialSchemaForPlanOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetServiceCredentialSchemaForPlanUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetServiceCredentialSchemaForPlanForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetServiceCredentialSchemaForPlanInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetServiceCredentialSchemaForPlanOK creates a GetServiceCredentialSchemaForPlanOK with default headers values
func NewGetServiceCredentialSchemaForPlanOK() *GetServiceCredentialSchemaForPlanOK {
	return &GetServiceCredentialSchemaForPlanOK{}
}

/*GetServiceCredentialSchemaForPlanOK handles this case with default header values.

Contains the service credential schema definition
*/
type GetServiceCredentialSchemaForPlanOK struct {
}

func (o *GetServiceCredentialSchemaForPlanOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/servicecredentialschemas/{kind}/{name}][%d] getServiceCredentialSchemaForPlanOK ", 200)
}

func (o *GetServiceCredentialSchemaForPlanOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetServiceCredentialSchemaForPlanUnauthorized creates a GetServiceCredentialSchemaForPlanUnauthorized with default headers values
func NewGetServiceCredentialSchemaForPlanUnauthorized() *GetServiceCredentialSchemaForPlanUnauthorized {
	return &GetServiceCredentialSchemaForPlanUnauthorized{}
}

/*GetServiceCredentialSchemaForPlanUnauthorized handles this case with default header values.

If not authenticated
*/
type GetServiceCredentialSchemaForPlanUnauthorized struct {
}

func (o *GetServiceCredentialSchemaForPlanUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/servicecredentialschemas/{kind}/{name}][%d] getServiceCredentialSchemaForPlanUnauthorized ", 401)
}

func (o *GetServiceCredentialSchemaForPlanUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetServiceCredentialSchemaForPlanForbidden creates a GetServiceCredentialSchemaForPlanForbidden with default headers values
func NewGetServiceCredentialSchemaForPlanForbidden() *GetServiceCredentialSchemaForPlanForbidden {
	return &GetServiceCredentialSchemaForPlanForbidden{}
}

/*GetServiceCredentialSchemaForPlanForbidden handles this case with default header values.

If authenticated but not authorized
*/
type GetServiceCredentialSchemaForPlanForbidden struct {
}

func (o *GetServiceCredentialSchemaForPlanForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/servicecredentialschemas/{kind}/{name}][%d] getServiceCredentialSchemaForPlanForbidden ", 403)
}

func (o *GetServiceCredentialSchemaForPlanForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetServiceCredentialSchemaForPlanInternalServerError creates a GetServiceCredentialSchemaForPlanInternalServerError with default headers values
func NewGetServiceCredentialSchemaForPlanInternalServerError() *GetServiceCredentialSchemaForPlanInternalServerError {
	return &GetServiceCredentialSchemaForPlanInternalServerError{}
}

/*GetServiceCredentialSchemaForPlanInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type GetServiceCredentialSchemaForPlanInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *GetServiceCredentialSchemaForPlanInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/servicecredentialschemas/{kind}/{name}][%d] getServiceCredentialSchemaForPlanInternalServerError  %+v", 500, o.Payload)
}

func (o *GetServiceCredentialSchemaForPlanInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *GetServiceCredentialSchemaForPlanInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}