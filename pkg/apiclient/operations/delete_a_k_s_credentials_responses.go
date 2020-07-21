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

// DeleteAKSCredentialsReader is a Reader for the DeleteAKSCredentials structure.
type DeleteAKSCredentialsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteAKSCredentialsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteAKSCredentialsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDeleteAKSCredentialsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteAKSCredentialsOK creates a DeleteAKSCredentialsOK with default headers values
func NewDeleteAKSCredentialsOK() *DeleteAKSCredentialsOK {
	return &DeleteAKSCredentialsOK{}
}

/*DeleteAKSCredentialsOK handles this case with default header values.

Contains the former AKS credentials
*/
type DeleteAKSCredentialsOK struct {
	Payload *models.V1alpha1AKSCredentials
}

func (o *DeleteAKSCredentialsOK) Error() string {
	return fmt.Sprintf("[DELETE /api/v1alpha1/teams/{team}/akscredentials/{name}][%d] deleteAKSCredentialsOK  %+v", 200, o.Payload)
}

func (o *DeleteAKSCredentialsOK) GetPayload() *models.V1alpha1AKSCredentials {
	return o.Payload
}

func (o *DeleteAKSCredentialsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1alpha1AKSCredentials)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteAKSCredentialsDefault creates a DeleteAKSCredentialsDefault with default headers values
func NewDeleteAKSCredentialsDefault(code int) *DeleteAKSCredentialsDefault {
	return &DeleteAKSCredentialsDefault{
		_statusCode: code,
	}
}

/*DeleteAKSCredentialsDefault handles this case with default header values.

A generic API error containing the cause of the error
*/
type DeleteAKSCredentialsDefault struct {
	_statusCode int

	Payload *models.ApiserverError
}

// Code gets the status code for the delete a k s credentials default response
func (o *DeleteAKSCredentialsDefault) Code() int {
	return o._statusCode
}

func (o *DeleteAKSCredentialsDefault) Error() string {
	return fmt.Sprintf("[DELETE /api/v1alpha1/teams/{team}/akscredentials/{name}][%d] DeleteAKSCredentials default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteAKSCredentialsDefault) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *DeleteAKSCredentialsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}