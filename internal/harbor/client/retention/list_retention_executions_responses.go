// Code generated by go-swagger; DO NOT EDIT.

package retention

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/cryptk/kubernetes-mimic/internal/harbor/models"
)

// ListRetentionExecutionsReader is a Reader for the ListRetentionExecutions structure.
type ListRetentionExecutionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListRetentionExecutionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListRetentionExecutionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewListRetentionExecutionsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewListRetentionExecutionsForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewListRetentionExecutionsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewListRetentionExecutionsOK creates a ListRetentionExecutionsOK with default headers values
func NewListRetentionExecutionsOK() *ListRetentionExecutionsOK {
	return &ListRetentionExecutionsOK{}
}

/* ListRetentionExecutionsOK describes a response with status code 200, with default header values.

Get a Retention execution successfully.
*/
type ListRetentionExecutionsOK struct {

	/* Link to previous page and next page
	 */
	Link string

	/* The total count of available items
	 */
	XTotalCount int64

	Payload []*models.RetentionExecution
}

func (o *ListRetentionExecutionsOK) Error() string {
	return fmt.Sprintf("[GET /retentions/{id}/executions][%d] listRetentionExecutionsOK  %+v", 200, o.Payload)
}
func (o *ListRetentionExecutionsOK) GetPayload() []*models.RetentionExecution {
	return o.Payload
}

func (o *ListRetentionExecutionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header Link
	hdrLink := response.GetHeader("Link")

	if hdrLink != "" {
		o.Link = hdrLink
	}

	// hydrates response header X-Total-Count
	hdrXTotalCount := response.GetHeader("X-Total-Count")

	if hdrXTotalCount != "" {
		valxTotalCount, err := swag.ConvertInt64(hdrXTotalCount)
		if err != nil {
			return errors.InvalidType("X-Total-Count", "header", "int64", hdrXTotalCount)
		}
		o.XTotalCount = valxTotalCount
	}

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListRetentionExecutionsUnauthorized creates a ListRetentionExecutionsUnauthorized with default headers values
func NewListRetentionExecutionsUnauthorized() *ListRetentionExecutionsUnauthorized {
	return &ListRetentionExecutionsUnauthorized{}
}

/* ListRetentionExecutionsUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type ListRetentionExecutionsUnauthorized struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *ListRetentionExecutionsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /retentions/{id}/executions][%d] listRetentionExecutionsUnauthorized  %+v", 401, o.Payload)
}
func (o *ListRetentionExecutionsUnauthorized) GetPayload() *models.Errors {
	return o.Payload
}

func (o *ListRetentionExecutionsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewListRetentionExecutionsForbidden creates a ListRetentionExecutionsForbidden with default headers values
func NewListRetentionExecutionsForbidden() *ListRetentionExecutionsForbidden {
	return &ListRetentionExecutionsForbidden{}
}

/* ListRetentionExecutionsForbidden describes a response with status code 403, with default header values.

Forbidden
*/
type ListRetentionExecutionsForbidden struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *ListRetentionExecutionsForbidden) Error() string {
	return fmt.Sprintf("[GET /retentions/{id}/executions][%d] listRetentionExecutionsForbidden  %+v", 403, o.Payload)
}
func (o *ListRetentionExecutionsForbidden) GetPayload() *models.Errors {
	return o.Payload
}

func (o *ListRetentionExecutionsForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewListRetentionExecutionsInternalServerError creates a ListRetentionExecutionsInternalServerError with default headers values
func NewListRetentionExecutionsInternalServerError() *ListRetentionExecutionsInternalServerError {
	return &ListRetentionExecutionsInternalServerError{}
}

/* ListRetentionExecutionsInternalServerError describes a response with status code 500, with default header values.

Internal server error
*/
type ListRetentionExecutionsInternalServerError struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *ListRetentionExecutionsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /retentions/{id}/executions][%d] listRetentionExecutionsInternalServerError  %+v", 500, o.Payload)
}
func (o *ListRetentionExecutionsInternalServerError) GetPayload() *models.Errors {
	return o.Payload
}

func (o *ListRetentionExecutionsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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
