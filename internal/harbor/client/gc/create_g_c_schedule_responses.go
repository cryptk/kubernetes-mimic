// Code generated by go-swagger; DO NOT EDIT.

package gc

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/cryptk/kubernetes-mimic/internal/harbor/models"
)

// CreateGCScheduleReader is a Reader for the CreateGCSchedule structure.
type CreateGCScheduleReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateGCScheduleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreateGCScheduleCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewCreateGCScheduleBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewCreateGCScheduleUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewCreateGCScheduleForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewCreateGCScheduleConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewCreateGCScheduleInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewCreateGCScheduleCreated creates a CreateGCScheduleCreated with default headers values
func NewCreateGCScheduleCreated() *CreateGCScheduleCreated {
	return &CreateGCScheduleCreated{}
}

/* CreateGCScheduleCreated describes a response with status code 201, with default header values.

Created
*/
type CreateGCScheduleCreated struct {

	/* The location of the resource
	 */
	Location string

	/* The ID of the corresponding request for the response
	 */
	XRequestID string
}

func (o *CreateGCScheduleCreated) Error() string {
	return fmt.Sprintf("[POST /system/gc/schedule][%d] createGCScheduleCreated ", 201)
}

func (o *CreateGCScheduleCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header Location
	hdrLocation := response.GetHeader("Location")

	if hdrLocation != "" {
		o.Location = hdrLocation
	}

	// hydrates response header X-Request-Id
	hdrXRequestID := response.GetHeader("X-Request-Id")

	if hdrXRequestID != "" {
		o.XRequestID = hdrXRequestID
	}

	return nil
}

// NewCreateGCScheduleBadRequest creates a CreateGCScheduleBadRequest with default headers values
func NewCreateGCScheduleBadRequest() *CreateGCScheduleBadRequest {
	return &CreateGCScheduleBadRequest{}
}

/* CreateGCScheduleBadRequest describes a response with status code 400, with default header values.

Bad request
*/
type CreateGCScheduleBadRequest struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *CreateGCScheduleBadRequest) Error() string {
	return fmt.Sprintf("[POST /system/gc/schedule][%d] createGCScheduleBadRequest  %+v", 400, o.Payload)
}
func (o *CreateGCScheduleBadRequest) GetPayload() *models.Errors {
	return o.Payload
}

func (o *CreateGCScheduleBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header X-Request-Id
	hdrXRequestID := response.GetHeader("X-Request-Id")

	if hdrXRequestID != "" {
		o.XRequestID = hdrXRequestID
	}

	o.Payload = new(models.Errors)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateGCScheduleUnauthorized creates a CreateGCScheduleUnauthorized with default headers values
func NewCreateGCScheduleUnauthorized() *CreateGCScheduleUnauthorized {
	return &CreateGCScheduleUnauthorized{}
}

/* CreateGCScheduleUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type CreateGCScheduleUnauthorized struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *CreateGCScheduleUnauthorized) Error() string {
	return fmt.Sprintf("[POST /system/gc/schedule][%d] createGCScheduleUnauthorized  %+v", 401, o.Payload)
}
func (o *CreateGCScheduleUnauthorized) GetPayload() *models.Errors {
	return o.Payload
}

func (o *CreateGCScheduleUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header X-Request-Id
	hdrXRequestID := response.GetHeader("X-Request-Id")

	if hdrXRequestID != "" {
		o.XRequestID = hdrXRequestID
	}

	o.Payload = new(models.Errors)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateGCScheduleForbidden creates a CreateGCScheduleForbidden with default headers values
func NewCreateGCScheduleForbidden() *CreateGCScheduleForbidden {
	return &CreateGCScheduleForbidden{}
}

/* CreateGCScheduleForbidden describes a response with status code 403, with default header values.

Forbidden
*/
type CreateGCScheduleForbidden struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *CreateGCScheduleForbidden) Error() string {
	return fmt.Sprintf("[POST /system/gc/schedule][%d] createGCScheduleForbidden  %+v", 403, o.Payload)
}
func (o *CreateGCScheduleForbidden) GetPayload() *models.Errors {
	return o.Payload
}

func (o *CreateGCScheduleForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header X-Request-Id
	hdrXRequestID := response.GetHeader("X-Request-Id")

	if hdrXRequestID != "" {
		o.XRequestID = hdrXRequestID
	}

	o.Payload = new(models.Errors)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateGCScheduleConflict creates a CreateGCScheduleConflict with default headers values
func NewCreateGCScheduleConflict() *CreateGCScheduleConflict {
	return &CreateGCScheduleConflict{}
}

/* CreateGCScheduleConflict describes a response with status code 409, with default header values.

Conflict
*/
type CreateGCScheduleConflict struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *CreateGCScheduleConflict) Error() string {
	return fmt.Sprintf("[POST /system/gc/schedule][%d] createGCScheduleConflict  %+v", 409, o.Payload)
}
func (o *CreateGCScheduleConflict) GetPayload() *models.Errors {
	return o.Payload
}

func (o *CreateGCScheduleConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header X-Request-Id
	hdrXRequestID := response.GetHeader("X-Request-Id")

	if hdrXRequestID != "" {
		o.XRequestID = hdrXRequestID
	}

	o.Payload = new(models.Errors)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateGCScheduleInternalServerError creates a CreateGCScheduleInternalServerError with default headers values
func NewCreateGCScheduleInternalServerError() *CreateGCScheduleInternalServerError {
	return &CreateGCScheduleInternalServerError{}
}

/* CreateGCScheduleInternalServerError describes a response with status code 500, with default header values.

Internal server error
*/
type CreateGCScheduleInternalServerError struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *CreateGCScheduleInternalServerError) Error() string {
	return fmt.Sprintf("[POST /system/gc/schedule][%d] createGCScheduleInternalServerError  %+v", 500, o.Payload)
}
func (o *CreateGCScheduleInternalServerError) GetPayload() *models.Errors {
	return o.Payload
}

func (o *CreateGCScheduleInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header X-Request-Id
	hdrXRequestID := response.GetHeader("X-Request-Id")

	if hdrXRequestID != "" {
		o.XRequestID = hdrXRequestID
	}

	o.Payload = new(models.Errors)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}