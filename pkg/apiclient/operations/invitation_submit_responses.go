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

// InvitationSubmitReader is a Reader for the InvitationSubmit structure.
type InvitationSubmitReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *InvitationSubmitReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewInvitationSubmitOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewInvitationSubmitInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewInvitationSubmitOK creates a InvitationSubmitOK with default headers values
func NewInvitationSubmitOK() *InvitationSubmitOK {
	return &InvitationSubmitOK{}
}

/*InvitationSubmitOK handles this case with default header values.

Indicates the generated link is valid and the user has been granted access
*/
type InvitationSubmitOK struct {
	Payload *models.TypesTeamInvitationResponse
}

func (o *InvitationSubmitOK) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/teams/invitation/{token}][%d] invitationSubmitOK  %+v", 200, o.Payload)
}

func (o *InvitationSubmitOK) GetPayload() *models.TypesTeamInvitationResponse {
	return o.Payload
}

func (o *InvitationSubmitOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.TypesTeamInvitationResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewInvitationSubmitInternalServerError creates a InvitationSubmitInternalServerError with default headers values
func NewInvitationSubmitInternalServerError() *InvitationSubmitInternalServerError {
	return &InvitationSubmitInternalServerError{}
}

/*InvitationSubmitInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type InvitationSubmitInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *InvitationSubmitInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /api/v1alpha1/teams/invitation/{token}][%d] invitationSubmitInternalServerError  %+v", 500, o.Payload)
}

func (o *InvitationSubmitInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *InvitationSubmitInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
