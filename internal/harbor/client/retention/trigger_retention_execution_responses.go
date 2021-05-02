// Code generated by go-swagger; DO NOT EDIT.

package retention

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/cryptk/kubernetes-mimic/internal/harbor/models"
)

// TriggerRetentionExecutionReader is a Reader for the TriggerRetentionExecution structure.
type TriggerRetentionExecutionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *TriggerRetentionExecutionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewTriggerRetentionExecutionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 201:
		result := NewTriggerRetentionExecutionCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewTriggerRetentionExecutionUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewTriggerRetentionExecutionForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewTriggerRetentionExecutionInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewTriggerRetentionExecutionOK creates a TriggerRetentionExecutionOK with default headers values
func NewTriggerRetentionExecutionOK() *TriggerRetentionExecutionOK {
	return &TriggerRetentionExecutionOK{}
}

/* TriggerRetentionExecutionOK describes a response with status code 200, with default header values.

Trigger a Retention job successfully.
*/
type TriggerRetentionExecutionOK struct {
}

func (o *TriggerRetentionExecutionOK) Error() string {
	return fmt.Sprintf("[POST /retentions/{id}/executions][%d] triggerRetentionExecutionOK ", 200)
}

func (o *TriggerRetentionExecutionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewTriggerRetentionExecutionCreated creates a TriggerRetentionExecutionCreated with default headers values
func NewTriggerRetentionExecutionCreated() *TriggerRetentionExecutionCreated {
	return &TriggerRetentionExecutionCreated{}
}

/* TriggerRetentionExecutionCreated describes a response with status code 201, with default header values.

Created
*/
type TriggerRetentionExecutionCreated struct {

	/* The location of the resource
	 */
	Location string

	/* The ID of the corresponding request for the response
	 */
	XRequestID string
}

func (o *TriggerRetentionExecutionCreated) Error() string {
	return fmt.Sprintf("[POST /retentions/{id}/executions][%d] triggerRetentionExecutionCreated ", 201)
}

func (o *TriggerRetentionExecutionCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewTriggerRetentionExecutionUnauthorized creates a TriggerRetentionExecutionUnauthorized with default headers values
func NewTriggerRetentionExecutionUnauthorized() *TriggerRetentionExecutionUnauthorized {
	return &TriggerRetentionExecutionUnauthorized{}
}

/* TriggerRetentionExecutionUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type TriggerRetentionExecutionUnauthorized struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *TriggerRetentionExecutionUnauthorized) Error() string {
	return fmt.Sprintf("[POST /retentions/{id}/executions][%d] triggerRetentionExecutionUnauthorized  %+v", 401, o.Payload)
}
func (o *TriggerRetentionExecutionUnauthorized) GetPayload() *models.Errors {
	return o.Payload
}

func (o *TriggerRetentionExecutionUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewTriggerRetentionExecutionForbidden creates a TriggerRetentionExecutionForbidden with default headers values
func NewTriggerRetentionExecutionForbidden() *TriggerRetentionExecutionForbidden {
	return &TriggerRetentionExecutionForbidden{}
}

/* TriggerRetentionExecutionForbidden describes a response with status code 403, with default header values.

Forbidden
*/
type TriggerRetentionExecutionForbidden struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *TriggerRetentionExecutionForbidden) Error() string {
	return fmt.Sprintf("[POST /retentions/{id}/executions][%d] triggerRetentionExecutionForbidden  %+v", 403, o.Payload)
}
func (o *TriggerRetentionExecutionForbidden) GetPayload() *models.Errors {
	return o.Payload
}

func (o *TriggerRetentionExecutionForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewTriggerRetentionExecutionInternalServerError creates a TriggerRetentionExecutionInternalServerError with default headers values
func NewTriggerRetentionExecutionInternalServerError() *TriggerRetentionExecutionInternalServerError {
	return &TriggerRetentionExecutionInternalServerError{}
}

/* TriggerRetentionExecutionInternalServerError describes a response with status code 500, with default header values.

Internal server error
*/
type TriggerRetentionExecutionInternalServerError struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *TriggerRetentionExecutionInternalServerError) Error() string {
	return fmt.Sprintf("[POST /retentions/{id}/executions][%d] triggerRetentionExecutionInternalServerError  %+v", 500, o.Payload)
}
func (o *TriggerRetentionExecutionInternalServerError) GetPayload() *models.Errors {
	return o.Payload
}

func (o *TriggerRetentionExecutionInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

/*TriggerRetentionExecutionBody trigger retention execution body
swagger:model TriggerRetentionExecutionBody
*/
type TriggerRetentionExecutionBody struct {

	// dry run
	DryRun bool `json:"dry_run,omitempty"`
}

// Validate validates this trigger retention execution body
func (o *TriggerRetentionExecutionBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this trigger retention execution body based on context it is used
func (o *TriggerRetentionExecutionBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *TriggerRetentionExecutionBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *TriggerRetentionExecutionBody) UnmarshalBinary(b []byte) error {
	var res TriggerRetentionExecutionBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
