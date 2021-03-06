// Code generated by go-swagger; DO NOT EDIT.

package costs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/appvia/kore/pkg/apiclient/models"
)

// GetCostReader is a Reader for the GetCost structure.
type GetCostReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCostReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetCostOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetCostUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetCostForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetCostNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetCostInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetCostOK creates a GetCostOK with default headers values
func NewGetCostOK() *GetCostOK {
	return &GetCostOK{}
}

/*GetCostOK handles this case with default header values.

Cost found
*/
type GetCostOK struct {
	Payload *models.V1beta1Cost
}

func (o *GetCostOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/costs/{name}][%d] getCostOK  %+v", 200, o.Payload)
}

func (o *GetCostOK) GetPayload() *models.V1beta1Cost {
	return o.Payload
}

func (o *GetCostOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1Cost)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCostUnauthorized creates a GetCostUnauthorized with default headers values
func NewGetCostUnauthorized() *GetCostUnauthorized {
	return &GetCostUnauthorized{}
}

/*GetCostUnauthorized handles this case with default header values.

If not authenticated
*/
type GetCostUnauthorized struct {
}

func (o *GetCostUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/costs/{name}][%d] getCostUnauthorized ", 401)
}

func (o *GetCostUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetCostForbidden creates a GetCostForbidden with default headers values
func NewGetCostForbidden() *GetCostForbidden {
	return &GetCostForbidden{}
}

/*GetCostForbidden handles this case with default header values.

If authenticated but not authorized
*/
type GetCostForbidden struct {
}

func (o *GetCostForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/costs/{name}][%d] getCostForbidden ", 403)
}

func (o *GetCostForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetCostNotFound creates a GetCostNotFound with default headers values
func NewGetCostNotFound() *GetCostNotFound {
	return &GetCostNotFound{}
}

/*GetCostNotFound handles this case with default header values.

Cost doesn't exist
*/
type GetCostNotFound struct {
}

func (o *GetCostNotFound) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/costs/{name}][%d] getCostNotFound ", 404)
}

func (o *GetCostNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetCostInternalServerError creates a GetCostInternalServerError with default headers values
func NewGetCostInternalServerError() *GetCostInternalServerError {
	return &GetCostInternalServerError{}
}

/*GetCostInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type GetCostInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *GetCostInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/costs/{name}][%d] getCostInternalServerError  %+v", 500, o.Payload)
}

func (o *GetCostInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *GetCostInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
