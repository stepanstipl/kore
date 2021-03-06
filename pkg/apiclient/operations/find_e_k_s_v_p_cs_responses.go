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

// FindEKSVPCsReader is a Reader for the FindEKSVPCs structure.
type FindEKSVPCsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *FindEKSVPCsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewFindEKSVPCsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewFindEKSVPCsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewFindEKSVPCsForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewFindEKSVPCsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewFindEKSVPCsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewFindEKSVPCsOK creates a FindEKSVPCsOK with default headers values
func NewFindEKSVPCsOK() *FindEKSVPCsOK {
	return &FindEKSVPCsOK{}
}

/*FindEKSVPCsOK handles this case with default header values.

Contains the former team definition from the kore
*/
type FindEKSVPCsOK struct {
	Payload *models.V1alpha1EKSVPCList
}

func (o *FindEKSVPCsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/teams/{team}/eksvpcs][%d] findEKSVPCsOK  %+v", 200, o.Payload)
}

func (o *FindEKSVPCsOK) GetPayload() *models.V1alpha1EKSVPCList {
	return o.Payload
}

func (o *FindEKSVPCsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1alpha1EKSVPCList)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewFindEKSVPCsUnauthorized creates a FindEKSVPCsUnauthorized with default headers values
func NewFindEKSVPCsUnauthorized() *FindEKSVPCsUnauthorized {
	return &FindEKSVPCsUnauthorized{}
}

/*FindEKSVPCsUnauthorized handles this case with default header values.

If not authenticated
*/
type FindEKSVPCsUnauthorized struct {
}

func (o *FindEKSVPCsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/teams/{team}/eksvpcs][%d] findEKSVPCsUnauthorized ", 401)
}

func (o *FindEKSVPCsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewFindEKSVPCsForbidden creates a FindEKSVPCsForbidden with default headers values
func NewFindEKSVPCsForbidden() *FindEKSVPCsForbidden {
	return &FindEKSVPCsForbidden{}
}

/*FindEKSVPCsForbidden handles this case with default header values.

If authenticated but not authorized
*/
type FindEKSVPCsForbidden struct {
}

func (o *FindEKSVPCsForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/teams/{team}/eksvpcs][%d] findEKSVPCsForbidden ", 403)
}

func (o *FindEKSVPCsForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewFindEKSVPCsInternalServerError creates a FindEKSVPCsInternalServerError with default headers values
func NewFindEKSVPCsInternalServerError() *FindEKSVPCsInternalServerError {
	return &FindEKSVPCsInternalServerError{}
}

/*FindEKSVPCsInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type FindEKSVPCsInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *FindEKSVPCsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/teams/{team}/eksvpcs][%d] findEKSVPCsInternalServerError  %+v", 500, o.Payload)
}

func (o *FindEKSVPCsInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *FindEKSVPCsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewFindEKSVPCsDefault creates a FindEKSVPCsDefault with default headers values
func NewFindEKSVPCsDefault(code int) *FindEKSVPCsDefault {
	return &FindEKSVPCsDefault{
		_statusCode: code,
	}
}

/*FindEKSVPCsDefault handles this case with default header values.

A generic API error containing the cause of the error
*/
type FindEKSVPCsDefault struct {
	_statusCode int

	Payload *models.ApiserverError
}

// Code gets the status code for the find e k s v p cs default response
func (o *FindEKSVPCsDefault) Code() int {
	return o._statusCode
}

func (o *FindEKSVPCsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/teams/{team}/eksvpcs][%d] findEKSVPCs default  %+v", o._statusCode, o.Payload)
}

func (o *FindEKSVPCsDefault) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *FindEKSVPCsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
