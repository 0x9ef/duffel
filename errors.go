// Copyright 2021-present Airheart, Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package duffel

import "fmt"

type ErrorType string

// ErrorCode represents the error code returned by the API.
type ErrorCode string

const (
	AuthenticationError ErrorType = "authentication_error"
	AirlineError        ErrorType = "airline_error"
	InvalidStateError   ErrorType = "invalid_state_error"
	RateLimitError      ErrorType = "rate_limit_error"
	ValidationError     ErrorType = "validation_error"
	InvalidRequestError ErrorType = "invalid_request_error"
	ApiError            ErrorType = "api_error"

	// The access token used is not recognized by our system
	AccessTokenNotFound ErrorCode = "access_token_not_found"

	// The airline has responded with an internal error, please contact support
	AirlineInternal ErrorCode = "airline_internal"

	// The airline responded with an unexpected error, please contact support
	AirlineUnknown ErrorCode = "airline_unknown"

	// Requested ancillary service item(s) (e.g. seats) are no longer available, please update your requested services or create a new offer request
	AncillaryServiceNotAvailable ErrorCode = "ancillary_service_not_available"

	// The provided order has already been cancelled
	AlreadyCancelled ErrorCode = "already_cancelled"

	// The request was unacceptable
	BadRequest ErrorCode = "bad_request"

	// A booking with the same details was already found for the selected itinerary, please select another offer
	DuplicateBooking ErrorCode = "duplicate_booking"

	// The order cannot contain more than one passenger with with the same name
	DuplicatePassengerName ErrorCode = "duplicate_passenger_name"

	// The provided access token has expired
	ExpiredAccessToken ErrorCode = "expired_access_token"

	// There wasn't enough balance in the wallet for the operation - for example, you booked a flight for £300 with only £200 available in the wallet
	InsufficientBalance ErrorCode = "insufficient_balance"

	// The provided token doesn't have sufficient permissions to perform the requested action
	InsufficientPermissions ErrorCode = "insufficient_permissions"

	// There was something wrong on our end, please contact support
	InternalServerError ErrorCode = "internal_server_error"

	// The Authorization header must conform to the following format: Bearer API_TOKEN
	InvalidAuthorizationHeader ErrorCode = "invalid_authorization_header"

	// The Content-Type should be set to application/json
	InvalidContentTypeHeader ErrorCode = "invalid_content_type_header"

	// The data in the request body should be a JSON object
	InvalidDataParam ErrorCode = "invalid_data_param"

	// The airline does not support the format of the email address provided
	InvalidEmailAddress ErrorCode = "invalid_email_address"

	// The request had an invalid combination of fields
	InvalidFieldsSet ErrorCode = "invalid_field_sets"

	// The airline did not recognise the loyalty programme account details for one or more of the passengers
	InvalidLoyaltyCard ErrorCode = "invalid_loyalty_card"

	// The title of one of the passengers is not valid
	InvalidPassengerTitle ErrorCode = "invalid_passenger_title"

	// The phone number is not valid
	InvalidPhoneNumber ErrorCode = "invalid_phone_number"

	// The Duffel-Version header must be a known version of our API as indicated in our Docs
	InvalidVersionHeader ErrorCode = "invalid_version_header"

	// The data in the request body is not valid
	MalformedDataParam ErrorCode = "malformed_data_param"

	// The Authorization header must be set and contain a valid API token
	MissingAuthorizationHeader ErrorCode = "missing_authorization_header"

	// The Content-Type header needs to be set to application/json
	MissingContentTypeHeader ErrorCode = "missing_content_type_header"

	// The data in the request body should be nested under the data key
	MissingDataParam ErrorCode = "missing_data_param"

	// The Duffel-Version header is required and must be a valid API version
	MissingVersionHeader ErrorCode = "missing_version_header"

	// The resource you are trying to access does not exist
	NotFound ErrorCode = "not_found"

	// The selected offer has already expired
	OfferExpired ErrorCode = "offer_expired"

	// The provided offer is no longer available, please select another offer or create a new offer request to get the latest availability
	OfferNoLongerAvailable ErrorCode = "offer_no_longer_available"

	// An offer from this offer request has already been booked; please perform a new search
	OfferRequestAlreadyBooked ErrorCode = "offer_request_already_booked"

	// The order change has already been actioned and cannot be actioned again
	OrderChangeAlreadyActioned ErrorCode = "order_change_already_actioned"

	// Order creation has already been attempted for the provided offer. You should not retry this request.
	OrderCreationAlreadyAttempted ErrorCode = "order_creation_already_attempted"

	// The request to create an order was not successful. You should not retry this request.
	OrderNotCreated ErrorCode = "order_not_created"

	// The amount provided in the payment does not match the total_amount of the order
	PaymentAmountDoesNotMatchOrderAmount ErrorCode = "payment_amount_does_not_match_order_amount"

	// The currency provided in the payment does not match the total_currency of the order
	PaymentCurrenctDoesNotMatchOrderCurrency ErrorCode = "payment_currency_does_not_match_order_currency"

	// The provided offer is no longer available for the same price, please retrieve the offer again to get the latest pricing information.
	PriceChanged ErrorCode = "price_changed"

	// Too many requests have hit the API too quickly. Please retry your request after the time specified in the ratelimit-reset header returned to you
	RateLimitExceeded ErrorCode = "rate_limit_exceeded"

	// The change you tried to accept is not the latest. Please retry the request with the latest one
	StaleAirlineInitiatedChangeAccept ErrorCode = "stale_airline_initiated_change_accept"

	// The change you tried to update is not the latest. Please retry the request with the latest one
	StaleAirlineInitiatedChangeUpdate ErrorCode = "stale_airline_initiated_change_update"

	// The feature you requested is not available. Please contact help@duffel.com if you are interested in getting access to it
	UnavailableFeature ErrorCode = "unavailable_feature"

	// The resource does not support the following action
	UnsupportedAction ErrorCode = "unsupported_action"

	// The API does not support the format set in the Accept header, please use a supported format
	UnsupportedFormat ErrorCode = "unsupported_format"

	// The version set to the Duffel-Version header is no longer supported by the API, please upgrade
	UnsupportedVersion ErrorCode = "unsupported_version"

	// The credit card number provided is not valid
	ValidationChecksum ErrorCode = "validation_checksum"

	// The field submitted has an invalid format
	ValidationFormat ErrorCode = "validation_format"

	// The field submitted must be one of a fixed set of values
	ValidationInclusion ErrorCode = "validation_inclusion"

	// The length of the submitted field is out of the boundaries for that field
	ValidationLength ErrorCode = "validation_length"

	// The field submitted cannot be blank
	ValidationRequired ErrorCode = "validation_required"

	// The field submitted has an invalid type
	ValidationType ErrorCode = "validation_type"

	// The field submitted must be unique
	ValidationUnique ErrorCode = "validation_unique"
)

// IsErrorCode is a concenience method to check if an error is a specific error code from Duffel.
// This simplifies error handling branches without needing to type cast multiple times in your code.
func IsErrorCode(err error, code ErrorCode) bool {
	if err, ok := err.(*DuffelError); ok {
		return err.IsCode(code)
	}
	return false
}

// IsErrorType is a concenience method to check if an error is a specific error type from Duffel.
// This simplifies error handling branches without needing to type cast multiple times in your code.
func IsErrorType(err error, typ ErrorType) bool {
	if err, ok := err.(*DuffelError); ok {
		return err.IsType(typ)
	}
	return false
}

// RequestIDFromError returns the request ID from the error. Use this when contacting Duffel support
// for non-retryable errors such as `AirlineInternal` or `AirlineUnknown`.
func RequestIDFromError(err error) (string, bool) {
	if err, ok := err.(*DuffelError); ok {
		return err.Meta.RequestID, true
	}
	return "", false
}

// ErrIsRetryable returns true if the request that generated this error is retryable.
func ErrIsRetryable(err error) bool {
	if err, ok := err.(*DuffelError); ok {
		return err.Retryable
	}
	return false
}

type DuffelError struct {
	Meta       ErrorMeta `json:"meta"`
	Errors     []Error   `json:"errors"`
	StatusCode int       `json:"-"`
	Retryable  bool      `json:"-"`
}

func (e *DuffelError) Error() string {
	return fmt.Sprintf("duffel: %s", e.Errors[0].Message)
}

func (e *DuffelError) IsType(t ErrorType) bool {
	for _, err := range e.Errors {
		if err.Type == t {
			return true
		}
	}
	return false
}

func (e *DuffelError) IsCode(t ErrorCode) bool {
	for _, err := range e.Errors {
		if err.Code == t {
			return true
		}
	}
	return false
}

type Error struct {
	Type             ErrorType `json:"type"`
	Title            string    `json:"title"`
	Message          string    `json:"message"`
	DocumentationURL string    `json:"documentation_url"`
	Code             ErrorCode `json:"code"`
}

type ErrorMeta struct {
	Status    int64  `json:"status"`
	RequestID string `json:"request_id"`
}
