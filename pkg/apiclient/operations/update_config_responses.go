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

// UpdateConfigReader is a Reader for the UpdateConfig structure.
type UpdateConfigReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateConfigReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateConfigOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdateConfigBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewUpdateConfigUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewUpdateConfigForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUpdateConfigNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateConfigInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUpdateConfigOK creates a UpdateConfigOK with default headers values
func NewUpdateConfigOK() *UpdateConfigOK {
	return &UpdateConfigOK{}
}

/*UpdateConfigOK handles this case with default header values.

Contains the config definition
*/
type UpdateConfigOK struct {
	Payload *models.V1Config
}

func (o *UpdateConfigOK) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/configs/{config}][%d] updateConfigOK  %+v", 200, o.Payload)
}

func (o *UpdateConfigOK) GetPayload() *models.V1Config {
	return o.Payload
}

func (o *UpdateConfigOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1Config)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateConfigBadRequest creates a UpdateConfigBadRequest with default headers values
func NewUpdateConfigBadRequest() *UpdateConfigBadRequest {
	return &UpdateConfigBadRequest{}
}

/*UpdateConfigBadRequest handles this case with default header values.

Validation error of supplied parameters/body
*/
type UpdateConfigBadRequest struct {
	Payload *models.ValidationError
}

func (o *UpdateConfigBadRequest) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/configs/{config}][%d] updateConfigBadRequest  %+v", 400, o.Payload)
}

func (o *UpdateConfigBadRequest) GetPayload() *models.ValidationError {
	return o.Payload
}

func (o *UpdateConfigBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateConfigUnauthorized creates a UpdateConfigUnauthorized with default headers values
func NewUpdateConfigUnauthorized() *UpdateConfigUnauthorized {
	return &UpdateConfigUnauthorized{}
}

/*UpdateConfigUnauthorized handles this case with default header values.

If not authenticated
*/
type UpdateConfigUnauthorized struct {
}

func (o *UpdateConfigUnauthorized) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/configs/{config}][%d] updateConfigUnauthorized ", 401)
}

func (o *UpdateConfigUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateConfigForbidden creates a UpdateConfigForbidden with default headers values
func NewUpdateConfigForbidden() *UpdateConfigForbidden {
	return &UpdateConfigForbidden{}
}

/*UpdateConfigForbidden handles this case with default header values.

If authenticated but not authorized
*/
type UpdateConfigForbidden struct {
}

func (o *UpdateConfigForbidden) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/configs/{config}][%d] updateConfigForbidden ", 403)
}

func (o *UpdateConfigForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateConfigNotFound creates a UpdateConfigNotFound with default headers values
func NewUpdateConfigNotFound() *UpdateConfigNotFound {
	return &UpdateConfigNotFound{}
}

/*UpdateConfigNotFound handles this case with default header values.

config does not exist
*/
type UpdateConfigNotFound struct {
}

func (o *UpdateConfigNotFound) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/configs/{config}][%d] updateConfigNotFound ", 404)
}

func (o *UpdateConfigNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateConfigInternalServerError creates a UpdateConfigInternalServerError with default headers values
func NewUpdateConfigInternalServerError() *UpdateConfigInternalServerError {
	return &UpdateConfigInternalServerError{}
}

/*UpdateConfigInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type UpdateConfigInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *UpdateConfigInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/configs/{config}][%d] updateConfigInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateConfigInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *UpdateConfigInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
