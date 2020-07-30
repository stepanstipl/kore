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

// ListUserIdentityReader is a Reader for the ListUserIdentity structure.
type ListUserIdentityReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListUserIdentityReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListUserIdentityOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewListUserIdentityUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewListUserIdentityForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewListUserIdentityNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewListUserIdentityInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewListUserIdentityOK creates a ListUserIdentityOK with default headers values
func NewListUserIdentityOK() *ListUserIdentityOK {
	return &ListUserIdentityOK{}
}

/*ListUserIdentityOK handles this case with default header values.

Contains the identities definitions from the kore
*/
type ListUserIdentityOK struct {
	Payload *models.V1IdentityList
}

func (o *ListUserIdentityOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/identities/{user}][%d] listUserIdentityOK  %+v", 200, o.Payload)
}

func (o *ListUserIdentityOK) GetPayload() *models.V1IdentityList {
	return o.Payload
}

func (o *ListUserIdentityOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1IdentityList)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListUserIdentityUnauthorized creates a ListUserIdentityUnauthorized with default headers values
func NewListUserIdentityUnauthorized() *ListUserIdentityUnauthorized {
	return &ListUserIdentityUnauthorized{}
}

/*ListUserIdentityUnauthorized handles this case with default header values.

If not authenticated
*/
type ListUserIdentityUnauthorized struct {
}

func (o *ListUserIdentityUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/identities/{user}][%d] listUserIdentityUnauthorized ", 401)
}

func (o *ListUserIdentityUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewListUserIdentityForbidden creates a ListUserIdentityForbidden with default headers values
func NewListUserIdentityForbidden() *ListUserIdentityForbidden {
	return &ListUserIdentityForbidden{}
}

/*ListUserIdentityForbidden handles this case with default header values.

If authenticated but not authorized
*/
type ListUserIdentityForbidden struct {
}

func (o *ListUserIdentityForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/identities/{user}][%d] listUserIdentityForbidden ", 403)
}

func (o *ListUserIdentityForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewListUserIdentityNotFound creates a ListUserIdentityNotFound with default headers values
func NewListUserIdentityNotFound() *ListUserIdentityNotFound {
	return &ListUserIdentityNotFound{}
}

/*ListUserIdentityNotFound handles this case with default header values.

User does not exist
*/
type ListUserIdentityNotFound struct {
}

func (o *ListUserIdentityNotFound) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/identities/{user}][%d] listUserIdentityNotFound ", 404)
}

func (o *ListUserIdentityNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewListUserIdentityInternalServerError creates a ListUserIdentityInternalServerError with default headers values
func NewListUserIdentityInternalServerError() *ListUserIdentityInternalServerError {
	return &ListUserIdentityInternalServerError{}
}

/*ListUserIdentityInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type ListUserIdentityInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *ListUserIdentityInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/identities/{user}][%d] listUserIdentityInternalServerError  %+v", 500, o.Payload)
}

func (o *ListUserIdentityInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *ListUserIdentityInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
