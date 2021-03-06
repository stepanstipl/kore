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

// ListIDPTypesReader is a Reader for the ListIDPTypes structure.
type ListIDPTypesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListIDPTypesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListIDPTypesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListIDPTypesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListIDPTypesOK creates a ListIDPTypesOK with default headers values
func NewListIDPTypesOK() *ListIDPTypesOK {
	return &ListIDPTypesOK{}
}

/*ListIDPTypesOK handles this case with default header values.

A list of all the possible identity provider types
*/
type ListIDPTypesOK struct {
	Payload []*models.V1IDPConfig
}

func (o *ListIDPTypesOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/idp/types][%d] listIdPTypesOK  %+v", 200, o.Payload)
}

func (o *ListIDPTypesOK) GetPayload() []*models.V1IDPConfig {
	return o.Payload
}

func (o *ListIDPTypesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListIDPTypesDefault creates a ListIDPTypesDefault with default headers values
func NewListIDPTypesDefault(code int) *ListIDPTypesDefault {
	return &ListIDPTypesDefault{
		_statusCode: code,
	}
}

/*ListIDPTypesDefault handles this case with default header values.

A generic API error containing the cause of the error
*/
type ListIDPTypesDefault struct {
	_statusCode int

	Payload *models.ApiserverError
}

// Code gets the status code for the list ID p types default response
func (o *ListIDPTypesDefault) Code() int {
	return o._statusCode
}

func (o *ListIDPTypesDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/idp/types][%d] ListIDPTypes default  %+v", o._statusCode, o.Payload)
}

func (o *ListIDPTypesDefault) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *ListIDPTypesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
