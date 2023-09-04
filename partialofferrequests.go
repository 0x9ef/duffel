package duffel

import (
	"context"
	"net/url"
)

const partialOsfferRequestIDPrefix = "prq_"

type (
	PartialOfferRequestClient interface {
		CreatePartialOfferRequest(ctx context.Context, requestInput PartialOfferRequestInput, requestOptions ...RequestOption) (*OfferRequest, error)
		GetPartialOfferRequest(ctx context.Context, partialOfferRequestID string, requestOptions ...RequestOption) (*OfferRequest, error)
		ListPartialOfferRequestFares(ctx context.Context, partialOfferRequestID string, options []ListPartialOfferRequestFaresParams, requestOptions ...RequestOption) (*OfferRequest, error)
	}

	PartialOfferRequestInput struct {
		// The passengers who want to travel. If you specify an age for a passenger,
		// the type may differ for the same passenger in different offers due to airline's different rules.
		// e.g. one airline may treat a 14 year old as an adult, and another as a young adult.
		// You may only specify an age or a type – not both.
		Passengers []OfferRequestPassenger `json:"passengers" url:"-"`
		// The slices that make up this offer request. One-way journeys can be expressed using one slice,
		// whereas return trips will need two.
		Slices []OfferRequestSlice `json:"slices" url:"-"`
		// The cabin that the passengers want to travel in
		CabinClass CabinClass `json:"cabin_class" url:"-"`
		// The maximum number of connections within any slice of the offer.
		// For example 0 means a direct flight which will have a single segment within each slice
		// and 1 means a maximum of two segments within each slice of the offer.
		MaxConnections *int `json:"max_connections,omitempty" url:"-"`
		// The maximum amount of time in milliseconds to wait for each airline to respond
		SupplierTimeout int `json:"-" url:"supplier_timeout,omitempty"`
	}

	ListPartialOfferRequestFaresParams struct {
		SelectedPartialOffer string `url:"selected_partial_offer"`
	}
)

func (a *API) CreatePartialOfferRequest(ctx context.Context, requestInput PartialOfferRequestInput, requestOptions ...RequestOption) (*OfferRequest, error) {
	return newRequestWithAPI[PartialOfferRequestInput, OfferRequest](a).
		Post("/air/partial_offer_requests", &requestInput).
		WithParams(requestInput).
		WithOptions(requestOptions...).
		Single(ctx)
}

func (a *API) GetPartialOfferRequest(ctx context.Context, partialOfferRequestID string, requestOptions ...RequestOption) (*OfferRequest, error) {
	if err := validateID(partialOfferRequestID, partialOsfferRequestIDPrefix); err != nil {
		return nil, err
	}

	return newRequestWithAPI[EmptyPayload, OfferRequest](a).
		Getf("/air/partial_offer_requests/%s", partialOfferRequestID).
		WithOptions(requestOptions...).
		Single(ctx)
}

func (a *API) ListPartialOfferRequestFares(ctx context.Context, partialOfferRequestID string, options []ListPartialOfferRequestFaresParams, requestOptions ...RequestOption) (*OfferRequest, error) {
	if err := validateID(partialOfferRequestID, partialOsfferRequestIDPrefix); err != nil {
		return nil, err
	}

	return newRequestWithAPI[ListPartialOfferRequestFaresParams, OfferRequest](a).
		Getf("/air/partial_offer_requests/%s", partialOfferRequestID).
		WithParams(normalizeParams(options)...).
		WithOptions(requestOptions...).
		Single(ctx)
}

var _ PartialOfferRequestClient = (*API)(nil)

func (o ListPartialOfferRequestFaresParams) Encode(q url.Values) error {
	q.Set("selected_partial_offer", o.SelectedPartialOffer)
	return nil
}

// Encode implements the ParamEncoder interface.
func (o PartialOfferRequestInput) Encode(q url.Values) error {
	return nil
}
