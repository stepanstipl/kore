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

// ListAWSAccountOUsReader is a Reader for the ListAWSAccountOUs structure.
type ListAWSAccountOUsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListAWSAccountOUsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListAWSAccountOUsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListAWSAccountOUsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListAWSAccountOUsOK creates a ListAWSAccountOUsOK with default headers values
func NewListAWSAccountOUsOK() *ListAWSAccountOUsOK {
	return &ListAWSAccountOUsOK{}
}

/*ListAWSAccountOUsOK handles this case with default header values.

Account OUs
*/
type ListAWSAccountOUsOK struct {
	Payload []string
}

func (o *ListAWSAccountOUsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/teams/{team}/awsorganizations/awsAccountOUs][%d] listAWSAccountOUsOK  %+v", 200, o.Payload)
}

func (o *ListAWSAccountOUsOK) GetPayload() []string {
	return o.Payload
}

func (o *ListAWSAccountOUsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListAWSAccountOUsDefault creates a ListAWSAccountOUsDefault with default headers values
func NewListAWSAccountOUsDefault(code int) *ListAWSAccountOUsDefault {
	return &ListAWSAccountOUsDefault{
		_statusCode: code,
	}
}

/*ListAWSAccountOUsDefault handles this case with default header values.

A generic API error containing the cause of the error
*/
type ListAWSAccountOUsDefault struct {
	_statusCode int

	Payload *models.ApiserverError
}

// Code gets the status code for the list a w s account o us default response
func (o *ListAWSAccountOUsDefault) Code() int {
	return o._statusCode
}

func (o *ListAWSAccountOUsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/teams/{team}/awsorganizations/awsAccountOUs][%d] ListAWSAccountOUs default  %+v", o._statusCode, o.Payload)
}

func (o *ListAWSAccountOUsDefault) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *ListAWSAccountOUsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
