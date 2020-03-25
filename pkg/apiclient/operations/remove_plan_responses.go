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

// RemovePlanReader is a Reader for the RemovePlan structure.
type RemovePlanReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RemovePlanReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewRemovePlanOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRemovePlanDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRemovePlanOK creates a RemovePlanOK with default headers values
func NewRemovePlanOK() *RemovePlanOK {
	return &RemovePlanOK{}
}

/*RemovePlanOK handles this case with default header values.

Contains the class definintion from the kore
*/
type RemovePlanOK struct {
	Payload *models.V1Plan
}

func (o *RemovePlanOK) Error() string {
	return fmt.Sprintf("[DELETE /api/v1alpha1/plans/{name}][%d] removePlanOK  %+v", 200, o.Payload)
}

func (o *RemovePlanOK) GetPayload() *models.V1Plan {
	return o.Payload
}

func (o *RemovePlanOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1Plan)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRemovePlanDefault creates a RemovePlanDefault with default headers values
func NewRemovePlanDefault(code int) *RemovePlanDefault {
	return &RemovePlanDefault{
		_statusCode: code,
	}
}

/*RemovePlanDefault handles this case with default header values.

A generic API error containing the cause of the error
*/
type RemovePlanDefault struct {
	_statusCode int

	Payload *models.ApiserverError
}

// Code gets the status code for the remove plan default response
func (o *RemovePlanDefault) Code() int {
	return o._statusCode
}

func (o *RemovePlanDefault) Error() string {
	return fmt.Sprintf("[DELETE /api/v1alpha1/plans/{name}][%d] RemovePlan default  %+v", o._statusCode, o.Payload)
}

func (o *RemovePlanDefault) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *RemovePlanDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}