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

// GetServicePlanReader is a Reader for the GetServicePlan structure.
type GetServicePlanReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetServicePlanReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetServicePlanOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetServicePlanUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetServicePlanForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetServicePlanNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetServicePlanInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetServicePlanOK creates a GetServicePlanOK with default headers values
func NewGetServicePlanOK() *GetServicePlanOK {
	return &GetServicePlanOK{}
}

/*GetServicePlanOK handles this case with default header values.

Contains the service plan definition
*/
type GetServicePlanOK struct {
	Payload *models.V1ServicePlan
}

func (o *GetServicePlanOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/serviceplans/{name}][%d] getServicePlanOK  %+v", 200, o.Payload)
}

func (o *GetServicePlanOK) GetPayload() *models.V1ServicePlan {
	return o.Payload
}

func (o *GetServicePlanOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1ServicePlan)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetServicePlanUnauthorized creates a GetServicePlanUnauthorized with default headers values
func NewGetServicePlanUnauthorized() *GetServicePlanUnauthorized {
	return &GetServicePlanUnauthorized{}
}

/*GetServicePlanUnauthorized handles this case with default header values.

If not authenticated
*/
type GetServicePlanUnauthorized struct {
}

func (o *GetServicePlanUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/serviceplans/{name}][%d] getServicePlanUnauthorized ", 401)
}

func (o *GetServicePlanUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetServicePlanForbidden creates a GetServicePlanForbidden with default headers values
func NewGetServicePlanForbidden() *GetServicePlanForbidden {
	return &GetServicePlanForbidden{}
}

/*GetServicePlanForbidden handles this case with default header values.

If authenticated but not authorized
*/
type GetServicePlanForbidden struct {
}

func (o *GetServicePlanForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/serviceplans/{name}][%d] getServicePlanForbidden ", 403)
}

func (o *GetServicePlanForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetServicePlanNotFound creates a GetServicePlanNotFound with default headers values
func NewGetServicePlanNotFound() *GetServicePlanNotFound {
	return &GetServicePlanNotFound{}
}

/*GetServicePlanNotFound handles this case with default header values.

the service plan with the given name doesn't exist
*/
type GetServicePlanNotFound struct {
}

func (o *GetServicePlanNotFound) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/serviceplans/{name}][%d] getServicePlanNotFound ", 404)
}

func (o *GetServicePlanNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetServicePlanInternalServerError creates a GetServicePlanInternalServerError with default headers values
func NewGetServicePlanInternalServerError() *GetServicePlanInternalServerError {
	return &GetServicePlanInternalServerError{}
}

/*GetServicePlanInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type GetServicePlanInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *GetServicePlanInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/serviceplans/{name}][%d] getServicePlanInternalServerError  %+v", 500, o.Payload)
}

func (o *GetServicePlanInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *GetServicePlanInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}