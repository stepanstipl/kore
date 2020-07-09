// Code generated by go-swagger; DO NOT EDIT.

package metadata

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/appvia/kore/pkg/apiclient/models"
)

// GetCloudRegionsReader is a Reader for the GetCloudRegions structure.
type GetCloudRegionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCloudRegionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetCloudRegionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetCloudRegionsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetCloudRegionsForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetCloudRegionsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetCloudRegionsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetCloudRegionsOK creates a GetCloudRegionsOK with default headers values
func NewGetCloudRegionsOK() *GetCloudRegionsOK {
	return &GetCloudRegionsOK{}
}

/*GetCloudRegionsOK handles this case with default header values.

A list of all the regions organised by continent
*/
type GetCloudRegionsOK struct {
	Payload *models.V1beta1ContinentList
}

func (o *GetCloudRegionsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/metadata/cloud/{cloud}/regions][%d] getCloudRegionsOK  %+v", 200, o.Payload)
}

func (o *GetCloudRegionsOK) GetPayload() *models.V1beta1ContinentList {
	return o.Payload
}

func (o *GetCloudRegionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.V1beta1ContinentList)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCloudRegionsUnauthorized creates a GetCloudRegionsUnauthorized with default headers values
func NewGetCloudRegionsUnauthorized() *GetCloudRegionsUnauthorized {
	return &GetCloudRegionsUnauthorized{}
}

/*GetCloudRegionsUnauthorized handles this case with default header values.

If not authenticated
*/
type GetCloudRegionsUnauthorized struct {
}

func (o *GetCloudRegionsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/metadata/cloud/{cloud}/regions][%d] getCloudRegionsUnauthorized ", 401)
}

func (o *GetCloudRegionsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetCloudRegionsForbidden creates a GetCloudRegionsForbidden with default headers values
func NewGetCloudRegionsForbidden() *GetCloudRegionsForbidden {
	return &GetCloudRegionsForbidden{}
}

/*GetCloudRegionsForbidden handles this case with default header values.

If authenticated but not authorized
*/
type GetCloudRegionsForbidden struct {
}

func (o *GetCloudRegionsForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/metadata/cloud/{cloud}/regions][%d] getCloudRegionsForbidden ", 403)
}

func (o *GetCloudRegionsForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetCloudRegionsNotFound creates a GetCloudRegionsNotFound with default headers values
func NewGetCloudRegionsNotFound() *GetCloudRegionsNotFound {
	return &GetCloudRegionsNotFound{}
}

/*GetCloudRegionsNotFound handles this case with default header values.

cloud doesn't exist
*/
type GetCloudRegionsNotFound struct {
}

func (o *GetCloudRegionsNotFound) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/metadata/cloud/{cloud}/regions][%d] getCloudRegionsNotFound ", 404)
}

func (o *GetCloudRegionsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetCloudRegionsInternalServerError creates a GetCloudRegionsInternalServerError with default headers values
func NewGetCloudRegionsInternalServerError() *GetCloudRegionsInternalServerError {
	return &GetCloudRegionsInternalServerError{}
}

/*GetCloudRegionsInternalServerError handles this case with default header values.

A generic API error containing the cause of the error
*/
type GetCloudRegionsInternalServerError struct {
	Payload *models.ApiserverError
}

func (o *GetCloudRegionsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /api/v1alpha1/metadata/cloud/{cloud}/regions][%d] getCloudRegionsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetCloudRegionsInternalServerError) GetPayload() *models.ApiserverError {
	return o.Payload
}

func (o *GetCloudRegionsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ApiserverError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
