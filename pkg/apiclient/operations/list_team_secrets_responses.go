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

// ListTeamSecretsReader is a Reader for the ListTeamSecrets structure.
type ListTeamSecretsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListTeamSecretsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListTeamSecretsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewListTeamSecretsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewListTeamSecretsOK creates a ListTeamSecretsOK with default headers values
func NewListTeamSecretsOK() *ListTeamSecretsOK {
	return &ListTeamSecretsOK{}
}

/*ListTeamSecretsOK handles this case with default header values.

Contains the definition for the resource
*/
type ListTeamSecretsOK struct {
	Payload *models.V1Secret
}

func (o *ListTeamSecretsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/teams/{team}/secrets][%d] listTeamSecretsOK  %+v", 200, o.Payload)
}

func (o *ListTeamSecretsOK) GetPayload() *models.V1Secret {
	return o.Payload
}

func (o *ListTeamSecretsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1Secret)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListTeamSecretsInternalServerError creates a ListTeamSecretsInternalServerError with default headers values
func NewListTeamSecretsInternalServerError() *ListTeamSecretsInternalServerError {
	return &ListTeamSecretsInternalServerError{}
}

/*ListTeamSecretsInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type ListTeamSecretsInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *ListTeamSecretsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/teams/{team}/secrets][%d] listTeamSecretsInternalServerError  %+v", 500, o.Payload)
}

func (o *ListTeamSecretsInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *ListTeamSecretsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
