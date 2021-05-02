// Code generated by go-swagger; DO NOT EDIT.

package preheat

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

// ListExecutionsReader is a Reader for the ListExecutions structure.
type ListExecutionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListExecutionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListExecutionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewListExecutionsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewListExecutionsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewListExecutionsForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewListExecutionsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewListExecutionsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewListExecutionsOK creates a ListExecutionsOK with default headers values
func NewListExecutionsOK() *ListExecutionsOK {
	return &ListExecutionsOK{}
}

/* ListExecutionsOK describes a response with status code 200, with default header values.

List executions success
*/
type ListExecutionsOK struct {

	/* Link refers to the previous page and next page
	 */
	Link string

	/* The total count of executions
	 */
	XTotalCount int64

	Payload []*models.Execution
}

func (o *ListExecutionsOK) Error() string {
	return fmt.Sprintf("[GET /projects/{project_name}/preheat/policies/{preheat_policy_name}/executions][%d] listExecutionsOK  %+v", 200, o.Payload)
}
func (o *ListExecutionsOK) GetPayload() []*models.Execution {
	return o.Payload
}

func (o *ListExecutionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewListExecutionsBadRequest creates a ListExecutionsBadRequest with default headers values
func NewListExecutionsBadRequest() *ListExecutionsBadRequest {
	return &ListExecutionsBadRequest{}
}

/* ListExecutionsBadRequest describes a response with status code 400, with default header values.

Bad request
*/
type ListExecutionsBadRequest struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *ListExecutionsBadRequest) Error() string {
	return fmt.Sprintf("[GET /projects/{project_name}/preheat/policies/{preheat_policy_name}/executions][%d] listExecutionsBadRequest  %+v", 400, o.Payload)
}
func (o *ListExecutionsBadRequest) GetPayload() *models.Errors {
	return o.Payload
}

func (o *ListExecutionsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewListExecutionsUnauthorized creates a ListExecutionsUnauthorized with default headers values
func NewListExecutionsUnauthorized() *ListExecutionsUnauthorized {
	return &ListExecutionsUnauthorized{}
}

/* ListExecutionsUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type ListExecutionsUnauthorized struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *ListExecutionsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /projects/{project_name}/preheat/policies/{preheat_policy_name}/executions][%d] listExecutionsUnauthorized  %+v", 401, o.Payload)
}
func (o *ListExecutionsUnauthorized) GetPayload() *models.Errors {
	return o.Payload
}

func (o *ListExecutionsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewListExecutionsForbidden creates a ListExecutionsForbidden with default headers values
func NewListExecutionsForbidden() *ListExecutionsForbidden {
	return &ListExecutionsForbidden{}
}

/* ListExecutionsForbidden describes a response with status code 403, with default header values.

Forbidden
*/
type ListExecutionsForbidden struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *ListExecutionsForbidden) Error() string {
	return fmt.Sprintf("[GET /projects/{project_name}/preheat/policies/{preheat_policy_name}/executions][%d] listExecutionsForbidden  %+v", 403, o.Payload)
}
func (o *ListExecutionsForbidden) GetPayload() *models.Errors {
	return o.Payload
}

func (o *ListExecutionsForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewListExecutionsNotFound creates a ListExecutionsNotFound with default headers values
func NewListExecutionsNotFound() *ListExecutionsNotFound {
	return &ListExecutionsNotFound{}
}

/* ListExecutionsNotFound describes a response with status code 404, with default header values.

Not found
*/
type ListExecutionsNotFound struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *ListExecutionsNotFound) Error() string {
	return fmt.Sprintf("[GET /projects/{project_name}/preheat/policies/{preheat_policy_name}/executions][%d] listExecutionsNotFound  %+v", 404, o.Payload)
}
func (o *ListExecutionsNotFound) GetPayload() *models.Errors {
	return o.Payload
}

func (o *ListExecutionsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewListExecutionsInternalServerError creates a ListExecutionsInternalServerError with default headers values
func NewListExecutionsInternalServerError() *ListExecutionsInternalServerError {
	return &ListExecutionsInternalServerError{}
}

/* ListExecutionsInternalServerError describes a response with status code 500, with default header values.

Internal server error
*/
type ListExecutionsInternalServerError struct {

	/* The ID of the corresponding request for the response
	 */
	XRequestID string

	Payload *models.Errors
}

func (o *ListExecutionsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /projects/{project_name}/preheat/policies/{preheat_policy_name}/executions][%d] listExecutionsInternalServerError  %+v", 500, o.Payload)
}
func (o *ListExecutionsInternalServerError) GetPayload() *models.Errors {
	return o.Payload
}

func (o *ListExecutionsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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
