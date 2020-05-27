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

// GetConfigReader is a Reader for the GetConfig structure.
type GetConfigReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetConfigReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetConfigOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetConfigUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetConfigForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetConfigNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetConfigInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetConfigOK creates a GetConfigOK with default headers values
func NewGetConfigOK() *GetConfigOK {
	return &GetConfigOK{}
}

/*GetConfigOK handles this case with default header values.

A list of all the config in the kore
*/
type GetConfigOK struct {
	Payload *models.V1Config
}

func (o *GetConfigOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/configs/{config}][%d] getConfigOK  %+v", 200, o.Payload)
}

func (o *GetConfigOK) GetPayload() *models.V1Config {
	return o.Payload
}

func (o *GetConfigOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1Config)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetConfigUnauthorized creates a GetConfigUnauthorized with default headers values
func NewGetConfigUnauthorized() *GetConfigUnauthorized {
	return &GetConfigUnauthorized{}
}

/*GetConfigUnauthorized handles this case with default header values.

If not authenticated
*/
type GetConfigUnauthorized struct {
}

func (o *GetConfigUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/configs/{config}][%d] getConfigUnauthorized ", 401)
}

func (o *GetConfigUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetConfigForbidden creates a GetConfigForbidden with default headers values
func NewGetConfigForbidden() *GetConfigForbidden {
	return &GetConfigForbidden{}
}

/*GetConfigForbidden handles this case with default header values.

If authenticated but not authorized
*/
type GetConfigForbidden struct {
}

func (o *GetConfigForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/configs/{config}][%d] getConfigForbidden ", 403)
}

func (o *GetConfigForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetConfigNotFound creates a GetConfigNotFound with default headers values
func NewGetConfigNotFound() *GetConfigNotFound {
	return &GetConfigNotFound{}
}

/*GetConfigNotFound handles this case with default header values.

config does not exist
*/
type GetConfigNotFound struct {
}

func (o *GetConfigNotFound) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/configs/{config}][%d] getConfigNotFound ", 404)
}

func (o *GetConfigNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetConfigInternalServerError creates a GetConfigInternalServerError with default headers values
func NewGetConfigInternalServerError() *GetConfigInternalServerError {
	return &GetConfigInternalServerError{}
}

/*GetConfigInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type GetConfigInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *GetConfigInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/configs/{config}][%d] getConfigInternalServerError  %+v", 500, o.Payload)
}

func (o *GetConfigInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *GetConfigInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
