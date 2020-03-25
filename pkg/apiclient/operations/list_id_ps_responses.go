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

// ListIDPsReader is a Reader for the ListIDPs structure.
type ListIDPsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListIDPsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListIDPsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListIDPsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListIDPsOK creates a ListIDPsOK with default headers values
func NewListIDPsOK() *ListIDPsOK {
	return &ListIDPsOK{}
}

/*ListIDPsOK handles this case with default header values.

A list of all the configured identity providers
*/
type ListIDPsOK struct {
	Payload []*models.V1IDP
}

func (o *ListIDPsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/idp/configured][%d] listIdPsOK  %+v", 200, o.Payload)
}

func (o *ListIDPsOK) GetPayload() []*models.V1IDP {
	return o.Payload
}

func (o *ListIDPsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListIDPsDefault creates a ListIDPsDefault with default headers values
func NewListIDPsDefault(code int) *ListIDPsDefault {
	return &ListIDPsDefault{
		_statusCode: code,
	}
}

/*ListIDPsDefault handles this case with default header values.

A generic API error containing the cause of the error
*/
type ListIDPsDefault struct {
	_statusCode int

	Payload *models.ApiserverError
}

// Code gets the status code for the list ID ps default response
func (o *ListIDPsDefault) Code() int {
	return o._statusCode
}

func (o *ListIDPsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/idp/configured][%d] ListIDPs default  %+v", o._statusCode, o.Payload)
}

func (o *ListIDPsDefault) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *ListIDPsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}