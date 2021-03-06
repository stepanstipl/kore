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

// UpdateAccountReader is a Reader for the UpdateAccount structure.
type UpdateAccountReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateAccountReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateAccountOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdateAccountBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewUpdateAccountUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewUpdateAccountForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateAccountInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUpdateAccountOK creates a UpdateAccountOK with default headers values
func NewUpdateAccountOK() *UpdateAccountOK {
	return &UpdateAccountOK{}
}

/*UpdateAccountOK handles this case with default header values.

Contains the account definition
*/
type UpdateAccountOK struct {
	Payload *models.V1beta1AccountManagement
}

func (o *UpdateAccountOK) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/accountmanagements/{name}][%d] updateAccountOK  %+v", 200, o.Payload)
}

func (o *UpdateAccountOK) GetPayload() *models.V1beta1AccountManagement {
	return o.Payload
}

func (o *UpdateAccountOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1AccountManagement)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateAccountBadRequest creates a UpdateAccountBadRequest with default headers values
func NewUpdateAccountBadRequest() *UpdateAccountBadRequest {
	return &UpdateAccountBadRequest{}
}

/*UpdateAccountBadRequest handles this case with default header values.

Validation error of supplied parameters/body
*/
type UpdateAccountBadRequest struct {
	Payload *models.ValidationError
}

func (o *UpdateAccountBadRequest) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/accountmanagements/{name}][%d] updateAccountBadRequest  %+v", 400, o.Payload)
}

func (o *UpdateAccountBadRequest) GetPayload() *models.ValidationError {
	return o.Payload
}

func (o *UpdateAccountBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateAccountUnauthorized creates a UpdateAccountUnauthorized with default headers values
func NewUpdateAccountUnauthorized() *UpdateAccountUnauthorized {
	return &UpdateAccountUnauthorized{}
}

/*UpdateAccountUnauthorized handles this case with default header values.

If not authenticated
*/
type UpdateAccountUnauthorized struct {
}

func (o *UpdateAccountUnauthorized) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/accountmanagements/{name}][%d] updateAccountUnauthorized ", 401)
}

func (o *UpdateAccountUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateAccountForbidden creates a UpdateAccountForbidden with default headers values
func NewUpdateAccountForbidden() *UpdateAccountForbidden {
	return &UpdateAccountForbidden{}
}

/*UpdateAccountForbidden handles this case with default header values.

If authenticated but not authorized
*/
type UpdateAccountForbidden struct {
}

func (o *UpdateAccountForbidden) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/accountmanagements/{name}][%d] updateAccountForbidden ", 403)
}

func (o *UpdateAccountForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateAccountInternalServerError creates a UpdateAccountInternalServerError with default headers values
func NewUpdateAccountInternalServerError() *UpdateAccountInternalServerError {
	return &UpdateAccountInternalServerError{}
}

/*UpdateAccountInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type UpdateAccountInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *UpdateAccountInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/accountmanagements/{name}][%d] updateAccountInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateAccountInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *UpdateAccountInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
