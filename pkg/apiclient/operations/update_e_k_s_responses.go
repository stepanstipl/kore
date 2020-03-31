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

// UpdateEKSReader is a Reader for the UpdateEKS structure.
type UpdateEKSReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateEKSReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateEKSOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewUpdateEKSDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateEKSOK creates a UpdateEKSOK with default headers values
func NewUpdateEKSOK() *UpdateEKSOK {
	return &UpdateEKSOK{}
}

/*UpdateEKSOK handles this case with default header values.

Contains the former team definition from the kore
*/
type UpdateEKSOK struct {
	Payload *models.V1alpha1EKS
}

func (o *UpdateEKSOK) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/teams/{team}/ekss/{name}][%d] updateEKSOK  %+v", 200, o.Payload)
}

func (o *UpdateEKSOK) GetPayload() *models.V1alpha1EKS {
	return o.Payload
}

func (o *UpdateEKSOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1alpha1EKS)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateEKSDefault creates a UpdateEKSDefault with default headers values
func NewUpdateEKSDefault(code int) *UpdateEKSDefault {
	return &UpdateEKSDefault{
		_statusCode: code,
	}
}

/*UpdateEKSDefault handles this case with default header values.

A generic API error containing the cause of the error
*/
type UpdateEKSDefault struct {
	_statusCode int

	Payload *models.ApiserverError
}

// Code gets the status code for the update e k s default response
func (o *UpdateEKSDefault) Code() int {
	return o._statusCode
}

func (o *UpdateEKSDefault) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/teams/{team}/ekss/{name}][%d] updateEKS default  %+v", o._statusCode, o.Payload)
}

func (o *UpdateEKSDefault) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *UpdateEKSDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}