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

// GetAKSReader is a Reader for the GetAKS structure.
type GetAKSReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAKSReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAKSOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetAKSDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetAKSOK creates a GetAKSOK with default headers values
func NewGetAKSOK() *GetAKSOK {
	return &GetAKSOK{}
}

/*GetAKSOK handles this case with default header values.

Contains the definition of the AKS cluster
*/
type GetAKSOK struct {
	Payload *models.V1alpha1AKS
}

func (o *GetAKSOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/teams/{team}/aks/{name}][%d] getAKSOK  %+v", 200, o.Payload)
}

func (o *GetAKSOK) GetPayload() *models.V1alpha1AKS {
	return o.Payload
}

func (o *GetAKSOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1alpha1AKS)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAKSDefault creates a GetAKSDefault with default headers values
func NewGetAKSDefault(code int) *GetAKSDefault {
	return &GetAKSDefault{
		_statusCode: code,
	}
}

/*GetAKSDefault handles this case with default header values.

A generic API error containing the cause of the error
*/
type GetAKSDefault struct {
	_statusCode int

	Payload *models.ApiserverError
}

// Code gets the status code for the get a k s default response
func (o *GetAKSDefault) Code() int {
	return o._statusCode
}

func (o *GetAKSDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/teams/{team}/aks/{name}][%d] getAKS default  %+v", o._statusCode, o.Payload)
}

func (o *GetAKSDefault) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *GetAKSDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}